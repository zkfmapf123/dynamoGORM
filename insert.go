package dynamodbgo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	RETRY_COUNT = 10
)

type TableParmas struct {
	tableName string

	primarykey  string // if already exists not used
	sortedKey   string // if already exists not used
	billingMode bool   // true : on-demain , false : provisioned
}

func Insert(params TableParmas, data map[string]any) error {
	orm := dynamoGORM(context.Background())

	attempt := 1
	for attempt = 1; attempt < RETRY_COUNT; attempt++ {

		if !isExistTableName(params.tableName) {
			createTable(&orm, params)
		}

		if isActiveDynamoTable(params.tableName) {
			break
		}

		time.Sleep(time.Second * 2)
		fmt.Println("테이블 생성 중... ")
	}

	item, err := SerializeToDynamoDB(data)
	if err != nil {
		return err
	}

	_, err = orm.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: &params.tableName,
		Item:      item,
	})

	return err
}

func createTable(orm *dynamoGORMParmas, params TableParmas) error {

	attrDefinition := []types.AttributeDefinition{
		{
			AttributeName: wrapString(params.primarykey),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}

	keySchema := []types.KeySchemaElement{
		{
			AttributeName: wrapString(params.primarykey),
			KeyType:       types.KeyTypeHash,
		},
	}

	// TOBE. 정렬 키는 추후 구성 예정
	_, err := orm.db.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName:            wrapString(params.tableName),
		AttributeDefinitions: attrDefinition,
		KeySchema:            keySchema,
		BillingMode:          getBillingMode(params.billingMode),
	})

	return err
}

func isActiveDynamoTable(tableName string) bool {
	ctx := context.Background()
	orm := dynamoGORM(ctx)

	res, err := orm.db.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: wrapString(tableName),
	})

	if err != nil {
		// log.Println(err)
		return false
	}

	return res.Table.TableStatus == types.TableStatusActive
}
