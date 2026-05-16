// This file defines pluggable authentication provider business error codes.

package authprovider

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeAuthProviderNotFound reports an unknown or disabled external provider.
	CodeAuthProviderNotFound = bizerr.MustDefine(
		"AUTH_PROVIDER_NOT_FOUND",
		"Authentication provider is not available",
		gcode.CodeNotFound,
	)
	// CodeAuthProviderUnavailable reports that provider metadata exists but its plugin handler is unavailable.
	CodeAuthProviderUnavailable = bizerr.MustDefine(
		"AUTH_PROVIDER_UNAVAILABLE",
		"Authentication provider plugin is not available",
		gcode.CodeNotAuthorized,
	)
	// CodeAuthStateInvalid reports invalid, expired, or already consumed authorization state.
	CodeAuthStateInvalid = bizerr.MustDefine(
		"AUTH_PROVIDER_STATE_INVALID",
		"Authentication provider state is invalid or expired",
		gcode.CodeNotAuthorized,
	)
	// CodeAuthIdentityNotBound reports an external identity without a local binding.
	CodeAuthIdentityNotBound = bizerr.MustDefine(
		"AUTH_IDENTITY_NOT_BOUND",
		"External identity is not bound to a local user",
		gcode.CodeNotAuthorized,
	)
	// CodeAuthIdentityAlreadyBound reports a duplicate external identity binding.
	CodeAuthIdentityAlreadyBound = bizerr.MustDefine(
		"AUTH_IDENTITY_ALREADY_BOUND",
		"External identity is already bound",
		gcode.CodeValidationFailed,
	)
	// CodeAuthIdentityNotFound reports a missing current-user identity binding.
	CodeAuthIdentityNotFound = bizerr.MustDefine(
		"AUTH_IDENTITY_NOT_FOUND",
		"External identity binding does not exist",
		gcode.CodeNotFound,
	)
)
