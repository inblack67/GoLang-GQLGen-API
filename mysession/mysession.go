package mysession

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/inblack67/GQLGenAPI/constants"
	"github.com/inblack67/GQLGenAPI/types"
)

var (
	// SessionStore ...
	SessionStore = sessions.NewCookieStore([]byte(constants.KSessionSecret))
)

// DestroySession ...
func DestroySession (res http.ResponseWriter, req *http.Request) error {

	session, err := SessionStore.Get(req, constants.KAuthSession)
	if err != nil {
		log.Println("session get err = ", err)
		return err
	}

	session.Values[constants.KCurrentUser] = nil

	session.Options.MaxAge = -1		// delete cookie

	err2 := session.Save(req, res)

	if err2 != nil {
		log.Println("session save err = ", err2)
		return err2
	}

	return nil
}

// SetSessionData ...
func SetSessionData (res http.ResponseWriter, req *http.Request, data types.SSession, maxAge int) error {

	session, err := SessionStore.Get(req, constants.KAuthSession)
	if err != nil {
		return err
	}

	marshelledData, marshallErr := json.Marshal(data)

	if marshallErr != nil {
		log.Fatal(marshallErr)
	}

	session.Values[constants.KCurrentUser] = marshelledData

	session.Options = &sessions.Options{
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge: maxAge,
		Secure: false,
	}

	err2 := session.Save(req, res)

	if err2 != nil {
		return err2
	}

	return nil
}

// GetSessionData ...
func GetSessionData (res http.ResponseWriter, req *http.Request, key string) (*types.SSession, error) {

	session, err := SessionStore.Get(req, constants.KAuthSession)
	if err != nil {
		log.Println("getting session err = ", err)
		return nil, err
	}

	data := session.Values[key]

	if data == nil {
		return nil, errors.New(constants.KNotAuthenticated)
	}

	byteData, ok := data.([]byte)

	if !ok {
		log.Println("byte err")
		return nil, errors.New(constants.InternalServerError)
	}

	var sessionData = new(types.SSession)

	unmarshallErr := json.Unmarshal([]byte(byteData), sessionData)

	if unmarshallErr != nil {
		log.Fatal(unmarshallErr)
		return nil, errors.New(constants.InternalServerError)
	}

	return sessionData, nil
}

