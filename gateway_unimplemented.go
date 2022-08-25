//go:build !darwin && !linux && !windows

package gateway

var (
	errNotImplemented = errors.New("not implemented for OS: " + runtime.GOOS)
)

func (inte *Interface) getGatewayOS() error {
	return errNotImplemented
}
