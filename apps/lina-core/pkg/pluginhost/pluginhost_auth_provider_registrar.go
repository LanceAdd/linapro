// This file implements the host-owned registrar used by source plugins to
// publish external authentication providers.

package pluginhost

import "strings"

// authProviderRegistrar is the host-owned AuthProviderRegistrar implementation.
type authProviderRegistrar struct {
	pluginID     string
	hostServices HostServices
	providers    []AuthProvider
}

// NewAuthProviderRegistrar creates one auth-provider registrar for a source plugin.
func NewAuthProviderRegistrar(pluginID string, hostServices HostServices) AuthProviderRegistrar {
	return &authProviderRegistrar{
		pluginID:     strings.TrimSpace(pluginID),
		hostServices: hostServices,
		providers:    make([]AuthProvider, 0),
	}
}

// Add registers one plugin-owned authentication provider.
func (r *authProviderRegistrar) Add(provider AuthProvider) {
	if r == nil || provider == nil {
		return
	}
	r.providers = append(r.providers, provider)
}

// HostServices returns host-published services available to the plugin.
func (r *authProviderRegistrar) HostServices() HostServices {
	if r == nil {
		return nil
	}
	return r.hostServices
}

// Providers returns a snapshot of providers registered in this session.
func (r *authProviderRegistrar) Providers() []AuthProvider {
	if r == nil || len(r.providers) == 0 {
		return nil
	}
	return append([]AuthProvider(nil), r.providers...)
}

// PluginID returns the plugin owning this registration session.
func (r *authProviderRegistrar) PluginID() string {
	if r == nil {
		return ""
	}
	return r.pluginID
}
