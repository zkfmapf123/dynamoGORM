package goddb

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	awsClient *aws.Config
	ddbClient *dynamodb.Client
)

func beforeHook() {
	awsClient, _ := config.LoadDefaultConfig(context.Background(), config.WithRegion("ap-northeast-2"))
	ddbClient = dynamodb.NewFromConfig(awsClient)
}

func Test_DDBStart(t *testing.T) {
	beforeHook()

	ddb := NewDDB(ddbClient)
	ddb.
		AddTable("user_logs", DDBTableParams{
			IsCreate:        true,
			IsPK:            true,
			PkAttributeType: types.ScalarAttributeTypeS,
			IsSK:            true,
			SkAttributeType: types.ScalarAttributeTypeS,
			BillingMode: DDBBillingMode{
				isOnDemand: true,
			},
		}).Start(true)
}
