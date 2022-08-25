//go:build linux

package gateway

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func (g *Gateway) parse(stdout string, family int) error {
	stdout = strings.Trim(stdout, " ")

	lines := regexp.MustCompile("\n").Split(stdout, -1)

	for _, line := range lines {
		matchs := regexp.MustCompile(`default( via .+?)?( dev .+?)( |$)`).FindAllStringSubmatch(line, -1)
		if len(matchs) == 0 {
			continue
		}
		data := matchs[0]
		gateway := data[1][5:]
		iface := data[2][5:]
		if len(gateway) > 0 && IsIP(gateway) > 0 {
			g.Gateway = gateway
			g.Interface = iface
			return nil
		} else if len(iface) > 0 && len(gateway) == 0 {
			inte, err := net.InterfaceByName(iface)
			if err != nil {
				return err
			}
			addresses, err := inte.Addrs()
			if err != nil {
				return nil
			}
			for _, address := range addresses {
				ip := address.String()
				if IsIP(ip) == family {
					g.Gateway = ip
					g.Interface = inte.Name
				}
			}

		}
	}
	if len(g.Gateway) == 0 && len(g.Interface) == 0 {
		return errNoGateway
	}
	return nil
}

func (g *Gateway) getGatewayOS(family int) error {
	stdout, err := execute("ip", fmt.Sprintf("-%d", family), "r")
	if err != nil {
		return err
	}
	return g.parse(stdout, family)
}
