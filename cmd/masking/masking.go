package masking

import (
	"regexp"
	"strings"

	"github.com/mziyabo/masking-proxy/cmd/shared"
)

var config shared.ProxyConfig

func init() {
	config = shared.Config
}

// Mask performs data-masking operation based off proxy config rules
func Mask(d []byte) []byte {

	data := string(d)

	for _, rule := range config.Rules {
		r := regexp.MustCompile(rule.Pattern)

		if r.Match([]byte(data)) {

			data = r.ReplaceAllString(data, rule.Replacement)

			// TODO: shorten this up
			// pad byte-array to maintain content-length
			g := []byte(data)
			g = append(g, make([]byte, (len(d)-len(data)))...)
			s := strings.ReplaceAll(string(g), "\x00", " ")

			d = []byte(s)
		}
	}

	return d
}
