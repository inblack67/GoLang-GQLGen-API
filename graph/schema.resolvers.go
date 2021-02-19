package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	"github.com/inblack67/GQLGenAPI/utils"
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

func (r *queryResolver) Hello(ctx context.Context) (*model.Hello, error) {

	myUser, err := middlewares.GetUserFromCtx(ctx)

	fmt.Println("myuser lol", myUser)

	if err != nil {
		log.Fatal("not auth")
	}

	return &model.Hello{
		Reply: myUser,
	}, nil

	// data  := ctx.Value(constants.KCurrentUser)
	// fmt.Println("woah", data)

	// reply, ok := data.(string)

	// if !ok{
	// 	return &model.Hello{
	// 	Reply: string("worlds"),
	// }, nil
	// }

	// return &model.Hello{
	// 	Reply: string(reply),
	// }, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {

	cachedMarshalledUsers, getErr := cache.RedisClient.Get(context.Background(), constants.KGetUsers).Result()

	// not in redis yet => query goes to db
	if getErr != nil {

		defer utils.Elapsed("db query => users")()

		var dbUsers []mymodels.User

		dbc := db.PgConn.Find(&dbUsers)

		var users []*model.User

		for _, v := range dbUsers {
			users = append(users, &model.User{
				Name:     v.Name,
				Email:    v.Email,
				Username: v.Username,
				CreatedAt: v.CreatedAt.String(),
				UpdatedAt: v.UpdatedAt.String(),
				DeletedAt: v.DeletedAt.Time.String(),
				UUID: v.UUID,
			})
		}

		marshalledUsers, marshallErr  := json.Marshal(users)

		if marshallErr != nil{
			log.Fatal("marshallErr", marshallErr)
		}

		setErr := cache.RedisClient.Set(context.Background(), constants.KGetUsers, marshalledUsers, time.Hour * 24).Err()

		if setErr != nil{
			log.Fatal("setErr", setErr)
		}

		return users, dbc.Error
	}

	defer utils.Elapsed("redis query => users")()

	// cached
	var cachedUsers []*model.User

	unmarshalErr := json.Unmarshal([]byte(cachedMarshalledUsers), &cachedUsers)

	if unmarshalErr != nil{
		log.Fatal("unmarshalErr", unmarshalErr)
	}

	return cachedUsers, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
