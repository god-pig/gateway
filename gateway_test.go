package gateway

import (
	"testing"
)

func TestGetGateway(t *testing.T) {
	g, err := GetGateway(IPv4)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(g.Interface)
	t.Log(g.Gateway)

}
