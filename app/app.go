package app

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/google/uuid"
	"github.com/rsora/luncher/app/web"
)

// App represents our server App.
type App struct {
	mux      *httptreemux.ContextMux
	mw       []web.Middleware
	shutdown chan os.Signal
}

// New creates an App value that handle a set of routes for the application.
func New(shutdown chan os.Signal, mw ...web.Middleware) *App {

	// Create an OpenTelemetry HTTP Handler which wraps our router. This will start
	// the initial span and annotate it with information about the request/response.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if an client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/

	mux := httptreemux.NewContextMux()

	return &App{
		mux:      mux,
		shutdown: shutdown,
		mw:       mw,
	}
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application router.
func (a *App) Handle(method string, path string, handler web.Handler, mw ...web.Middleware) {
	// First wrap handler specific middleware around this handler.
	handler = web.WrapMiddleware(mw, handler)
	// Add the application's general middleware to the handler chain.
	handler = web.WrapMiddleware(a.mw, handler)

	// The function to execute for each request.
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pull the context from the request and
		// use it as a separate parameter.
		ctx := r.Context()

		// Set the context with the required values to
		// process the request.
		// The Now field has to be used downstream (in the data CRUD layer as well)
		// as a consistent timestamp for the request handling start.
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}
		ctx = context.WithValue(ctx, key, &v)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			// Some bad and unrecoverable error happened.
			a.SignalShutdown()
			return
		}
	})

	a.mux.Handle(method, path, h)
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// ServeHTTP implements the http.Handler interface. It's the entry point for
// all http traffic and allows the opentelemetry mux to run first to handle
// tracing. The opentelemetry mux then calls the application mux to handle
// application traffic. This was setup on line 44 in the NewApp function.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
