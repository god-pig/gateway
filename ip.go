package gateway

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	IPv4 int = 4
	IPv6 int = 6
)

// IPv4 Segment
var v4Seg string = `(?:[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])`

var v4Str string = fmt.Sprintf(`(%s[.]){3}%s`, v4Seg, v4Seg)

var v4Reg *regexp.Regexp = regexp.MustCompile(fmt.Sprintf(`^%s$`, v4Str))

// // IPv6 Segment
var v6Seg string = `(?:[0-9a-fA-F]{1,4})`

var v6Reg *regexp.Regexp = regexp.MustCompile(`^(` +
	fmt.Sprintf(`(?:%s:){7}(?:%s|:)|`, v6Seg, v6Seg) +
	fmt.Sprintf(`(?:%s:){6}(?:%s|:%s|:)|`, v6Seg, v4Seg, v6Seg) +
	fmt.Sprintf(`(?:%s:){5}(?::%s|(:%s){1,2}|:)|`, v6Seg, v4Seg, v6Seg) +
	fmt.Sprintf(`(?:%s:){4}(?:(:%s){0,1}:%s|(:%s){1,3}|:)|`, v6Seg, v6Seg, v4Seg, v6Seg) +
	fmt.Sprintf(`(?:%s:){3}(?:(:%s){0,2}:%s|(:%s){1,4}|:)|`, v6Seg, v6Seg, v4Seg, v6Seg) +
	fmt.Sprintf(`(?:%s:){2}(?:(:%s){0,3}:%s|(:%s){1,5}|:)|`, v6Seg, v6Seg, v4Seg, v6Seg) +
	fmt.Sprintf(`(?:%s:){1}(?:(:%s){0,4}:%s|(:%s){1,6}|:)|`, v6Seg, v6Seg, v4Seg, v6Seg) +
	fmt.Sprintf(`(?::((?::%s){0,5}:%s|(?::%s){1,7}|:))`, v6Seg, v4Seg, v6Seg) +
	`)(%[0-9a-zA-Z-.:]{1,})?$`)

var (
	errNoGateway = errors.New("gateway no found")
)

func IsIPv4(s string) bool {
	return v4Reg.Match([]byte(s))
}

func IsIPv6(s string) bool {
	return v6Reg.Match([]byte(s))
}

func IsIP(s string) int {
	if len(s) == 0 {
		return 0
	}
	if IsIPv4(s) {
		return 4
	}
	if IsIPv6(s) {
		return 6
	}
	return 0
}
