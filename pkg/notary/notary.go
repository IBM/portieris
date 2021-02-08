// Copyright 2018,2021 Portieris Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package notary

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/IBM/portieris/internal/info"
	"github.com/docker/distribution/registry/client/transport"
	notaryclient "github.com/theupdateframework/notary/client"
	"github.com/theupdateframework/notary/trustpinning"
	"github.com/theupdateframework/notary/tuf/data"
)

// Client .
type Client struct {
	trustDir string
	rootCAs  *x509.CertPool
}

// Interface .
type Interface interface {
	GetNotaryRepo(server, image, notaryToken string) (notaryclient.Repository, error)
}

// NewClient creates and initializes the client
func NewClient(trustDir string, customCA []byte) (Interface, error) {
	// Create a trust directory
	err := createTrustDir(trustDir)
	if err != nil {
		return nil, err
	}
	rootCA, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	if customCA != nil {
		rootCA.AppendCertsFromPEM(customCA)
	}
	return &Client{trustDir: trustDir, rootCAs: rootCA}, nil
}

// GetNotaryRepo .
func (c Client) GetNotaryRepo(server, image, notaryToken string) (notaryclient.Repository, error) {
	return notaryclient.NewFileCachedRepository(
		c.trustDir,
		data.GUN(image),
		server,
		c.makeHubTransport(notaryToken),
		nil,
		trustpinning.TrustPinConfig{},
	)
}

func (c Client) makeHubTransport(notaryToken string) http.RoundTripper {
	base := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			// Avoid fallback by default to SSL protocols < TLS1.2
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			RootCAs:                  c.rootCAs,
		},
		DisableKeepAlives: true,
	}

	modifiers := []transport.RequestModifier{
		transport.NewHeaderRequestModifier(http.Header{
			"User-Agent": []string{"portieris/" + info.Version},
		}),
	}

	if notaryToken != "" {
		modifiers = []transport.RequestModifier{
			transport.NewHeaderRequestModifier(http.Header{
				"User-Agent":    []string{"portieris/" + info.Version},
				"Authorization": []string{fmt.Sprintf("Bearer %s", notaryToken)},
			}),
		}
	}

	return transport.NewTransport(base, modifiers...)
}

func createTrustDir(trustDir string) error {
	// Create a new directory only if it doesn't exist
	if !fileExists(trustDir) {
		if err := os.MkdirAll(trustDir, 0700); err != nil {
			return err
		}
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
