package tstree

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

// mergeProp совмещает два узла одного поля при повторном обходе с другим
// живым примером: объекты мержатся, разные формы (объект vs массив) — union.
func mergeProp(prev, child Node) Node {
	// примитив any уступает конкретике: тип мог заполняться без живого примера
	// (RawMessage пуст / другой роут обошёл раньше) — порядок обхода роутов
	// случаен (map), поэтому any и конкретика должны стабильно сливаться
	if prev.Kind == KindPrimitive && prev.Type == "any" && child.Kind != KindPrimitive {
		return child
	}
	if child.Kind == KindPrimitive && child.Type == "any" {
		return prev // конкретика важнее any
	}
	// union растёт вариантами
	if prev.Kind == KindUnion {
		for i := range prev.Of {
			if prev.Of[i].Kind == child.Kind {
				prev.Of[i] = mergeProp(prev.Of[i], child)
				return prev
			}
		}
		prev.Of = append(prev.Of, child)
		return prev
	}
	// объект+объект: объединяем ключи (разные формы RawJSON — всё равно один интерфейс)
	if prev.Kind == KindObject && child.Kind == KindObject {
		mergeNode(&prev, &child)
		// счётчики: total суммируется, counts ключей child прибавляются
		if prev.meta == nil {
			prev.meta = newMeta()
		}
		if child.meta != nil {
			for k, c := range child.meta.counts {
				prev.meta.counts[k] += c
			}
		}
		prev.meta.total += child.metaTotal()
		return prev
	}
	// массив+массив: мерж элементы
	if prev.Kind == KindArray && child.Kind == KindArray {
		if child.Item != nil {
			if prev.Item == nil {
				prev.Item = child.Item
			} else {
				m := mergeProp(*prev.Item, *child.Item)
				prev.Item = &m
			}
		}
		return prev
	}
	// разные формы (объект vs массив) -> union
	if (prev.Kind == KindObject || prev.Kind == KindArray) && (child.Kind == KindObject || child.Kind == KindArray) {
		return Node{Kind: KindUnion, Of: []Node{prev, child}}
	}
	// ref/примитивы не трогаем: дедупликация по типу сама
	return prev
}

// fillingGuard выполняет fn с защитой от рекурсии самоссылающихся типов.
func (w *walker) fillingGuard(t reflect.Type, fn func()) {
	if w.filling == nil {
		w.filling = map[reflect.Type]bool{}
	}
	if w.filling[t] {
		return // уже заполняется выше по стеку — рекурсия
	}
	w.filling[t] = true
	defer delete(w.filling, t)
	fn()
}

// walker обходит Go-типы и строит Node'ы, регистрируя именованные в реестре.
type walker struct {
	reg      *Registry
	path     []string // путь полей: для имён анонимных типов (RootInput$user.address)
	depth    int
	filling  map[reflect.Type]bool // защита от рекурсии самоссылающихся типов
}

const maxDepth = 32 // защита от циклов в аномалиях; самоссылки через byType не растут

// pushPath/popPath — контекст пути при обходе полей.
func (w *walker) pushPath(part string) func() {
	w.path = append(w.path, part)
	return func() { w.path = w.path[:len(w.path)-1] }
}

// mergeNode мержит объект-поля и элементы массива при повторном обходе типа
// с другим живым примером (разные cl.Execute наполняют разные ключи).
func mergeNode(dst, src *Node) {
	switch dst.Kind {
	case KindObject:
		if src.Kind != KindObject {
			return
		}
		for k, v := range src.Props {
			if prev, ok := dst.Props[k]; ok {
				dst.Props[k] = mergeProp(prev, v)
				continue
			}
			dst.Props[k] = v
		}
		// счётчики примеров: total суммируется, counts прибавляются
		if dst.meta == nil {
			dst.meta = newMeta()
		}
		dst.meta.total += src.metaTotal()
		if src.meta != nil {
			for k, c := range src.meta.counts {
				dst.meta.counts[k] += c
			}
		}
	case KindArray:
		if src.Kind != KindArray || src.Item == nil {
			return
		}
		if dst.Item == nil {
			dst.Item = src.Item
			return
		}
		if dst.Item.Kind == src.Item.Kind && (src.Item.Kind == KindObject || src.Item.Kind == KindArray) {
			mergeNode(dst.Item, src.Item)
		}
	}
}

// metaTotal — число примеров, слитых в узел (0 если меты нет).
func (n *Node) metaTotal() int {
	if n.meta == nil {
		return 0
	}
	return n.meta.total
}

// rawJSONNode декодирует живой пример json.RawMessage и строит узел
// из реальных данных: объект -> props по ключам, массив -> массив, скаляры.
func (w *walker) rawJSONNode(val reflect.Value) (Node, bool) {
	if !val.IsValid() || val.Kind() != reflect.Slice || val.Len() == 0 {
		return Node{}, false
	}
	var decoded any
	if err := json.Unmarshal(val.Bytes(), &decoded); err != nil {
		return Node{}, false
	}
	n := decodedNode(decoded, w)
	if n.Kind == "" {
		return Node{}, false
	}
	return n, true
}

// decodedNode строит узел из декодированного JSON (map/slice/скаляры).
func decodedNode(v any, w *walker) Node {
	switch tv := v.(type) {
	case nil:
		return Node{Kind: KindPrimitive, Type: "any"}
	case map[string]any:
		node := Node{Kind: KindObject, Props: map[string]Node{}, meta: newMeta()}
		keys := make([]string, 0, len(tv))
		for k := range tv {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			w.pushPath(k)
			node.Props[k] = decodedNode(tv[k], w)
			w.path = w.path[:len(w.path)-1]
			node.meta.counts[k] = 1 // ключ присутствует в этом примере
		}
		node.meta.total = 1
		return node
	case []any:
		item := Node{Kind: KindPrimitive, Type: "any"}
		if len(tv) > 0 {
			item = decodedNode(tv[0], w)
		}
		return Node{Kind: KindArray, Item: &item}
	case string:
		return Node{Kind: KindPrimitive, Type: "string"}
	case float64:
		return Node{Kind: KindPrimitive, Type: "number"}
	case bool:
		return Node{Kind: KindPrimitive, Type: "boolean"}
	default:
		return Node{Kind: KindPrimitive, Type: "any"}
	}
}

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

	// uuid.UUID: [16]byte под капотом, на фронте ждём строку
	switch t.PkgPath() {
	case "github.com/google/uuid":
		if t.Name() == "UUID" {
			return Node{Kind: KindPrimitive, Type: "string", Docs: "uuid"}
		}
	}

	// json.RawMessage — сырой JSON: если есть живой пример, декодируем
	// и обходим как обычные данные (реальные ключи попадут в типы).
	if t == rawMessageType {
		if n, ok := w.rawJSONNode(val); ok {
			return n
		}
		return Node{Kind: KindPrimitive, Type: "any", Docs: "raw JSON"}
	}

	// []byte/[]uint8 — base64-строка в JSON
	if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return Node{Kind: KindArray, Item: &Node{Kind: KindPrimitive, Type: "number"}}
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

	// тип уже в реестре (обойдён ранее или зарезервирован через ReserveName —
	// тег `ts:"..."`) — просто ссылка; узел заполняем при первой встрече
	named, ok := w.reg.byType[t]
	if !ok {
		named = &Named{Name: w.reg.uniqueName(typeName(t, w))}
		w.reg.byType[t] = named
		w.reg.byName[named.Name] = named
	}

	if named.Node == nil || !w.filling[t] {
		// тип либо ещё не заполнялся, либо заполнялся по другому сэмплу —
		// повторно проходим и мержим новые ключи из текущего живого примера
		if w.filling == nil {
			w.filling = map[reflect.Type]bool{}
		}
		if w.filling[t] {
			return Node{Kind: KindRef, Name: named.Name} // рекурсия самоссылки
		}
		w.filling[t] = true
		defer delete(w.filling, t)

		if named.Node == nil {
			named.Node = &Node{Kind: KindObject, Props: map[string]Node{}, Optional: map[string]bool{}}
		}
		w.fillProps(named.Node, t, val)
	}

	return Node{Kind: KindRef, Name: named.Name}
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

		tag := f.Tag.Get("json")
		if tag == "-" {
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

		tagOverride := f.Tag.Get("ts")
		// тег `ts:"Name"` на поле-структуру — резервируем ДО обхода, чтобы
		// structNode взял имя из реестра. Для инлайн-узлов (RawMessage и др.)
		// имя назначаем после обхода — дедупликации по типу там нет.
		if tagOverride != "" && f.Type.Kind() == reflect.Struct {
			w.reg.ReserveName(f.Type, tagOverride)
		}

		w.pushPath(name)
		child := w.typeof(f.Type, fv)
		w.path = w.path[:len(w.path)-1] // popPath (не defer: цикл)

		// поле-указатель или omitempty -> optional
		if f.Type.Kind() == reflect.Pointer || omit {
			node.Optional[name] = true
		}

		// RawMessage без примера уже развёрнут по живым данным;
		// при нескольких сэмплах мержим или строим union (объект vs массив)
		if tagOverride != "" && f.Type.Kind() != reflect.Struct {
			// тег `ts:"Name"` на не-структурном поле (RawMessage и т.п.):
			// инлайн-узел выносится в именованный тип, повторные сэмплы мержатся в него
			node.Props[name] = w.reg.WrapName(child, tagOverride)
		} else if prev, ok := node.Props[name]; ok {
			node.Props[name] = mergeProp(prev, child)
		} else {
			node.Props[name] = child
		}
	}
}

// arrayNode — слайс/массив.
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

// anonPathName — осмысленное имя анонимного типа из пути полей:
// поле "user.address" при корне "UsersListInput" -> "UsersListInput$user$address".
// $ — валидный символ TS-идентификатора, отличает сгенерированные имена.
func (w *walker) anonPathName() string {
	if len(w.path) == 0 {
		return w.anonName() // нет пути — fallback на AnonN
	}
	prefix := sanitizeName(w.path[0])
	if prefix == "" {
		prefix = "Anon"
	}
	parts := []string{prefix}
	for _, p := range w.path[1:] {
		if s := sanitizeName(p); s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, "$")
}

// typeName возвращает имя Go-типа или сгенерированное из пути полей для анонимных.
func typeName(t reflect.Type, w *walker) string {
	if name := t.Name(); name != "" {
		return name
	}
	return w.anonPathName()
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
var rawMessageType = reflect.TypeOf(json.RawMessage{})

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
