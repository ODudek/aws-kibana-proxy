package main

import (
	"net/http"
	"testing"
)

func Test_copyHeader(t *testing.T) {
	type args struct {
		dst http.Header
		src http.Header
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test copy header to destination",
			args: struct {
				dst http.Header
				src http.Header
			}{dst: map[string][]string{
				"key2": {"value"},
			}, src: map[string][]string{
				"key2": {"value"},
			}},
		},
		{
			name: "test don't copy header to destination",
			args: struct {
				dst http.Header
				src http.Header
			}{dst: map[string][]string{
				"key":  {"value"},
				"key2": {"value"},
			}, src: map[string][]string{
				"key2": {"value"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copyHeader(tt.args.dst, tt.args.src)
		})
	}
}
