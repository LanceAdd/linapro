// This file verifies the host-governed plugin consumer surface projection
// without requiring plugin lifecycle database fixtures or concrete plugins.

package plugin

import (
	"testing"

	"lina-core/internal/model/entity"
	"lina-core/internal/service/plugin/internal/catalog"
	"lina-core/pkg/pluginhost"
)

// TestBuildConsumerSurfaceSnapshotAggregatesHostInputs verifies C-side API,
// frontend, enablement, version, and tenant governance metadata are grouped by plugin.
func TestBuildConsumerSurfaceSnapshotAggregatesHostInputs(t *testing.T) {
	supportsMultiTenant := true
	manifests := []*catalog.Manifest{
		{
			ID:                  "mall",
			Version:             "0.1.0",
			Type:                catalog.TypeSource.String(),
			ScopeNature:         catalog.ScopeNatureTenantAware.String(),
			SupportsMultiTenant: &supportsMultiTenant,
			DefaultInstallMode:  catalog.InstallModeTenantScoped.String(),
		},
		{
			ID:                 "admin-only",
			Version:            "0.1.0",
			Type:               catalog.TypeSource.String(),
			ScopeNature:        catalog.ScopeNaturePlatformOnly.String(),
			DefaultInstallMode: catalog.InstallModeGlobal.String(),
		},
	}
	registries := []*entity.SysPlugin{
		{
			PluginId:    "mall",
			Version:     "0.2.0",
			ScopeNature: catalog.ScopeNatureTenantAware.String(),
			InstallMode: catalog.InstallModeTenantScoped.String(),
		},
	}
	bindings := []pluginhost.SourceRouteBinding{
		{
			PluginID:     "mall",
			Method:       "GET",
			Path:         "/api/c/v1/mall/products",
			Surface:      pluginhost.SurfaceConsumer,
			Documentable: true,
		},
		{
			PluginID: "mall",
			Method:   "POST",
			Path:     "/api/v1/mall/internal",
			Surface:  pluginhost.SurfaceAdmin,
		},
	}
	frontendIndex := &sourceConsumerFrontendResourceIndex{
		mounts: []*sourceConsumerFrontendMountEntry{
			{
				pluginID:    "mall",
				version:     "0.2.0",
				mountPath:   "/mall",
				index:       "index.html",
				spaFallback: true,
				assets: map[string]struct{}{
					"frontend/consumer/index.html":     {},
					"frontend/consumer/assets/app.js":  {},
					"frontend/consumer/assets/app.css": {},
				},
			},
		},
	}

	snapshot := buildConsumerSurfaceSnapshot(
		manifests,
		registries,
		bindings,
		frontendIndex,
		map[string]bool{"mall": true, "admin-only": true},
	)

	if snapshot == nil {
		t.Fatalf("expected consumer surface snapshot")
	}
	if len(snapshot.Plugins) != 1 {
		t.Fatalf("expected only C-side capable plugin, got %#v", snapshot.Plugins)
	}
	got := snapshot.Plugins[0]
	if got.PluginID != "mall" || got.Version != "0.2.0" || !got.Enabled {
		t.Fatalf("unexpected plugin identity or enablement: %#v", got)
	}
	if !got.TenantAware ||
		got.ScopeNature != catalog.ScopeNatureTenantAware.String() ||
		got.DefaultInstallMode != catalog.InstallModeTenantScoped.String() {
		t.Fatalf("unexpected tenant governance projection: %#v", got)
	}
	if got.ConsumerAPIRouteCount != 1 || len(got.ConsumerAPIRoutes) != 1 {
		t.Fatalf("expected one C-side route, got %#v", got.ConsumerAPIRoutes)
	}
	if got.ConsumerAPIRoutes[0].Path != "/api/c/v1/mall/products" ||
		got.ConsumerAPIRoutes[0].Method != "GET" ||
		!got.ConsumerAPIRoutes[0].Documentable {
		t.Fatalf("unexpected route snapshot: %#v", got.ConsumerAPIRoutes[0])
	}
	if got.ConsumerFrontend == nil ||
		got.ConsumerFrontend.MountPath != "/mall" ||
		got.ConsumerFrontend.AssetCount != 3 ||
		!got.ConsumerFrontend.SPAFallback {
		t.Fatalf("unexpected frontend snapshot: %#v", got.ConsumerFrontend)
	}
}

// TestBuildConsumerSurfaceSnapshotDerivesPluginIDFromConsumerPath verifies
// route bindings remain governable when the binding lacks an explicit plugin ID.
func TestBuildConsumerSurfaceSnapshotDerivesPluginIDFromConsumerPath(t *testing.T) {
	snapshot := buildConsumerSurfaceSnapshot(
		nil,
		nil,
		[]pluginhost.SourceRouteBinding{
			{
				Method:       "GET",
				Path:         "/api/c/v1/portal/home",
				Surface:      pluginhost.SurfaceConsumer,
				Documentable: true,
			},
		},
		nil,
		nil,
	)

	if snapshot == nil || len(snapshot.Plugins) != 1 {
		t.Fatalf("expected one derived plugin snapshot, got %#v", snapshot)
	}
	if snapshot.Plugins[0].PluginID != "portal" {
		t.Fatalf("expected plugin id from route path, got %#v", snapshot.Plugins[0])
	}
}

// TestBuildConsumerSurfaceSnapshotSortsPluginsAndRoutes verifies governance
// output remains deterministic for review and future API exposure.
func TestBuildConsumerSurfaceSnapshotSortsPluginsAndRoutes(t *testing.T) {
	snapshot := buildConsumerSurfaceSnapshot(
		nil,
		nil,
		[]pluginhost.SourceRouteBinding{
			{PluginID: "portal", Method: "POST", Path: "/api/c/v1/portal/login", Surface: pluginhost.SurfaceConsumer},
			{PluginID: "mall", Method: "GET", Path: "/api/c/v1/mall/products", Surface: pluginhost.SurfaceConsumer},
			{PluginID: "portal", Method: "GET", Path: "/api/c/v1/portal/home", Surface: pluginhost.SurfaceConsumer},
		},
		nil,
		nil,
	)

	if len(snapshot.Plugins) != 2 {
		t.Fatalf("expected two plugin snapshots, got %#v", snapshot.Plugins)
	}
	if snapshot.Plugins[0].PluginID != "mall" || snapshot.Plugins[1].PluginID != "portal" {
		t.Fatalf("expected sorted plugins, got %#v", snapshot.Plugins)
	}
	portal := snapshot.Plugins[1]
	if portal.ConsumerAPIRoutes[0].Path != "/api/c/v1/portal/home" ||
		portal.ConsumerAPIRoutes[1].Path != "/api/c/v1/portal/login" {
		t.Fatalf("expected sorted portal routes, got %#v", portal.ConsumerAPIRoutes)
	}
}
