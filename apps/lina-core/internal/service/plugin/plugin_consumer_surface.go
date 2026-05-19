// This file builds host-governed source-plugin consumer-facing projections by
// combining manifest declarations, frontend resource indexes, and current
// enablement state.

package plugin

import (
	"context"
	"sort"
	"strings"

	"lina-core/internal/model/entity"
	"lina-core/internal/service/plugin/internal/catalog"
)

// BuildConsumerSurfaceSnapshot builds an on-demand host governance snapshot
// for source-plugin frontend mounts, enablement, version, and tenant
// governance declarations.
func (s *serviceImpl) BuildConsumerSurfaceSnapshot(ctx context.Context) (*ConsumerSurfaceSnapshot, error) {
	if err := s.ensureRuntimeCacheFresh(ctx); err != nil {
		return nil, err
	}
	manifests, err := s.catalogSvc.ScanEmbeddedSourceManifests()
	if err != nil {
		return nil, err
	}
	registries, err := s.catalogSvc.ListAllRegistries(ctx)
	if err != nil {
		return nil, err
	}
	frontendIndex, err := s.loadSourceConsumerFrontendResourceIndex(ctx)
	if err != nil {
		return nil, err
	}
	enabledByID := s.buildConsumerSurfaceEnabledMap(ctx, manifests)
	snapshot := buildConsumerSurfaceSnapshot(
		manifests,
		registries,
		frontendIndex,
		enabledByID,
	)
	return snapshot, nil
}

// buildConsumerSurfaceEnabledMap resolves business-entry enablement for every
// discovered source plugin while reusing registry rows already loaded for the
// governance snapshot.
func (s *serviceImpl) buildConsumerSurfaceEnabledMap(
	ctx context.Context,
	manifests []*catalog.Manifest,
) map[string]bool {
	enabledByID := make(map[string]bool, len(manifests))
	for _, manifest := range manifests {
		pluginID := strings.TrimSpace(manifestID(manifest))
		if pluginID == "" {
			continue
		}
		enabledByID[pluginID] = s.IsEnabled(ctx, pluginID)
	}
	return enabledByID
}

// buildConsumerSurfaceSnapshot assembles the read-only plugin consumer surface
// projection from already-loaded host governance inputs.
func buildConsumerSurfaceSnapshot(
	manifests []*catalog.Manifest,
	registries []*entity.SysPlugin,
	frontendIndex *sourceConsumerFrontendResourceIndex,
	enabledByID map[string]bool,
) *ConsumerSurfaceSnapshot {
	pluginSnapshots := make(map[string]*ConsumerSurfacePluginSnapshot)
	registryByID := buildRegistryByPluginID(registries)
	for _, manifest := range manifests {
		item := buildConsumerSurfacePluginSnapshot(manifest, registryByID[strings.TrimSpace(manifestID(manifest))], enabledByID)
		if item != nil {
			pluginSnapshots[item.PluginID] = item
		}
	}

	if frontendIndex != nil {
		for _, mount := range frontendIndex.mounts {
			if mount == nil || strings.TrimSpace(mount.pluginID) == "" {
				continue
			}
			item := pluginSnapshots[mount.pluginID]
			if item == nil {
				item = &ConsumerSurfacePluginSnapshot{PluginID: mount.pluginID, Version: mount.version}
				pluginSnapshots[mount.pluginID] = item
			}
			item.ConsumerFrontend = &ConsumerSurfaceFrontendSnapshot{
				MountPath:   mount.mountPath,
				Index:       mount.index,
				SPAFallback: mount.spaFallback,
				AssetCount:  len(mount.assets),
			}
		}
	}

	items := make([]*ConsumerSurfacePluginSnapshot, 0, len(pluginSnapshots))
	for _, item := range pluginSnapshots {
		if item == nil || !consumerSurfacePluginHasCapability(item) {
			continue
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i int, j int) bool {
		if items[i] == nil {
			return false
		}
		if items[j] == nil {
			return true
		}
		return items[i].PluginID < items[j].PluginID
	})
	return &ConsumerSurfaceSnapshot{
		Plugins: items,
	}
}

// buildConsumerSurfacePluginSnapshot projects stable plugin governance metadata
// from manifest and registry state without requiring any concrete business plugin.
func buildConsumerSurfacePluginSnapshot(
	manifest *catalog.Manifest,
	registry *entity.SysPlugin,
	enabledByID map[string]bool,
) *ConsumerSurfacePluginSnapshot {
	if manifest == nil || catalog.NormalizeType(manifest.Type) != catalog.TypeSource {
		return nil
	}
	pluginID := strings.TrimSpace(manifest.ID)
	if pluginID == "" {
		return nil
	}
	scopeNature := catalog.NormalizeScopeNature(manifest.ScopeNature).String()
	installMode := catalog.NormalizeInstallMode(manifest.DefaultInstallMode).String()
	version := strings.TrimSpace(manifest.Version)
	if registry != nil {
		if strings.TrimSpace(registry.Version) != "" {
			version = strings.TrimSpace(registry.Version)
		}
		if strings.TrimSpace(registry.ScopeNature) != "" {
			scopeNature = catalog.NormalizeScopeNature(registry.ScopeNature).String()
		}
		if strings.TrimSpace(registry.InstallMode) != "" {
			installMode = catalog.NormalizeInstallMode(registry.InstallMode).String()
		}
	}
	return &ConsumerSurfacePluginSnapshot{
		PluginID:           pluginID,
		Version:            version,
		Enabled:            enabledByID[pluginID],
		TenantAware:        manifest.SupportsTenantGovernance(),
		ScopeNature:        scopeNature,
		DefaultInstallMode: installMode,
	}
}

// consumerSurfacePluginHasCapability reports whether a plugin contributes any
// consumer-facing host surface that should appear in the governance snapshot.
func consumerSurfacePluginHasCapability(item *ConsumerSurfacePluginSnapshot) bool {
	return item != nil && item.ConsumerFrontend != nil
}

// manifestID safely returns one manifest identifier for map lookups.
func manifestID(manifest *catalog.Manifest) string {
	if manifest == nil {
		return ""
	}
	return manifest.ID
}
