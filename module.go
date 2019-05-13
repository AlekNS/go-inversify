package inversify

// Module .
type Module struct {
	name string

	registerCallback   func(ContainerBinder) error
	unRegisterCallback func(ContainerBinder) error
}

// Register .
func (mdl *Module) Register(callback func(ContainerBinder) error) *Module {
	mdl.registerCallback = callback
	return mdl
}

// UnRegister .
func (mdl *Module) UnRegister(callback func(ContainerBinder) error) *Module {
	mdl.unRegisterCallback = callback
	return mdl
}

func dummyModuleCallback(ContainerBinder) error { return nil }

// NewModule .
func NewModule(name string) *Module {
	return &Module{
		name:               name,
		registerCallback:   dummyModuleCallback,
		unRegisterCallback: dummyModuleCallback,
	}
}
