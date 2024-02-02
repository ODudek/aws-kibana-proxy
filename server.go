package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	config     *AppConfig
	httpClient *http.Client
}

func (s *Server) createProxy() *httputil.ReverseProxy {
	remote, err := url.Parse(s.config.EsEndpoint)
	if err != nil {
		// Crash app because there is a problem with parse endpoint url
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(remote)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (s *Server) Handle(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		rw.WriteHeader(http.StatusOK)
		return
	}
	req, err := SignRequest(s.config, r)
	if err != nil {
		log.Error("Error signing request:", err.Error())
		return
	}
	originalResponse, err := s.httpClient.Do(req)

	if err != nil {
		log.Error("Error sending request:", err.Error())
		return
	}
	copyHeader(rw.Header(), originalResponse.Header)
	rw.WriteHeader(originalResponse.StatusCode)
	io.Copy(rw, originalResponse.Body)
	originalResponse.Body.Close()
}

func (s *Server) Start() {
	log.Info("Starting server on port", s.config.Port)
	err := http.ListenAndServe(s.config.Port, http.HandlerFunc(s.Handle))
	if err != nil {
		// Crash app because there is a problem with starting server
		panic(err)
	}
}

func NewServer(config *AppConfig) *Server {
	return &Server{
		config,
		http.DefaultClient,
	}
}
