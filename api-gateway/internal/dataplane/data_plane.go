package dataplane

// setup reverse proxy for data plane
import (
	"net/http/httputil"
	"net/url"
)

func GetGatewayRouter(targetURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy, nil
}