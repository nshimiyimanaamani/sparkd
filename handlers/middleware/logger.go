package middleware

import (
	"net/http"

	"github.com/iradukunda1/firecrackerland/internal/render"
	lgg "github.com/sirupsen/logrus"
)

type Middleware func(h http.Handler) http.Handler

// include context with logger in http server for downstream use
func SetLoggerCtx(lg *lgg.Logger) Middleware {

	return func(next http.Handler) http.Handler {

		f := func(w http.ResponseWriter, r *http.Request) {

			r = r.WithContext(render.SetLogger(r.Context(), lg))

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(f)
	}
}
