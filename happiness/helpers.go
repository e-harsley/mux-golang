package happiness

import "encoding/json"

func MapDump(data interface{}) (map[string]interface{}, error) {
	var mod map[string]interface{}
	jsonString, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonString, &mod); err != nil {
		return mod, err
	}
	return mod, nil

}
