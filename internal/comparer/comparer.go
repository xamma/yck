package comparer

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlKeyComparer struct {
	SourceFilePath      string
	CompareFilePath     string
	ValueMismatchEnabled bool
	SourceData          map[string]interface{}
	CompareData         map[string]interface{}
}

func LoadYAML(filePath string) (map[string]interface{}, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %v", err)
	}
	return data, nil
}

func (c *YamlKeyComparer) CompareMaps(source, compare map[string]interface{}, path string) {
	for key, compareValue := range compare {
		currentPath := key
		if path != "" {
			currentPath = path + "." + key
		}

		sourceValue, exists := source[key]
		if !exists {
			fmt.Printf("\033[31m[MISSING] %s is missing in %s\033[0m\n", currentPath, c.SourceFilePath)
			continue
		}

		if sourceValue == nil {
			// fmt.Printf("\033[33m[PRESENT BUT NIL] %s exists but is nil in %s\033[0m\n", currentPath, c.SourceFilePath)
		}

		switch compareValueTyped := compareValue.(type) {
		case map[string]interface{}:
			if sourceMap, ok := sourceValue.(map[string]interface{}); ok {
				if len(sourceMap) == 0 && len(compareValueTyped) > 0 {
					fmt.Printf("\033[33m[VALUE MISMATCH] %s exists but is empty in %s\033[0m\n", currentPath, c.SourceFilePath)
				} else {
					c.CompareMaps(sourceMap, compareValueTyped, currentPath)
				}
			} else {
				fmt.Printf("\033[31m[TYPE MISMATCH] %s expected map, but got different type\033[0m\n", currentPath)
			}

		case []interface{}:
			if sourceList, ok := sourceValue.([]interface{}); ok {
				c.CompareLists(sourceList, compareValueTyped, currentPath)
			} else {
				fmt.Printf("\033[31m[TYPE MISMATCH] %s expected list, but got different type\033[0m\n", currentPath)
			}

		default:
			if c.ValueMismatchEnabled && sourceValue != compareValue {
				fmt.Printf("\033[33m[VALUE MISMATCH] %s\n%s: %v <---> %s: %v\033[0m\n",
					currentPath, c.SourceFilePath, sourceValue, c.CompareFilePath, compareValue)
			}
		}
	}
}

func (c *YamlKeyComparer) CompareLists(source, compare []interface{}, path string) {
	maxLen := len(source)
	if len(compare) > maxLen {
		maxLen = len(compare)
	}

	for i := 0; i < maxLen; i++ {
		currentPath := fmt.Sprintf("%s[%d]", path, i)

		if i >= len(source) {
			fmt.Printf("\033[31m[EXTRA ITEM] %s: %v\033[0m\n", currentPath, compare[i])
			continue
		}
		if i >= len(compare) {
			fmt.Printf("\033[31m[MISSING ITEM] %s: %v\033[0m\n", currentPath, source[i])
			continue
		}

		sourceValue := source[i]
		compareValue := compare[i]

		switch compareValueTyped := compareValue.(type) {
		case map[string]interface{}:
			if sourceMap, ok := sourceValue.(map[string]interface{}); ok {
				c.CompareMaps(sourceMap, compareValueTyped, currentPath)
			} else {
				fmt.Printf("\033[31m[TYPE MISMATCH] %s expected map, but got different type\033[0m\n", currentPath)
			}

		case []interface{}:
			if sourceList, ok := sourceValue.([]interface{}); ok {
				c.CompareLists(sourceList, compareValueTyped, currentPath)
			} else {
				fmt.Printf("\033[31m[TYPE MISMATCH] %s expected list, but got different type\033[0m\n", currentPath)
			}

		default:
			if c.ValueMismatchEnabled && sourceValue != compareValue {
				fmt.Printf("\033[33m[VALUE MISMATCH] %s\n%s: %v <---> %s: %v\033[0m\n",
					currentPath, c.SourceFilePath, sourceValue, c.CompareFilePath, compareValue)
			}
		}
	}
}

func (c *YamlKeyComparer) PrintYAMLDebug(data map[string]interface{}, prefix string) {
	for k, v := range data {
		fmt.Printf("%s%s: %T\n", prefix, k, v)
		if nested, ok := v.(map[string]interface{}); ok {
			c.PrintYAMLDebug(nested, prefix+"  ")
		}
	}
}