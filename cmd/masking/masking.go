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
func Mask(data string) string {

	data = strings.ReplaceAll(data, `\\\`, "\\")

	// GitHub OAUTH client details:
	r := regexp.MustCompile(`(client(_?Secret|ID)):((\s?\\?\"?)([a-z0-9]*)(\\\"|\"|\s)?)`)
	data = r.ReplaceAllString(data, "$1:$4******$6")

	return data
}
