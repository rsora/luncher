// Package handlers contains all routes and related handler functions supported by this service
package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/rsora/mywebapp/app/services/mywebapp/handlers/debug/checkgrp"
	"github.com/rsora/mywebapp/app/services/mywebapp/handlers/v1/usergrp"
	"github.com/rsora/mywebapp/business/core/user"
	"github.com/rsora/mywebapp/business/sys/metrics"
	"github.com/rsora/mywebapp/business/web/mid"
	"github.com/rsora/mywebapp/foundation/web"
	"go.uber.org/zap"
)

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux.
// Using the DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it!
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// DebugMux registers all the debug standard library routes and then custom
// debug application routes for the service. This bypassing the use of the
// DefaultServerMux. Using the DefaultServerMux would be a security risk since
// a dependency could inject a handler into our service without us knowing it.
func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	DB       *sqlx.DB
	Metrics  *metrics.Metrics
}

// APIMux constructs an http.Handler with all the application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(cfg.Metrics),
		mid.ExpvarMetrics(),
		mid.Panics(),
	)

	// Register v1 endpoints.
	v1(app, cfg)

	return app
}

// v1 binds 1 routes.
func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	// Register user management and authentication endpoints.
	ugh := usergrp.Handlers{
		User: user.NewCore(cfg.Log, cfg.DB),
	}
	app.Handle(http.MethodGet, version, "/users/token", ugh.Token)
	app.Handle(http.MethodGet, version, "/users/:page/:rows", ugh.Query)
	app.Handle(http.MethodGet, version, "/users/:id", ugh.QueryByID)
	app.Handle(http.MethodPost, version, "/users", ugh.Create)
	app.Handle(http.MethodPut, version, "/users/:id", ugh.Update)
	app.Handle(http.MethodDelete, version, "/users/:id", ugh.Delete)

}
