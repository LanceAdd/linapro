// This file defines the shared guest host-service client contracts used by
// both wasip1 host-call implementations and non-WASI unsupported stubs.

package guest

import "time"

// RuntimeHostService exposes guest-side helpers for the runtime host service.
type RuntimeHostService interface {
	// Log writes one structured runtime log entry through the host.
	Log(level int, message string, fields map[string]string) error
	// StateGet reads one plugin-scoped runtime state value by key.
	StateGet(key string) (string, bool, error)
	// StateSet writes one plugin-scoped runtime state value.
	StateSet(key string, value string) error
	// StateDelete removes one plugin-scoped runtime state value.
	StateDelete(key string) error
	// StateGetInt reads one integer runtime state value.
	StateGetInt(key string) (int, bool, error)
	// StateSetInt writes one integer runtime state value.
	StateSetInt(key string, value int) error
	// Now returns the current host time string.
	Now() (string, error)
	// UUID returns one host-generated unique identifier string.
	UUID() (string, error)
	// Node returns the current host node identity string.
	Node() (string, error)
}

// StorageHostService exposes guest-side helpers for the governed storage host service.
type StorageHostService interface {
	// Put writes one governed storage object under the given logical path.
	Put(objectPath string, body []byte, contentType string, overwrite bool) (*HostServiceStorageObject, error)
	// PutText writes one UTF-8 text object under the given logical path.
	PutText(objectPath string, content string, contentType string, overwrite bool) (*HostServiceStorageObject, error)
	// Get reads one governed storage object under the given logical path.
	Get(objectPath string) ([]byte, *HostServiceStorageObject, bool, error)
	// GetText reads one UTF-8 text object under the given logical path.
	GetText(objectPath string) (string, *HostServiceStorageObject, bool, error)
	// Delete removes one governed storage object under the given logical path.
	Delete(objectPath string) error
	// List lists governed storage objects under one logical path prefix.
	List(prefix string, limit uint32) ([]*HostServiceStorageObject, error)
	// Stat reads metadata for one governed storage object under the given logical path.
	Stat(objectPath string) (*HostServiceStorageObject, bool, error)
}

// HTTPHostService exposes guest-side helpers for the governed outbound HTTP host service.
type HTTPHostService interface {
	// Request executes one governed outbound HTTP request through the host.
	Request(targetURL string, request *HostServiceNetworkRequest) (*HostServiceNetworkResponse, error)
}

// DataHostService exposes the compatibility guest-side helpers for the governed
// structured data host service. New guest code should prefer plugindb.
type DataHostService interface {
	// List executes one governed structured data list request.
	List(table string, filters map[string]string, pageNum int32, pageSize int32) (*DataListResult, error)
	// ListRequest executes one governed structured data list request with the raw host-service request payload.
	ListRequest(table string, request *HostServiceDataListRequest) (*HostServiceDataListResponse, error)
	// Get reads one governed record by key from an authorized table.
	Get(table string, key any) (map[string]any, bool, error)
	// GetRequest executes one governed data get request with the raw host-service request payload.
	GetRequest(table string, request *HostServiceDataGetRequest) (*DataGetResult, error)
	// Create creates one governed record in an authorized table.
	Create(table string, record map[string]any) (*DataMutationResult, error)
	// Update updates one governed record in an authorized table.
	Update(table string, key any, record map[string]any) (*DataMutationResult, error)
	// Delete deletes one governed record in an authorized table.
	Delete(table string, key any) (*DataMutationResult, error)
	// Transaction executes one governed structured data transaction.
	Transaction(table string, operations []*DataTransactionInput) (*DataTransactionResult, error)
}

// DataListResult is the decoded guest-side result of one data list request.
type DataListResult struct {
	// Records is the ordered JSON-decoded result set.
	Records []map[string]any
	// Total is the total number of matching rows before pagination.
	Total int32
}

// DataMutationResult is the decoded guest-side result of one data mutation.
type DataMutationResult struct {
	// AffectedRows is the number of rows affected by the mutation.
	AffectedRows int64
	// Key is the JSON-decoded resource key returned by the host when available.
	Key any
	// Record is the optional JSON-decoded record snapshot returned by the host.
	Record map[string]any
}

// DataTransactionInput describes one guest-side transaction step.
type DataTransactionInput struct {
	// Method is one structured mutation method such as create/update/delete.
	Method string
	// Key is the optional resource key used by update/delete.
	Key any
	// Record is the optional input document used by create/update.
	Record map[string]any
}

// DataTransactionResult is the decoded guest-side result of one data transaction.
type DataTransactionResult struct {
	// Results is the ordered list of per-step mutation results.
	Results []*DataMutationResult
	// AffectedRows is the aggregate affected row count across all steps.
	AffectedRows int64
}

// DataGetResult is the decoded guest-side result of one data get request.
type DataGetResult struct {
	// Found reports whether the requested record exists.
	Found bool
	// Record is the optional JSON-decoded record snapshot returned by the host.
	Record map[string]any
}

// CacheHostService exposes guest-side helpers for the governed distributed cache host service.
type CacheHostService interface {
	// Get reads one governed cache value from the authorized namespace.
	Get(namespace string, key string) (*HostServiceCacheValue, bool, error)
	// Set writes one governed cache value into the authorized namespace.
	Set(namespace string, key string, value string, expireSeconds int64) (*HostServiceCacheValue, error)
	// Delete removes one governed cache value from the authorized namespace.
	Delete(namespace string, key string) error
	// Incr increments one governed cache integer value inside the authorized namespace.
	Incr(namespace string, key string, delta int64, expireSeconds int64) (*HostServiceCacheValue, error)
	// Expire updates one governed cache expiration policy inside the authorized namespace.
	Expire(namespace string, key string, expireSeconds int64) (bool, string, error)
}

// LockHostService exposes guest-side helpers for the governed distributed lock host service.
type LockHostService interface {
	// Acquire attempts to acquire one governed distributed lock.
	Acquire(lockName string, leaseMillis int64) (*HostServiceLockAcquireResponse, error)
	// Renew extends one governed distributed lock using the issued ticket.
	Renew(lockName string, ticket string) (*HostServiceLockRenewResponse, error)
	// Release releases one governed distributed lock using the issued ticket.
	Release(lockName string, ticket string) error
}

// ConfigHostService exposes guest-side helpers for the read-only config host service.
type ConfigHostService interface {
	// Get reads one plugin-scoped configuration value as JSON.
	Get(key string) (string, bool, error)
	// Exists reports whether one configuration key exists.
	Exists(key string) (bool, error)
	// String reads one configuration value as a string.
	String(key string) (string, bool, error)
	// Bool reads one configuration value as a bool.
	Bool(key string) (bool, bool, error)
	// Int reads one configuration value as an int.
	Int(key string) (int, bool, error)
	// Duration reads one configuration value as a duration.
	Duration(key string) (time.Duration, bool, error)
}

// NotifyHostService exposes guest-side helpers for the governed unified notify host service.
type NotifyHostService interface {
	// Send sends one governed notification through the authorized channel.
	Send(channelKey string, request *HostServiceNotifySendRequest) (*HostServiceNotifySendResponse, error)
}

// CronHostService exposes guest-side helpers for the cron registration host
// service.
type CronHostService interface {
	// Register submits one dynamic-plugin cron declaration to the current
	// host-side discovery collector.
	Register(contract *CronContract) error
}

// HostConfigHostService exposes guest-side helpers for whitelisted public host config.
type HostConfigHostService interface {
	// Get reads one whitelisted public host config value as JSON.
	Get(key string) (string, bool, error)
	// String reads one whitelisted public host config value as a string.
	String(key string) (string, bool, error)
	// Bool reads one whitelisted public host config value as a bool.
	Bool(key string) (bool, bool, error)
	// Int reads one whitelisted public host config value as an int.
	Int(key string) (int, bool, error)
	// Duration reads one whitelisted public host config value as a duration.
	Duration(key string) (time.Duration, bool, error)
}

// ManifestHostService exposes guest-side helpers for plugin-scoped manifest resources.
type ManifestHostService interface {
	// Get reads one manifest resource as bytes.
	Get(path string) ([]byte, bool, error)
	// GetText reads one manifest resource as UTF-8 text.
	GetText(path string) (string, bool, error)
	// Scan decodes a YAML manifest resource or nested key into target.
	Scan(path string, key string, target any) (bool, error)
}

// HostDBQueryResult preserves the previous guest-side result shape for callers
// that have not yet migrated to the structured data service SDK.
type HostDBQueryResult struct {
	// Columns lists the result set column names in row order.
	Columns []string
	// Rows stores the tabular result values as string slices per row.
	Rows [][]string
	// RowCount reports the total number of rows returned by the legacy query.
	RowCount int
}
