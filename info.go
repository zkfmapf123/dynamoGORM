package goddb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DDBTableInfoParams struct {
	TableName             string
	TableSizeBytes        int64
	ItemCount             int64
	CreationDateTime      time.Time
	TableStatus           types.TableStatus
	ProvisionedThroughput struct {
		ReadCapacityUnits  int64
		WriteCapacityUnits int64
	}
	TableArn string
	TableId  string
}

func (c DDBClient) GetTables() ([]string, error) {

	output, err := c.client.ListTables(context.Background(), &dynamodb.ListTablesInput{})
	if err != nil {
		c.trace(ERROR, "DDBClient.GetTables.ListTables.Error", map[string]any{
			"error": err,
		})

		return nil, err
	}

	return output.TableNames, nil
}

func (c DDBClient) GetTable(tableName string) (DDBTableInfoParams, error) {

	output, err := c.client.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		c.trace(ERROR, "DDBClient.GetTable.DescribeTable.Error", map[string]any{
			"error": err,
		})
		return DDBTableInfoParams{}, err
	}

	return DDBTableInfoParams{
		TableName:        *output.Table.TableName,
		TableSizeBytes:   *output.Table.TableSizeBytes,
		ItemCount:        *output.Table.ItemCount,
		CreationDateTime: *output.Table.CreationDateTime,
		TableStatus:      output.Table.TableStatus,
		ProvisionedThroughput: struct {
			ReadCapacityUnits  int64
			WriteCapacityUnits int64
		}{
			ReadCapacityUnits:  *output.Table.ProvisionedThroughput.ReadCapacityUnits,
			WriteCapacityUnits: *output.Table.ProvisionedThroughput.WriteCapacityUnits,
		},
		TableArn: *output.Table.TableArn,
		TableId:  *output.Table.TableId,
	}, nil

}
