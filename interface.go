package goddb

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DDB

type DDBBillingMode struct {
	isOnDemand    bool
	isProvisioned struct {
		ReadCapacityUnits  int
		WriteCapacityUnits int
	}
}

type DDBTableParams struct {
	IsCreate bool // Table 생성 유무

	IsPK            bool                      // primary key
	PkAttributeType types.ScalarAttributeType // Hash

	IsSK            bool // sort key
	SkAttributeType types.ScalarAttributeType

	BillingMode DDBBillingMode
}

type DDBClient struct {
	client *dynamodb.Client
	tables map[string]DDBTableParams
}

// Log
type CustomLogParmas struct {
	ph  string // placeholder
	msg map[string]any
}
