package inversify

import (
	"reflect"
)

// IContainer .
type IContainer interface {
	// Bind .
	Bind(Any) *Binding
	// BindByType .
	BindByType(reflect.Type) *Binding

	// Get .
	Get(Any) (Any, error)
	// GetByType .
	GetByType(reflect.Type) (Any, error)
	// IsBound .
	IsBound(Any) bool

	// Build .
	Build()

	// Merge with another container
	Merge(IContainer) IContainer
	// SetParent supports for hierarchical DI systems
	SetParent(IContainer)

	// GetParent .
	GetParent() IContainer

	// Snapshot .
	// Snapshot() IContainer
}

type optionalBind struct {
	dep Any
}

// Optional .
func Optional(dep Any) Any {
	return optionalBind{dep}
}

type containerImpl struct {
	parent *containerImpl

	factories map[Any]*Binding
}

func (c *containerImpl) Bind(symbol Any) *Binding {
	b := &Binding{}
	asVal := reflect.ValueOf(symbol)
	if asVal.Kind() == reflect.Ptr && asVal.IsNil() {
		symbol = asVal.Interface()
	}
	c.factories[symbol] = b
	return b
}

func (c *containerImpl) BindByType(reflectedType reflect.Type) *Binding {
	b := &Binding{}
	c.factories[reflectedType] = b
	return b
}

func (c *containerImpl) findFactory(symbol Any) (*Binding, bool) {
	factory, ok := c.factories[symbol]

	if !ok {
		if c.parent != nil {
			return c.parent.findFactory(symbol)
		}

		return nil, false
	}

	return factory, true
}

func (c *containerImpl) Build() {
	buildContainerImpl(c)
}

func (c *containerImpl) Get(symbol Any) (Any, error) {
	if c.factories[symbol].resolves == nil {
		c.Build()
	}
	return c.factories[symbol].factory()
}

func (c *containerImpl) GetByType(reflectedType reflect.Type) (Any, error) {
	return c.Get(reflectedType)
}

func (c *containerImpl) IsBound(symbol Any) bool {
	_, ok := c.factories[symbol]
	return ok
}

func (c *containerImpl) Merge(other IContainer) IContainer {
	container := newContainer()

	otherImpl, ok := other.(*containerImpl)
	if !ok {
		panic("container is not compatible")
	}

	for symbol, factory := range c.factories {
		container.factories[symbol] = factory
	}

	for symbol, factory := range otherImpl.factories {
		container.factories[symbol] = factory
	}

	return container
}

func (c *containerImpl) SetParent(parent IContainer) {
	parentImpl, ok := parent.(*containerImpl)
	if !ok {
		panic("container is not compatible")
	}

	c.parent = parentImpl
}

func (c *containerImpl) GetParent() IContainer {
	return c.parent
}

func newContainer() *containerImpl {
	return &containerImpl{
		factories: make(map[Any]*Binding),
	}
}

func Container() IContainer {
	return newContainer()
}
