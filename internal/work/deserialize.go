package work

import (
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

// Deserialize unmarshals a yaml file into an array of tasks
func Deserialize(pathToYaml string) ([]ParsedTask, error) {
	data, readErr := ioutil.ReadFile(pathToYaml)
	if readErr != nil {
		return nil, readErr
	}

	var unknownTasks []map[string]interface{}
	if err := yaml.Unmarshal(data, &unknownTasks); err != nil {
		return nil, err
	}

	parsedTasks := make([]ParsedTask, 0)
	for _, t := range unknownTasks {
		if t["type"] == "ensure file" {
			var task EnsureFileTask
			if err := mapstructure.Decode(t, &task); err != nil {
				return nil, err
			}

			parsedTasks = append(parsedTasks, ParsedTask{
				Original: t,
				Task:     &task,
			})
		}
	}

	return parsedTasks, nil
}
