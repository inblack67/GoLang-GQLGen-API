package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gofrs/uuid"
	"github.com/inblack67/GQLGenAPI/cache"
	"github.com/inblack67/GQLGenAPI/constants"
	"github.com/inblack67/GQLGenAPI/db"
	"github.com/inblack67/GQLGenAPI/graph/generated"
	"github.com/inblack67/GQLGenAPI/graph/model"
	"github.com/inblack67/GQLGenAPI/middlewares"
	"github.com/inblack67/GQLGenAPI/mymodels"
	"github.com/inblack67/GQLGenAPI/types"
	"github.com/inblack67/GQLGenAPI/utils"
	"gorm.io/gorm"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterParams) (bool, error) {
	var newUser = new(mymodels.User)

	newUser = &mymodels.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Username: input.Username,
	}

	validationErr := newUser.ValidateMe()

	if validationErr != nil {
		return false, validationErr
	}

	hashedPassword, hashErr := argon2id.CreateHash(input.Password, argon2id.DefaultParams)

	if hashErr != nil {
		log.Fatalf(hashErr.Error())
	}

	myuuid, errUUID := uuid.NewV4()

	if errUUID != nil {
		log.Fatalf(errUUID.Error())
	}

	strUUID := myuuid.String()

	newUser.Password = hashedPassword
	newUser.UUID = strUUID

	userCreationErr := db.PgConn.Create(&newUser).Error

	if userCreationErr != nil {
		return false, userCreationErr
	}

	return true, nil
}

func (r *mutationResolver) LoginUser(ctx context.Context, input model.LoginParams) (bool, error) {
	isAuth := middlewares.IsAuthenticated(ctx)

	if isAuth {
		return false, errors.New(constants.KNotAuthorized)
	}

	var user = new(mymodels.User)

	err := db.PgConn.Find(&user, mymodels.User{Username: input.Username}).Error

	if err != nil {
		return false, errors.New(err.Error())
	}

	notFoundErr := errors.Is(err, gorm.ErrRecordNotFound)
	if notFoundErr || (user.Username == "") {
		return false, errors.New(constants.KInvalidCredentials)
	}

	isValidPassword, argonErr := argon2id.ComparePasswordAndHash(input.Password, user.Password)

	if argonErr != nil {
		return false, errors.New(argonErr.Error())
	}

	if !isValidPassword {
		return false, errors.New(constants.KInvalidCredentials)
	}

	sessionData := new(types.SSession)

	sessionData.ID = user.ID
	sessionData.Username = user.Username
	sessionData.UUID = user.UUID

	marshalledSessionData, marshallErr := json.Marshal(sessionData)

	if marshallErr != nil {
		return false, errors.New(marshallErr.Error())
	}

	setErr := cache.RedisClient.Set(ctx, constants.KCurrentUser, marshalledSessionData, time.Hour*24).Err()

	uuid, uuidErr := uuid.NewV4()

	token := fmt.Sprint(uuid)

	if uuidErr != nil {
		log.Fatal(uuidErr.Error())
	}

	setErr2 := cache.RedisClient.Set(ctx, constants.KAuthSession, token, time.Hour*24).Err()

	if setErr != nil {
		log.Fatal(setErr.Error())
	}

	if setErr2 != nil {
		log.Fatal(setErr2.Error())
	}

	myCtx := ctx.Value(constants.KMyContext).(types.MyCtx)

	http.SetCookie(myCtx.ResponseWriter, &http.Cookie{
		Name:     constants.KAuthSession,
		Value:    token,
		Expires:  time.Now().AddDate(0, 0, 1),
		MaxAge:   int(time.Hour) * 24,
		SameSite: http.SameSiteLaxMode,
	})

	return true, nil
}

func (r *mutationResolver) LogoutUser(ctx context.Context) (bool, error) {
	isAuth := middlewares.IsAuthenticated(ctx)

	if !isAuth {
		return false, errors.New(constants.KNotAuthenticated)
	}

	delErr := cache.RedisClient.Del(ctx, constants.KAuthSession).Err()

	delErr2 := cache.RedisClient.Del(ctx, constants.KCurrentUser).Err()

	if delErr != nil {
		log.Fatal(delErr)
	}

	if delErr2 != nil {
		log.Fatal(delErr2)
	}

	myCtx := ctx.Value(constants.KMyContext).(types.MyCtx)

	// delete cookie
	http.SetCookie(myCtx.ResponseWriter, &http.Cookie{
		Name:     constants.KAuthSession,
		Value:    "nil",
		Expires:  time.Now(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	return true, nil
}

func (r *queryResolver) Hello(ctx context.Context) (*model.Hello, error) {
	return &model.Hello{
		Reply: "myUser",
	}, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	cachedMarshalledUsers, getErr := cache.RedisClient.Get(ctx, constants.KGetUsers).Result()

	// not in redis yet => query goes to db
	if getErr != nil {

		defer utils.Elapsed("db query => users")()

		var dbUsers []mymodels.User

		dbc := db.PgConn.Find(&dbUsers)

		var users []*model.User

		for _, v := range dbUsers {
			users = append(users, &model.User{
				Name:      v.Name,
				Email:     v.Email,
				Username:  v.Username,
				CreatedAt: v.CreatedAt.String(),
				UpdatedAt: v.UpdatedAt.String(),
				DeletedAt: v.DeletedAt.Time.String(),
				UUID:      v.UUID,
			})
		}

		marshalledUsers, marshallErr := json.Marshal(users)

		if marshallErr != nil {
			log.Fatal("marshallErr", marshallErr)
		}

		setErr := cache.RedisClient.Set(ctx, constants.KGetUsers, marshalledUsers, time.Hour*24).Err()

		if setErr != nil {
			log.Fatal("setErr", setErr)
		}

		return users, dbc.Error
	}

	defer utils.Elapsed("redis query => users")()

	// cached
	var cachedUsers []*model.User

	unmarshalErr := json.Unmarshal([]byte(cachedMarshalledUsers), &cachedUsers)

	if unmarshalErr != nil {
		log.Fatal("unmarshalErr", unmarshalErr)
	}

	return cachedUsers, nil
}

func (r *queryResolver) GetMe(ctx context.Context) (*model.GetMeResponse, error) {
	isAuth := middlewares.IsAuthenticated(ctx)

	if !isAuth {
		return nil, errors.New(constants.KNotAuthenticated)
	}

	ctxUser, err := middlewares.GetUserFromCtx(ctx)

	if err != nil {
		return nil, err
	}

	sendUser := new(model.GetMeResponse)

	sendUser.ID = ctxUser.UUID
	sendUser.Username = ctxUser.Username

	return sendUser, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
