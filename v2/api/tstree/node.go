// Package tstree строит рекурсивное дерево типов на основе reflect.Type.
//
// В отличие от сбора "живых примеров" (Inputs/Outputs), дерево строится из
// самих Go-типов, поэтому:
//   - один и тот же тип, встреченный в разных роутах, получает одно имя
//     и один интерфейс на фронте (а не element1Type/element2Type);
//   - nil-значения и пустые слайсы не ломают вывод типов;
//   - опциональность полей берётся из указателей и omitempty.
//
// Реестр (Registry) — глобальный на весь дамп: создавайте его один раз в
// getDocs и передавайте во все роуты, тогда унификация имён сквозная.
package tstree

// Node — элемент дерева типов.
type Node struct {
	Kind     Kind            `json:"kind,omitempty"`
	Name     string          `json:"name,omitempty"`     // для KindRef: имя интерфейса; для KindMap: имя типа ключа
	Type     string          `json:"type,omitempty"`     // для KindPrimitive: "string" | "number" | "boolean" | "any" | "Date"
	Props    map[string]Node `json:"props,omitempty"`    // для KindObject: поля
	Optional map[string]bool `json:"optional,omitempty"` // какие поля объекта опциональны
	Item     *Node           `json:"item,omitempty"`     // для KindArray / KindMap: тип элемента/значения
	Of       []Node          `json:"of,omitempty"`       // для KindUnion / KindEnum: варианты
	Docs     string          `json:"docs,omitempty"`     // комментарий-описание
}

// Kind — вид узла дерева.
type Kind string

const (
	KindRef       Kind = "ref"       // ссылка на именованный тип (интерфейс из Registry)
	KindObject    Kind = "object"    // объект { ... }
	KindArray     Kind = "array"     // массив Item[]
	KindMap       Kind = "map"       // словарь { [key: string]: Item }
	KindUnion     Kind = "union"     // объединение Of (A | B)
	KindPrimitive Kind = "primitive" // Type: string | number | boolean | any | Date
	KindEnum      Kind = "enum"      // набор строк в Of, на фронте -> string-union
)
