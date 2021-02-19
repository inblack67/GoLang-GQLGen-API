package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/inblack67/GQLGenAPI/cache"
	"github.com/inblack67/GQLGenAPI/constants"
	"github.com/inblack67/GQLGenAPI/types"
)

// AuthMiddleware ...
func AuthMiddleware() func(handler http.Handler) http.Handler{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			var sessionStore  = sessions.NewCookieStore([]byte(constants.KSessionSecret))

			session, err := sessionStore.Get(req, constants.KAuthSession)

			if err != nil{
				log.Fatalf(err.Error())
			}

			cachedMarshalledSessionData, getErr := cache.RedisClient.Get(context.Background(), constants.KAuthSession).Result()

			if getErr != nil {

				session.Options = &sessions.Options{
					// Domain: ".domain.com",	// for prod
					// Secure: true,	// for prod
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					MaxAge: -1,	// MaxAge<0 means delete cookie immediately.
				}

				next.ServeHTTP(res, req.WithContext(req.Context()))
				return
			}

			session.Options = &sessions.Options{
				// Domain: ".domain.com",	// for prod
				// Secure: true,	// for prod
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				MaxAge: int(time.Hour) * 24,
			}

			var sessionData = new(types.SSession)

			unmarshallErr := json.Unmarshal([]byte(cachedMarshalledSessionData), sessionData)

			if unmarshallErr != nil {
				next.ServeHTTP(res, req.WithContext(req.Context()))
				return
			}

			marshall, err := json.Marshal(sessionData)

			session.Values[constants.KCurrentUser] = marshall

			sessionSaveErr := session.Save(req, res)
			if sessionSaveErr != nil {
				next.ServeHTTP(res, req.WithContext(req.Context()))
				return
			}

			newCtx := context.WithValue(req.Context(), constants.KCurrentUser, marshall)

			next.ServeHTTP(res, req.WithContext(newCtx))
		})
	}
}

// GetUserFromCtx ...
func GetUserFromCtx (ctx context.Context) (*types.SSession, error) {

	user := new(types.SSession)

	ctxUser := ctx.Value(constants.KCurrentUser)

	str, ok := ctxUser.([]byte)

	if !ok {
		return nil, errors.New(constants.KNotAuthenticated)
	}

	err := json.Unmarshal([]byte(str), &user)

	if err != nil {
		return nil, errors.New(constants.KNotAuthenticated)
	}

	return user, nil
}

// SETSession ...
func SETSession (data interface {}) {

}