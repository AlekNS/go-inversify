package inversify

import (
	"container/list"
)

// getStrongConnectedDependencyList based on Tarjan algorithm (finding strong connected components).
func getStrongConnectedDependencyList(factories map[Any]*Binding) []NAny {
	index := 0
	indexes := make(map[Any]int, len(factories))
	lowIndexes := make(map[Any]int, len(factories))
	stack := list.New()
	scDeps := make([]NAny, 0, len(factories))

	for symbol, bind := range factories {
		if _, ok := indexes[symbol]; !ok {
			getStrongConnectedDependency(factories, symbol, bind, &index, indexes, lowIndexes, stack, &scDeps)
		}
	}

	return scDeps
}

// @TODO: Use stack instead of recursion
func getStrongConnectedDependency(factories map[Any]*Binding,
	symbol Any,
	bind *Binding,
	index *int,
	indexes map[Any]int,
	lowIndexes map[Any]int,
	stack *list.List,
	scDeps *[]NAny) {

	stack.PushBack(symbol)
	indexes[symbol] = *index
	lowIndexes[symbol] = *index
	*index++

	for _, dep := range bind.dependencies {
		if symbol == dep {
			*scDeps = append(*scDeps, NAny{symbol, dep})
		}
		if _, ok := indexes[dep]; !ok {
			getStrongConnectedDependency(factories, dep, factories[dep], index, indexes, lowIndexes, stack, scDeps)
			lowIndexes[symbol] = minInt(lowIndexes[symbol], lowIndexes[dep])
		} else if listFindValue(stack, dep) != nil {
			lowIndexes[symbol] = minInt(lowIndexes[symbol], indexes[dep])
		}
	}

	if lowIndexes[symbol] == indexes[symbol] {
		scdLocal := NAny{}

		for {
			w := stack.Back()
			stack.Remove(w)
			scdLocal = append(scdLocal, w.Value)
			if symbol == w.Value {
				break
			}
		}
		*scDeps = append(*scDeps, scdLocal)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func listFindValue(l *list.List, value interface{}) *list.Element {
	elem := l.Back()
	for elem != nil {
		if elem.Value == value {
			return elem
		}
		elem = elem.Prev()
	}
	return nil
}
