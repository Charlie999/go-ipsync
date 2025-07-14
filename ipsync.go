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

func shouldKeepAddr(newAddrList []*netlink.Addr, existing *netlink.Addr) *netlink.Addr {
	for _, addr := range newAddrList {
		if (existing.Equal(*addr)) {
			return addr
		}
	}

	return nil
}

func shouldAddAddr(newaddr *netlink.Addr, intaddr4 []netlink.Addr, intaddr6 []netlink.Addr) bool {
	for _, addr4 := range intaddr4 {
		if (addr4.Equal(*newaddr)) {
			return false;
		}
	}

	for _, addr6 := range intaddr6 {
		if (addr6.Equal(*newaddr)) {
			return false;
		}
	}

	return true;
}

/*
SyncAddrOnInterface(addrlist []*netlink.Addr, ifname String)

Apply all addresses in addrlist to interface, and remove all that are not present.
Must be done as root.
*/
func SyncAddrOnInterface(addrlist []*netlink.Addr, ifname string) error {
	link, err := netlink.LinkByName(ifname)
	if (err != nil) {
		return err;
	}

	// Get addresses
	addrList4, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if (err != nil) {
		return err;
	}
	addrList6, err := netlink.AddrList(link, netlink.FAMILY_V6)
	if (err != nil) {
		return err;
	}

	actions := []AddrSyncAction{}

	// Check for removals
	for _, addr := range addrList4 {
		match := shouldKeepAddr(addrlist, &addr)
		if (match == nil) {
			actions = append(actions, AddrSyncAction{addr, ActionRemove}) 
		}
	}

	for _, addr := range addrList6 {
		match := shouldKeepAddr(addrlist, &addr)
		if (match == nil) {
			actions = append(actions, AddrSyncAction{addr, ActionRemove}) 
		}
	}

	// Check for additions
	for _, addr := range addrlist {
		if (shouldAddAddr(addr, addrList4, addrList6)) {
			actions = append(actions, AddrSyncAction{*addr, ActionAdd})
		}
	}

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