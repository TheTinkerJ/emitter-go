package bootstrap

type Application struct {
	env *Env
}

func NewApp() Application {
	app := Application{
		env: NewEnv(),
	}
	return app
}
