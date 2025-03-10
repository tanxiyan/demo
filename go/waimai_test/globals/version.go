package globals

import (
	"fmt"
	"strconv"
	"strings"
)

const VERSION = "2.1.13"

func Version() string {
	return VERSION
}

// `
func VersionToNumber(version string) int64 {
	arr := strings.Split(version, ".")
	var n = ""
	for _, d := range arr {
		n += fmt.Sprintf("%04s", d)
	}
	d, _ := strconv.ParseInt(n, 10, 64)
	return d
}
