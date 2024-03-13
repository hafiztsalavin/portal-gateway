package service

import (
	"net/http"
	"time"
)

type BackendService struct {
	Name            string   `json:"name" yaml:"name"`
	Scheme          string   `json:"scheme" yaml:"scheme"`
	UpstreamTargets []string `json:"upstreamTargets" yaml:"upstreamTargets"`
	Path            string   `json:"path,omitempty" yaml:"path,omitempty"`
	Domain          string   `json:"domain" yaml:"domain"`

	MaxIdleConns int           `json:"maxIdleConns,omitempty" yaml:"maxIdleConns,omitempty"`
	MaxIdleTime  time.Duration `json:"maxIdleTime" yaml:"maxIdleTime"`
	Timeout      time.Duration `json:"timeout" yaml:"timeout"`

	httpClient *http.Client
}

func (bs *BackendService) setHttpClient() {
	transport := &http.Transport{
		MaxIdleConns:        bs.MaxIdleConns,
		IdleConnTimeout:     bs.MaxIdleTime * time.Second,
		TLSHandshakeTimeout: bs.Timeout * time.Second,
	}

	bs.httpClient = &http.Client{Transport: transport}
}

func (bs *BackendService) Init() {
	bs.setHttpClient()
}
