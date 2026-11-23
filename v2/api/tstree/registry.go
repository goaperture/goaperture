package tstree

import (
	"encoding/json"
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
// Узлы канонизируются: props и union-варианты сортируются, чтобы результат
// не зависел от случайного порядка обхода map и сэмплов.
func (r *Registry) Named() []Named {
	out := make([]Named, 0, len(r.byName))
	for _, n := range r.byName {
		if n.Node != nil {
			canonNode(n.Node, map[*Node]bool{})
			out = append(out, *n)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// ReserveName резервирует имя типа до первого обхода (структурный тег `ts:"..."`).
// Вызывать ДО обхода типа: structNode увидит тип в byType и возьмёт
// зарезервированное имя, все ref-узлы сразу получат верное имя.
func (r *Registry) ReserveName(t reflect.Type, name string) {
	if t == nil || name == "" {
		return
	}
	if _, exists := r.byType[t]; exists {
		return // тип уже обойдён — имя уже зафиксировано
	}
	name = r.uniqueName(name)
	named := &Named{Name: name}
	r.byType[t] = named
	r.byName[name] = named
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
// (скобки, точки, звёздочки) вырезаются. $ разрешён — валиден в TS.
func sanitizeName(name string) string {
	var b strings.Builder
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '$' {
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

// canonNode канонизирует узел: сортирует props и union-варианты.
// Рекурсивно, с защитой от циклов (visited — узлы уже в обработке).
func canonNode(n *Node, visited map[*Node]bool) {
	if n == nil || visited[n] {
		return
	}
	visited[n] = true
	defer delete(visited, n)

	switch n.Kind {
	case KindObject:
		keys := make([]string, 0, len(n.Props))
		for k := range n.Props {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			pn := n.Props[k]
			canonNode(&pn, visited)
			n.Props[k] = pn
		}
		// обязательность по счётчикам примеров: ключ обязателен, только если
		// есть во ВСЕХ слитых примерах (и примеров было больше одного).
		// Один пример = всё обязательно (счётчиков нет — данные не делились)
		if n.meta != nil && n.meta.total > 1 {
			for _, k := range keys {
				if n.meta.counts[k] < n.meta.total {
					if n.Optional == nil {
						n.Optional = map[string]bool{}
					}
					n.Optional[k] = true
				}
			}
		}
	case KindArray, KindMap:
		if n.Item != nil {
			canonNode(n.Item, visited)
		}
	case KindUnion:
		for i := range n.Of {
			canonNode(&n.Of[i], visited)
		}
		// канонический ключ узла: строковое представление для стабильной сортировки
		sort.Slice(n.Of, func(i, j int) bool {
			return unionKey(n.Of[i]) < unionKey(n.Of[j])
		})
	}
}

// unionKey — каноническое строковое представление узла для сортировки union.
func unionKey(n Node) string {
	b, _ := json.Marshal(n)
	return string(b)
}

// WrapName выносит инлайн-узел (union, RawJSON-объект и т.п.) в именованный
// тип (структурный тег `ts:"..."` на не-структурном поле). При повторной
// встрече с тем же тегом узлы мержатся — варианты union сохраняются.
func (r *Registry) WrapName(node Node, name string) Node {
	named, ok := r.byName[name]
	if !ok {
		named = &Named{Name: name, Node: &node}
		r.byName[name] = named
		return Node{Kind: KindRef, Name: name}
	}
	// уже есть: мержим (например, другой сэмпл добавил вариант в union)
	if named.Node != nil {
		m := mergeProp(*named.Node, node)
		named.Node = &m
	}
	return Node{Kind: KindRef, Name: name}
}
