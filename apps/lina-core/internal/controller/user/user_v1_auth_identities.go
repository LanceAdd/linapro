// This file handles current-user external authentication identity endpoints.

package user

import (
	"context"

	v1 "lina-core/api/user/v1"
	authprovidersvc "lina-core/internal/service/authprovider"
)

const authIdentityBindPurpose = "bind"

// AuthIdentityList lists current-user external authentication identities.
func (c *ControllerV1) AuthIdentityList(ctx context.Context, _ *v1.AuthIdentityListReq) (res *v1.AuthIdentityListRes, err error) {
	userID := 0
	if c.bizCtxSvc != nil {
		if localCtx := c.bizCtxSvc.Get(ctx); localCtx != nil {
			userID = localCtx.UserId
		}
	}
	items, err := c.authProviderSvc.ListCurrentUserIdentities(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &v1.AuthIdentityListRes{List: toAuthIdentityEntities(items)}, nil
}

// AuthIdentityUnbind deletes one current-user external authentication identity.
func (c *ControllerV1) AuthIdentityUnbind(ctx context.Context, req *v1.AuthIdentityUnbindReq) (res *v1.AuthIdentityUnbindRes, err error) {
	userID := 0
	if c.bizCtxSvc != nil {
		if localCtx := c.bizCtxSvc.Get(ctx); localCtx != nil {
			userID = localCtx.UserId
		}
	}
	if err = c.authProviderSvc.UnbindCurrentUserProvider(ctx, userID, req.ProviderKey); err != nil {
		return nil, err
	}
	return &v1.AuthIdentityUnbindRes{}, nil
}

// AuthIdentityBind starts current-user external authentication identity binding.
func (c *ControllerV1) AuthIdentityBind(ctx context.Context, req *v1.AuthIdentityBindReq) (res *v1.AuthIdentityBindRes, err error) {
	userID := 0
	if c.bizCtxSvc != nil {
		if localCtx := c.bizCtxSvc.Get(ctx); localCtx != nil {
			userID = localCtx.UserId
		}
	}
	output, err := c.authProviderSvc.Authorize(ctx, authprovidersvc.AuthorizeInput{
		ProviderKey: req.ProviderKey,
		Purpose:     authIdentityBindPurpose,
		RedirectURI: req.RedirectUri,
		UserID:      userID,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AuthIdentityBindRes{
		RedirectUrl: output.RedirectURL,
		State:       output.State,
	}, nil
}

// toAuthIdentityEntities converts service identity items into API DTOs.
func toAuthIdentityEntities(items []authprovidersvc.IdentityItem) []*v1.AuthIdentityEntity {
	if len(items) == 0 {
		return []*v1.AuthIdentityEntity{}
	}
	out := make([]*v1.AuthIdentityEntity, 0, len(items))
	for _, item := range items {
		entity := &v1.AuthIdentityEntity{
			ProviderKey:      item.ProviderKey,
			ProviderType:     item.ProviderType,
			Subject:          item.Subject,
			ExternalTenantId: item.ExternalTenantID,
			Email:            item.Email,
			Mobile:           item.Mobile,
			DisplayName:      item.DisplayName,
			Avatar:           item.Avatar,
		}
		if item.LastLoginAt != nil {
			entity.LastLoginAt = item.LastLoginAt.String()
		}
		if item.BoundAt != nil {
			entity.BoundAt = item.BoundAt.String()
		}
		out = append(out, entity)
	}
	return out
}
