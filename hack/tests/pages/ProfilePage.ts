import type { Locator, Page } from '@playwright/test';

import { waitForRouteReady } from '../support/ui';

export class ProfilePage {
  constructor(private page: Page) {}

  get profilePanel(): Locator {
    return this.page.locator('.ant-card').first();
  }

  get nicknameInput(): Locator {
    return this.page.getByTestId('profile-base-form').locator('input').first();
  }

  get passwordTab(): Locator {
    return this.page.getByRole('tab').nth(1);
  }

  get securityTab(): Locator {
    return this.page.getByRole('tab').nth(2);
  }

  get authIdentityList(): Locator {
    return this.page.getByTestId('profile-auth-identities');
  }

  get authIdentityEmpty(): Locator {
    return this.page.getByTestId('profile-auth-identities-empty');
  }

  panelDisplayName(name: string): Locator {
    return this.profilePanel.getByText(name, { exact: true }).first();
  }

  async goto() {
    await this.page.goto('/profile');
    await waitForRouteReady(this.page);
    await this.nicknameInput.waitFor({ state: 'visible', timeout: 10000 });
  }

  async openPasswordTab() {
    await this.passwordTab.click();
    await waitForRouteReady(this.page);
  }

  async openSecurityTab() {
    await this.securityTab.click();
    await waitForRouteReady(this.page);
    await this.authIdentityList.waitFor({ state: 'visible', timeout: 10000 });
  }

  authProvider(providerKey: string): Locator {
    return this.page.getByTestId(`profile-auth-provider-${providerKey}`);
  }

  authProviderBind(providerKey: string): Locator {
    return this.page.getByTestId(`profile-auth-bind-${providerKey}`);
  }

  authProviderUnbind(providerKey: string): Locator {
    return this.page.getByTestId(`profile-auth-unbind-${providerKey}`);
  }
}
