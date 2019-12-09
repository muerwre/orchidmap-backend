package app

import "github.com/sirupsen/logrus"

import "github.com/muerwre/orchidgo/db"

type App struct {
	Config *Config
	DB     *db.DB
}

func (a *App) NewContext() *Context {
	return &Context{
		Logger: logrus.New(),
	}
}

func New() (app *App, err error) {
	app = &App{}
	app.Config, err = InitConfig()

	if err != nil {
		return nil, err
	}

	app.DB, err = db.New()

	if err != nil {
		return nil, err
	}

	return app, err
}

func (a *App) Close() error {
	return a.DB.Close()
}
