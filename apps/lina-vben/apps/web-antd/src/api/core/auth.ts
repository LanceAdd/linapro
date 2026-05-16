import type { TenantAwareLoginResult } from '#/api/tenant/model';

import { requestClient } from '#/api/request';

export namespace AuthApi {
  /** 登录接口参数 */
  export interface LoginParams {
    password?: string;
    username?: string;
  }

  /** 登录接口返回值 */
  export interface LoginResult extends TenantAwareLoginResult {
    redirectUri?: string;
  }

  /** 刷新 token 接口参数 */
  export interface RefreshTokenParams {
    refreshToken: string;
  }

  /** 刷新 token 接口返回值 */
  export interface RefreshTokenResult {
    accessToken: string;
    refreshToken?: string;
  }

  /** External authentication provider shown on login and profile pages. */
  export interface AuthProvider {
    icon?: string;
    name: string;
    providerKey: string;
    providerType: string;
    sort: number;
  }

  /** External authentication provider list response. */
  export interface AuthProviderListResult {
    list: AuthProvider[];
  }

  /** External authentication authorization query parameters. */
  export interface AuthProviderAuthorizeParams {
    purpose?: 'bind' | 'login';
    redirectUri?: string;
  }

  /** External authentication authorization response. */
  export interface AuthProviderAuthorizeResult {
    redirectUrl: string;
    state: string;
  }

  /** External authentication callback query parameters. */
  export interface AuthProviderCallbackParams {
    code?: string;
    error?: string;
    state: string;
  }
}

/**
 * 登录
 */
export async function loginApi(data: AuthApi.LoginParams) {
  return requestClient.post<AuthApi.LoginResult>('/auth/login', data);
}

/**
 * 退出登录
 */
export async function logoutApi() {
  return requestClient.post('/auth/logout');
}

/**
 * 刷新 access token
 */
export async function refreshTokenApi(data: AuthApi.RefreshTokenParams) {
  return requestClient.post<AuthApi.RefreshTokenResult>('/auth/refresh', data);
}

/**
 * List external authentication providers available to the current frontend.
 */
export async function getAuthProvidersApi() {
  const res = await requestClient.get<AuthApi.AuthProviderListResult>(
    '/auth/providers',
  );
  return Array.isArray(res.list) ? res.list : [];
}

/**
 * Create an authorization state for one external authentication provider.
 */
export async function authorizeAuthProviderApi(
  providerKey: string,
  params?: AuthApi.AuthProviderAuthorizeParams,
) {
  return requestClient.get<AuthApi.AuthProviderAuthorizeResult>(
    `/auth/providers/${providerKey}/authorize`,
    { params },
  );
}

/**
 * Consume an external authentication callback.
 */
export async function callbackAuthProviderApi(
  providerKey: string,
  params: AuthApi.AuthProviderCallbackParams,
) {
  return requestClient.get<AuthApi.LoginResult>(
    `/auth/providers/${providerKey}/callback`,
    { params },
  );
}
