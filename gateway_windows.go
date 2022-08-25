//go:build windows

package gateway

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var gwArgs []string = strings.Split("path Win32_NetworkAdapterConfiguration where IPEnabled=true get DefaultIPGateway,GatewayCostMetric,IPConnectionMetric,Index /format:table", " ")

func ifArgs(id int) []string {
	return strings.Split(fmt.Sprintf(`path Win32_NetworkAdapter where Index=%d get NetConnectionID,MACAddress /format:table`, id), " ")
}

func parseGwTable(gwTable string, family int) (string, int) {
	var bestGw string
	var bestMetric int
	var bestId int
	gwTable = strings.Trim(gwTable, " ")

	lines := regexp.MustCompile("\r*\n").Split(gwTable, -1)
	lines = lines[1:]

	for _, line := range lines {
		line = strings.Trim(line, " ")
		matchs := regexp.MustCompile(`({.+?}) +?({.+?}) +?([0-9]+) +?([0-9]+)`).FindAllStringSubmatch(line, -1)
		if len(matchs) == 0 {
			continue
		}
		data := matchs[0]
		if len(data) < 5 {
			continue
		}
		gwArr := data[1]
		gwCostsArr := data[2]
		id, _ := strconv.Atoi(data[3])
		ipMetric, _ := strconv.Atoi(data[4])
		gateways := regexp.MustCompile(`"(.+?)"`).FindAllString(gwArr, -1)
		for i, gateway := range gateways {
			gateways[i] = gateway[1 : len(gateway)-1]
		}
		gatewayCosts := regexp.MustCompile(`[0-9]+`).FindAllString(gwCostsArr, -1)

		for i, gateway := range gateways {
			if len(gateway) == 0 || IsIP(gateway) != family {
				continue
			}
			cost, _ := strconv.Atoi(gatewayCosts[i])
			metric := cost + ipMetric

			if len(bestGw) == 0 || metric < bestMetric {
				bestGw, bestMetric, bestId = gateway, metric, id
			}
		}
	}

	return bestGw, bestId
}

func parseIfTable(ifTable string) (string, error) {
	ifTable = strings.Trim(ifTable, " ")
	lines := regexp.MustCompile("\n").Split(ifTable, -1)
	line := strings.Trim(lines[1], " ")
	data := regexp.MustCompile(`\s+`).Split(line, -1)
	mac := strings.ToLower(data[0])
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var name string
	for _, inte := range interfaces {
		if strings.ToLower(inte.HardwareAddr.String()) == mac {
			name = inte.Name
			break
		}
	}
	return name, nil
}

func (g *Gateway) getGatewayOS(family int) error {
	out, err := execute("wmic", gwArgs...)
	if err != nil {
		return err
	}

	gateway, id := parseGwTable(out, family)
	if len(gateway) == 0 {
		return errNoGateway
	}

	var name string
	out, err = execute("wmic", ifArgs(id)...)
	if err != nil {
		return err
	}

	name, err = parseIfTable(out)
	if err != nil {
		return err
	}

	g.Gateway = gateway
	g.Interface = name
	return nil
}
