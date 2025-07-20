package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	// "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dynamoDBInterface interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	// We'll add more methods later for Get/Delete, but keep minimal for now.
}

type Client struct {
	db dynamoDBInterface
	Table string
}

func (c *Client) Set(ctx context.Context, key, val string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.Table),
		Item: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: key},
			"value": &types.AttributeValueMemberS{Value: val},
		},
	}

	_, err := c.db.PutItem(ctx, input)
	if err != nil {
		return err
	}

	return nil
}