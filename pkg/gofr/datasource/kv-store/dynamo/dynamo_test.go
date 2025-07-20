package dynamo

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestClient_Set_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDB := NewMockdynamoDBInterface(ctrl)

	client := &Client{
		db: mockDB,
		Table: "test-table",
	}

	ctx := context.Background()
	key := "test-key"
	value := "test-value"

	// Set expectations for PutItem call
	expectedInput := &dynamodb.PutItemInput{
		TableName: aws.String(client.Table),
		Item: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: key},
			"value": &types.AttributeValueMemberS{Value: value},
		},
	}

	mockDB.EXPECT().PutItem(ctx, expectedInput, gomock.Any()).Return(&dynamodb.PutItemOutput{},nil)

	err := client.Set(ctx, key, value)

	require.NoError(t, err)
}

