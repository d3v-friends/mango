package mtest

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/aggregate"
	"github.com/d3v-friends/mango/mvars"
	"github.com/d3v-friends/pure-go/fnEnv"
	"github.com/d3v-friends/pure-go/fnPanic"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestQueryBuilder(test *testing.T) {
	client := fnPanic.OnPointer(mango.NewClient(&mango.ClientOpt{
		Host:     fnEnv.Read("HOST"),
		Username: fnEnv.Read("USERNAME"),
		Password: fnEnv.Read("PASSWORD"),
		Database: fnEnv.Read("DATABASE"),
	}))

	fnPanic.On(client.Migrate(context.TODO(), modelAll...))

	var now = time.Now()
	var account = &Account{
		Id:            primitive.NewObjectID(),
		AccountDataId: primitive.NewObjectID(),
		UserType:      UserTypeGeneral,
		UpdatedAt:     now,
	}

	var accountData = &AccountData{
		Id:        account.AccountDataId,
		AccountId: account.Id,
		Name:      gofakeit.Username(),
		CreatedAt: now,
	}

	fnPanic.OnPointer(
		client.
			Database().
			Collection(account.GetCollectionNm()).
			InsertOne(context.TODO(), account),
	)

	fnPanic.OnPointer(
		client.
			Database().
			Collection(accountData.GetCollectionNm()).
			InsertOne(context.TODO(), accountData),
	)

	test.Run("lookup one", func(t *testing.T) {
		var ctx = context.TODO()
		var builder, err = aggregate.NewQueryBuilder[AccountWithData](aggregate.LookUp)
		if err != nil {
			t.Fatal(err)
		}

		var res *AccountWithData
		if res, err = builder.FindOne(ctx, client.Database(), bson.M{
			mvars.FID: account.Id,
		}); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, account.Id, res.Account.Id)
		assert.Equal(t, account.AccountDataId, res.Data.Id)
		assert.Equal(t, accountData.Name, res.Data.Name)
	})
}
