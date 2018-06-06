package appengine

import (
	"context"
	"net/http"

	"google.golang.org/appengine/urlfetch"

	"google.golang.org/appengine"
)

// ContextHelper is a helper.ContextHelper implementation specific to appengine.
type ContextHelper struct {
}

// Context provides appengine context from http.Request
func (h *ContextHelper) Context(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

// ClientHelper is a helper.ClientHelper implementation specific to appengine.
type ClientHelper struct {
}

// Client provides http.Client from appengine context.
func (h *ClientHelper) Client(ctx context.Context) *http.Client {
	return urlfetch.Client(ctx)
}
