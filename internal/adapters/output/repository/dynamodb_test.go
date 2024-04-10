package repository_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/posilva/simpleboards/internal/adapters/output/repository"
	"github.com/posilva/simpleboards/internal/testutil"
	testmocks "github.com/posilva/simpleboards/internal/testutil/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDynamoDBRepository_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	client := testmocks.NewMockDynamoDBClient(ctrl)

	attributes := make(map[string]types.AttributeValue)
	attributes["score"] = &types.AttributeValueMemberN{Value: "2"}

	client.EXPECT().UpdateItem(ctx, gomock.Any()).Return(
		&dynamodb.UpdateItemOutput{
			Attributes: attributes,
		}, nil).AnyTimes()

	settings := testutil.NewMockDefaultDynamoDBSettings(client)
	r, err := repository.NewDynamoDBRepository(settings)
	assert.NoError(t, err)

	entry := testutil.NewID()
	leaderboard := testutil.NewUnique(testutil.Name(t))
	score := float64(1)

	_, err = r.Add(entry, leaderboard, score)
	v1, err := r.Add(entry, leaderboard, score)

	assert.NoError(t, err)
	assert.Equal(t, uint64(2*score), uint64(v1))
}
