package main_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gtirloni/prometheus-json-exporter"
)

type kvPair struct {
	key   string
	value float64
}

type receiver struct {
	received []kvPair
}

func (r *receiver) Receive(key string, value float64) {
	r.received = append(r.received, kvPair{key, value})
}

func TestWalkJSON(t *testing.T) {
	testData := []struct {
		name     string
		bytes    []byte
		expected []kvPair
	}{
		{
			name:  "float value",
			bytes: []byte(`{"x": 1.0}`),
			expected: []kvPair{
				kvPair{key: "x", value: 1.0},
			},
		},
		{
			name:  "int value",
			bytes: []byte(`{"x": 1}`),
			expected: []kvPair{
				kvPair{key: "x", value: 1.0},
			},
		},
		{
			name:  "bool value",
			bytes: []byte(`{"x": true}`),
			expected: []kvPair{
				kvPair{key: "x", value: 1.0},
			},
		},
		{
			name:     "string value",
			bytes:    []byte(`{"x": "ok"}`),
			expected: nil,
		},
		{
			name:     "null value",
			bytes:    []byte(`{"x": null}`),
			expected: nil,
		},
		{
			name:  "array value",
			bytes: []byte(`{"x": [1, 2, 3]}`),
			expected: []kvPair{
				kvPair{key: "x__0", value: 1},
				kvPair{key: "x__1", value: 2},
				kvPair{key: "x__2", value: 3},
			},
		},
		{
			name:  "nested value",
			bytes: []byte(`{"x": {"y": 1}}`),
			expected: []kvPair{
				kvPair{key: "x.y", value: 1.0},
			},
		},
		{
			name:  "nested^2 value",
			bytes: []byte(`{"x": {"y": {"z": 1}}}`),
			expected: []kvPair{
				kvPair{key: "x.y.z", value: 1.0},
			},
		},
		{
			name:  "array in nested value",
			bytes: []byte(`{"x": {"y": [1, 2, 3]}}`),
			expected: []kvPair{
				kvPair{key: "x.y__0", value: 1},
				kvPair{key: "x.y__1", value: 2},
				kvPair{key: "x.y__2", value: 3},
			},
		},
		{
			name:  "array in array value",
			bytes: []byte(`{"x": [[1, 2], [3, 4]]}`),
			expected: []kvPair{
				kvPair{key: "x__0__0", value: 1},
				kvPair{key: "x__0__1", value: 2},
				kvPair{key: "x__1__0", value: 3},
				kvPair{key: "x__1__1", value: 4},
			},
		},
		{
			name:  "array at root",
			bytes: []byte(`[1, 2, 3]`),
			expected: []kvPair{
				kvPair{key: "__0", value: 1},
				kvPair{key: "__1", value: 2},
				kvPair{key: "__2", value: 3},
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var jsonData interface{}
			err := json.Unmarshal(tt.bytes, &jsonData)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			r := &receiver{}
			main.WalkJSON("", jsonData, r)
			if !reflect.DeepEqual(r.received, tt.expected) {
				t.Errorf("Got: %#v, expected: %#v", r.received, tt.expected)
			}
		})
	}
}
