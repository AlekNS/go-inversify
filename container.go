package inversify

import (
	"fmt"
	"reflect"
)

// Container .
type Container interface {
	// Bind .
	Bind(Any, ...string) *Binding
	// Rebind .
	Rebind(Any, ...string) *Binding
	// Unbind .
	Unbind(Any, ...string) Container

	// Get .
	Get(Any, ...string) (Any, error)
	// IsBound .
	IsBound(Any, ...string) bool

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

type namedBind struct {
	dependency Any
	name       string
}

// Named .
func Named(dep Any, name string) Any {
	if _, isOptional := dep.(optionalBind); isOptional {
		panic(fmt.Sprintf("Optional binding couldn't be embedded into named"))
	}
	return namedBind{dep, name}
}

type containerDefault struct {
	parent *containerDefault

	factories map[Any]map[string]*Binding
}

func reflectInterfacePointers(symbol Any) Any {
	asValue := reflect.ValueOf(symbol)
	if asValue.Kind() == reflect.Ptr && asValue.IsNil() {
		symbol = asValue.Interface()
	}
	return symbol
}

func (c *containerDefault) Bind(symbol Any, names ...string) *Binding {
	return c.bindInternal(false, symbol, names...)
}

func (c *containerDefault) Rebind(symbol Any, names ...string) *Binding {
	return c.bindInternal(true, symbol, names...)
}

func (c *containerDefault) bindInternal(isRebinding bool, symbol Any, names ...string) *Binding {
	name := getFirstStringArgumentOrEmpty(names)
	isBindingExists := c.IsBound(symbol, name)

	if isRebinding && !isBindingExists {
		panic(fmt.Sprintf(`binding "%+v[%s]" is not exists for re-declaration`, symbol, name))
	} else if !isRebinding && isBindingExists {
		panic(fmt.Sprintf(`binding "%+v[%s]" is already registered, use Rebind to replace binding`, symbol, name))
	}

	binding := &Binding{}
	bindings, exists := c.factories[reflectInterfacePointers(symbol)]
	if !exists {
		bindings = make(map[string]*Binding)
		c.factories[reflectInterfacePointers(symbol)] = bindings
	}

	bindings[name] = binding
	return binding
}

func (c *containerDefault) Unbind(symbol Any, names ...string) Container {
	bindings, exists := c.factories[reflectInterfacePointers(symbol)]
	if exists {
		name := getFirstStringArgumentOrEmpty(names)
		delete(bindings, name)
		if len(bindings) == 0 {
			delete(c.factories, reflectInterfacePointers(symbol))
		}
	}
	// else panic!?
	return c
}

func (c *containerDefault) findFactory(symbol Any, name string) (*Binding, bool) {
	if !c.hasFactory(symbol, name) {
		if c.parent != nil {
			return c.parent.findFactory(symbol, name)
		}

		return nil, false
	}

	return c.factories[symbol][name], true
}

func (c *containerDefault) Build() {
	err := resolveContainerDependencies(c)
	if err != nil {
		panic(err.Error())
	}
}

func (c *containerDefault) Get(symbol Any, names ...string) (Any, error) {
	if c.parent == nil {
		return c.factories[symbol][getFirstStringArgumentOrEmpty(names)].factory()
	}
	bining, _ := c.findFactory(symbol, getFirstStringArgumentOrEmpty(names))
	return bining.factory()
}

func (c *containerDefault) IsBound(symbol Any, names ...string) bool {
	_, ok := c.findFactory(symbol, getFirstStringArgumentOrEmpty(names))
	return ok
}

func (c *containerDefault) hasFactory(symbol Any, names ...string) bool {
	_, ok := c.factories[symbol]
	if !ok {
		return false
	}
	_, ok = c.factories[symbol][getFirstStringArgumentOrEmpty(names)]
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
		factories: make(map[Any]map[string]*Binding),
	}
}

// NewContainer .
func NewContainer() Container {
	return newDefaultContainer()
}
