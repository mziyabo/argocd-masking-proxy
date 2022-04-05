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

	// pad := strings.Join([]string{"%-", fmt.Sprint((len(d) - len(data))), "s"}, "")
	// g := fmt.Sprintf(pad, data)
	g := []byte(data)
	g = append(g, make([]byte, (len(d)-len(data)))...)

	s := strings.ReplaceAll(string(g), "\x00", " ")

	return ([]byte(s))
}
