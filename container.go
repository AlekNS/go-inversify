package inversify

import (
	"reflect"
)

// Container .
type Container interface {
	// Bind .
	Bind(Any) *Binding
	// Unbind .
	Unbind(Any) Container

	// Get .
	Get(Any) (Any, error)
	// IsBound .
	IsBound(Any) bool

	// Rebuild .
	Rebuild()

	// Merge with another container
	Merge(Container) Container
	// SetParent supports for hierarchical DI systems
	SetParent(Container)

	// GetParent .
	GetParent() Container

	// Load .
	Load(*Module) error

	// UnLoad .
	UnLoad(*Module) error

	// Snapshot @TODO
	// Snapshot() Container
}

type optionalBind struct {
	dep Any
}

// Optional .
func Optional(dep Any) Any {
	return optionalBind{dep}
}

type containerDefault struct {
	parent *containerDefault

	factories map[Any]*Binding
}

func (c *containerDefault) Bind(symbol Any) *Binding {
	b := &Binding{}
	asVal := reflect.ValueOf(symbol)
	if asVal.Kind() == reflect.Ptr && asVal.IsNil() {
		symbol = asVal.Interface()
	}
	c.factories[symbol] = b
	return b
}

func (c *containerDefault) Unbind(symbol Any) Container {
	asVal := reflect.ValueOf(symbol)
	if asVal.Kind() == reflect.Ptr && asVal.IsNil() {
		symbol = asVal.Interface()
	}
	delete(c.factories, symbol)
	return c
}

func (c *containerDefault) findFactory(symbol Any) (*Binding, bool) {
	factory, ok := c.factories[symbol]

	if !ok {
		if c.parent != nil {
			return c.parent.findFactory(symbol)
		}

		return nil, false
	}

	return factory, true
}

func (c *containerDefault) Rebuild() {
	resolveContainerDependencies(c)
}

func (c *containerDefault) Get(symbol Any) (Any, error) {
	if c.factories[symbol].resolves == nil {
		c.Rebuild()
	}
	return c.factories[symbol].factory()
}

func (c *containerDefault) IsBound(symbol Any) bool {
	_, ok := c.factories[symbol]
	return ok
}

func (c *containerDefault) Merge(other Container) Container {
	container := newDefaultContainer()

	otherImpl, ok := other.(*containerDefault)
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

func (c *containerDefault) SetParent(parent Container) {
	parentImpl, ok := parent.(*containerDefault)
	if !ok {
		panic("container is not compatible")
	}

	c.parent = parentImpl
}

func (c *containerDefault) GetParent() Container {
	return c.parent
}

func (c *containerDefault) Load(module *Module) error {
	return module.registerCallback(newContainerBinderProxy(c))
}

func (c *containerDefault) UnLoad(module *Module) error {
	return module.unRegisterCallback(newContainerBinderProxy(c))
}

func newDefaultContainer() *containerDefault {
	return &containerDefault{
		factories: make(map[Any]*Binding),
	}
}

// NewContainer .
func NewContainer() Container {
	return newDefaultContainer()
}
