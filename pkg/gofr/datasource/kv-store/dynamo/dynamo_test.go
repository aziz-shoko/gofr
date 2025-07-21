package dynamo

import (
	"context"
	"testing"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_ClientSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockdynamoDBInterface(ctrl)
	mockLogger := NewMockLogger(ctrl)
	mockMetrics := NewMockMetrics(ctrl)

	client := &Client{
		db:      mockDB,
		configs: &Configs{Table: "test-table", Region: "us-east-1"},
		logger:  mockLogger,
		metrics: mockMetrics,
	}

	ctx := context.Background()
	key := "test-key"
	val := "test-value"

	expectedInput := &dynamodb.PutItemInput{
		TableName: aws.String("test-table"),
		Item: map[string]types.AttributeValue{
			"pk":    &types.AttributeValueMemberS{Value: key},
			"value": &types.AttributeValueMemberS{Value: val},
		},
	}

	mockDB.EXPECT().PutItem(ctx, expectedInput, gomock.Any()).Return(&dynamodb.PutItemOutput{}, nil)

	mockLogger.EXPECT().Debug(gomock.Any())

	mockMetrics.EXPECT().RecordHistogram(
		gomock.Any(),
		"app_dynamodb_stats",
		gomock.Any(),
		"table", client.configs.Table,
		"type", "SET",
	)

	err := client.Set(ctx, key, val)

	require.NoError(t, err)
}

func Test_ClientSetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockdynamoDBInterface(ctrl)
	mockLogger := NewMockLogger(ctrl)
	mockMetrics := NewMockMetrics(ctrl)

	client := &Client{
		db:      mockDB,
		configs: &Configs{Table: "test-table", Region: "us-east-1"},
		logger:  mockLogger,
		metrics: mockMetrics,
	}

	ctx := context.Background()
	key := "test-key"
	val := "test-value"

	expectedError := errors.New("dynamodb error")

	mockDB.EXPECT().PutItem(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedError)

	mockLogger.EXPECT().Debugf("error while setting data for key: %v, error: %v", key, expectedError)

	mockLogger.EXPECT().Debug(gomock.Any())

	mockMetrics.EXPECT().RecordHistogram(
		gomock.Any(),
		"app_dynamodb_stats",
		gomock.Any(),
		"table", "test-table",
		"type", "SET",
	)

	err := client.Set(ctx, key, val)

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}
