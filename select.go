package dynamodbgo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func FindByUsePK[T any](tableName, key string, value string) (T, error) {

	context := context.Background()
	orm := dynamoGORM(context)

	output, err := orm.db.GetItem(context, &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			key: &types.AttributeValueMemberS{Value: value},
		},
	})

	var zero T
	if err != nil {
		return zero, nil
	}

	data, err := deserialize[T](output.Item)
	if err != nil {
		return zero, nil
	}

	return data, nil
}

func SelectAll[T any](tableName string) ([]T, error) {

	orm := dynamoGORM(context.Background())

	outputs, err := orm.db.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: wrapString(tableName),
	})

	if err != nil {
		return nil, err
	}

	values := make([]T, len(outputs.Items))
	for i, output := range outputs.Items {
		v, err := deserialize[T](output)
		if err != nil {
			return values, err
		}

		values[i] = v
	}

	return values, nil
}
