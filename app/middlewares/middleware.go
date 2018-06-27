package middlewares

import (
	"context"
	"net/http"

	"github.com/sarulabs/dingo"
	"github.com/sarulabs/dingo-example/var/lib/services/dic"
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
