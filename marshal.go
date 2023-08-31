package flatted

import "encoding/json"

func FlattedFromStruct(input interface{}) (string, error) {
	str, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return Flatted(string(str)), nil
}

func UnFlattedToStruct(str string, output interface{}) error {
	err := json.Unmarshal([]byte(UnFlatted(str)), &output)
	if err != nil {
		return err
	}
	return nil
}
