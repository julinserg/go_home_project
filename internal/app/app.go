package app

type App struct { // TODO
}

type Logger interface { // TODO
}

func New(logger Logger) *App {
	return &App{}
}
