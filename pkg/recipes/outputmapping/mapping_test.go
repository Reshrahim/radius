/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package outputmapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApply(t *testing.T) {
	tests := []struct {
		name            string
		rawOutputs      map[string]OutputValue
		outputsMap      map[string]string
		expectedValues  map[string]any
		expectedSecrets map[string]any
	}{
		{
			name:            "nil outputs",
			rawOutputs:      nil,
			outputsMap:      nil,
			expectedValues:  map[string]any{},
			expectedSecrets: map[string]any{},
		},
		{
			name:            "empty outputs",
			rawOutputs:      map[string]OutputValue{},
			outputsMap:      nil,
			expectedValues:  map[string]any{},
			expectedSecrets: map[string]any{},
		},
		{
			name: "pass-through when no mapping (non-sensitive)",
			rawOutputs: map[string]OutputValue{
				"fqdn": {Value: "db.example.com", Sensitive: false},
				"port": {Value: 5432, Sensitive: false},
			},
			outputsMap:      nil,
			expectedValues:  map[string]any{"fqdn": "db.example.com", "port": 5432},
			expectedSecrets: map[string]any{},
		},
		{
			name: "pass-through routes sensitive to secrets",
			rawOutputs: map[string]OutputValue{
				"host":     {Value: "db.example.com", Sensitive: false},
				"password": {Value: "secret123", Sensitive: true},
			},
			outputsMap:      nil,
			expectedValues:  map[string]any{"host": "db.example.com"},
			expectedSecrets: map[string]any{"password": "secret123"},
		},
		{
			name: "mapping renames outputs to resource properties",
			rawOutputs: map[string]OutputValue{
				"fqdn":        {Value: "db.example.com", Sensitive: false},
				"listen_port": {Value: 5432, Sensitive: false},
				"extra":       {Value: "not-mapped", Sensitive: false},
			},
			outputsMap: map[string]string{
				"host": "fqdn",
				"port": "listen_port",
			},
			expectedValues:  map[string]any{"host": "db.example.com", "port": 5432},
			expectedSecrets: map[string]any{},
		},
		{
			name: "mapping with sensitive output",
			rawOutputs: map[string]OutputValue{
				"connection_string": {Value: "Server=...;Password=...", Sensitive: true},
				"host":              {Value: "db.example.com", Sensitive: false},
			},
			outputsMap: map[string]string{
				"connectionString": "connection_string",
				"hostname":         "host",
			},
			expectedValues:  map[string]any{"hostname": "db.example.com"},
			expectedSecrets: map[string]any{"connectionString": "Server=...;Password=..."},
		},
		{
			name: "mapping with literal value (no matching module output)",
			rawOutputs: map[string]OutputValue{
				"host": {Value: "db.example.com", Sensitive: false},
			},
			outputsMap: map[string]string{
				"hostname": "host",
				"port":     "5432", // not a module output → literal
			},
			expectedValues:  map[string]any{"hostname": "db.example.com", "port": "5432"},
			expectedSecrets: map[string]any{},
		},
		{
			name: "mixed module outputs and literal values",
			rawOutputs: map[string]OutputValue{
				"fqdn": {Value: "pg-postgresql.postgres.database.azure.com", Sensitive: false},
				"name": {Value: "pg-postgresql", Sensitive: false},
			},
			outputsMap: map[string]string{
				"host":     "fqdn",      // module output reference
				"database": "name",      // module output reference
				"username": "pgadmin",   // literal value
				"password": "test123",   // literal value
			},
			expectedValues: map[string]any{
				"host":     "pg-postgresql.postgres.database.azure.com",
				"database": "pg-postgresql",
				"username": "pgadmin",
				"password": "test123",
			},
			expectedSecrets: map[string]any{},
		},
		{
			name: "empty mapping with outputs (pass-through)",
			rawOutputs: map[string]OutputValue{
				"host": {Value: "db.example.com", Sensitive: false},
			},
			outputsMap:      map[string]string{},
			expectedValues:  map[string]any{"host": "db.example.com"},
			expectedSecrets: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, secrets := Apply(tt.rawOutputs, tt.outputsMap)
			assert.Equal(t, tt.expectedValues, values)
			assert.Equal(t, tt.expectedSecrets, secrets)
		})
	}
}
