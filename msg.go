package gostorm

import "encoding/json"

type TaskIds []int
type Values map[string]interface{}

func ParseValues(data []byte) (Values, error) {
	var vals Values
    err := json.Unmarshal(data, &vals)
	return vals, err
}

func (vals Values) Set(name string, val interface{}) {
	vals[name] = val
}

func (vals Values) GetString(name string) (string, bool) {
	if val, ok := vals[name]; ok {
		return val.(string), true
	}
	return "", false
}

func (vals Values) GetStringList(name string) ([]string, bool) {
	if val, ok := vals[name]; ok {
		items := val.([]interface{})
		strs := make([]string, 0, len(items))
		for _, item := range items {
			strs = append(strs, item.(string))
		}
		return strs, true
	}
	return nil, false
}

func (vals Values) GetInt(name string) (int, bool) {
	if val, ok := vals[name]; ok {
		return int(val.(float64)), true
	}
	return 0, false
}

func (vals Values) GetValues(name string) (Values, bool) {
	if val, ok := vals[name]; ok {
		fieldVals := val.(map[string]interface{})
		returnVals := make(Values)
		for key, fVal := range fieldVals {
			returnVals.Set(key, fVal)
		}
		return returnVals, true
	}
	return nil, false
}

func (vals Values) String() string {
	data, err := json.Marshal(vals)
	if err != nil {
		panic(err.Error())
	}
	return string(data)
}

func ParseTaskIds(data []byte) (TaskIds, error) {
	var tids TaskIds
	err := json.Unmarshal(data, &tids)
	return tids, err
}