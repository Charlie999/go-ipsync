# go-ipsync
Synchronise a slice of IP addresses to an interface, using netlink.

## Usage
```go
package main

import (
	"github.com/charlie999/go-ipsync"
	"github.com/vishvananda/netlink"
	"log"
)

func main() {
	addrsList := []*netlink.Addr{}
	addr, _ := netlink.ParseAddr("1.2.3.4/32")
	addrsList = append(addrsList, addr)
	// ...

	err := ipsync.SyncAddrOnInterface(addrsList, "dummy0");
	if (err != nil) {
		log.Fatal(err);
	}
}
```
