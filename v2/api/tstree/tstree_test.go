package tstree

import (
	"encoding/json"
	"reflect"
	"testing"
)

// === вспомогательные типы для тестов (имитируют типы реального приложения) ===

type ProductType string
type GlassShape string

type Size struct {
	Size  int `json:"size"`
	Price int `json:"price"`
}

type GlassCard struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Sizes []Size `json:"sizes"`
}

type RawUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// === Build / дедупликация ===

// TestBuildBasic: структура обходит в ref + именованный интерфейс.
func TestBuildBasic(t *testing.T) {
	reg := NewRegistry()

	type Input struct {
		Page int `json:"page"`
	}
	type Output struct {
		Items []string `json:"items"`
	}

	res := Build(reg, "XInput", "XOutput",
		reflect.TypeOf(Input{}), reflect.TypeOf(Output{}), nil, nil)

	if res.Input.Kind != KindRef || res.Input.Name != "XInput" {
		t.Fatalf("input = %+v, want ref XInput", res.Input)
	}
	if res.Output.Kind != KindRef || res.Output.Name != "XOutput" {
		t.Fatalf("output = %+v, want ref XOutput", res.Output)
	}
	if len(reg.Named()) != 2 {
		t.Fatalf("named = %d, want 2", len(reg.Named()))
	}
}

// TestExpectArrayRef: ExpectArray оборачивает сэмпл в []any — корень должен
// уточняться до []GlassCard с ref-элементом, а не превращаться в any[].
func TestExpectArrayRef(t *testing.T) {
	reg := NewRegistry()
	outType := reflect.TypeOf((*interface{ any })(nil)).Elem() // Output interface{ any }

	outputs := []any{[]any{GlassCard{ID: 1, Title: "x", Sizes: []Size{{Size: 1, Price: 1}}}}}
	res := Build(reg, "XInput", "GlassCards", nil, outType, nil, outputs)

	b, _ := json.Marshal(res.Output)
	want := `{"kind":"ref","name":"GlassCards"}`
	if string(b) != want {
		t.Fatalf("output = %s, want %s", b, want)
	}

	named := findNamed(t, reg, "GlassCards")
	if named.Kind != KindArray {
		t.Fatalf("GlassCards kind = %s, want array", named.Kind)
	}
	if named.Item == nil || named.Item.Kind != KindRef || named.Item.Name != "GlassCard" {
		t.Fatalf("item = %+v, want ref GlassCard", named.Item)
	}
}

// TestAnonNames: анонимные структуры получают имена из пути полей.
func TestAnonNames(t *testing.T) {
	reg := NewRegistry()

	type Input struct {
		Filter struct {
			Q string `json:"q"`
		} `json:"filter"`
	}

	Build(reg, "XInput", "XOutput", reflect.TypeOf(Input{}), nil, nil, nil)

	n := findNamed(t, reg, "XInput")
	child := n.Props["filter"]
	if child.Kind != KindRef || child.Name != "XInput$filter" {
		t.Fatalf("filter = %+v, want ref XInput$filter", child)
	}
}

// TestTsTagOnStruct: тег ts на структурном поле переименовывает тип.
func TestTsTagOnStruct(t *testing.T) {
	reg := NewRegistry()

	type Card struct {
		ID int `json:"id"`
	}
	type Order struct {
		Card Card `json:"card" ts:"GlassCard"`
	}

	Build(reg, "OrderInput", "OrderOutput", reflect.TypeOf(Order{}), nil, nil, nil)

	n := findNamed(t, reg, "OrderInput")
	if n.Props["card"].Kind != KindRef || n.Props["card"].Name != "GlassCard" {
		t.Fatalf("card = %+v, want ref GlassCard", n.Props["card"])
	}
	c := findNamed(t, reg, "GlassCard")
	if _, ok := c.Props["id"]; !ok {
		t.Fatal("GlassCard interface missing")
	}
}

// TestTsTagOnRawMessage: тег ts на json.RawMessage выносит инлайн-узел
// в именованный тип, повторные сэмплы мержатся в него.
func TestTsTagOnRawMessage(t *testing.T) {
	reg := NewRegistry()

	type Data struct {
		Element json.RawMessage `json:"element" ts:"RawCalc"`
		Type    string          `json:"type"`
	}

	samples := []any{
		Data{Type: "custom", Element: json.RawMessage(`{"id":"x","name":"n"}`)},
		Data{Type: "glass", Element: json.RawMessage(`{"type":"glass","view":"clear"}`)},
		Data{Type: "triplex", Element: json.RawMessage(`[{"id":"l"}]`)},
		Data{Type: "empty"}, // RawMessage пуст -> any
	}

	Build(reg, "XInput", "XOutput", reflect.TypeOf(Data{}), nil, samples, nil)

	n := findNamed(t, reg, "XInput")
	if n.Props["element"].Kind != KindRef || n.Props["element"].Name != "RawCalc" {
		t.Fatalf("element = %+v, want ref RawCalc", n.Props["element"])
	}

	raw := findNamed(t, reg, "RawCalc")
	// custom+glass слились в один объект, triplex — отдельный вариант-массив
	if raw.Kind != KindUnion || len(raw.Of) != 2 {
		t.Fatalf("RawCalc = %+v, want union of 2", raw)
	}
}

// TestRawMessageWithoutTag: RawMessage с живым примером разворачивается
// в структуру по реальным ключам.
func TestRawMessageWithoutTag(t *testing.T) {
	reg := NewRegistry()

	type Data struct {
		Element json.RawMessage `json:"element"`
	}

	samples := []any{Data{Element: json.RawMessage(`[{"id":"l1","element":{"view":"clear"}}]`)}}
	Build(reg, "XInput", "XOutput", reflect.TypeOf(Data{}), nil, samples, nil)

	n := findNamed(t, reg, "XInput")
	el := n.Props["element"]
	if el.Kind != KindArray || el.Item == nil {
		t.Fatalf("element = %+v, want array", el)
	}
	if _, ok := el.Item.Props["id"]; !ok {
		t.Fatalf("item props = %+v, want id+element", el.Item.Props)
	}
	if _, ok := el.Item.Props["element"]; !ok {
		t.Fatalf("item props = %+v, want id+element", el.Item.Props)
	}
}

// TestRawMessageWithoutSample: RawMessage без примера -> any + docs.
func TestRawMessageWithoutSample(t *testing.T) {
	reg := NewRegistry()

	type Data struct {
		Element json.RawMessage `json:"element"`
	}

	Build(reg, "XInput", "XOutput", reflect.TypeOf(Data{}), nil, nil, nil)

	n := findNamed(t, reg, "XInput")
	el := n.Props["element"]
	if el.Kind != KindPrimitive || el.Type != "any" || el.Docs != "raw JSON" {
		t.Fatalf("element = %+v, want any/raw JSON", el)
	}
}

// TestMergeSamplesUnion: несколько сэмплов одного типа мержатся: одинаковые
// контейнеры сливаются, разные формы (объект vs массив) -> union.
func TestMergeSamplesUnion(t *testing.T) {
	reg := NewRegistry()

	type Data struct {
		Element json.RawMessage `json:"element"`
	}

	samples := []any{
		Data{Element: json.RawMessage(`{"id":"x","name":"n"}`)},
		Data{Element: json.RawMessage(`{"type":"glass"}`)},
		Data{Element: json.RawMessage(`[{"id":"l"}]`)},
	}
	Build(reg, "XInput", "XOutput", reflect.TypeOf(Data{}), nil, samples, nil)

	n := findNamed(t, reg, "XInput")
	el := n.Props["element"]
	if el.Kind != KindUnion {
		t.Fatalf("element kind = %s, want union", el.Kind)
	}
	// сэмплы-объекты слились в один вариант, массив — второй
	if len(el.Of) != 2 {
		t.Fatalf("union len = %d, want 2", len(el.Of))
	}
	obj := el.Of[0]
	if obj.Kind != KindObject {
		obj = el.Of[1]
	}
	for _, k := range []string{"id", "name", "type"} {
		if _, ok := obj.Props[k]; !ok {
			t.Fatalf("merged object missing key %s", k)
		}
	}
}

// TestMergePropPriorityAny: any уступает конкретике при слиянии
// (порядок роутов случаен — результат не должен зависеть от него).
func TestMergePropPriorityAny(t *testing.T) {
	anyFirst := mergeProp(
		Node{Kind: KindPrimitive, Type: "any"},
		Node{Kind: KindRef, Name: "A"},
	)
	if anyFirst.Kind != KindRef {
		t.Fatalf("any+ref = %+v, want ref", anyFirst)
	}
	refFirst := mergeProp(
		Node{Kind: KindRef, Name: "A"},
		Node{Kind: KindPrimitive, Type: "any"},
	)
	if refFirst.Kind != KindRef {
		t.Fatalf("ref+any = %+v, want ref", refFirst)
	}
}

// TestStabilityAcrossPermutations: результат не зависит от порядка сэмплов
// и от того, был ли тип зарегистрирован другим роутом раньше.
func TestStabilityAcrossPermutations(t *testing.T) {
	type Data struct {
		Element json.RawMessage `json:"element"`
		Type    string          `json:"type"`
	}
	samples := []any{
		Data{Type: "custom", Element: json.RawMessage(`{"id":"x","name":"n"}`)},
		Data{Type: "glass", Element: json.RawMessage(`{"type":"glass","view":"clear"}`)},
		Data{Type: "triplex", Element: json.RawMessage(`[{"id":"l"}]`)},
		Data{Type: "empty"},
	}

	var first string
	for iter := 0; iter < 24; iter++ {
		reg := NewRegistry()
		dt := reflect.TypeOf(Data{})
		// в половине случаев тип уже встречался (другой роут)
		if iter%2 == 0 {
			named := &Named{Name: "CalcData"}
			named.Node = &Node{Kind: KindObject, Props: map[string]Node{}, Optional: map[string]bool{}}
			reg.byType[dt] = named
			reg.byName["CalcData"] = named
		}

		shuffled := make([]any, len(samples))
		for i, idx := range []int{iter % 4, (iter + 1) % 4, (iter + 2) % 4, (iter + 3) % 4} {
			shuffled[i] = samples[idx]
		}

		Build(reg, "XInput", "XOutput", dt, nil, shuffled, nil)

		el := reg.byType[dt].Node.Props["element"]
		canonNode(&el, map[*Node]bool{})
		b, _ := json.Marshal(el)
		if first == "" {
			first = string(b)
		} else if string(b) != first {
			t.Fatalf("iter %d UNSTABLE:\n%s\nvs\n%s", iter, first, b)
		}
	}
}

// TestRequiredFromCounts: при слиянии примеров ключ обязателен только если
// есть во всех примерах (иначе optional).
func TestRequiredFromCounts(t *testing.T) {
	reg := NewRegistry()

	type Data struct {
		Element json.RawMessage `json:"element" ts:"RawCalc"`
	}

	samples := []any{
		Data{Element: json.RawMessage(`{"id":"x","name":"n"}`)},
		Data{Element: json.RawMessage(`{"id":"y","services":{}}`)},
	}
	Build(reg, "XInput", "XOutput", reflect.TypeOf(Data{}), nil, samples, nil)

	raw := findNamed(t, reg, "RawCalc")
	if raw.Optional["id"] {
		t.Fatal("id должен быть required (есть во всех примерах)")
	}
	if !raw.Optional["name"] || !raw.Optional["services"] {
		t.Fatalf("name/services должны быть optional: %+v", raw.Optional)
	}
}

// TestExplicitRootAliasForKnownType: если тип уже зарегистрирован другим
// роутом, явное имя корня (Types.Output) всё равно появляется в named
// как алиас-ссылка на реальный тип.
func TestExplicitRootAliasForKnownType(t *testing.T) {
	reg := NewRegistry()
	// другой роут первым обошёл RawUser под дефолтным именем
	Build(reg, "OtherInput", "OtherOutput", nil,
		reflect.TypeOf((*RawUser)(nil)).Elem(), nil, []any{RawUser{}})

	// GetUsers: Types.Output="Users", handler возвращает []*RawUser
	usersType := reflect.TypeOf((*[]*RawUser)(nil)).Elem()
	res := Build(reg, "GetUsersInput", "Users", nil, usersType, nil, []any{[]*RawUser{{ID: 1}}})

	if res.Output.Name != "Users" {
		t.Fatalf("output = %+v, want ref Users", res.Output)
	}
	u := findNamed(t, reg, "Users")
	if u.Kind != KindArray || u.Item == nil || u.Item.Name == "" {
		t.Fatalf("Users = %+v, want array of ref", u)
	}
}

// TestUnionSortStable: канонизация делает порядок union-вариантов и props
// стабильным независимо от порядка сэмплов.
func TestUnionSortStable(t *testing.T) {
	type Data struct {
		Element json.RawMessage `json:"element"`
	}

	s1 := Data{Element: json.RawMessage(`[{"a":1}]`)}
	s2 := Data{Element: json.RawMessage(`{"b":"s"}`)}

	var first string
	for iter := 0; iter < 4; iter++ {
		samples := []any{s1, s2}
		if iter%2 == 1 {
			samples = []any{s2, s1}
		}
		reg := NewRegistry()
		Build(reg, "XInput", "XOutput", reflect.TypeOf(Data{}), nil, samples, nil)
		var got string
		for _, n := range reg.Named() {
			nb, _ := json.Marshal(n.Node)
			got += string(nb)
		}
		if first == "" {
			first = got
		} else if got != first {
			t.Fatalf("iter %d UNSTABLE", iter)
		}
	}
}

// === helpers ===

func findNamed(t *testing.T, reg *Registry, name string) Node {
	t.Helper()
	for _, n := range reg.Named() {
		if n.Name == name {
			if n.Node == nil {
				t.Fatalf("named %s has nil node", name)
			}
			return *n.Node
		}
	}
	t.Fatalf("named %s not found in %v", name, reg.Named())
	return Node{}
}
