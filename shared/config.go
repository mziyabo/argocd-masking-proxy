package shared

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var Config ProxyConfig

func init() {
	readConfig()
	tokenDir := "/var/run/secrets/kubernetes.io/serviceaccount/token"

	// Initialize config object
	Config.ApiURL, _ = url.Parse(viper.GetString("target"))
	Config.Port = viper.GetInt("serve.port")
	Config.Host = viper.GetString("serve.host")
	tokenBytes, err := os.ReadFile(tokenDir)
	if err != nil {
		log.Printf("%s\n", err)
	} else {
		Config.Token = string(tokenBytes)
	}

	Config.TLSConfig = TLSClientConfig{
		Enabled: viper.GetBool("serve.tls.enabled"),
		Cert:    viper.GetString("serve.tls.cert"),
		Key:     viper.GetString("serve.tls.certKey"),
	}

	if Config.TLSConfig.Enabled {
		Config.ProxyURL, _ = url.Parse(strings.Join([]string{"https://", Config.Host, ":", fmt.Sprint(Config.Port)}, ""))
	} else {
		Config.ProxyURL, _ = url.Parse(strings.Join([]string{"http://", Config.Host, ":", fmt.Sprint(Config.Port)}, ""))
	}
}

// Proxy Config
type ProxyConfig struct {
	Rules     []ProxyRule
	ApiURL    *url.URL // Kubernetes API server url
	ProxyURL  *url.URL
	Port      int
	Host      string
	TLSConfig TLSClientConfig
	Token     string // ServiceAccount token
}

type TLSClientConfig struct {
	Enabled bool
	Cert    string
	Key     string
}

// Masking Rule
type ProxyRule struct {
	// Target Kubernetes api object type
	// Kind string

	// Regex pattern to mask
	RegexPattern string

	// Regex replacement string
	Replacement string

	// Fields to include otherwise Include ALL
	// Meant to stop targeting certain things like namespaces
	// IncludeFields []string
}

// Read proxy.conf.json file
func readConfig() {

	// Name of config file (without extension)
	viper.SetConfigName("proxy.conf")
	// REQUIRED if the config file does not have the extension in the name
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/kubemask")
	viper.AddConfigPath("$HOME/.kubemask")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
