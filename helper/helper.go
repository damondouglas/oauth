package helper

import (
	"context"
	"net/http"
)

// ContextHelper provides context.Context from http.Request
type ContextHelper interface {
	Context(*http.Request) context.Context
}

// ClientHelper provides http.Client from context.Context
type ClientHelper interface {
	Client(context.Context) *http.Client
}
