package tstree

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// walker обходит Go-типы и строит Node'ы, регистрируя именованные в реестре.
type walker struct {
	reg      *Registry
	path     []string // путь для имён анонимных типов
	depth    int
	maxDepth int
}

const maxDepth = 32 // защита от циклов в аномалиях; самоссылки через byType не растут

// Value строит узел по живому значению (пример из коллектора).
func (w *walker) Value(v any) Node {
	if v == nil {
		return Node{Kind: KindPrimitive, Type: "any"}
	}
	return w.typeof(reflect.TypeOf(v), reflect.ValueOf(v))
}

// Type строит узел по типу (без значения).
func (w *walker) Type(t reflect.Type) Node {
	return w.typeof(t, reflect.Value{})
}

func (w *walker) typeof(t reflect.Type, val reflect.Value) Node {
	w.depth++
	defer func() { w.depth-- }()

	if w.depth > maxDepth {
		return Node{Kind: KindPrimitive, Type: "any"}
	}

	// указатель: разыменовываем (опциональность решает поле-родитель)
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	// uuid.UUID и совместимые: [16]byte под капотом, на фронте ждём строку
	switch t.PkgPath() {
	case "github.com/google/uuid":
		if t.Name() == "UUID" {
			return Node{Kind: KindPrimitive, Type: "string", Docs: "uuid"}
		}
	}

	// uuid.UUID и совместимые: [16]byte под капотом, на фронте ждём строку
	switch t.PkgPath() {
	case "github.com/google/uuid":
		if t.Name() == "UUID" {
			return Node{Kind: KindPrimitive, Type: "string", Docs: "uuid"}
		}
	}

	switch t.Kind() {
	case reflect.Struct:
		return w.structNode(t, val)

	case reflect.Map:
		item := w.typeof(t.Elem(), mapValue(val))
		return Node{Kind: KindMap, Name: mapKeyName(t), Item: &item}

	case reflect.Slice, reflect.Array:
		return w.arrayNode(t, val)

	case reflect.Interface:
		// any: используем живой пример, если он есть
		if val.IsValid() && val.Kind() != reflect.Pointer && !val.IsZero() {
			vt := val.Elem().Type()
			for vt.Kind() == reflect.Pointer {
				vt = vt.Elem()
			}
			return w.typeof(vt, val.Elem())
		}
		return Node{Kind: KindPrimitive, Type: "any"}

	default:
		return primitive(t, val)
	}
}

// structNode — именованный объект. Все структуры (в т.ч. анонимные) получают
// имя в реестре и ссылаются через KindRef: дедупликация по reflect.Type.
func (w *walker) structNode(t reflect.Type, val reflect.Value) Node {
	// time.Time и наследники -> строка ISO
	if t == timeType || (t.Name() == "Time" && strings.HasSuffix(t.String(), "time.Time")) {
		return Node{Kind: KindPrimitive, Type: "string", Docs: "ISO datetime"}
	}

	// уже в реестре — просто ссылка (сквозная дедупликация по типу)
	if named, ok := w.reg.byType[t]; ok {
		return Node{Kind: KindRef, Name: named.Name}
	}

	name := w.reg.uniqueName(typeName(t, w))
	node := &Node{Kind: KindObject, Props: map[string]Node{}, Optional: map[string]bool{}}
	named := &Named{Name: name, Node: node}
	w.reg.byType[t] = named
	w.reg.byName[name] = named

	w.fillProps(node, t, val)
	return Node{Kind: KindRef, Name: name}
}

// fillProps заполняет Props структуры, разворачивая встроенные (embedded) поля.
func (w *walker) fillProps(node *Node, t reflect.Type, val reflect.Value) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			// embedded struct: разворачиваем поля в родителя
			var fv reflect.Value
			if val.IsValid() && val.Kind() == reflect.Struct {
				fv = val.Field(i)
			}
			w.fillProps(node, f.Type, fv)
			continue
		}

		// unexported поля json не сериализует — пропускаем
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}

		// unexported поля json не сериализует — пропускаем
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}

		tag := f.Tag.Get("json")
		if tag == "-" {
			continue
		}
		// пользовательский фильтр (например, системные поля ent)
		if w.reg.SkipField != nil && w.reg.SkipField(t, f) {
			continue
		}
		// пользовательский фильтр (например, системные поля ent)
		if w.reg.SkipField != nil && w.reg.SkipField(t, f) {
			continue
		}
		name := strings.Split(tag, ",")[0]
		omit := strings.Contains(tag, "omitempty")
		if name == "" {
			name = f.Name
		}

		var fv reflect.Value
		if val.IsValid() && val.Kind() == reflect.Struct {
			fv = val.Field(i)
		}

		child := w.typeof(f.Type, fv)
		// поле-указатель или omitempty -> optional
		if f.Type.Kind() == reflect.Pointer || omit {
			node.Optional[name] = true
		}
		node.Props[name] = child
	}
}

func (w *walker) arrayNode(t reflect.Type, val reflect.Value) Node {
	// слияние примеров нескольких элементов слайса: берём общий знаменатель.
	// Для слайса структур elemType один и тот же -> один ref, этого достаточно.
	item := w.typeof(t.Elem(), sliceFirst(val))
	return Node{Kind: KindArray, Item: &item}
}

func (w *walker) anonName() string {
	w.reg.anon++
	return fmt.Sprintf("Anon%d", w.reg.anon)
}

// typeName возвращает имя Go-типа или сгенерированное для анонимных.
func typeName(t reflect.Type, w *walker) string {
	if name := t.Name(); name != "" {
		return name
	}
	return w.anonName()
}

//
// helpers
//

func sliceFirst(val reflect.Value) reflect.Value {
	if val.IsValid() && (val.Kind() == reflect.Slice || val.Kind() == reflect.Array) && val.Len() > 0 {
		return val.Index(0)
	}
	return reflect.Value{}
}

func mapValue(val reflect.Value) reflect.Value {
	if val.IsValid() && val.Kind() == reflect.Map {
		for _, k := range val.MapKeys() {
			return val.MapIndex(k)
		}
	}
	return reflect.Value{}
}

// mapKeyName — имя типа ключа map для TS ({ [k: string]: T }).
func mapKeyName(t reflect.Type) string {
	if t.Key().Kind() == reflect.String {
		if k := t.Key().Name(); k != "" && k != "string" {
			return k
		}
	}
	return "string"
}

var timeType = reflect.TypeOf(time.Time{})

// primitive — маппинг Go-примитивов на TS.
func primitive(t reflect.Type, val reflect.Value) Node {
	n := Node{Kind: KindPrimitive}
	switch t.Kind() {
	case reflect.String:
		n.Type = "string"
	case reflect.Bool:
		n.Type = "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n.Type = "number"
	case reflect.Float32, reflect.Float64:
		n.Type = "number"
	default:
		n.Type = "any"
	}
	return n
}
