package schema

import "net/url"

// Proxy Config
type ProxyConfig struct {
	Rules        []ProxyRule
	ApiServerUrl *url.URL
	Port         int
	Host         string
}

// Masking Rule
type ProxyRule struct {
	// Target Kubernetes api object type
	Kind string

	// Regex pattern to mask
	Pattern string

	// Fields to include otherwise Include ALL
	// Meant to stop targeting certain things like namespaces
	IncludeFields []string
}
