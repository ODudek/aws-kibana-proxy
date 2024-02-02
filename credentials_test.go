package main

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestPrepareRemoteUrl(t *testing.T) {
	rr := http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme:   "http",
			Path:     "/1234",
			RawPath:  "/1234",
			RawQuery: "q=1",
		},
	}
	type args struct {
		endpoint string
		r        *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name: "test prepare remote url",
			args: args{
				endpoint: "http://localhost:9200",
				r:        &rr,
			},
			want: &url.URL{
				Scheme:   "http",
				Host:     "localhost:9200",
				Path:     "/1234",
				RawPath:  "/1234",
				RawQuery: "q=1",
			},
			wantErr: false,
		},
		{
			name: "test prepare remote url with error",
			args: args{
				endpoint: "ht://localhost:eqwe12",
				r:        &rr,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrepareRemoteUrl(tt.args.endpoint, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrepareRemoteUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrepareRemoteUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignRequest(t *testing.T) {
	requiredHeaders := []string{"Authorization", "X-Amz-Date", "X-Amz-Security-Token"}
	c := AppConfig{
		Port:       ":9200",
		EsEndpoint: "http://localhost:9200",
	}
	rGet := http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme:   "http",
			Path:     "/1234",
			RawPath:  "/1234",
			RawQuery: "q=1",
		},
		Body: io.NopCloser(strings.NewReader("test")),
	}
	//rPost := http.Request{
	//	Method: "POST",
	//	URL: &url.URL{
	//		Scheme:   "http",
	//		Path:     "/1234",
	//		RawPath:  "/1234",
	//		RawQuery: "q=1",
	//	},
	//	Body: nil,
	//}
	type args struct {
		config *AppConfig
		r      *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name: "test sign get request",
			args: args{
				config: &c,
				r:      &rGet,
			},
			want: &http.Request{
				Header: map[string][]string{
					"Authorization": {"AWS4-HMAC-SHA256 Credential=, SignedHeaders=, Signature="},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SignRequest(tt.args.config, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false && (tt.want == nil || tt.want.Header == nil) {
				t.Errorf("Missing required headers")
				return
			}
			for _, header := range requiredHeaders {
				if _, ok := got.Header[header]; !ok {
					t.Errorf("SignRequest() missing required header: %s", header)
					return
				}
				if got.Header.Get(header) == "" && tt.want.Header.Get(header) == "" {
					t.Errorf("SignRequest() got = %v, want %v", got.Header.Get(header), tt.want.Header.Get(header))
					return
				}
			}

		})
	}
}
