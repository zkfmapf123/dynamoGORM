package goddb

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gookit/assert"
)

type Message struct {
	PK   string `dynamodbav:"PK"`
	SK   string `dynamodbav:"SK"`
	Name string `dynamodbav:"Name"`
	Age  int    `dynamodbav:"Age"`
}

var (
	ddbClient *dynamodb.Client
	client    *DDBClient
)

func scenarioBeforeHook() {
	awsClient, _ := config.LoadDefaultConfig(context.Background(), config.WithRegion("ap-northeast-2"))
	ddbClient = dynamodb.NewFromConfig(awsClient)
	client = NewDDB(ddbClient)
}

func scenarioAfterHook() {
	client.dropTable("user_logs_1")
	client.dropTable("user_logs_2")
	time.Sleep(10 * time.Second)
}

func Test_DDBCreate(t *testing.T) {

	scenarioBeforeHook()

	t.Run("1. 테이블 생성 확인 여부", func(t *testing.T) {

		err := client.
			AddTable("user_logs_1", DDBTableParams{
				IsCreate:        true,
				IsPK:            true,
				PkAttributeType: types.ScalarAttributeTypeS,
				IsSK:            true,
				SkAttributeType: types.ScalarAttributeTypeS,
				BillingMode: DDBBillingMode{
					isOnDemand: true,
				},
			}).
			AddTable("user_logs_2", DDBTableParams{
				IsCreate:        true,
				IsPK:            true,
				PkAttributeType: types.ScalarAttributeTypeS,
				IsSK:            true,
				SkAttributeType: types.ScalarAttributeTypeS,
				BillingMode: DDBBillingMode{
					isOnDemand: true,
				},
			}).Start(true)

		assert.NoError(t, err)
		time.Sleep(10 * time.Second)
	})

	// DDB 테이블 생성 Wait

	t.Run("2. 테이블 중복 생성 에러 여부", func(t *testing.T) {

		err := client.
			AddTable("user_logs_1", DDBTableParams{
				IsCreate:        true,
				IsPK:            true,
				PkAttributeType: types.ScalarAttributeTypeS,
				IsSK:            true,
				SkAttributeType: types.ScalarAttributeTypeS,
				BillingMode: DDBBillingMode{
					isOnDemand: true,
				},
			}).Start(true)

		assert.Err(t, err)

	})

	t.Run("3. row 단건 추가", func(t *testing.T) {

		err := client.Insert("user_logs_1", Message{
			PK:   "1",
			SK:   "1",
			Name: "test",
			Age:  10,
		})

		assert.NoError(t, err)
	})

	t.Run("4. row 중복 추가할때 에러 여부", func(t *testing.T) {

		err := client.Insert("user_logs_1", Message{
			PK:   "1",
			SK:   "1",
			Name: "test",
			Age:  10,
		})

		assert.Err(t, err)
	})

	t.Run("5. row batch 추가 여부", func(t *testing.T) {

		err := client.InsertBatch("user_logs_1", []any{
			Message{
				PK:   "10",
				SK:   "10",
				Name: "test",
				Age:  10,
			},
			Message{
				PK:   "11",
				SK:   "11",
				Name: "test",
				Age:  10,
			},
		})

		assert.NoError(t, err)
	})

	client.dropTable("user_logs_1")
	client.dropTable("user_logs_2")
	time.Sleep(10 * time.Second)
}

func Test_DDBInfo(t *testing.T) {

	scenarioBeforeHook()
	defer scenarioAfterHook()

	t.Run("0.테이블 생성 확인 여부", func(t *testing.T) {
		err := client.
			AddTable("user_logs_1", DDBTableParams{
				IsCreate:        true,
				IsPK:            true,
				PkAttributeType: types.ScalarAttributeTypeS,
				IsSK:            true,
				SkAttributeType: types.ScalarAttributeTypeS,
				BillingMode: DDBBillingMode{
					isOnDemand: true,
				},
			}).Start(true)
		assert.NoError(t, err)
		time.Sleep(10 * time.Second)
	})

	t.Run("1. 테이블 목록 조회 여부", func(t *testing.T) {
		tables, err := client.GetTables()

		assert.NoError(t, err)
		assert.Contains(t, tables, "user_logs_1")
	})

	t.Run("2. 테이블 상세조회", func(t *testing.T) {
		table, err := client.GetTable("user_logs_1")

		assert.NoError(t, err)
		assert.Eq(t, table.TableName, "user_logs_1")
		assert.Eq(t, table.TableSizeBytes, int64(0))
		assert.Eq(t, table.ItemCount, int64(0))
		assert.Eq(t, table.TableStatus, types.TableStatusActive)
		assert.Eq(t, table.ProvisionedThroughput.ReadCapacityUnits, int64(0))
		assert.Eq(t, table.ProvisionedThroughput.WriteCapacityUnits, int64(0))
		assert.NotNil(t, table.TableArn)
		assert.NotNil(t, table.TableId)
	})

	/*
		시나리오

		현재 회사에 3개의 부서가 존재
		- DEV
		- DESIGN

		3개의 부서의 정보를 기재 함

		- Query 시나리오
			- 특정 사용자 추출
			- 특정 팀의 사용자들 추출 (나이 순서)
	*/

	t.Run("3. Data Insert & Select > 개발그룹 사용자의 데이터 추가 (단일 추가)", func(t *testing.T) {

		var err error

		err = client.Insert("user_logs_1", Message{
			PK:   "USER#1",
			SK:   "#PROFILE",
			Name: "tom",
			Age:  32,
		})

		assert.NoError(t, err)

		err = client.Insert("user_logs_1", Message{
			PK:   "USER#2",
			SK:   "#PROFILE",
			Name: "jerry",
			Age:  24,
		})

		assert.NoError(t, err)

		err = client.Insert("user_logs_1", Message{
			PK:   "USER#3",
			SK:   "#PROFILE",
			Name: "harry",
			Age:  20,
		})

		assert.NoError(t, err)

		err = client.Insert("user_logs_1", Message{
			PK:   "GROUP#DEV",
			SK:   "USER#1",
			Name: "harry",
		})

		assert.NoError(t, err)

		err = client.Insert("user_logs_1", Message{
			PK:   "GROUP#DEV",
			SK:   "USER#2",
			Name: "jerry",
		})

		assert.NoError(t, err)

		err = client.Insert("user_logs_1", Message{
			PK:   "GROUP#DEV",
			SK:   "USER#3",
			Name: "tom",
		})

		assert.NoError(t, err)
	})

	t.Run("4. Data BatchWrite & Select", func(t *testing.T) {
		err := client.InsertBatch("user_logs_1", []any{
			Message{
				PK:   "USER#4",
				SK:   "#PROFILE",
				Name: "james",
				Age:  35,
			},
			Message{
				PK:   "GROUP#DESIGN",
				SK:   "USER#4",
				Name: "james",
			},
			Message{
				PK:   "USER#5",
				SK:   "#PROFILE",
				Name: "john",
				Age:  50,
			},
			Message{
				PK:   "GROUP#DESIGN",
				SK:   "USER#6",
				Name: "john",
			},
		})
		assert.NoError(t, err)
	})

	t.Run("5. Data 단건 조회 pk = 1", func(t *testing.T) {
		item, err := client.FindByKey("user_logs_1", "USER#1", "#PROFILE")
		assert.NoError(t, err)

		result := MarshalMap[Message](item)

		assert.Eq(t, result.PK, "USER#1")
		assert.Eq(t, result.SK, "#PROFILE")
		assert.Eq(t, result.Name, "tom")
		assert.Eq(t, result.Age, 32)
	})

	t.Run("6. 개발팀 조회", func(t *testing.T) {

		items, err := client.FindByKeyUseExpression("user_logs_1", "GROUP#DEV", "", 25, RangeParams{
			KeyConditionExpression: "PK = :pk",
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "GROUP#DEV"},
			},
		})

		assert.NoError(t, err)

		results := MarshalMaps[Message](items)

		assert.Eq(t, len(results), 3)
	})

	t.Run("7. 디자인 팀 조회", func(t *testing.T) {

		items, err := client.FindByKeyUseExpression("user_logs_1", "GROUP#DESIGN", "", 25, RangeParams{
			KeyConditionExpression: "PK = :pk",
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "GROUP#DESIGN"},
			},
		})

		assert.NoError(t, err)

		results := MarshalMaps[Message](items)

		assert.Eq(t, len(results), 2)
	})

	t.Run("8. 개발팀 중 나이가 제일 어린 사람 조회 (비효율적이지만 테스트를 위해...)", func(t *testing.T) {

		items, err := client.FindByKeyUseExpression("user_logs_1", "GROUP#DEV", "", 25, RangeParams{
			KeyConditionExpression: "PK = :pk",
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "GROUP#DEV"},
			},
		})

		assert.NoError(t, err)

		results := MarshalMaps[Message](items)

		name, age := "", math.MaxInt
		for _, result := range results {

			userItme, err := client.FindByKey("user_logs_1", result.SK, "#PROFILE")
			assert.NoError(t, err)

			userResult := MarshalMap[Message](userItme)

			if userResult.Age < age {
				age = userResult.Age
				name = userResult.Name
			}
		}

		assert.Eq(t, name, "harry")
		assert.Eq(t, age, 20)
	})
}
