package tstree

import "reflect"

// Result — результат построения дерева для одного роута.
type Result struct {
	// Input/Output — корневые узлы. Для структур это KindRef с именем
	// из реестра; сам интерфейс лежит в Registry.Named().
	Input  Node `json:"input,omitempty"`
	Output Node `json:"output,omitempty"`
}

// Build строит дерево для корней роута.
//
// inputs/outputs — живые примеры из коллектора ([]any): по ним извлекается
// reflect.Type; если примеров нет, берутся типы из generic-параметров.
// rootInput/rootOutput — имена корневых интерфейсов (обычно XxInput/XxOutput).
//
// reg один на весь дамп — дедупликация типов сквозная по всем роутам.
func Build(
	reg *Registry,
	rootInput, rootOutput string,
	inputType, outputType reflect.Type,
	inputs, outputs []any,
) Result {
	// резервируем корневые имена до обхода полей
	scope := reg.Declare([]Root{
		{Name: rootInput},
		{Name: rootOutput},
	})

	w := &walker{reg: reg}

	res := Result{}
	res.Input = w.rootNode(scope.roots[0].Name, inputType, inputs)
	res.Output = w.rootNode(scope.roots[1].Name, outputType, outputs)
	return res
}

// rootNode строит корневой узел: именованный (ref) либо анонимный инлайн.
func (w *walker) rootNode(name string, t reflect.Type, samples []any) Node {
	// вытаскиваем тип из примеров, если generic-тип не передан (interface{})
	if t == nil || (t.Kind() == reflect.Interface && t.NumMethod() == 0) {
		for _, s := range samples {
			if s != nil {
				st := reflect.TypeOf(s)
				for st.Kind() == reflect.Pointer {
					st = st.Elem()
				}
				if st.Kind() != reflect.Interface {
					t = st
					break
				}
			}
		}
	}

	if t == nil {
		return Node{}
	}

	// именованный корень: регистрируем под зарезервированным именем
	if named, ok := w.reg.byType[t]; ok {
		return Node{Kind: KindRef, Name: named.Name}
	}

	// узел — по реальной форме типа (структура, map, слайс, примитив)
	var node *Node
	switch t.Kind() {
	case reflect.Struct:
		node = &Node{Kind: KindObject, Props: map[string]Node{}, Optional: map[string]bool{}}
		w.fillProps(node, t, firstValue(samples))
	case reflect.Map:
		item := w.typeof(t.Elem(), reflect.Value{})
		node = &Node{Kind: KindMap, Name: mapKeyName(t), Item: &item}
	case reflect.Slice, reflect.Array:
		item := w.typeof(t.Elem(), reflect.Value{})
		node = &Node{Kind: KindArray, Item: &item}
	default:
		pn := primitive(t, reflect.Value{})
		node = &pn
	}

	named := &Named{Name: name, Node: node}
	w.reg.byType[t] = named
	w.reg.byName[name] = named

	return Node{Kind: KindRef, Name: name}
}

func firstValue(samples []any) reflect.Value {
	for _, s := range samples {
		if s != nil {
			return reflect.ValueOf(s)
		}
	}
	return reflect.Value{}
}

func isPrimitiveKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String:
		return true
	}
	return false
}
