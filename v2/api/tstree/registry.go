package tstree

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Registry — реестр именованных типов на весь дамп.
// Ключ дедупликации — reflect.Type, а не форма JSON-примера:
// один и тот же Go-тип, встреченный в разных роутах и полях,
// получит одно имя и один интерфейс на фронте.
type Registry struct {
	byType map[reflect.Type]*Named
	byName map[string]*Named
	anon   int

	// SkipField — опциональный хук: вернуть true, чтобы исключить поле
	// из дерева (например, системные поля ent). t — тип-владелец, f — поле.
	SkipField func(t reflect.Type, f reflect.StructField) bool
}

// Root — входная точка одного роута: input или output.
type Root struct {
	Name string
	Node *Node
}

// Named — именованный тип (единый на все встречи). Node вида object
// рендерится на фронте как interface, остальные — как type alias.
type Named struct {
	Name string `json:"name"`
	Node *Node  `json:"node"`
}

// NewRegistry создаёт реестр. Один на весь getDocs.
func NewRegistry() *Registry {
	return &Registry{
		byType: map[reflect.Type]*Named{},
		byName: map[string]*Named{},
	}
}

// Declare резервирует корневые имена роута ДО обхода полей: защищает от
// коллизий (вложенный Go-тип с именем XxInput) и фиксирует порядок.
func (r *Registry) Declare(roots []Root) *RouteScope {
	for i := range roots {
		roots[i].Name = r.uniqueName(roots[i].Name)
		r.byName[roots[i].Name] = &Named{Name: roots[i].Name}
	}
	return &RouteScope{reg: r, roots: roots}
}

// Named возвращает все именованные типы (заполненные), отсортированные по имени.
func (r *Registry) Named() []Named {
	out := make([]Named, 0, len(r.byName))
	for _, n := range r.byName {
		if n.Node != nil {
			out = append(out, *n)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (r *Registry) uniqueName(name string) string {
	name = sanitizeName(name)
	if name == "" {
		r.anon++
		name = fmt.Sprintf("Anon%d", r.anon)
	}
	if _, taken := r.byName[name]; !taken {
		return name
	}
	base := name
	for i := 2; ; i++ {
		name = fmt.Sprintf("%s%d", base, i)
		if _, taken := r.byName[name]; !taken {
			return name
		}
	}
}

// sanitizeName приводит имя Go-типа к валидному TS-идентификатору:
// инстанцированные дженерики ("Group[*app/ent.LeadTime]") и прочий мусор
// (скобки, точки, звёздочки) вырезаются.
func sanitizeName(name string) string {
	var b strings.Builder
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			b.WriteRune(c)
		}
	}
	return b.String()
}

// RouteScope — регистратор одного роута.
type RouteScope struct {
	reg   *Registry
	roots []Root
}

// Roots возвращает корни роута (узлы заполняются при обходе).
func (s *RouteScope) Roots() []Root {
	return s.roots
}
