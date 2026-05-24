// This file verifies the guest-side pluginbridge helper contracts exposed to
// dynamic plugins stay interface-based rather than leaking concrete structs.

package guest

import (
	"reflect"
	"testing"
)

// TestGuestContractsUseInterfaces verifies guest-facing behavioral helpers are
// published as interfaces.
func TestGuestContractsUseInterfaces(t *testing.T) {
	assertGuestInterfaceType(t, (*GuestRuntime)(nil), "GuestRuntime")
	assertGuestInterfaceType(t, (*GuestControllerRouteDispatcher)(nil), "GuestControllerRouteDispatcher")
	assertGuestInterfaceType(t, (*DynamicRouteRegistrar)(nil), "DynamicRouteRegistrar")
	assertGuestInterfaceType(t, (*RuntimeHostService)(nil), "RuntimeHostService")
	assertGuestInterfaceType(t, (*StorageHostService)(nil), "StorageHostService")
	assertGuestInterfaceType(t, (*HTTPHostService)(nil), "HTTPHostService")
	assertGuestInterfaceType(t, (*DataHostService)(nil), "DataHostService")
	assertGuestInterfaceType(t, (*CacheHostService)(nil), "CacheHostService")
	assertGuestInterfaceType(t, (*LockHostService)(nil), "LockHostService")
	assertGuestInterfaceType(t, (*ConfigHostService)(nil), "ConfigHostService")
	assertGuestInterfaceType(t, (*NotifyHostService)(nil), "NotifyHostService")
	assertGuestInterfaceType(t, (*CronHostService)(nil), "CronHostService")
	assertGuestInterfaceType(t, (*HostConfigHostService)(nil), "HostConfigHostService")
	assertGuestInterfaceType(t, (*ManifestHostService)(nil), "ManifestHostService")
}

// assertGuestInterfaceType verifies the reflected type under test is an
// interface.
func assertGuestInterfaceType(t *testing.T, value interface{}, name string) {
	t.Helper()

	if reflect.TypeOf(value).Elem().Kind() != reflect.Interface {
		t.Fatalf("expected %s to be declared as interface", name)
	}
}
