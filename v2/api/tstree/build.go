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

	// ExpectArray оборачивает сэмпл в []any: уточняем до конкретного типа
	// элемента, иначе корень становится any[] (элемент-интерфейс без значения)
	t = refineAnyWrapper(t, samples)

	// именованный корень: если тип уже встречался — мержим текущие сэмплы
	// в существующий узел (порядок роутов случаен, данные могут быть только здесь)
	if named, ok := w.reg.byType[t]; ok {
		if named.Node != nil && t.Kind() == reflect.Struct && len(samples) > 0 {
			w.fillingGuard(t, func() {
				for _, s := range samples {
					if s == nil {
						continue
					}
					w.fillProps(named.Node, t, reflect.ValueOf(s))
				}
			})
		}
		return Node{Kind: KindRef, Name: named.Name}
	}

	// узел — по реальной форме типа (структура, map, слайс, примитив)
	var node *Node
	sample := firstValue(samples)
	switch t.Kind() {
	case reflect.Struct:
		node = &Node{Kind: KindObject, Props: map[string]Node{}, Optional: map[string]bool{}}
		defer w.pushPath(name)() // путь анонимных типов начинается с имени корня
		// мержим ВСЕ сэмплы: разные Execute могут наполнять разные части типа
		for _, s := range samples {
			if s == nil {
				continue
			}
			w.fillProps(node, t, reflect.ValueOf(s))
		}
	case reflect.Map:
		item := w.typeof(t.Elem(), mapValue(sample))
		node = &Node{Kind: KindMap, Name: mapKeyName(t), Item: &item}
	case reflect.Slice, reflect.Array:
		item := w.typeof(t.Elem(), sliceFirst(sample))
		node = &Node{Kind: KindArray, Item: &item}
	default:
		pn := primitive(t, reflect.Value{})
		node = &pn
	}

	named, ok := w.reg.byType[t]
	if !ok {
		named = &Named{Name: name}
		w.reg.byType[t] = named
		w.reg.byName[name] = named
	}

	if named.Node == nil {
		named.Node = node
	} else {
		// тип уже встречался: мержим построенный узел в существующий
		mergeNode(named.Node, node)
	}

	return Node{Kind: KindRef, Name: named.Name}
}

// refineAnyWrapper уточняет тип-обёртку ([]any, map[string]any) до конкретного
// типа из живого примера. ExpectArray кладёт сэмпл в []any, а вывод по
// пустому интерфейсу без значения даёт any — без уточнения корень any[].
func refineAnyWrapper(t reflect.Type, samples []any) reflect.Type {
	if t == nil || len(samples) == 0 {
		return t
	}

	// слайс с элементом-пустым интерфейсом: тип первого элемента примера
	if t.Kind() == reflect.Slice && isEmptyInterface(t.Elem()) {
		if et := sampleElemType(sliceFirst(firstValue(samples))); et != nil {
			return reflect.SliceOf(et)
		}
		return t
	}

	// map со значением-пустым интерфейсом
	if t.Kind() == reflect.Map && isEmptyInterface(t.Elem()) {
		if et := sampleElemType(mapValue(firstValue(samples))); et != nil {
			return reflect.MapOf(t.Key(), et)
		}
	}

	return t
}

// sampleElemType — конкретный тип значения (интерфейс и указатель разворачиваются).
func sampleElemType(val reflect.Value) reflect.Type {
	if !val.IsValid() {
		return nil
	}
	if val.Kind() == reflect.Interface {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	for val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if !val.IsValid() || isEmptyInterface(val.Type()) {
		return nil
	}
	return val.Type()
}

func isEmptyInterface(t reflect.Type) bool {
	return t != nil && t.Kind() == reflect.Interface && t.NumMethod() == 0
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
