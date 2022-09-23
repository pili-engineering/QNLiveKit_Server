package application

func RegisterModule(name string, module Module) error {
	return moduleManager.RegisterModule(name, module)
}
