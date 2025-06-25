package dynamodbgo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

/*
@desc pk를 제외한 값들을 정의합니다.
*/
func UpdatePartial(tableName, pkName, pkValue string, updates map[string]any) error {
	orm := dynamoGORM(context.Background())

	updateExpression := "SET "
	expressionAttributeNames := make(map[string]string)
	expressionAttributeValues := make(map[string]types.AttributeValue)

	i := 1
	for key, value := range updates {
		placeholder := fmt.Sprintf("#field%d", i)
		valuePlaceholder := fmt.Sprintf(":val%d", i)

		updateExpression += fmt.Sprintf("%s = %s, ", placeholder, valuePlaceholder)
		expressionAttributeNames[placeholder] = key
		expressionAttributeValues[valuePlaceholder] = serializeValue(value)

		i++
	}
	updateExpression = updateExpression[:len(updateExpression)-2]

	_, err := orm.db.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			pkName: &types.AttributeValueMemberS{Value: pkValue},
		},
		UpdateExpression:          &updateExpression,
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	})

	if err != nil {
		fmt.Println("UpdateItem error:", err)
	}

	return err
}

func serializeValue(value any) types.AttributeValue {
	switch v := value.(type) {
	case string:
		return &types.AttributeValueMemberS{Value: v}
	case int:
		return &types.AttributeValueMemberN{Value: strconv.Itoa(v)}
	case bool:
		return &types.AttributeValueMemberBOOL{Value: v}
	default:
		return &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", v)}
	}
}
