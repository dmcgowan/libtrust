package client

import (
	"crypto/x509"
	"errors"
	"net/http"

	"github.com/docker/libtrust/trustapi"
	"github.com/docker/libtrust/trustgraph"
	"github.com/gorilla/mux"
)

var (
	ErrBaseGraphDoesNotExist = errors.New("base graph does not exist")
)

type TrustClient struct {
	client    *http.Client
	router    *mux.Router
	endpoint  string
	authority *x509.CertPool
}

func NewTrustClient(server string, authority *x509.CertPool) *TrustClient {
	return &TrustClient{
		client:    &http.Client{},
		router:    trustapi.NewRouter(server),
		endpoint:  server,
		authority: authority,
	}
}

func (c *TrustClient) GetBaseGraph(name string) (*trustgraph.Statement, error) {
	u, err := c.router.Get("graphbase").URL("graphname", name)
	if err != nil {
		return nil, err
	}

	u.Host = c.endpoint
	req := &http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Body:       nil,
		Host:       c.endpoint,
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, ErrBaseGraphDoesNotExist
	}

	defer resp.Body.Close()

	return trustgraph.LoadStatement(resp.Body, c.authority)
}
