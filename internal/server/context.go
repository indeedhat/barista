package server

import "context"

type Context struct {
	context.Context
	cfg ServerConfig
}

// NewContext creates a new instance of the servermux context from the provided ResponseWriter and Request params
func NewContext(ctx context.Context, cfg ServerConfig) Context {
	return Context{
		Context: ctx,
		cfg:     cfg,
	}
}

func (c Context) Config() ServerConfig {
	return c.cfg
}

func (c Context) WithData(key, val any) Context {
	return Context{
		Context: context.WithValue(c.Context, key, val),
		cfg:     c.cfg,
	}
}
