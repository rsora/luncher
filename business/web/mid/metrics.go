package mid

import (
	"context"
	"net/http"

	"github.com/rsora/mywebapp/business/sys/metrics"
	"github.com/rsora/mywebapp/foundation/web"
)

// Metrics updates information about api calls.
func Metrics(mt *metrics.Metrics) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			t := mt.TraceAPI(r.URL.Path)
			err := handler(ctx, w, r)
			t.Mark()

			return err
		}
		return h
	}
	return m
}
