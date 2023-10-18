package controller

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"
	httptransport "github.com/go-git/go-git/v5/plumbing/transport/http"
	ctrl "sigs.k8s.io/controller-runtime"
)

// HttpTransportWithCustomCerts returns an HTTP transport with custom certificates.
// If proxyStr is provided, it will be used as the proxy URL.
// If not, it tries to fetch the proxy from an environment variable.
func HttpTransportwithCustomCerts(tlsConfig *tls.Config, proxyStr *transport.ProxyOptions, ctx context.Context) (transport.Transport, error) {
	log := ctrl.LoggerFrom(ctx)
	var message string

	var (
		proxyUrl *url.URL
		err      error
	)
	if proxyStr != nil {
		proxyUrl, err = url.Parse(proxyStr.URL)
		if err != nil {
			message = fmt.Sprintf("failed to parse proxy url: %s", proxyStr.URL)
			log.Info(message)

			proxyUrl = &url.URL{}
		}
	} else {
		proxyUrl, err = GetHTTPSProxy()
		if err != nil {
			message = fmt.Sprintf("https_proxy environment variable is not set or invalid: %v", err)
			log.Info(message)

			proxyUrl = &url.URL{}
		}
	}
	return httptransport.NewClient(&http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			TLSClientConfig: tlsConfig,
		},
	}), nil

}

// GetHTTPSProxy returns the value of the https_proxy environment variable.
func GetHTTPSProxy() (*url.URL, error) {
	proxy := os.Getenv("https_proxy")
	if proxy == "" {
		proxy = os.Getenv("HTTPS_PROXY")
	}
	if proxy == "" {
		return nil, fmt.Errorf("no https_proxy environment variable set")
	}
	parsedURL, err := url.Parse(proxy)
	if err != nil {
		return nil, fmt.Errorf("invalid https_proxy URL: %w", err)
	}
	return parsedURL, nil
}
