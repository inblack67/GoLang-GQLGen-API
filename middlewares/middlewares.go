package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/inblack67/GQLGenAPI/constants"
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

			data := "foobar"

			session.Values[constants.KCurrentUser] = data

			currentUser := session.Values[constants.KCurrentUser]

			fmt.Println("currentUser",currentUser)

			sessionSaveErr := session.Save(req, res)
			if sessionSaveErr != nil {
				http.Error(res, err.Error(), http.		StatusInternalServerError)
				return
			}

			newCtx := context.WithValue(req.Context(), constants.KCurrentUser, data)

			next.ServeHTTP(res, req.WithContext(newCtx))

		})
	}
}

// GetUserFromCtx ...
func GetUserFromCtx (ctx context.Context) (string, error) {
	user := ctx.Value(constants.KCurrentUser)
	if user == nil {
		return "nil", errors.New("not auth")
	}
	strUser, ok := user.(string)

	if !ok {
		return "type err", nil
	}

	return strUser, nil
}