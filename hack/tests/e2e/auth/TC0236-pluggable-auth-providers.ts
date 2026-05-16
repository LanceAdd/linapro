import { Buffer } from 'node:buffer';

import type { Page, Route } from '@playwright/test';

import { test, expect } from '../../fixtures/auth';
import { config } from '../../fixtures/config';
import { LoginPage } from '../../pages/LoginPage';
import { ProfilePage } from '../../pages/ProfilePage';

const providerList = [
  {
    icon: 'brand-google',
    name: 'Google',
    providerKey: 'google',
    providerType: 'oidc',
    sort: 10,
  },
  {
    icon: 'brand-microsoft',
    name: 'Microsoft',
    providerKey: 'microsoft',
    providerType: 'oidc',
    sort: 20,
  },
];

const boundIdentity = {
  boundAt: '2026-05-16 10:00:00',
  displayName: 'Alex Google',
  email: 'alex@example.com',
  providerKey: 'google',
  providerType: 'oidc',
  subject: 'google-subject-1',
};

function ok(data: unknown) {
  return { code: 0, data, message: 'ok' };
}

function callbackPayload(data: unknown) {
  return Buffer.from(JSON.stringify(data))
    .toString('base64')
    .replaceAll('+', '-')
    .replaceAll('/', '_')
    .replace(/=+$/u, '');
}

async function fulfillJson(route: Route, data: unknown) {
  await route.fulfill({
    contentType: 'application/json',
    json: ok(data),
    status: 200,
  });
}

async function mockProviderList(page: Page, providers = providerList) {
  await page.route('**/api/v1/auth/providers', async (route) => {
    if (route.request().method() !== 'GET') {
      await route.fallback();
      return;
    }
    await fulfillJson(route, { list: providers });
  });
}

test.describe('TC-236 pluggable authentication providers', () => {
  test('TC-236a: login page renders enabled external providers and starts authorization', async ({
    page,
  }) => {
    const loginPage = new LoginPage(page);
    let authorizeCalled = false;

    await mockProviderList(page);
    await page.route(
      '**/api/v1/auth/providers/google/authorize**',
      async (route) => {
        authorizeCalled = true;
        await fulfillJson(route, {
          redirectUrl: `${config.baseURL}/auth/login?external=google`,
          state: 'mock-login-state',
        });
      },
    );

    await loginPage.goto();
    await expect(loginPage.externalAuthProviderList).toBeVisible();
    await expect(loginPage.externalAuthProviderButton('google')).toBeVisible();
    await expect(
      loginPage.externalAuthProviderButton('microsoft'),
    ).toBeVisible();

    await loginPage.clickExternalAuthProvider('google');
    await expect
      .poll(() => authorizeCalled)
      .toBeTruthy();
    await expect(page).toHaveURL(/external=google/);
  });

  test('TC-236b: empty provider list keeps password login available without external entry', async ({
    page,
  }) => {
    const loginPage = new LoginPage(page);

    await mockProviderList(page, []);
    await loginPage.goto();

    await expect(loginPage.usernameInput).toBeVisible();
    await expect(loginPage.passwordInput).toBeVisible();
    await expect(loginPage.loginButton).toBeVisible();
    await expect(loginPage.externalAuthProviderList).toHaveCount(0);
    await expect(loginPage.externalAuthProviderEmpty).toHaveCount(1);
  });

  test('TC-236c: callback bridge consumes fragment payload without a callback API request', async ({
    page,
  }) => {
    let callbackCalled = false;

    await page.route(
      '**/api/v1/auth/providers/google/callback**',
      async (route) => {
        callbackCalled = true;
        await route.fallback();
      },
    );

    const target = '/auth/login?externalBind=google';
    const payload = callbackPayload({
      redirectUri: target,
      tenants: [],
    });
    await page.goto(
      `/auth/providers/callback#/auth/providers/callback?providerKey=google&payload=${payload}&target=${encodeURIComponent(target)}`,
      { waitUntil: 'domcontentloaded' },
    );

    await expect.poll(() => callbackCalled).toBeFalsy();
    await expect(page).toHaveURL(/externalBind=google/);
  });

  test('TC-236d: profile security tab lists, binds, and unbinds external identities', async ({
    adminPage,
  }) => {
    const profilePage = new ProfilePage(adminPage);
    let identities = [boundIdentity];
    let bindCalled = false;
    let unbindCalled = false;

    await mockProviderList(adminPage);
    await adminPage.route('**/api/v1/user/auth-identities', async (route) => {
      if (route.request().method() === 'GET') {
        await fulfillJson(route, { list: identities });
        return;
      }
      await route.fallback();
    });
    await adminPage.route(
      '**/api/v1/user/auth-identities/microsoft/bind',
      async (route) => {
        bindCalled = true;
        await fulfillJson(route, {
          redirectUrl: `${config.baseURL}/profile?externalBind=microsoft`,
          state: 'mock-bind-state',
        });
      },
    );
    await adminPage.route(
      '**/api/v1/user/auth-identities/google',
      async (route) => {
        if (route.request().method() !== 'DELETE') {
          await route.fallback();
          return;
        }
        unbindCalled = true;
        identities = [];
        await fulfillJson(route, {});
      },
    );

    await profilePage.goto();
    await profilePage.openSecurityTab();

    await expect(profilePage.authProvider('google')).toContainText(
      'Alex Google',
    );
    await expect(profilePage.authProviderUnbind('google')).toBeVisible();
    await expect(profilePage.authProviderBind('microsoft')).toBeVisible();

    await profilePage.authProviderBind('microsoft').click();
    await expect
      .poll(() => bindCalled)
      .toBeTruthy();
    await expect(adminPage).toHaveURL(/externalBind=microsoft/);

    await profilePage.goto();
    await profilePage.openSecurityTab();
    await profilePage.authProviderUnbind('google').click();
    await adminPage
      .locator('.ant-popconfirm-buttons .ant-btn-primary')
      .last()
      .click();
    await expect
      .poll(() => unbindCalled)
      .toBeTruthy();
    await expect(profilePage.authProviderBind('google')).toBeVisible();
  });
});
