package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sarulabs/dingo"
	"github.com/sarulabs/dingo-example/app/models/helpers"
	"github.com/sarulabs/dingo-example/var/lib/services/dic"
	"go.uber.org/zap"
)

// DingoMiddleware adds a container in each request context.
func DingoMiddleware(h http.HandlerFunc, app *dic.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create a request container from tha app container
		ctn, err := app.SubContainer()
		if err != nil {
			panic(err)
		}
		defer ctn.Delete()

		// store the container in the http.Request
		req := r.WithContext(context.WithValue(r.Context(), dingo.ContainerKey("dingo"), ctn))

		h(w, req)
	}
}

// PanicRecoveryMiddleware handles the panic in the handlers.
func PanicRecoveryMiddleware(h http.HandlerFunc, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// log the error
				logger.Error(fmt.Sprint(rec))

				// write the error response
				helpers.JSONResponse(w, 500, map[string]interface{}{
					"error": "Internal Error",
				})
			}
		}()

		h(w, r)
	}
}
