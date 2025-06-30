package main

import "encoding/json"

// Predict function that adds two numbers
// This function will be called by gin-mcp when the model is executed
func Predict(input []byte) ([]byte, error) {
	// Parse the input JSON
	var data map[string]int
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, err
	}

	// Extract the numbers to add
	a, existsA := data["a"]
	b, existsB := data["b"]

	if !existsA || !existsB {
		errorResult := map[string]string{
			"error": "Both 'a' and 'b' parameters are required",
		}
		return json.Marshal(errorResult)
	}

	// Perform the addition
	result := map[string]int{
		"sum": a + b,
		"a":   a,
		"b":   b,
	}

	// Return the result as JSON
	return json.Marshal(result)
}
