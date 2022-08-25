# default gateway

Get the default network gateway, cross-platform.

## Usage

```go
package main

import (
	"fmt"
	"github.com/god-pig/gateway"
)

func main() {
	data := gateway.GetGateway(gateway.IPv4)
	fmt.Println(data.Gateway)
	// -> 192.168.1.1

	fmt.Println(data.Interface)
	// -> 以太网
}
```