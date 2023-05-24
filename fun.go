package winterMVC

import "regexp"

const patternProjectUrl = `^v\d+$`

// MustCompileVersion path 为 version 则返回 true
func MustCompileVersion(path string) bool {
	reg := regexp.MustCompile(patternProjectUrl)
	return reg.MatchString(path)
}
