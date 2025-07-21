package dynamo

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// in client_test.go

func setupTest(t *testing.T) (
	ctx context.Context,
	client *Client,
	mockDB *MockdynamoDBInterface,
	mockLogger *MockLogger,
	mockMetrics *MockMetrics,
	finish func(),
) {
	ctrl := gomock.NewController(t)
	mockDB = NewMockdynamoDBInterface(ctrl)
	mockLogger = NewMockLogger(ctrl)
	mockMetrics = NewMockMetrics(ctrl)
	client = &Client{
		db:      mockDB,
		configs: &Configs{Table: "test-table", Region: "us-east-1"},
		logger:  mockLogger,
		metrics: mockMetrics,
	}
	ctx = context.Background()
	finish = func() { ctrl.Finish() }
	return
}

func Test_ClientSet(t *testing.T) {
	ctx, client, mockDB, mockLogger, mockMetrics, finish := setupTest(t)
	defer finish()

	key, val := "test-key", "test-value"
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

	require.NoError(t, client.Set(ctx, key, val))
}

func Test_ClientSetError(t *testing.T) {
	ctx, client, mockDB, mockLogger, mockMetrics, finish := setupTest(t)
	defer finish()

	key, val := "test-key", "test-value"
	expectedErr := errors.New("dynamodb error")

	mockDB.EXPECT().PutItem(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedErr)
	mockLogger.EXPECT().Debugf("error while setting data for key: %v, error: %v", key, expectedErr)
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
	assert.Equal(t, expectedErr, err)
}
