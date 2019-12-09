package app

import "github.com/sirupsen/logrus"

import "github.com/muerwre/orchidgo/db"

type Context struct {
	Logger        logrus.FieldLogger
	RemoteAddress string
	DB            *db.DB
	Config        *Config
}

func (ctx *Context) WithLogger(logger logrus.FieldLogger) *Context {
	ret := *ctx
	ret.Logger = logger
	return &ret
}

func (ctx *Context) WithRemoteAddress(address string) *Context {
	ret := *ctx
	ret.RemoteAddress = address
	return &ret
}
