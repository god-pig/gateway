package gateway

import (
	"bytes"
	"os/exec"
)

type Gateway struct {
	Gateway   string
	Interface string
}

func GetGateway(family int) (*Gateway, error) {
	inte := &Gateway{}
	if family != IPv6 {
		family = IPv4
	}

	err := inte.getGatewayOS(family)
	if err != nil {
		return nil, err
	}
	return inte, nil
}

func execute(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	out := stdout.String()
	if err != nil {
		return "", err
	}
	return out, nil
}
