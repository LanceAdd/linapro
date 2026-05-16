// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package auth

import (
	"context"

	"lina-core/api/auth/v1"
)

type IAuthV1 interface {
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
	Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error)
	ProviderList(ctx context.Context, req *v1.ProviderListReq) (res *v1.ProviderListRes, err error)
	ProviderAuthorize(ctx context.Context, req *v1.ProviderAuthorizeReq) (res *v1.ProviderAuthorizeRes, err error)
	ProviderCallback(ctx context.Context, req *v1.ProviderCallbackReq) (res *v1.ProviderCallbackRes, err error)
	ProviderCallbackPost(ctx context.Context, req *v1.ProviderCallbackPostReq) (res *v1.ProviderCallbackPostRes, err error)
	Refresh(ctx context.Context, req *v1.RefreshReq) (res *v1.RefreshRes, err error)
}
