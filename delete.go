package dynamodbgo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Delete(tableName, pk, pkValue string) error {
	orm := dynamoGORM(context.Background())

	_, err := orm.db.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			pk: &types.AttributeValueMemberS{Value: pkValue},
		},
	})

	return err
}
