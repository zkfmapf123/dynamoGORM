package dynamodbgo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dynamoGORMParmas struct {
	db dynamodb.Client
}

func dynamoGORM(c context.Context) dynamoGORMParmas {

	client, err := config.LoadDefaultConfig(c)
	if err != nil {
		panic(err)
	}

	db := dynamodb.NewFromConfig(client)
	return dynamoGORMParmas{
		db: *db,
	}
}

func isExistTableName(tableName string) bool {
	orm := dynamoGORM(context.Background())

	_, err := orm.db.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: wrapString(tableName),
	})

	// not found
	if err != nil {
		// log.Println(err)
		return false
	}

	return true
}

func wrapString(v string) *string {
	return aws.String(v)
}

func wrapInt(v int) *int {
	return aws.Int(v)
}

func wrapBool(v bool) *bool {
	return aws.Bool(v)
}

func getBillingMode(v bool) types.BillingMode {
	if v {
		return types.BillingModePayPerRequest
	}

	return types.BillingModeProvisioned

}
