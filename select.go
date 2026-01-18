package goddb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// 단건 조회
func (c DDBClient) FindByKey(tableName, pk, sk string) (map[string]types.AttributeValue, error) {

	c.trace(DEBUG, "DDBClient.FindByKey", map[string]any{
		"tableName": tableName,
		"pk":        pk,
		"sk":        sk,
	})

	output, err := c.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})

	if err != nil {
		c.trace(ERROR, "DDBClient.FindByKey.GetItem.Error", map[string]any{
			"tableName": tableName,
			"pk":        pk,
			"sk":        sk,
			"error":     err,
		})

		return nil, err
	}

	if output.Item == nil {

		c.trace(DEBUG, "DDBClient.FindByKey.GetItem.Item.Error", map[string]any{
			"tableName": tableName,
			"pk":        pk,
			"sk":        sk,
			"error":     "item not found",
		})

		return nil, errors.New("item not found")
	}

	return output.Item, nil
}

// 조회 - Range
type RangeParams struct {
	KeyConditionExpression    string
	ExpressionAttributeValues map[string]types.AttributeValue
}

// Expression 을 사용하여 조회
func (c DDBClient) FindByKeyUseExpression(tableName, pk, sk string, limit int, params RangeParams) ([]map[string]types.AttributeValue, error) {

	c.trace(DEBUG, "DDBClient.FindByKeyUseRange", map[string]any{
		"tableName":  tableName,
		"pk":         pk,
		"sk":         sk,
		"limit":      limit,
		"expression": params,
	})

	res, err := c.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(params.KeyConditionExpression),
		ExpressionAttributeValues: params.ExpressionAttributeValues,
		ScanIndexForward:          aws.Bool(true), // 최신 순
		Limit:                     aws.Int32(int32(limit)),
	})

	if err != nil {
		c.trace(ERROR, "DDBClient.FindByKeyUseRange.Query.Error", map[string]any{
			"tableName":  tableName,
			"pk":         pk,
			"sk":         sk,
			"limit":      limit,
			"expression": params,
			"error":      err,
		})
		return nil, err
	}

	return res.Items, nil
}
