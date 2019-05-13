package inversify

// IContainer .
type IContainer interface {
	// Bind .
	Bind(Any) *Binding

	// Get .
	Get(Any) (Any, error)
	// IsBound .
	IsBound(Any) bool

	Build()

	// Merge with another container
	Merge(IContainer) IContainer
	// SetParent supports for hierarchical DI systems
	SetParent(IContainer)

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
	parent    *containerImpl
	factories map[Any]*Binding
}

func (c *containerImpl) Bind(symbol Any) *Binding {
	b := &Binding{}
	c.factories[symbol] = b
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
	return c.factories[symbol].factory()
}

func (c *containerImpl) IsBound(symbol Any) bool {
	_, ok := c.factories[symbol]
	return ok
}

func (c *containerImpl) Merge(other IContainer) IContainer {
	container := newContainer()

	otherImpl, ok := other.(*containerImpl)
	if !ok {
		panic("container is not containerImpl")
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
		panic("container is not containerImpl")
	}

	c.parent = parentImpl
}

func newContainer() *containerImpl {
	return &containerImpl{
		factories: make(map[Any]*Binding),
	}
}

func Container() IContainer {
	return newContainer()
}
