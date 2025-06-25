package dynamodbgo

import "testing"

func Test_Delete(t *testing.T) {
	err := Delete("users", "user_id", "user_14")
	if err != nil {
		panic(err)
	}
}
