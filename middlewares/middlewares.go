package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/inblack67/GQLGenAPI/cache"
	"github.com/inblack67/GQLGenAPI/constants"
	"github.com/inblack67/GQLGenAPI/types"
)

// AuthMiddleware ...
func AuthMiddleware() func(handler http.Handler) http.Handler{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

			newCtx := context.WithValue(req.Context(), constants.KMyContext, types.MyCtx{ Request: req, ResponseWriter: res })

			next.ServeHTTP(res, req.WithContext(newCtx))
		})
	}
}

// GetUserFromCtx ...
func GetUserFromCtx (ctx context.Context) (*types.SSession, error) {

	user := new(types.SSession)

	ctxUser, err := cache.RedisClient.Get(ctx, constants.KCurrentUser).Result()

	if err != nil {
		log.Fatal(err)
	}

	err2 := json.Unmarshal([]byte(ctxUser), &user)

	if err2 != nil {
		log.Fatal(err2)
	}

	return user, nil
}

// IsAuthenticated ...
func IsAuthenticated (ctx context.Context) bool {

	myCtx := ctx.Value(constants.KMyContext).(types.MyCtx)

	cookie, cookieErr := myCtx.Request.Cookie(constants.KAuthSession)

	// no cookie of ours
	if cookie == nil || cookieErr != nil {
		return false
	}

	ourToken, getErr := cache.RedisClient.Get(ctx, constants.KAuthSession).Result()

	if getErr != nil {
		// no session => not auth
		return false
	}

	recievedToken := string(cookie.Value)

	if recievedToken != ourToken {
		http.Error(myCtx.ResponseWriter, "Invalid cookie", http.StatusForbidden)
		return false
	}

	return true
}