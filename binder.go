package inversify

import "sync"

// Binding holds factory and dependencies
type Binding struct {
	once         sync.Once
	factory      FactoryFunc
	resolves     NAny
	dependencies NAny
}

// To binds to any object that can converted interface{}
func (b *Binding) To(obj Any) *Binding {
	b.factory = func() (Any, error) {
		return obj, nil
	}

	return b
}

// ToFactory binds to abstract function with specified dependencies
func (b *Binding) ToFactory(factoryMethod Any, dependencies ...Any) *Binding {
	return b.toFactoryMethod(wrapAbstractApplyFuncAsSlice(factoryMethod), dependencies...)
}

// ToTypedFactory binds to typed function with specified dependencies
func (b *Binding) ToTypedFactory(factoryMethod Any, dependencies ...Any) *Binding {
	return b.toFactoryMethod(wrapTypedApplyFuncAsSlice(factoryMethod), dependencies...)
}

func (b *Binding) toFactoryMethod(factoryMethod func([]Any) (Any, error), dependencies ...Any) *Binding {
	b.dependencies = dependencies
	emptyDeps := []Any{}
	depsCount := len(dependencies)

	b.factory = func() (Any, error) {
		if depsCount == 0 {
			return factoryMethod(emptyDeps)
		}
		var err error
		resolvedDeps := make([]Any, depsCount, depsCount)
		for inx, dep := range b.resolves {
			resolvedDeps[inx], err = dep.(FactoryFunc)()
			if err != nil {
				return nil, err
			}
		}
		return factoryMethod(resolvedDeps)
	}

	return b
}

// InSingletonScope declares dependency as singleton
func (b *Binding) InSingletonScope() {
	var instance Any
	var err error

	originalFactory := b.factory

	b.factory = func() (Any, error) {
		if instance == nil {
			b.once.Do(func() {
				instance, err = originalFactory()
			})
		}
		return instance, err
	}
}
