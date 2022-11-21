package operate_log

import (
	"errors"
	"fmt"
	"strings"

	"im/pkg/utils"

	"github.com/mitchellh/mapstructure"
)

type MappingField map[string]string

const (
	KEY_LOG           = "operatelog"
	KEY_BEFORE        = "operatelog-fields-before"
	KEY_AFTER         = "operatelog-fields-after"
	KEY_NEED_FIELDS   = "operatelog-fields-need"
	KEY_EXPECT_FIELDS = "operatelog-fields-expect"
)

var (
	defaultLang  = "cn"
	fieldMapping map[string]MappingField
	joinMapping  map[string]MappingField
	lang         = defaultLang
)

func SetLang(val string) {
	lang = val
}

func SetFieldMapping(val map[string]MappingField) {
	fieldMapping = val
}

func SetJoinMapping(val map[string]MappingField) {
	joinMapping = val
}

func CreateFieldsLog(fields interface{}, needFields, expectFields []string) (string, error) {
	result := ""

	fieldsMap, err := decodeFields(fields)
	if err != nil {
		return result, errors.New("decodeFields " + err.Error())
	}
	fieldsMapRes := filtrationFields(fieldsMap, needFields, expectFields)
	fieldsRes, err := translateFields(fieldsMapRes)
	if err != nil {
		return result, errors.New("Translate beforeFields " + err.Error())
	}

	for key := range fieldsRes {
		if fmt.Sprintf("%v", fieldsRes[key]) == "" {
			continue
		}
		result += fmt.Sprintf("%v%v%v%v%v",
			key,
			joinMapping["fieldBefore"][lang],
			fieldsRes[key],
			joinMapping["fieldAfter"][lang],
			joinMapping["split"][lang],
		)
	}
	result = strings.TrimRight(result, joinMapping["split"][lang])
	result += joinMapping["end"][lang]
	return result, nil
}

func UpdateFieldsLog(beforeFields interface{}, afterFields interface{}, needFields, expectFields []string) (string, error) {
	result := ""

	fieldsMap, err := decodeFields(beforeFields)
	if err != nil {
		return result, errors.New("decodeFields " + err.Error())
	}
	fieldsMapRes := filtrationFields(fieldsMap, needFields, expectFields)
	before, err := translateFields(fieldsMapRes)
	if err != nil {
		return result, errors.New("Translate beforeFields " + err.Error())
	}

	fieldsMap, err = decodeFields(afterFields)
	if err != nil {
		return result, errors.New("decodeFields " + err.Error())
	}
	fieldsMapRes = filtrationFields(fieldsMap, needFields, expectFields)
	after, err := translateFields(fieldsMapRes)
	if err != nil {
		return result, errors.New("Translate afterFields " + err.Error())
	}

	for key := range before {
		if _, ok := after[key]; !ok {
			continue
		}
		if fmt.Sprintf("%v", before[key]) == fmt.Sprintf("%v", after[key]) {
			continue
		}
		result += fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v",
			key,
			joinMapping["before"][lang],
			joinMapping["fieldBefore"][lang],
			before[key],
			joinMapping["fieldAfter"][lang],
			joinMapping["after"][lang],
			joinMapping["fieldBefore"][lang],
			after[key],
			joinMapping["fieldAfter"][lang],
			joinMapping["split"][lang],
		)
	}
	result = strings.TrimRight(result, joinMapping["split"][lang])
	result += joinMapping["end"][lang]
	return result, nil
}

func GetFieldsLogSlice(beforeFields interface{}, afterFields interface{}, needFields, expectFields []string) ([]string, []string, []string, error) {
	fieldsSlice, beforeSlice, afterSlice := []string{}, []string{}, []string{}
	fieldsMap, err := decodeFields(beforeFields)
	if err != nil {
		return fieldsSlice, beforeSlice, afterSlice, errors.New("decodeFields " + err.Error())
	}
	fieldsMapRes := filtrationFields(fieldsMap, needFields, expectFields)
	before, err := translateFields(fieldsMapRes)
	if err != nil {
		return fieldsSlice, beforeSlice, afterSlice, errors.New("Translate beforeFields " + err.Error())
	}

	fieldsMap, err = decodeFields(afterFields)
	if err != nil {
		return fieldsSlice, beforeSlice, afterSlice, errors.New("decodeFields " + err.Error())
	}
	fieldsMapRes = filtrationFields(fieldsMap, needFields, expectFields)
	after, err := translateFields(fieldsMapRes)
	if err != nil {
		return fieldsSlice, beforeSlice, afterSlice, errors.New("Translate afterFields " + err.Error())
	}

	fields := []string{}
	for k := range before {
		fields = append(fields, k)
	}
	for k := range after {
		fields = append(fields, k)
	}
	fields = utils.RemoveRepeat(fields)
	for _, v := range fields {
		beforeVal, beforeOK := before[v]
		afterVal, afterOK := after[v]
		if !beforeOK && afterOK {
			if afterVal != "" {
				afterValStr := fmt.Sprintf("%v", afterVal)
				afterSlice = append(afterSlice, afterValStr)
				fieldsSlice = append(fieldsSlice, v)
			}
		} else if beforeOK && afterOK {
			if fmt.Sprintf("%v", beforeVal) != fmt.Sprintf("%v", afterVal) {
				beforeValStr := fmt.Sprintf("%v", beforeVal)
				beforeSlice = append(beforeSlice, beforeValStr)
				afterValStr := fmt.Sprintf("%v", afterVal)
				afterSlice = append(afterSlice, afterValStr)
				fieldsSlice = append(fieldsSlice, v)
			}
		}
	}

	return fieldsSlice, beforeSlice, afterSlice, nil
}

func decodeFields(fields interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{}, 0)
	err := mapstructure.Decode(fields, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func translateFields(fields map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for key, value := range fields {
		key = utils.ToSnakeCase(key)
		if object, ok := fieldMapping[key][lang]; ok {
			result[object] = value
		} else {
			result[key] = value
		}
	}
	return result, nil
}

func filtrationFields(fieldsMap map[string]interface{}, needFields, expectFields []string) map[string]interface{} {
	fieldsMapRes := make(map[string]interface{})
	for key, val := range fieldsMap {
		if len(needFields) > 0 && utils.Contains(needFields, utils.ToSnakeCase(key)) == -1 {
			continue
		}
		if len(expectFields) > 0 && utils.Contains(expectFields, utils.ToSnakeCase(key)) > -1 {
			continue
		}
		fieldsMapRes[key] = val
	}
	return fieldsMapRes
}

// init default mapping
func init() {
	lang = "cn"
	joinMapping = map[string]MappingField{
		"before": {
			"cn": "由",
		},
		"after": {
			"cn": "改为",
		},
		"split": {
			"cn": "，",
		},
		"end": {
			"cn": "",
		},
		"fieldBefore": {
			"cn": "[",
		},
		"fieldAfter": {
			"cn": "]",
		},
	}
	fieldMapping = map[string]MappingField{
		"username": {
			"cn": "用户名",
		},
		"mobile": {
			"cn": "联系电话",
		},
		"email": {
			"cn": "电子邮件",
		},
	}
}
