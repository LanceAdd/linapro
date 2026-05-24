//go:build !wasip1

// This file verifies the root pluginbridge facade exposes guest host-service
// clients during ordinary Go builds by forwarding to the framework-owned guest
// stubs.

package pluginbridge

import (
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/pluginbridge/guest"
)

// TestFacadeGuestHostServiceStubs verifies dynamic plugins can import the root
// facade in non-WASI tests without defining plugin-local unsupported clients.
func TestFacadeGuestHostServiceStubs(t *testing.T) {
	t.Parallel()

	if _, err := Runtime().Now(); !gerror.Is(err, guest.ErrHostCallsUnavailable) {
		t.Fatalf("expected runtime unavailable sentinel, got %v", err)
	}
	if _, _, err := Config().String("demo.greeting"); !gerror.Is(err, guest.ErrHostCallsUnavailable) {
		t.Fatalf("expected config unavailable sentinel, got %v", err)
	}
	if _, _, err := HostConfig().String("workspace.basePath"); !gerror.Is(err, guest.ErrHostCallsUnavailable) {
		t.Fatalf("expected host runtime unavailable sentinel, got %v", err)
	}
	if _, _, err := Manifest().GetText("metadata.yaml"); !gerror.Is(err, guest.ErrHostCallsUnavailable) {
		t.Fatalf("expected manifest unavailable sentinel, got %v", err)
	}
}
