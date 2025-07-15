package ipsync

import (
	"github.com/vishvananda/netlink"
)

type AddrSyncActionType int
const (
	ActionAdd AddrSyncActionType = iota
	ActionRemove
)


type AddrSyncAction struct {
	addr netlink.Addr
	action AddrSyncActionType
}

func sliceContainsAddr(needle netlink.Addr, haystack []netlink.Addr) bool {
	for _, addr := range haystack {
		if (addr.Equal(needle)) {
			return true
		}
	}

	return false
}

/*
SyncAddrOnInterface(addrlist []*netlink.Addr, ifname String)

Apply all addresses in addrlist to interface, and remove all that are not present.
Must be done as root.
*/
func SyncAddrOnInterface(addrlist []netlink.Addr, ifname string) error {
	link, err := netlink.LinkByName(ifname)
	if (err != nil) {
		return err;
	}

	// Get addresses
	addrList, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if (err != nil) {
		return err;
	}

	actions := []AddrSyncAction{}

	// Check for removals
	for _, addr := range addrList {
		if (!sliceContainsAddr(addr, addrlist)) {
			actions = append(actions, AddrSyncAction{addr, ActionRemove}) 
		}
	}

	// Check for additions
	for _, addr := range addrlist {
		if (!sliceContainsAddr(addr, addrList)) {
			actions = append(actions, AddrSyncAction{addr, ActionAdd})
		}
	}

	// Opportunity for dry-run here?

	for _, action := range actions {
		if (action.action == ActionAdd) { // Add address
			err = netlink.AddrAdd(link, &action.addr)
			if (err != nil) {
				return err;
			}
		} else if (action.action == ActionRemove) {
			err = netlink.AddrDel(link, &action.addr)
			if (err != nil) {
				return err;
			}
		}
	}

	return nil;
}