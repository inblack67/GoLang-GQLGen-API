package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gofrs/uuid"
	"github.com/inblack67/GQLGenAPI/cache"
	"github.com/inblack67/GQLGenAPI/constants"
	"github.com/inblack67/GQLGenAPI/db"
	"github.com/inblack67/GQLGenAPI/graph/generated"
	"github.com/inblack67/GQLGenAPI/graph/model"
	"github.com/inblack67/GQLGenAPI/mymodels"
	"github.com/inblack67/GQLGenAPI/mysession"
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
		log.Println("hashErr = ", hashErr)
		return false, errors.New(constants.InternalServerError)
	}

	myuuid, errUUID := uuid.NewV4()

	if errUUID != nil {
		log.Println("errUUID = ", errUUID)
		return false, errors.New(constants.InternalServerError)
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

	myCtx := ctx.Value(constants.KMyContext).(types.MyCtx)

	sessionData, sessionErr := mysession.GetSessionData(myCtx.ResponseWriter, myCtx.Request, constants.KCurrentUser)

	// already logged in
	if sessionErr == nil && sessionData != nil {
		return false, errors.New(constants.KNotAuthorized)
	}

	var user = new(mymodels.User)

	err := db.PgConn.Find(&user, mymodels.User{Username: input.Username}).Error

	if err != nil {
		return false, err
	}

	notFoundErr := errors.Is(err, gorm.ErrRecordNotFound)
	if notFoundErr || (user.Username == "") {
		return false, errors.New(constants.KInvalidCredentials)
	}

	isValidPassword, argonErr := argon2id.ComparePasswordAndHash(input.Password, user.Password)

	if argonErr != nil {
		log.Println("argonErr = ", argonErr)
		return false, errors.New(constants.InternalServerError)
	}

	if !isValidPassword {
		return false, errors.New(constants.KInvalidCredentials)
	}

	maxAge := int(time.Hour) * 24

	newSessionData := types.SSession{
		ID: user.ID,
		Username: user.Username,
		UUID: user.UUID,
	}

	sessErr := mysession.SetSessionData(myCtx.ResponseWriter, myCtx.Request, newSessionData, maxAge)

	if sessErr != nil {
		log.Println("sessErr = ", sessErr)
		return false, errors.New(constants.InternalServerError)
	}

	return true, nil
}

func (r *mutationResolver) LogoutUser(ctx context.Context) (bool, error) {

	myCtx := ctx.Value(constants.KMyContext).(types.MyCtx)

	_, sessionErr := mysession.GetSessionData(myCtx.ResponseWriter, myCtx.Request, constants.KCurrentUser)

	if sessionErr != nil {
		return false, sessionErr
	}

	err := mysession.DestroySession(myCtx.ResponseWriter, myCtx.Request)

	if err != nil {
		log.Println("destroy session err = ", err)
		return false, errors.New(constants.InternalServerError)
	}

	return true, nil
}

func (r *queryResolver) Hello(ctx context.Context) (*model.Hello, error) {
	return &model.Hello{
		Reply: "worlds",
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

	myCtx := ctx.Value(constants.KMyContext).(types.MyCtx)

	sessionData, sessionErr := mysession.GetSessionData(myCtx.ResponseWriter, myCtx.Request, constants.KCurrentUser)

	if sessionErr != nil {
		return nil, sessionErr
	}

	sendUser := new(model.GetMeResponse)

	sendUser.ID = sessionData.UUID
	sendUser.Username = sessionData.Username

	return sendUser, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
