package dynamodbgo

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type Users struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	IsDev bool   `json:"is_dev"`
}

type ComplexUser struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Age      int            `json:"age"`
	IsDev    bool           `json:"is_dev"`
	Skills   []string       `json:"skills"`
	Metadata map[string]any `json:"metadata"`
}

func Test_serialize(t *testing.T) {
	t.Run("기본 구조체 직렬화", func(t *testing.T) {
		u := Users{
			Name:  "leedonggyu",
			Age:   32,
			IsDev: false,
		}

		// 구조체를 map[string]any로 변환 후 직렬화
		data := map[string]any{
			"name":   u.Name,
			"age":    u.Age,
			"is_dev": u.IsDev,
		}

		result := serialize(data)

		// 결과 검증
		assert.NotNil(t, result)
		assert.Equal(t, "leedonggyu", result["name"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "32", result["age"].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, false, result["is_dev"].(*types.AttributeValueMemberBOOL).Value)
	})

	t.Run("복잡한 구조체 직렬화", func(t *testing.T) {
		u := ComplexUser{
			ID:     "user123",
			Name:   "leedonggyu",
			Age:    32,
			IsDev:  true,
			Skills: []string{"Go", "Python", "JavaScript"},
			Metadata: map[string]any{
				"department": "Engineering",
				"level":      3,
			},
		}

		data := map[string]any{
			"id":       u.ID,
			"name":     u.Name,
			"age":      u.Age,
			"is_dev":   u.IsDev,
			"skills":   u.Skills,
			"metadata": u.Metadata,
		}

		result := serialize(data)

		// 결과 검증
		assert.NotNil(t, result)
		assert.Equal(t, "user123", result["id"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "leedonggyu", result["name"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "32", result["age"].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, true, result["is_dev"].(*types.AttributeValueMemberBOOL).Value)

		// 슬라이스 검증
		skills := result["skills"].(*types.AttributeValueMemberSS).Value
		assert.Equal(t, []string{"Go", "Python", "JavaScript"}, skills)

		// 중첩 맵 검증
		metadata := result["metadata"].(*types.AttributeValueMemberM).Value
		assert.Equal(t, "Engineering", metadata["department"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "3", metadata["level"].(*types.AttributeValueMemberN).Value)
	})
}

func Test_deserialize(t *testing.T) {
	t.Run("기본 구조체 역직렬화", func(t *testing.T) {
		dynamoData := map[string]types.AttributeValue{
			"name":   &types.AttributeValueMemberS{Value: "leedonggyu"},
			"age":    &types.AttributeValueMemberN{Value: "32"},
			"is_dev": &types.AttributeValueMemberBOOL{Value: false},
		}

		result, err := deserialize[Users](dynamoData)

		assert.NoError(t, err)
		assert.Equal(t, "leedonggyu", result.Name)
		assert.Equal(t, 32, result.Age)
		assert.Equal(t, false, result.IsDev)
	})

	t.Run("복잡한 구조체 역직렬화", func(t *testing.T) {
		dynamoData := map[string]types.AttributeValue{
			"id":     &types.AttributeValueMemberS{Value: "user123"},
			"name":   &types.AttributeValueMemberS{Value: "leedonggyu"},
			"age":    &types.AttributeValueMemberN{Value: "32"},
			"is_dev": &types.AttributeValueMemberBOOL{Value: true},
			"skills": &types.AttributeValueMemberSS{Value: []string{"Go", "Python", "JavaScript"}},
			"metadata": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"department": &types.AttributeValueMemberS{Value: "Engineering"},
				"level":      &types.AttributeValueMemberN{Value: "3"},
			}},
		}

		result, err := deserialize[ComplexUser](dynamoData)

		assert.NoError(t, err)
		assert.Equal(t, "user123", result.ID)
		assert.Equal(t, "leedonggyu", result.Name)
		assert.Equal(t, 32, result.Age)
		assert.Equal(t, true, result.IsDev)
		assert.Equal(t, []string{"Go", "Python", "JavaScript"}, result.Skills)
		assert.Equal(t, "Engineering", result.Metadata["department"])
		assert.Equal(t, float64(3), result.Metadata["level"])
	})

	t.Run("nil 입력 처리", func(t *testing.T) {
		result, err := deserialize[Users](nil)

		assert.Error(t, err)
		assert.Equal(t, "input is nil", err.Error())
		assert.Equal(t, Users{}, result)
	})
}

func Test_SerializeToDynamoDB(t *testing.T) {
	t.Run("구조체를 DynamoDB 형식으로 변환", func(t *testing.T) {
		u := Users{
			Name:  "leedonggyu",
			Age:   32,
			IsDev: false,
		}

		result, err := serializeToDynamoDB(u)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "leedonggyu", result["name"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "32", result["age"].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, false, result["is_dev"].(*types.AttributeValueMemberBOOL).Value)
	})
}

func Test_DeserializeFromDynamoDB(t *testing.T) {
	t.Run("DynamoDB 데이터를 구조체로 변환", func(t *testing.T) {
		dynamoData := map[string]types.AttributeValue{
			"name":   &types.AttributeValueMemberS{Value: "leedonggyu"},
			"age":    &types.AttributeValueMemberN{Value: "32"},
			"is_dev": &types.AttributeValueMemberBOOL{Value: false},
		}

		result, err := deserializeFromDynamoDB[Users](dynamoData)

		assert.NoError(t, err)
		assert.Equal(t, "leedonggyu", result.Name)
		assert.Equal(t, 32, result.Age)
		assert.Equal(t, false, result.IsDev)
	})
}

func Test_ConvertFunctions(t *testing.T) {
	t.Run("문자열 변환", func(t *testing.T) {
		av := &types.AttributeValueMemberS{Value: "test"}
		result, err := convertToString(av)
		assert.NoError(t, err)
		assert.Equal(t, "test", result)
	})

	t.Run("정수 변환", func(t *testing.T) {
		av := &types.AttributeValueMemberN{Value: "123"}
		result, err := convertToInt(av)
		assert.NoError(t, err)
		assert.Equal(t, 123, result)
	})

	t.Run("실수 변환", func(t *testing.T) {
		av := &types.AttributeValueMemberN{Value: "123.45"}
		result, err := convertToFloat64(av)
		assert.NoError(t, err)
		assert.Equal(t, 123.45, result)
	})

	t.Run("불린 변환", func(t *testing.T) {
		av := &types.AttributeValueMemberBOOL{Value: true}
		result, err := convertToBool(av)
		assert.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("문자열 슬라이스 변환", func(t *testing.T) {
		av := &types.AttributeValueMemberSS{Value: []string{"a", "b", "c"}}
		result, err := convertToStringSlice(av)
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("정수 슬라이스 변환", func(t *testing.T) {
		av := &types.AttributeValueMemberNS{Value: []string{"1", "2", "3"}}
		result, err := convertToIntSlice(av)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func Test_deserializeToMap(t *testing.T) {
	t.Run("DynamoDB 데이터를 map으로 변환", func(t *testing.T) {
		dynamoData := map[string]types.AttributeValue{
			"name":   &types.AttributeValueMemberS{Value: "leedonggyu"},
			"age":    &types.AttributeValueMemberN{Value: "32"},
			"is_dev": &types.AttributeValueMemberBOOL{Value: true},
			"skills": &types.AttributeValueMemberSS{Value: []string{"Go", "Python"}},
		}

		result := deserializeToMap(dynamoData)

		assert.NotNil(t, result)
		assert.Equal(t, "leedonggyu", result["name"])
		assert.Equal(t, 32, result["age"])
		assert.Equal(t, true, result["is_dev"])
		assert.Equal(t, []string{"Go", "Python"}, result["skills"])
	})
}

func Test_deserializeSingleValueToAny(t *testing.T) {
	t.Run("단일 값 변환 테스트", func(t *testing.T) {
		tests := []struct {
			name     string
			input    types.AttributeValue
			expected any
		}{
			{
				name:     "문자열",
				input:    &types.AttributeValueMemberS{Value: "test"},
				expected: "test",
			},
			{
				name:     "정수",
				input:    &types.AttributeValueMemberN{Value: "123"},
				expected: 123,
			},
			{
				name:     "실수",
				input:    &types.AttributeValueMemberN{Value: "123.45"},
				expected: 123.45,
			},
			{
				name:     "불린",
				input:    &types.AttributeValueMemberBOOL{Value: true},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := deserializeSingleValueToAny(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}
