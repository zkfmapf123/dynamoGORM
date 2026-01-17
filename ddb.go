package goddb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	PrimaryKey = "PK"
	SortKey    = "SK"
)

func NewDDB(dynamoDBClient *dynamodb.Client) *DDBClient {
	return &DDBClient{
		client: dynamoDBClient,
		tables: map[string]DDBTableParams{},
	}
}

func (c *DDBClient) AddTable(tableName string, table DDBTableParams) *DDBClient {

	c.tables[tableName] = table
	return c
}

func (c *DDBClient) Start(isCreateTable bool) error {

	InfoLog(CustomLogParmas{
		ph: "DDBClient.Start",
		msg: map[string]any{
			"totalTableCount ": len(c.tables),
			"isCreate":         isCreateTable,
		},
	})

	for tableName, params := range c.tables {

		if params.IsCreate {

			InfoLog(CustomLogParmas{
				ph: "DDBClient.Start.CreateTable",
				msg: map[string]any{
					"tableName": tableName,
				},
			})

			keySchema, keyAttribute := getPKandSK(params)

			createTableInput := &dynamodb.CreateTableInput{
				TableName:            aws.String(tableName),
				KeySchema:            keySchema,
				AttributeDefinitions: keyAttribute,
			}

			// ondemand
			if params.BillingMode.isOnDemand {
				createTableInput.BillingMode = getBillingMode(params.BillingMode)
			}

			if !params.BillingMode.isOnDemand {
				createTableInput.BillingMode = getBillingMode(params.BillingMode)
				createTableInput.ProvisionedThroughput = &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(int64(params.BillingMode.isProvisioned.ReadCapacityUnits)),
					WriteCapacityUnits: aws.Int64(int64(params.BillingMode.isProvisioned.WriteCapacityUnits)),
				}
			}

			_, err := c.client.CreateTable(context.Background(), createTableInput)

			if err != nil {
				ErrorLog(CustomLogParmas{
					ph: "DDBClient.Start.CreateTable.Error",
					msg: map[string]any{
						"tableName": tableName,
						"error":     err,
					},
				})

				return err
			}

			InfoLog(CustomLogParmas{
				ph: "DDBClient.Start.CreateTable.Success",
				msg: map[string]any{
					"tableName": tableName,
				},
			})
		}
	}

	return nil
}

func getPKandSK(params DDBTableParams) ([]types.KeySchemaElement, []types.AttributeDefinition) {
	keySchema := []types.KeySchemaElement{}
	keyAttribute := []types.AttributeDefinition{}

	// use pk
	if params.IsPK {

		keySchema = append(keySchema, types.KeySchemaElement{
			AttributeName: aws.String(PrimaryKey),
			KeyType:       types.KeyTypeHash,
		})

		keyAttribute = append(keyAttribute, types.AttributeDefinition{
			AttributeName: aws.String(PrimaryKey),
			AttributeType: params.PkAttributeType,
		})

	}

	// use sk
	if params.IsSK {

		keySchema = append(keySchema, types.KeySchemaElement{
			AttributeName: aws.String(SortKey),
			KeyType:       types.KeyTypeRange,
		})

		keyAttribute = append(keyAttribute, types.AttributeDefinition{
			AttributeName: aws.String(SortKey),
			AttributeType: params.SkAttributeType,
		})

	}

	return keySchema, keyAttribute
}

func getBillingMode(billingMode DDBBillingMode) types.BillingMode {
	if billingMode.isOnDemand {
		return types.BillingModePayPerRequest
	}

	return types.BillingModeProvisioned
}
