package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	ErrAppendCACert     = errors.New("failed to append custom CA cert")
	ErrCACertNotExist   = errors.New("custom CA certificate path does not exist")
	ErrInvalidTimeout   = errors.New("timeout must be greater than 0")
	ErrReadCustomCACert = errors.New("failed to read custom CA cert")
	ErrSystemCACertPool = errors.New("failed to get system CA cert pool")
)

type HTTPClient struct {
	CustomCACertPath   string        `envconfig:"CUSTOM_CA_CERT_PATH" json:"custom_ca_cert_path"`
	InsecureSkipVerify bool          `default:"false"                 envconfig:"INSECURE_SKIP_VERIFY" json:"insecure_skip_verify"` //nolint:lll
	Timeout            time.Duration `default:"30s"                   envconfig:"TIMEOUT"              json:"timeout"`
}

func (cfg *HTTPClient) FromConfig() (*http.Client, error) {
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSystemCACertPool, err)
	}

	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	if cfg.CustomCACertPath != "" {
		caCert, err := os.ReadFile(cfg.CustomCACertPath)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrReadCustomCACert, err)
		}

		if !rootCAs.AppendCertsFromPEM(caCert) {
			return nil, ErrAppendCACert
		}
	}

	return &http.Client{ //nolint:exhaustruct
		Timeout: cfg.Timeout,
		Transport: &http.Transport{ //nolint:exhaustruct
			TLSClientConfig: &tls.Config{ //nolint:exhaustruct
				InsecureSkipVerify: cfg.InsecureSkipVerify,
				RootCAs:            rootCAs,
			},
		},
	}, nil
}

func (cfg *HTTPClient) Validate() error {
	if cfg.CustomCACertPath != "" {
		if _, err := os.Stat(cfg.CustomCACertPath); os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrCACertNotExist, cfg.CustomCACertPath)
		}
	}

	if cfg.Timeout <= 0 {
		return ErrInvalidTimeout
	}

	return nil
}
