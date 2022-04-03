package masking

import (
	"regexp"
	"strings"
)

func init() {
	// Load config
	// Read manifest
}

// perform masking operation
func Mask(d []byte) []byte {

	data := string(d)
	data = strings.ReplaceAll(data, `\\\`, "\\")

	// GitHub OAUTH client details:
	r := regexp.MustCompile(`(client(_?Secret|ID)):((\s?\\?\"?)([a-z0-9]*)(\\\"|\"|\s)?)`)
	data = r.ReplaceAllString(data, "$1:$4******$6")

	// r := regexp.MustCompile(`192.168.49.2`)
	// data = r.ReplaceAllString(data, "******")

	return ([]byte(data))
}
