package dynamodbgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsExistsTable(t *testing.T) {
	isOk := isExistTableName("not-found")
	assert.Equal(t, isOk, false)
}

// func Test_isInsertTable(t *testing.T) {
// 	value := map[string]any{
// 		"user_id": "user_1",
// 		"bb":      true,
// 		"cc":      true,
// 	}

// 	for i := 10; i < 20; i++ {

// 		value["user_id"] = fmt.Sprintf("user_%d", i)
// 		err := Insert(TableParmas{
// 			tableName:   "users",
// 			primarykey:  "user_id",
// 			billingMode: true,
// 		}, value)

// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 	}
// }
