package tool

import "fmt"

func BodyRequestHelper(requestBody []ToolInteractionElement, id string) (any, error) {
	for _, entry := range requestBody {
		if entry.Interface_id == id {
			return entry.Content, nil
		}
	}
	return nil, fmt.Errorf("missisng required field: %s", id)
}
