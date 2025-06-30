package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Execute function that performs mathematical calculations
// This function will be called by gin-mcp when the tool is executed
func Execute(input []byte) ([]byte, error) {
	// Parse the input JSON
	var data map[string]interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	// Extract the expression from arguments
	arguments, ok := data["arguments"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("arguments field is required")
	}

	expression, ok := arguments["expression"].(string)
	if !ok {
		return nil, fmt.Errorf("expression argument is required")
	}

	// Evaluate the expression
	result, err := evaluateExpression(expression)
	if err != nil {
		errorResult := map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Error: %v", err),
				},
			},
		}
		return json.Marshal(errorResult)
	}

	// Return the result in MCP format
	response := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Result: %v", result),
			},
		},
	}

	return json.Marshal(response)
}

// evaluateExpression evaluates a simple mathematical expression
func evaluateExpression(expr string) (float64, error) {
	// Remove whitespace
	expr = strings.ReplaceAll(expr, " ", "")

	// Simple expression evaluator for basic arithmetic
	// This is a simplified version - in production you'd want a proper parser

	// Handle basic operations
	if strings.Contains(expr, "+") {
		parts := strings.Split(expr, "+")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid expression format")
		}
		a, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", parts[0])
		}
		b, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", parts[1])
		}
		return a + b, nil
	}

	if strings.Contains(expr, "-") {
		parts := strings.Split(expr, "-")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid expression format")
		}
		a, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", parts[0])
		}
		b, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", parts[1])
		}
		return a - b, nil
	}

	if strings.Contains(expr, "*") {
		parts := strings.Split(expr, "*")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid expression format")
		}
		a, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", parts[0])
		}
		b, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", parts[1])
		}
		return a * b, nil
	}

	if strings.Contains(expr, "/") {
		parts := strings.Split(expr, "/")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid expression format")
		}
		a, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", parts[0])
		}
		b, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", parts[1])
		}
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	}

	// If no operator found, try to parse as a single number
	result, err := strconv.ParseFloat(expr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid expression: %s", expr)
	}

	return result, nil
}
