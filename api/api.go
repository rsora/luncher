package api

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux"
	"github.com/rsora/luncher/handler"
	"go.uber.org/zap"
)

// A Handler is a type that handles a http request, it differs from http.Handler
// because it returns an error and the context is explicitly passed.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// APIMuxConfig contains all the mandatory dependencies required by handlers.
type APIMuxConfig struct {
	Log *zap.SugaredLogger
}

// api represents our server api.
type api struct {
	*httptreemux.ContextMux
	// mw  []web.Middleware
	log      *zap.SugaredLogger
	shutdown chan os.Signal
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) http.Handler {
	a := &api{
		ContextMux: httptreemux.NewContextMux(),
		log:        cfg.Log,
	}

	// // Setup the middleware common to each handler.
	// a.mw = append(a.mw, middleware.RequestID())
	// a.mw = append(a.mw, middleware.NoCache())

	// // Don't log endpoints from a middleware.
	// // Each endpoint will log its own requests with more pertinent data fields.
	// // a.mw = append(a.mw, middleware.Logger(cfg.Log))

	// a.mw = append(a.mw, middleware.Errors(cfg.Log))
	// a.mw = append(a.mw, middleware.Panics())

	a.Handle(http.MethodGet, "/daily", handler.GetSuggestion)
	return a.ContextMux
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application router.
//func (a *api) Handle(method string, path string, handler web.Handler, mw ...web.Middleware) {
func (a *api) Handle(method string, path string, handler Handler) {
	// // First wrap handler specific middleware around this handler.
	// handler = web.WrapMiddleware(mw, handler)
	// // Add the application's general middleware to the handler chain.
	// handler = web.WrapMiddleware(a.mw, handler)

	// The function to execute for each request.
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pull the context from the request and
		// use it as a separate parameter.
		ctx := r.Context()

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			// Some bad and unrecoverable error happened.
			a.log.Errorw("handler unrecoverable", "error", err)
			a.SignalShutdown()
			return
		}
	})

	a.ContextMux.Handle(method, path, h)
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *api) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
