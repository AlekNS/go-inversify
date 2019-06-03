package inversify

import (
	"fmt"
	"reflect"
)

// Container holds dependencies graph
type Container interface {
	// Bind declares dependency (make a panic if already binded)
	Bind(Any, ...string) *Binding
	// Rebind re-declares dependency (make a panic if not exists)
	Rebind(Any, ...string) *Binding
	// Unbind removes dependency
	Unbind(Any, ...string) Container

	// Get resolves dependency
	Get(Any, ...string) (Any, error)
	// Get resolves dependency and make a panic if error was produced
	MustGet(Any, ...string) Any
	// IsBound check existences of dependency
	IsBound(Any, ...string) bool

	// Build and resolves dependencies
	Build()

	// Merge with another container and returns new container
	Merge(container Container, name string) Container
	// SetParent supports for hierarchical DI systems
	SetParent(Container)

	// GetParent gets parent of container
	GetParent() Container

	// Load binds module
	Load(*Module) error

	// UnLoad unbinds module
	UnLoad(*Module) error

	// Snapshot @TODO
	// Snapshot() Container
}

type optionalBind struct {
	dependency Any
}

// Optional declares dependency as optional (no errors occurred when dependency not resolved)
func Optional(dep Any) Any {
	return optionalBind{dep}
}

type namedBind struct {
	dependency Any
	name       string
}

// Named declares dependency that could have multiply resolves with distinct names
func Named(dep Any, name string) Any {
	if _, isOptional := dep.(optionalBind); isOptional {
		panic(fmt.Sprintf("Optional binding couldn't be embedded into named"))
	}
	return namedBind{dep, name}
}

type containerDefault struct {
	name    string
	parent  *containerDefault
	isBuilt bool

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
		panic(fmt.Sprintf(`binding "%+v[%v][%s]" is not exists for re-declaration`, symbol, reflect.TypeOf(symbol), name))
	} else if !isRebinding && isBindingExists {
		panic(fmt.Sprintf(`binding "%+v[%v][%s]" is already registered, use Rebind to replace binding`, symbol, reflect.TypeOf(symbol), name))
	}

	symbol = reflectInterfacePointers(symbol)

	binding := &Binding{}
	bindings, exists := c.factories[symbol]
	if !exists {
		bindings = make(map[string]*Binding)
		c.factories[symbol] = bindings
	}

	bindings[name] = binding
	c.isBuilt = false
	return binding
}

func (c *containerDefault) Unbind(symbol Any, names ...string) Container {
	symbol = reflectInterfacePointers(symbol)
	bindings, exists := c.factories[symbol]
	if exists {
		name := getFirstStringArgumentOrEmpty(names)
		delete(bindings, name)
		if len(bindings) == 0 {
			delete(c.factories, symbol)
		}
	}
	// else panic!?
	c.isBuilt = false
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
	if c.parent != nil {
		c.parent.Build()
	}

	if c.isBuilt {
		return
	}

	err := resolveContainerDependencies(c)
	if err != nil {
		panic(err.Error())
	}

	c.isBuilt = true
}

func (c *containerDefault) Get(symbol Any, names ...string) (Any, error) {
	name := getFirstStringArgumentOrEmpty(names)
	if c.parent == nil {
		return c.factories[symbol][name].factory()
	}
	bining, _ := c.findFactory(symbol, name)
	return bining.factory()
}

func (c *containerDefault) MustGet(symbol Any, names ...string) Any {
	resolved, err := c.Get(symbol, names...)
	if err != nil {
		panic(fmt.Sprintf("error was occurred when getting %#v[%v][%#v]: %v", symbol, reflect.TypeOf(symbol), names, err))
	}
	return resolved
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

func (c *containerDefault) Merge(other Container, name string) Container {
	container := newDefaultContainer(name)

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
	c.isBuilt = true
	return module.registerCallback(newContainerBinderProxy(c))
}

func (c *containerDefault) UnLoad(module *Module) error {
	c.isBuilt = true
	return module.unRegisterCallback(newContainerBinderProxy(c))
}

func (c *containerDefault) String() string {
	return c.name
}

func newDefaultContainer(name string) *containerDefault {
	return &containerDefault{
		name:      name,
		factories: make(map[Any]map[string]*Binding),
	}
}

// NewContainer .
func NewContainer(name string) Container {
	return newDefaultContainer(name)
}
