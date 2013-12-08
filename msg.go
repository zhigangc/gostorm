package gostorm

import "bytes"
import "strconv"
import "encoding/json"

type TaskIds []int
type Values map[string]interface{}

func ParseValues(data []byte) (Values, error) {
	var vals Values
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	err := dec.Decode(&vals)
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

func (vals Values) GetList(name string) ([]interface{}, bool) {
	if val, ok := vals[name]; ok {
		return val.([]interface{}), true
	}
	return nil, false
}

func (vals Values) GetInt(name string) (int, bool) {
	if val, ok := vals[name]; ok {
		num, err := strconv.ParseInt(string(val.(json.Number)), 10, 64)
		if err == nil {
			return int(num), true
		}
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