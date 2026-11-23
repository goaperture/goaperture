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

// nodeMeta — служебные счётчики слияния живых JSON-примеров (RawMessage и
// инлайн-данные). Не сериализуются: по ним canonNode помечает optional'ами
// ключи, отсутствующие в части примеров (обязательное = есть во всех).
// Указатель — чтобы счётчики жили сквозь копии Node (merge по значению).
type nodeMeta struct {
	counts map[string]int // сколько примеров содержат ключ
	total  int            // сколько примеров слилось в узел
}

// Node — элемент дерева типов.
type Node struct {
	Kind     Kind            `json:"kind,omitempty"`
	Name     string          `json:"name,omitempty"`
	Type     string          `json:"type,omitempty"`
	Props    map[string]Node `json:"props,omitempty"`
	Optional map[string]bool `json:"optional,omitempty"`
	Item     *Node           `json:"item,omitempty"`
	Of       []Node          `json:"of,omitempty"`
	Docs     string          `json:"docs,omitempty"`

	meta *nodeMeta `json:"-"`
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

// newMeta — счётчики примеров для слияния живых JSON-данных.
func newMeta() *nodeMeta {
	return &nodeMeta{counts: map[string]int{}}
}
