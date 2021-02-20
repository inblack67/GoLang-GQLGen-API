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

	cache.PopulateUsers()

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
		ID:       user.ID,
		Username: user.Username,
		UUID:     user.UUID,
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

func (r *mutationResolver) CreateStory(ctx context.Context, input model.CreateStoryParams) (bool, error) {
	myCtx := ctx.Value(constants.KMyContext).(types.MyCtx)

	sessionData, sessionErr := mysession.GetSessionData(myCtx.ResponseWriter, myCtx.Request, constants.KCurrentUser)

	if sessionErr != nil {
		return false, sessionErr
	}

	myuuid, errUUID := uuid.NewV4()

	if errUUID != nil {
		log.Println("errUUID = ", errUUID)
		return false, errors.New(constants.InternalServerError)
	}

	var newStory = new(mymodels.Story)

	newStory = &mymodels.Story{
		Title:    input.Title,
		UUID:     myuuid,
		UserID:   sessionData.ID,
		UserUUID: sessionData.UUID,
	}

	validationErr := newStory.ValidateStory()

	if validationErr != nil {
		return false, validationErr
	}

	storyCreationErr := db.PgConn.Create(&newStory).Error

	if storyCreationErr != nil {
		return false, storyCreationErr
	}

	// updating cache
	cache.PopulateStories()
	cache.PopulateUsers()

	return true, nil
}

func (r *queryResolver) Hello(ctx context.Context) (*model.Hello, error) {
	return &model.Hello{
		Reply: "worlds",
	}, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*mymodels.User, error) {
	cachedMarshalledUsers, getErr := cache.RedisClient.Get(ctx, constants.KUsers).Result()

	if getErr != nil {
		log.Println("redis get err = ", getErr)
		return nil, errors.New(constants.InternalServerError)
	}

	defer utils.Elapsed("redis query => users")()

	// cached
	var cachedUsers []*mymodels.User

	unmarshalErr := json.Unmarshal([]byte(cachedMarshalledUsers), &cachedUsers)

	if unmarshalErr != nil {
		log.Fatal("unmarshalErr", unmarshalErr)
		return nil, errors.New(constants.InternalServerError)
	}

	return cachedUsers, nil
}

func (r *queryResolver) GetMe(ctx context.Context) (*model.GetMeResponse, error) {
	defer utils.Elapsed("redis => getMe query")()

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

func (r *queryResolver) Stories(ctx context.Context) ([]*mymodels.Story, error) {
	defer utils.Elapsed("redis query => stories")()

	marshalledStories, err := cache.RedisClient.Get(ctx, constants.KStories).Result()

	if err != nil {
		log.Println("redis get err", err)
		return nil, errors.New(constants.InternalServerError)
	}

	var stories []*mymodels.Story

	unmarshallErr := json.Unmarshal([]byte(marshalledStories), &stories)

	if err != nil {
		log.Println("redis unmarshallErr", unmarshallErr)
		return nil, errors.New(constants.InternalServerError)
	}

	return stories, nil
}

func (r *storyResolver) UserID(ctx context.Context, obj *mymodels.Story) (string, error) {
	return obj.UserUUID, nil
}

func (r *storyResolver) CreatedAt(ctx context.Context, obj *mymodels.Story) (string, error) {
	return obj.CreatedAt.String(), nil
}

func (r *storyResolver) UpdatedAt(ctx context.Context, obj *mymodels.Story) (string, error) {
	return obj.UpdatedAt.String(), nil
}

func (r *storyResolver) DeletedAt(ctx context.Context, obj *mymodels.Story) (string, error) {
	return obj.DeletedAt.Time.String(), nil
}

func (r *storyResolver) UUID(ctx context.Context, obj *mymodels.Story) (string, error) {
	return obj.UUID.String(), nil
}

func (r *userResolver) CreatedAt(ctx context.Context, obj *mymodels.User) (string, error) {
	return obj.CreatedAt.String(), nil
}

func (r *userResolver) UpdatedAt(ctx context.Context, obj *mymodels.User) (string, error) {
	return obj.UpdatedAt.String(), nil
}

func (r *userResolver) DeletedAt(ctx context.Context, obj *mymodels.User) (string, error) {
	return obj.DeletedAt.Time.String(), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Story returns generated.StoryResolver implementation.
func (r *Resolver) Story() generated.StoryResolver { return &storyResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type storyResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
