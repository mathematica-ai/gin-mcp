package handlers

import (
	"encoding/json"
	"testing"
)

func TestMCPHandler_ValidateInput(t *testing.T) {
	handler := NewMCPHandler()

	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "valid JSON",
			input:   []byte(`{"key": "value"}`),
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{"key": "value"`),
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "null input",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMCPHandler_ValidateOutput(t *testing.T) {
	handler := NewMCPHandler()

	tests := []struct {
		name    string
		output  []byte
		wantErr bool
	}{
		{
			name:    "valid JSON",
			output:  []byte(`{"result": "success"}`),
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			output:  []byte(`{"result": "success"`),
			wantErr: true,
		},
		{
			name:    "empty output",
			output:  []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateOutput(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMCPHandler_ValidateAndFormatOutput(t *testing.T) {
	handler := NewMCPHandler()

	tests := []struct {
		name     string
		output   []byte
		wantJSON bool
	}{
		{
			name:     "valid JSON output",
			output:   []byte(`{"result": "success"}`),
			wantJSON: true,
		},
		{
			name:     "plain text output",
			output:   []byte("Hello, World!"),
			wantJSON: false,
		},
		{
			name:     "empty output",
			output:   []byte{},
			wantJSON: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handler.validateAndFormatOutput(tt.output)
			if err != nil {
				t.Errorf("validateAndFormatOutput() error = %v", err)
				return
			}

			// Check if result is valid JSON
			var jsonData interface{}
			jsonErr := json.Unmarshal(result, &jsonData)

			if tt.wantJSON && jsonErr != nil {
				t.Errorf("Expected valid JSON but got error: %v", jsonErr)
			}

			if !tt.wantJSON && jsonErr == nil {
				// Should be formatted as MCP content
				var content map[string]interface{}
				if err := json.Unmarshal(result, &content); err != nil {
					t.Errorf("Expected MCP content format but got error: %v", err)
				}

				if _, exists := content["content"]; !exists {
					t.Errorf("Expected 'content' field in MCP format")
				}
			}
		})
	}
}
