package inversify

// ContainerBinder holds interface to encapsulate bindings
type ContainerBinder interface {
	// Bind declares dependency (make a panic if already binded)
	Bind(symbol Any) *Binding

	// Unbind removes dependency
	Unbind(symbol Any)

	// IsBound check existences of dependency
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
