package middlewares

import (
	"context"
	"net/http"

	"github.com/inblack67/GQLGenAPI/constants"
	"github.com/inblack67/GQLGenAPI/types"
)

// AuthMiddleware ...
func AuthMiddleware() func(handler http.Handler) http.Handler{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

			// populating ctx with req and res
			newCtx := context.WithValue(req.Context(), constants.KMyContext, types.MyCtx{ Request: req, ResponseWriter: res })

			next.ServeHTTP(res, req.WithContext(newCtx))
		})
	}
}