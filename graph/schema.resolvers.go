package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"github.com/alexedwards/argon2id"
	"github.com/gofrs/uuid"
	"github.com/inblack67/GQLGenAPI/db"
	"github.com/inblack67/GQLGenAPI/graph/generated"
	"github.com/inblack67/GQLGenAPI/graph/model"
	"github.com/inblack67/GQLGenAPI/mymodels"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterParams) (bool, error) {

				var newUser = new(mymodels.User)

				newUser = &mymodels.User{
					Name: input.Name,
					Email: input.Email,
					Password: input.Password,
					Username: input.Username,
				}

				validationErr := newUser.ValidateMe()

				if validationErr != nil{
					return false, validationErr
				}

				hashedPassword, hashErr := argon2id.CreateHash(input.Password, argon2id.DefaultParams)

				if hashErr != nil{
					log.Fatalf(hashErr.Error())
				}

				myuuid, errUUID := uuid.NewV4()

				// strUUID := fmt.Sprintf("%v", myuuid)

				if errUUID != nil{
					log.Fatalf(errUUID.Error())
				}

				newUser.Password = hashedPassword
				newUser.UUID = myuuid

				userCreationErr := db.PgConn.Create(&newUser).Error

				if userCreationErr != nil {
					return false, userCreationErr
				}

				return true, nil
}

func (r *queryResolver) Hello(ctx context.Context) (*model.Hello, error) {
	return &model.Hello{
		Reply: "worlds",
	}, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
