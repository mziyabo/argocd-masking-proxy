package masking

import (
	"regexp"
)

// perform masking operation
func Mask(data string) string {

	// load config
	// read manifest

	//r := regexp.MustCompile(`client(_?Secret|ID):\s?\\?\"?[a-z0-9]*(\\\"|\"|\s)?`)
	r := regexp.MustCompile(`client(_?Secret|ID)`)
	data = r.ReplaceAllString(data, "private.******")
	// data = strings.Replace(data, "AIzaSyDol1GJR0myfynTGuHLnf2HITViVHpYfqg", "private.******", -1)

	return data
}
