package inversify

// ContainerBinder holds interface to encapsulate bindings
type ContainerBinder interface {
	// Bind .
	Bind(symbol Any) *Binding

	// Unbind .
	Unbind(symbol Any)

	// IsBound .
	IsBound(symbol Any) bool
}

type containerBinderProxy struct {
	container Container
}

func (proxy *containerBinderProxy) Bind(symbol Any) *Binding {
	return proxy.container.Bind(symbol)
}

func (proxy *containerBinderProxy) Unbind(symbol Any) {
	proxy.container.Unbind(symbol)
}

func (proxy *containerBinderProxy) IsBound(symbol Any) bool {
	return proxy.container.IsBound(symbol)
}

func newContainerBinderProxy(container Container) ContainerBinder {
	return &containerBinderProxy{
		container: container,
	}
}
