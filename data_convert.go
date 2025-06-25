package dynamodbgo

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDB AttributeValue를 제네릭 타입 T로 변환
func deserialize[T any](v map[string]types.AttributeValue) (T, error) {
	var result T

	if v == nil {
		return result, fmt.Errorf("input is nil")
	}

	resultType := reflect.TypeOf(result)
	resultValue := reflect.ValueOf(&result).Elem()

	if resultType.Kind() == reflect.Struct {
		for i := 0; i < resultType.NumField(); i++ {
			field := resultType.Field(i)
			fieldValue := resultValue.Field(i)

			if !fieldValue.CanSet() {
				continue
			}

			jsonTag := field.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = field.Name
			}

			if attrValue, exists := v[jsonTag]; exists {

				convertedValue, err := convertToFieldType(attrValue, field.Type)
				if err != nil {
					return result, fmt.Errorf("failed to convert field %s: %w", field.Name, err)
				}

				if convertedValue != nil {

					convertedReflectValue := reflect.ValueOf(convertedValue)
					if convertedReflectValue.Type().AssignableTo(field.Type) {
						fieldValue.Set(convertedReflectValue)
					} else {

						if fieldValue.CanAddr() {
							fieldValue.Set(convertedReflectValue.Convert(field.Type))
						}
					}
				}
			}
		}
	} else {
		return deserializeViaJSON[T](v)
	}

	return result, nil
}

// AttributeValue를 특정 타입으로 변환
func convertToFieldType(av types.AttributeValue, targetType reflect.Type) (any, error) {
	switch v := av.(type) {
	case *types.AttributeValueMemberS:
		switch targetType.Kind() {
		case reflect.String:
			return v.Value, nil
		case reflect.Int:
			if intVal, err := strconv.Atoi(v.Value); err == nil {
				return intVal, nil
			}
		case reflect.Int64:
			if intVal, err := strconv.ParseInt(v.Value, 10, 64); err == nil {
				return intVal, nil
			}
		case reflect.Float64:
			if floatVal, err := strconv.ParseFloat(v.Value, 64); err == nil {
				return floatVal, nil
			}
		case reflect.Bool:
			if boolVal, err := strconv.ParseBool(v.Value); err == nil {
				return boolVal, nil
			}
		}
		return v.Value, nil

	case *types.AttributeValueMemberN:
		switch targetType.Kind() {
		case reflect.String:
			return v.Value, nil
		case reflect.Int:
			if intVal, err := strconv.Atoi(v.Value); err == nil {
				return intVal, nil
			}
		case reflect.Int64:
			if intVal, err := strconv.ParseInt(v.Value, 10, 64); err == nil {
				return intVal, nil
			}
		case reflect.Float64:
			if floatVal, err := strconv.ParseFloat(v.Value, 64); err == nil {
				return floatVal, nil
			}
		}
		return v.Value, nil

	case *types.AttributeValueMemberBOOL:
		return v.Value, nil

	case *types.AttributeValueMemberB:
		return v.Value, nil

	case *types.AttributeValueMemberSS:
		return v.Value, nil

	case *types.AttributeValueMemberNS:
		numbers := make([]int, len(v.Value))
		for i, numStr := range v.Value {
			if num, err := strconv.Atoi(numStr); err == nil {
				numbers[i] = num
			}
		}
		return numbers, nil

	case *types.AttributeValueMemberM:
		nestedMap := deserializeToMap(v.Value)
		return nestedMap, nil

	case *types.AttributeValueMemberL:
		list := make([]any, len(v.Value))
		for i, val := range v.Value {
			list[i] = deserializeSingleValueToAny(val)
		}
		return list, nil

	default:
		return fmt.Sprintf("%v", av), nil
	}
}

// JSON을 통한 변환 (구조체가 아닌 경우)
func deserializeViaJSON[T any](v map[string]types.AttributeValue) (T, error) {
	var result T

	intermediate := deserializeToMap(v)

	jsonBytes, err := json.Marshal(intermediate)
	if err != nil {
		return result, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal to type %T: %w", result, err)
	}

	return result, nil
}

// DynamoDB AttributeValue를 map[string]any로 변환
func deserializeToMap(v map[string]types.AttributeValue) map[string]any {
	result := make(map[string]any)

	for key, value := range v {
		result[key] = deserializeSingleValueToAny(value)
	}

	return result
}

// 단일 AttributeValue를 any로 변환
func deserializeSingleValueToAny(value types.AttributeValue) any {
	switch v := value.(type) {
	case *types.AttributeValueMemberS:
		return v.Value
	case *types.AttributeValueMemberN:
		// 숫자로 파싱 시도
		if intVal, err := strconv.Atoi(v.Value); err == nil {
			return intVal
		} else if floatVal, err := strconv.ParseFloat(v.Value, 64); err == nil {
			return floatVal
		}
		return v.Value
	case *types.AttributeValueMemberBOOL:
		return v.Value
	case *types.AttributeValueMemberB:
		return v.Value
	case *types.AttributeValueMemberSS:
		return v.Value
	case *types.AttributeValueMemberNS:
		numbers := make([]int, len(v.Value))
		for i, numStr := range v.Value {
			if num, err := strconv.Atoi(numStr); err == nil {
				numbers[i] = num
			}
		}
		return numbers
	case *types.AttributeValueMemberM:
		return deserializeToMap(v.Value)
	case *types.AttributeValueMemberL:
		list := make([]any, len(v.Value))
		for i, val := range v.Value {
			list[i] = deserializeSingleValueToAny(val)
		}
		return list
	default:
		return fmt.Sprintf("%v", value)
	}
}

// Go 데이터를 DynamoDB AttributeValue로 변환
func serialize(data map[string]any) map[string]types.AttributeValue {
	item := make(map[string]types.AttributeValue)

	for key, value := range data {
		switch v := value.(type) {
		case string:
			item[key] = &types.AttributeValueMemberS{Value: v}
		case int:
			item[key] = &types.AttributeValueMemberN{Value: strconv.Itoa(v)}
		case int64:
			item[key] = &types.AttributeValueMemberN{Value: strconv.FormatInt(v, 10)}
		case float64:
			item[key] = &types.AttributeValueMemberN{Value: strconv.FormatFloat(v, 'f', -1, 64)}
		case bool:
			item[key] = &types.AttributeValueMemberBOOL{Value: v}
		case []byte:
			item[key] = &types.AttributeValueMemberB{Value: v}
		case []string:
			item[key] = &types.AttributeValueMemberSS{Value: v}
		case []int:
			numbers := make([]string, len(v))
			for i, num := range v {
				numbers[i] = strconv.Itoa(num)
			}
			item[key] = &types.AttributeValueMemberNS{Value: numbers}
		case map[string]any:
			nestedMap := serialize(v)
			item[key] = &types.AttributeValueMemberM{Value: nestedMap}
		case []any:
			list := make([]types.AttributeValue, len(v))
			for i, val := range v {
				switch listVal := val.(type) {
				case string:
					list[i] = &types.AttributeValueMemberS{Value: listVal}
				case int:
					list[i] = &types.AttributeValueMemberN{Value: strconv.Itoa(listVal)}
				case bool:
					list[i] = &types.AttributeValueMemberBOOL{Value: listVal}
				default:
					list[i] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", val)}
				}
			}
			item[key] = &types.AttributeValueMemberL{Value: list}
		default:
			item[key] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", value)}
		}
	}

	return item
}

// 편의 함수들
func serializeToDynamoDB[T any](data T) (map[string]types.AttributeValue, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct: %w", err)
	}

	var mapData map[string]any
	err = json.Unmarshal(jsonBytes, &mapData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal to map: %w", err)
	}

	return serialize(mapData), nil
}

func deserializeFromDynamoDB[T any](v map[string]types.AttributeValue) (T, error) {
	return deserialize[T](v)
}

func convertToString(av types.AttributeValue) (string, error) {
	switch v := av.(type) {
	case *types.AttributeValueMemberS:
		return v.Value, nil
	case *types.AttributeValueMemberN:
		return v.Value, nil
	case *types.AttributeValueMemberBOOL:
		return fmt.Sprintf("%t", v.Value), nil
	default:
		return fmt.Sprintf("%v", av), nil
	}
}

func convertToInt(av types.AttributeValue) (int, error) {
	switch v := av.(type) {
	case *types.AttributeValueMemberS:
		return strconv.Atoi(v.Value)
	case *types.AttributeValueMemberN:
		return strconv.Atoi(v.Value)
	default:
		return 0, fmt.Errorf("cannot convert to int: %v", av)
	}
}

func convertToFloat64(av types.AttributeValue) (float64, error) {
	switch v := av.(type) {
	case *types.AttributeValueMemberS:
		return strconv.ParseFloat(v.Value, 64)
	case *types.AttributeValueMemberN:
		return strconv.ParseFloat(v.Value, 64)
	default:
		return 0, fmt.Errorf("cannot convert to float64: %v", av)
	}
}

func convertToBool(av types.AttributeValue) (bool, error) {
	switch v := av.(type) {
	case *types.AttributeValueMemberBOOL:
		return v.Value, nil
	case *types.AttributeValueMemberS:
		return strconv.ParseBool(v.Value)
	default:
		return false, fmt.Errorf("cannot convert to bool: %v", av)
	}
}

func convertToStringSlice(av types.AttributeValue) ([]string, error) {
	switch v := av.(type) {
	case *types.AttributeValueMemberSS:
		return v.Value, nil
	case *types.AttributeValueMemberL:
		result := make([]string, len(v.Value))
		for i, val := range v.Value {
			if str, err := convertToString(val); err == nil {
				result[i] = str
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert to string slice: %v", av)
	}
}

func convertToIntSlice(av types.AttributeValue) ([]int, error) {
	switch v := av.(type) {
	case *types.AttributeValueMemberNS:
		result := make([]int, len(v.Value))
		for i, numStr := range v.Value {
			if num, err := strconv.Atoi(numStr); err == nil {
				result[i] = num
			}
		}
		return result, nil
	case *types.AttributeValueMemberL:
		result := make([]int, len(v.Value))
		for i, val := range v.Value {
			if num, err := convertToInt(val); err == nil {
				result[i] = num
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert to int slice: %v", av)
	}
}
