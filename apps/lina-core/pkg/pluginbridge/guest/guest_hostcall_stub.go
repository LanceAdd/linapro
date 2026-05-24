//go:build !wasip1

// This file provides host-build stubs for guest-side host-service clients.
// The stubs keep ordinary Go builds and unit tests compilable while making it
// explicit that real host calls are only available in wasip1 guest builds.

package guest

import (
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
)

// ErrHostCallsUnavailable reports that guest host calls are unavailable in
// non-WASI builds.
var ErrHostCallsUnavailable = gerror.New(
	"pluginbridge guest host services are only available for wasip1 builds",
)

// Host-build unsupported client implementations used by package-level helpers.
type (
	unsupportedRuntimeHostService    struct{}
	unsupportedStorageHostService    struct{}
	unsupportedHTTPHostService       struct{}
	unsupportedDataHostService       struct{}
	unsupportedCacheHostService      struct{}
	unsupportedLockHostService       struct{}
	unsupportedConfigHostService     struct{}
	unsupportedNotifyHostService     struct{}
	unsupportedCronHostService       struct{}
	unsupportedHostConfigHostService struct{}
	unsupportedManifestHostService   struct{}
)

var (
	defaultRuntimeHostService    RuntimeHostService    = unsupportedRuntimeHostService{}
	defaultStorageHostService    StorageHostService    = unsupportedStorageHostService{}
	defaultHTTPHostService       HTTPHostService       = unsupportedHTTPHostService{}
	defaultDataHostService       DataHostService       = unsupportedDataHostService{}
	defaultCacheHostService      CacheHostService      = unsupportedCacheHostService{}
	defaultLockHostService       LockHostService       = unsupportedLockHostService{}
	defaultConfigHostService     ConfigHostService     = unsupportedConfigHostService{}
	defaultNotifyHostService     NotifyHostService     = unsupportedNotifyHostService{}
	defaultCronHostService       CronHostService       = unsupportedCronHostService{}
	defaultHostConfigHostService HostConfigHostService = unsupportedHostConfigHostService{}
	defaultManifestHostService   ManifestHostService   = unsupportedManifestHostService{}
)

// Runtime returns the runtime host service guest client.
func Runtime() RuntimeHostService {
	return defaultRuntimeHostService
}

// Storage returns the storage host service guest client.
func Storage() StorageHostService {
	return defaultStorageHostService
}

// HTTP returns the outbound HTTP host service guest client.
func HTTP() HTTPHostService {
	return defaultHTTPHostService
}

// Data returns the compatibility structured data host service guest client.
// New guest code should prefer plugindb.Open().
func Data() DataHostService {
	return defaultDataHostService
}

// Cache returns the distributed cache host service guest client.
func Cache() CacheHostService {
	return defaultCacheHostService
}

// Lock returns the distributed lock host service guest client.
func Lock() LockHostService {
	return defaultLockHostService
}

// Config returns the read-only config host service guest client.
func Config() ConfigHostService {
	return defaultConfigHostService
}

// Notify returns the unified notify host service guest client.
func Notify() NotifyHostService {
	return defaultNotifyHostService
}

// Cron returns the cron host service guest client.
func Cron() CronHostService {
	return defaultCronHostService
}

// HostConfig returns the public host config guest client.
func HostConfig() HostConfigHostService {
	return defaultHostConfigHostService
}

// Manifest returns the plugin manifest-resource guest client.
func Manifest() ManifestHostService {
	return defaultManifestHostService
}

// Log reports that guest runtime host calls are unavailable.
func (unsupportedRuntimeHostService) Log(_ int, _ string, _ map[string]string) error {
	return ErrHostCallsUnavailable
}

// StateGet reports that guest runtime host calls are unavailable.
func (unsupportedRuntimeHostService) StateGet(_ string) (string, bool, error) {
	return "", false, ErrHostCallsUnavailable
}

// StateSet reports that guest runtime host calls are unavailable.
func (unsupportedRuntimeHostService) StateSet(_ string, _ string) error {
	return ErrHostCallsUnavailable
}

// StateDelete reports that guest runtime host calls are unavailable.
func (unsupportedRuntimeHostService) StateDelete(_ string) error {
	return ErrHostCallsUnavailable
}

// StateGetInt reports that guest runtime host calls are unavailable.
func (unsupportedRuntimeHostService) StateGetInt(_ string) (int, bool, error) {
	return 0, false, ErrHostCallsUnavailable
}

// StateSetInt reports that guest runtime host calls are unavailable.
func (unsupportedRuntimeHostService) StateSetInt(_ string, _ int) error {
	return ErrHostCallsUnavailable
}

// Now reports that guest runtime host calls are unavailable.
func (unsupportedRuntimeHostService) Now() (string, error) {
	return "", ErrHostCallsUnavailable
}

// UUID reports that guest runtime host calls are unavailable.
func (unsupportedRuntimeHostService) UUID() (string, error) {
	return "", ErrHostCallsUnavailable
}

// Node reports that guest runtime host calls are unavailable.
func (unsupportedRuntimeHostService) Node() (string, error) {
	return "", ErrHostCallsUnavailable
}

// HostLog writes one runtime log entry through the host.
func HostLog(level int, message string, fields map[string]string) error {
	return Runtime().Log(level, message, fields)
}

// HostStateGet reads one plugin-scoped runtime state value.
func HostStateGet(key string) (string, bool, error) {
	return Runtime().StateGet(key)
}

// HostStateSet writes one plugin-scoped runtime state value.
func HostStateSet(key string, value string) error {
	return Runtime().StateSet(key, value)
}

// HostStateDelete removes one plugin-scoped runtime state value.
func HostStateDelete(key string) error {
	return Runtime().StateDelete(key)
}

// HostStateGetInt reads one integer plugin-scoped runtime state value.
func HostStateGetInt(key string) (int, bool, error) {
	return Runtime().StateGetInt(key)
}

// HostStateSetInt writes one integer plugin-scoped runtime state value.
func HostStateSetInt(key string, value int) error {
	return Runtime().StateSetInt(key, value)
}

// Put reports that guest storage host calls are unavailable.
func (unsupportedStorageHostService) Put(
	_ string,
	_ []byte,
	_ string,
	_ bool,
) (*HostServiceStorageObject, error) {
	return nil, ErrHostCallsUnavailable
}

// PutText reports that guest storage host calls are unavailable.
func (unsupportedStorageHostService) PutText(
	_ string,
	_ string,
	_ string,
	_ bool,
) (*HostServiceStorageObject, error) {
	return nil, ErrHostCallsUnavailable
}

// Get reports that guest storage host calls are unavailable.
func (unsupportedStorageHostService) Get(
	_ string,
) ([]byte, *HostServiceStorageObject, bool, error) {
	return nil, nil, false, ErrHostCallsUnavailable
}

// GetText reports that guest storage host calls are unavailable.
func (unsupportedStorageHostService) GetText(
	_ string,
) (string, *HostServiceStorageObject, bool, error) {
	return "", nil, false, ErrHostCallsUnavailable
}

// Delete reports that guest storage host calls are unavailable.
func (unsupportedStorageHostService) Delete(_ string) error {
	return ErrHostCallsUnavailable
}

// List reports that guest storage host calls are unavailable.
func (unsupportedStorageHostService) List(
	_ string,
	_ uint32,
) ([]*HostServiceStorageObject, error) {
	return nil, ErrHostCallsUnavailable
}

// Stat reports that guest storage host calls are unavailable.
func (unsupportedStorageHostService) Stat(
	_ string,
) (*HostServiceStorageObject, bool, error) {
	return nil, false, ErrHostCallsUnavailable
}

// Request reports that guest outbound HTTP host calls are unavailable.
func (unsupportedHTTPHostService) Request(
	_ string,
	_ *HostServiceNetworkRequest,
) (*HostServiceNetworkResponse, error) {
	return nil, ErrHostCallsUnavailable
}

// List reports that guest structured data host calls are unavailable.
func (unsupportedDataHostService) List(
	_ string,
	_ map[string]string,
	_ int32,
	_ int32,
) (*DataListResult, error) {
	return nil, ErrHostCallsUnavailable
}

// ListRequest reports that guest structured data host calls are unavailable.
func (unsupportedDataHostService) ListRequest(
	_ string,
	_ *HostServiceDataListRequest,
) (*HostServiceDataListResponse, error) {
	return nil, ErrHostCallsUnavailable
}

// Get reports that guest structured data host calls are unavailable.
func (unsupportedDataHostService) Get(_ string, _ any) (map[string]any, bool, error) {
	return nil, false, ErrHostCallsUnavailable
}

// GetRequest reports that guest structured data host calls are unavailable.
func (unsupportedDataHostService) GetRequest(_ string, _ *HostServiceDataGetRequest) (*DataGetResult, error) {
	return nil, ErrHostCallsUnavailable
}

// Create reports that guest structured data host calls are unavailable.
func (unsupportedDataHostService) Create(_ string, _ map[string]any) (*DataMutationResult, error) {
	return nil, ErrHostCallsUnavailable
}

// Update reports that guest structured data host calls are unavailable.
func (unsupportedDataHostService) Update(_ string, _ any, _ map[string]any) (*DataMutationResult, error) {
	return nil, ErrHostCallsUnavailable
}

// Delete reports that guest structured data host calls are unavailable.
func (unsupportedDataHostService) Delete(_ string, _ any) (*DataMutationResult, error) {
	return nil, ErrHostCallsUnavailable
}

// Transaction reports that guest structured data host calls are unavailable.
func (unsupportedDataHostService) Transaction(_ string, _ []*DataTransactionInput) (*DataTransactionResult, error) {
	return nil, ErrHostCallsUnavailable
}

// Get reports that guest cache host calls are unavailable.
func (unsupportedCacheHostService) Get(_ string, _ string) (*HostServiceCacheValue, bool, error) {
	return nil, false, ErrHostCallsUnavailable
}

// Set reports that guest cache host calls are unavailable.
func (unsupportedCacheHostService) Set(_ string, _ string, _ string, _ int64) (*HostServiceCacheValue, error) {
	return nil, ErrHostCallsUnavailable
}

// Delete reports that guest cache host calls are unavailable.
func (unsupportedCacheHostService) Delete(_ string, _ string) error {
	return ErrHostCallsUnavailable
}

// Incr reports that guest cache host calls are unavailable.
func (unsupportedCacheHostService) Incr(_ string, _ string, _ int64, _ int64) (*HostServiceCacheValue, error) {
	return nil, ErrHostCallsUnavailable
}

// Expire reports that guest cache host calls are unavailable.
func (unsupportedCacheHostService) Expire(_ string, _ string, _ int64) (bool, string, error) {
	return false, "", ErrHostCallsUnavailable
}

// Acquire reports that guest lock host calls are unavailable.
func (unsupportedLockHostService) Acquire(_ string, _ int64) (*HostServiceLockAcquireResponse, error) {
	return nil, ErrHostCallsUnavailable
}

// Renew reports that guest lock host calls are unavailable.
func (unsupportedLockHostService) Renew(_ string, _ string) (*HostServiceLockRenewResponse, error) {
	return nil, ErrHostCallsUnavailable
}

// Release reports that guest lock host calls are unavailable.
func (unsupportedLockHostService) Release(_ string, _ string) error {
	return ErrHostCallsUnavailable
}

// Get reports that guest config host calls are unavailable.
func (unsupportedConfigHostService) Get(_ string) (string, bool, error) {
	return "", false, ErrHostCallsUnavailable
}

// Exists reports that guest config host calls are unavailable.
func (unsupportedConfigHostService) Exists(_ string) (bool, error) {
	return false, ErrHostCallsUnavailable
}

// String reports that guest config host calls are unavailable.
func (unsupportedConfigHostService) String(_ string) (string, bool, error) {
	return "", false, ErrHostCallsUnavailable
}

// Bool reports that guest config host calls are unavailable.
func (unsupportedConfigHostService) Bool(_ string) (bool, bool, error) {
	return false, false, ErrHostCallsUnavailable
}

// Int reports that guest config host calls are unavailable.
func (unsupportedConfigHostService) Int(_ string) (int, bool, error) {
	return 0, false, ErrHostCallsUnavailable
}

// Duration reports that guest config host calls are unavailable.
func (unsupportedConfigHostService) Duration(_ string) (time.Duration, bool, error) {
	return 0, false, ErrHostCallsUnavailable
}

// Send reports that guest notify host calls are unavailable.
func (unsupportedNotifyHostService) Send(
	_ string,
	_ *HostServiceNotifySendRequest,
) (*HostServiceNotifySendResponse, error) {
	return nil, ErrHostCallsUnavailable
}

// Register reports that guest cron host calls are unavailable.
func (unsupportedCronHostService) Register(_ *CronContract) error {
	return ErrHostCallsUnavailable
}

// Get reports that guest host config calls are unavailable.
func (unsupportedHostConfigHostService) Get(_ string) (string, bool, error) {
	return "", false, ErrHostCallsUnavailable
}

// String reports that guest host config calls are unavailable.
func (unsupportedHostConfigHostService) String(_ string) (string, bool, error) {
	return "", false, ErrHostCallsUnavailable
}

// Bool reports that guest host config calls are unavailable.
func (unsupportedHostConfigHostService) Bool(_ string) (bool, bool, error) {
	return false, false, ErrHostCallsUnavailable
}

// Int reports that guest host config calls are unavailable.
func (unsupportedHostConfigHostService) Int(_ string) (int, bool, error) {
	return 0, false, ErrHostCallsUnavailable
}

// Duration reports that guest host config calls are unavailable.
func (unsupportedHostConfigHostService) Duration(_ string) (time.Duration, bool, error) {
	return 0, false, ErrHostCallsUnavailable
}

// Get reports that guest manifest host calls are unavailable.
func (unsupportedManifestHostService) Get(_ string) ([]byte, bool, error) {
	return nil, false, ErrHostCallsUnavailable
}

// GetText reports that guest manifest host calls are unavailable.
func (unsupportedManifestHostService) GetText(_ string) (string, bool, error) {
	return "", false, ErrHostCallsUnavailable
}

// Scan reports that guest manifest host calls are unavailable.
func (unsupportedManifestHostService) Scan(_ string, _ string, _ any) (bool, error) {
	return false, ErrHostCallsUnavailable
}

// HostDBQuery is no longer part of the public host service protocol.
func HostDBQuery(_ string, _ []string, _ int) (*HostDBQueryResult, error) {
	return nil, gerror.New("HostDBQuery has been removed; use the structured pluginbridge.Data() service instead")
}

// HostDBExecute is no longer part of the public host service protocol.
func HostDBExecute(_ string, _ []string) (int64, int64, error) {
	return 0, 0, gerror.New("HostDBExecute has been removed; use the structured pluginbridge.Data() service instead")
}
