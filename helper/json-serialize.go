package helper

import "encoding/json"

func JsonSerialize(state interface{}) string {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	return string(stateJSON)
}
