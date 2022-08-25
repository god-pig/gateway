//go:build darwin

package gateway

import (
	"regexp"
	"strings"
)

var dests []string = []string{"default", "0.0.0.0", "0.0.0.0/0", "::", "::/0"}

func has(arr []string, item string) bool {
	for _, val := range arr {
		if val == item {
			return true
		}
	}
	return false
}

func (g *Gateway) parse(stdout string, family int) error {
	stdout = strings.Trim(stdout, " ")

	lines := regexp.MustCompile("\n").Split(stdout, -1)

	for _, line := range lines {
		matchs := regexp.MustCompile(` +`).Split(line, -1)
		if len(matchs) == 0 {
			continue
		}
		data := matchs[0]
		target := data[0]
		gateway := data[1]
		iface := data[3]
		if has(dests, target) && IsIP(gateway) > 0 {
			g.Gateway = gateway
			g.Interface = iface
		}
	}
	if len(g.Gateway) == 0 && len(g.Interface) == 0 {
		return errNoGateway
	}
	return nil
}

func (g *Gateway) getGatewayOS(family int) error {
	inet := "inet"
	if family == IPv6 {
		inet = "inet6"
	}
	stdout, err := execute("netstat", "-rn", "-f", inet)
	if err != nil {
		return err
	}
	return g.parse(stdout, family)
}
