package dynamodbgo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ResponseParams struct {
	UserId string `json:"user_id"`
	BB     bool   `json:"bb"`
	CC     bool   `json:"cc"`
}

func Test_FindById(t *testing.T) {

	data, err := FindByUsePK[ResponseParams]("users", "user_id", "user_14")
	if err != nil {
		panic(err)
	}
	fmt.Println(data)

	assert.Equal(t, data.UserId, "user_14")
	assert.Equal(t, data.BB, true)
	assert.Equal(t, data.CC, true)
}

func Test_FindByList(t *testing.T) {

	datas, err := SelectAll[ResponseParams]("users")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(datas) > 0, true)
}
