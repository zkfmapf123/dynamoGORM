package goddb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func MarshalMap[T any](item map[string]types.AttributeValue) T {

	var result T
	attributevalue.UnmarshalMap(item, &result)

	return result
}

func MarshalMaps[T any](item []map[string]types.AttributeValue) []T {

	var result []T

	for _, v := range item {
		result = append(result, MarshalMap[T](v))
	}

	return result
}
