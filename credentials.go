package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"io"
	"net/http"
	"net/url"
	"time"
)

func GetAWSConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO())
}

func CalculateSHA256Hash(reader io.Reader) (string, error) {

	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)
	hashHex := hex.EncodeToString(hashBytes)

	return hashHex, nil
}

func PrepareRemoteUrl(endpoint string, r *http.Request) (*url.URL, error) {
	remote, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	remote.Path = r.URL.Path
	remote.RawPath = r.URL.RawPath
	remote.RawQuery = r.URL.RawQuery
	return remote, nil
}

func SignRequest(config *AppConfig, r *http.Request) (*http.Request, error) {
	remote, err := PrepareRemoteUrl(config.EsEndpoint, r)
	if err != nil {
		return nil, err
	}

	cfg, err := GetAWSConfig()
	if err != nil {
		return nil, err
	}

	creds, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	rr, err := http.NewRequest(r.Method, remote.String(), bytes.NewBuffer(body))
	rr.URL = remote

	if err != nil {
		return nil, err
	}

	hash, err := CalculateSHA256Hash(bytes.NewBuffer(body))

	err = v4.NewSigner().SignHTTP(rr.Context(), creds, rr, hash, "es", cfg.Region, time.Now())

	if err != nil {
		return nil, err
	}

	rr.Header.Set("kbn-version", r.Header.Get("kbn-version"))
	rr.Header.Set("User-Agent", r.Header.Get("User-Agent"))

	return rr, nil
}
