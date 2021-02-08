package read

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type ServiceCfg struct {
	GRPC ServiceGRPCCfg `json:"grpc"`
}

type ServiceGRPCCfg struct {
	Server ServerCfg `json:"server"`
}

type ServerCfg struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

func yamlValueToJsonValue(value interface{}) (interface{}, error) {
	var result interface{}

	switch v := value.(type) {
	case map[interface{}]interface{}:
		// JSON objects

		mapValue := make(map[string]interface{}, len(v))

		for key, entryValue := range v {
			keyString, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("key %q is not a string", key)
			}

			entryValue2, err := yamlValueToJsonValue(entryValue)
			if err != nil {
				return nil, err
			}

			mapValue[keyString] = entryValue2
		}

		result = mapValue

	case []interface{}:
		// JSON arrays
		sliceValue := make([]interface{}, len(v))

		for i, elementValue := range v {
			elementValue2, err := yamlValueToJsonValue(elementValue)
			if err != nil {
				return nil, err
			}

			sliceValue[i] = elementValue2
		}

		result = sliceValue

	default:
		Result = Value
	}

	return result, nil
}

func LoadCfg(filePath string, cfg *ServiceCfg) error {
	// Load the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open %s: %w", filePath, err)
	}
	defer file.Close()

	tplData, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", filePath, err)
	}

	// Parse YAML data
	var yamlValue interface{}
	err1 := yaml.Unmarshal(tplData, &yamlValue)
	if err1 != nil {
		return fmt.Errorf("cannot parse %s: %w", filePath, err1)
	}

	// Convert to valid JSON
	jsonValue, err := yamlValueToJsonValue(yamlValue)
	if err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}

	// Generate and parse JSON data
	readCfg := func(value interface{}, dest interface{}) error {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("cannot genrerate json: %w", err)
		}

		if err := json.Unmarshal(data, dest); err != nil {
			return fmt.Errorf("cannot parse json: %w", err)
		}

		return nil
	}

	jsonValue, ok := jsonValue.(map[string]interface{})
	if !ok {
		return fmt.Errorf("top-level element is not an object")
	}

	if err := readCfg(jsonValue, &cfg); err != nil {
		return fmt.Errorf("cannot read service configuration: %w", err)
	}

	return nil
}
