package inversify

import (
	"container/list"
	"fmt"
)

// PairDependencyName holds dependency and name
type PairDependencyName struct {
	Symbol Any
	Name   string
}

func (pair PairDependencyName) isEqual(symbol Any, name string) bool {
	if pair.Symbol == symbol && pair.Name == name {
		return true
	}
	return false
}

func (pair PairDependencyName) String() string {
	return fmt.Sprintf("%+v[%s]", pair.Symbol, pair.Name)
}

type stronglyConnectedContext struct {
	factories  map[Any]map[string]*Binding
	index      int
	indexes    map[Any]map[string]int
	lowIndexes map[Any]map[string]int
	stack      *list.List
	scDeps     []NAny
}

func (sc *stronglyConnectedContext) setIndexes(symbol Any, name string, val int) {
	sub, ok := sc.indexes[symbol]
	if !ok {
		sc.indexes[symbol] = make(map[string]int)
		sub = sc.indexes[symbol]
	}
	sub[name] = val
}

func (sc *stronglyConnectedContext) isIndexes(symbol Any, name string) bool {
	sub, ok := sc.indexes[symbol]
	if !ok {
		return false
	}
	_, ok = sub[name]
	return ok
}

func (sc *stronglyConnectedContext) setLowIndexes(symbol Any, name string, val int) {
	sub, ok := sc.lowIndexes[symbol]
	if !ok {
		sc.lowIndexes[symbol] = make(map[string]int)
		sub = sc.lowIndexes[symbol]
	}
	sub[name] = val
}

// getStronglyConnectedDependencyList based on Tarjan algorithm (finding strong connected components).
func getStronglyConnectedDependencyList(factories map[Any]map[string]*Binding) []NAny {
	ctx := &stronglyConnectedContext{
		factories:  factories,
		index:      0,
		indexes:    make(map[Any]map[string]int, len(factories)),
		lowIndexes: make(map[Any]map[string]int, len(factories)),
		stack:      list.New(),
		scDeps:     make([]NAny, 0, len(factories)),
	}

	for rootSymbol, bindings := range factories {
		for rootName := range bindings {
			if _, ok := ctx.indexes[rootSymbol]; !ok {
				if _, ok := ctx.indexes[rootSymbol][rootName]; !ok {
					getStronglyConnectedDependency(ctx, rootSymbol, rootName)
				}
			}
		}
	}

	return ctx.scDeps
}

// @TODO: Use stack instead of recursion
func getStronglyConnectedDependency(ctx *stronglyConnectedContext,
	rootSymbol Any,
	rootName string) {

	ctx.stack.PushBack(PairDependencyName{rootSymbol, rootName})
	ctx.setIndexes(rootSymbol, rootName, ctx.index)
	ctx.setLowIndexes(rootSymbol, rootName, ctx.index)
	ctx.index++

	bindings, ok := ctx.factories[rootSymbol]
	if !ok {
		return
	}
	binding, ok := bindings[rootName]
	if !ok {
		return
	}

	for _, dependency := range binding.dependencies {
		dependency, name, _ := unwrapDependency(dependency)
		// self dependency
		if rootSymbol == dependency && rootName == name {
			ctx.scDeps = append(ctx.scDeps, NAny{PairDependencyName{
				rootSymbol, rootName}, PairDependencyName{dependency, name}})
		}
		// skip optional and hierarchy
		if _, ok := ctx.factories[dependency]; !ok {
			continue
		}
		if _, ok := ctx.factories[dependency][name]; !ok {
			continue
		}
		if !ctx.isIndexes(dependency, name) {
			getStronglyConnectedDependency(ctx, dependency, name)
			ctx.setLowIndexes(rootSymbol, rootName,
				minInt(ctx.lowIndexes[rootSymbol][rootName], ctx.lowIndexes[dependency][name]))
		} else if isListContainsSymbolPair(ctx.stack, dependency, name) {
			ctx.setLowIndexes(rootSymbol, rootName,
				minInt(ctx.lowIndexes[rootSymbol][rootName], ctx.indexes[dependency][name]))
		}
	}

	if ctx.lowIndexes[rootSymbol][rootName] == ctx.indexes[rootSymbol][rootName] {
		scdLocal := NAny{}

		for {
			w := ctx.stack.Back()
			ctx.stack.Remove(w)
			scdLocal = append(scdLocal, w.Value)
			pairValue := w.Value.(PairDependencyName)
			if pairValue.isEqual(rootSymbol, rootName) {
				break
			}
		}

		ctx.scDeps = append(ctx.scDeps, scdLocal)
	}
}

func isListContainsSymbolPair(l *list.List, symbol Any, name string) bool {
	elem := l.Back()
	for elem != nil {
		if elem.Value != nil {
			if elem.Value.(PairDependencyName).isEqual(symbol, name) {
				return true
			}
		}
		elem = elem.Prev()
	}
	return false
}
