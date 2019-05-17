package inversify

import (
	"container/list"
)

// getStronglyConnectedDependencyList based on Tarjan algorithm (finding strong connected components).
func getStronglyConnectedDependencyList(factories map[Any]*Binding) []NAny {
	index := 0
	indexes := make(map[Any]int, len(factories))
	lowIndexes := make(map[Any]int, len(factories))
	stack := list.New()
	scDeps := make([]NAny, 0, len(factories))

	for symbol, binding := range factories {
		if _, ok := indexes[symbol]; !ok {
			getStronglyConnectedDependency(factories, symbol, binding, &index, indexes, lowIndexes, stack, &scDeps)
		}
	}

	return scDeps
}

// @TODO: Use stack instead of recursion
func getStronglyConnectedDependency(factories map[Any]*Binding,
	symbol Any,
	binding *Binding,
	index *int,
	indexes map[Any]int,
	lowIndexes map[Any]int,
	stack *list.List,
	scDeps *[]NAny) {

	// it's optional binding
	if binding == nil {
		return
	}

	stack.PushBack(symbol)
	indexes[symbol] = *index
	lowIndexes[symbol] = *index
	*index++

	for _, dependency := range binding.dependencies {
		if symbol == dependency {
			*scDeps = append(*scDeps, NAny{symbol, dependency})
		}
		if _, ok := indexes[dependency]; !ok {
			getStronglyConnectedDependency(factories,
				dependency, factories[dependency], index, indexes, lowIndexes, stack, scDeps)
			lowIndexes[symbol] = minInt(lowIndexes[symbol], lowIndexes[dependency])
		} else if listFindValue(stack, dependency) != nil {
			lowIndexes[symbol] = minInt(lowIndexes[symbol], indexes[dependency])
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
