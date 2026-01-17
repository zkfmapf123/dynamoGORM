package goddb

import (
	"context"
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

	client.DropTable("user_logs_1")
	client.DropTable("user_logs_2")
	time.Sleep(10 * time.Second)
}
