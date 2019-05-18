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

	// Build .
	Build()

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
	dependency Any
}

// Optional .
func Optional(dep Any) Any {
	return optionalBind{dep}
}

type containerDefault struct {
	parent *containerDefault

	factories map[Any]*Binding
}

func reflectInterfacePointers(symbol Any) Any {
	asValue := reflect.ValueOf(symbol)
	if asValue.Kind() == reflect.Ptr && asValue.IsNil() {
		symbol = asValue.Interface()
	}
	return symbol
}

func (c *containerDefault) Bind(symbol Any) *Binding {
	binding := &Binding{}
	c.factories[reflectInterfacePointers(symbol)] = binding
	return binding
}

func (c *containerDefault) Unbind(symbol Any) Container {
	delete(c.factories, reflectInterfacePointers(symbol))
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

func (c *containerDefault) Build() {
	err := resolveContainerDependencies(c)
	if err != nil {
		panic(err.Error())
	}
}

func (c *containerDefault) Get(symbol Any) (Any, error) {
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
