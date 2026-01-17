package goddb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	BATCH_SIZE = 25
)

func (c DDBClient) Insert(tableName string, item any) error {
	c.trace(DEBUG, "DDBClient.Insert", map[string]any{
		"tableName": tableName,
		"item":      item,
	})

	marshalItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		c.trace(ERROR, "DDBClient.Insert.MarshalMap.Error", map[string]any{
			"tableName": tableName,
			"item":      item,
			"error":     err,
		})
		return err
	}

	_, err = c.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      marshalItem,

		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})

	if err != nil {
		var condFailed *types.ConditionalCheckFailedException
		if errors.As(err, &condFailed) {
			c.trace(ERROR, "DDBClient.Insert.ConditionCheckFailed.Error", map[string]any{
				"tableName": tableName,
				"item":      condFailed.Item,
				"error":     err,
			})
			return err
		}
	}

	c.trace(INFO, "DDBClient.Insert.Success", map[string]any{
		"tableName": tableName,
		"item":      marshalItem,
	})

	return nil
}

// insert batch (conditino 없음)
func (c DDBClient) InsertBatch(tableName string, items []any) error {

	c.trace(DEBUG, "DDBClient.InsertBatch", map[string]any{
		"tableName": tableName,
		"itemCount": len(items),
	})

	for i := 0; i < len(items); i += BATCH_SIZE {
		end := i + BATCH_SIZE

		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]

		var writeRequests []types.WriteRequest
		for _, v := range batch {

			marsharV, err := attributevalue.MarshalMap(v)
			if err != nil {
				c.trace(ERROR, "DDBClient.InsertBatch.MarshalMap.Error", map[string]any{
					"tableName": tableName,
					"item":      v,
					"error":     err,
				})
				return err
			}

			writeRequests = append(writeRequests, types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: marsharV,
				},
			})

		}

		results, err := c.client.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: writeRequests,
			},
		})

		if err != nil {
			c.trace(ERROR, "DDBClient.InsertBatch.BatchWriteItem.Error", map[string]any{
				"tableName": tableName,
				"error":     err,
			})
			return err
		}

		// retry
		retryCount := 0
		for len(results.UnprocessedItems) > 0 && retryCount < 3 {
			retryCount++
			results, err = c.client.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
				RequestItems: results.UnprocessedItems,
			})
			if err != nil {
				c.trace(ERROR, "DDBClient.InsertBatch.BatchWriteItem.Error", map[string]any{
					"tableName":  tableName,
					"error":      err,
					"retryCount": retryCount,
				})
				return err
			}
		}

		if len(results.UnprocessedItems) > 0 {
			c.trace(ERROR, "DDBClient.InsertBatch.BatchWriteItem.UnprocessedItems", map[string]any{
				"tableName":        tableName,
				"unprocessedItems": results.UnprocessedItems,
			})
			return errors.New("unprocessed items")
		}

		c.trace(INFO, "DDBClient.InsertBatch.Success", map[string]any{
			"tableName":      tableName,
			"itemCount":      len(items),
			"unprocessItems": len(results.UnprocessedItems),
		})
	}

	return nil
}
