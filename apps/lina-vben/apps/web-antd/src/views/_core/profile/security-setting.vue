<script setup lang="ts">
import type { AuthApi } from '#/api/core/auth';
import type { AuthIdentity } from '#/api/system/user';

import { computed, onMounted, ref } from 'vue';

import { IconifyIcon } from '@vben/icons';

import { Button, message, Popconfirm, Spin, Tag } from 'ant-design-vue';

import { getAuthProvidersApi } from '#/api/core/auth';
import {
  bindAuthIdentity,
  getAuthIdentities,
  unbindAuthIdentity,
} from '#/api/system/user';
import { $t } from '#/locales';

defineOptions({ name: 'ProfileSecuritySetting' });

type ProviderView = AuthApi.AuthProvider & {
  identity?: AuthIdentity;
};

const loading = ref(false);
const actionProviderKey = ref('');
const providers = ref<AuthApi.AuthProvider[]>([]);
const identities = ref<AuthIdentity[]>([]);

const providerViews = computed<ProviderView[]>(() => {
  const identityByProvider = new Map(
    identities.value.map((item) => [item.providerKey, item]),
  );
  return providers.value.map((provider) => ({
    ...provider,
    identity: identityByProvider.get(provider.providerKey),
  }));
});

function resolveProviderIcon(provider: AuthApi.AuthProvider) {
  const icon = provider.icon?.trim();
  if (!icon) {
    return 'lucide:key-round';
  }
  if (icon.includes(':') || icon.startsWith('svg:')) {
    return icon;
  }
  const builtInIcons: Record<string, string> = {
    'brand-github': 'mdi:github',
    'brand-google': 'svg:google',
    'brand-microsoft': 'mdi:microsoft',
    github: 'mdi:github',
    google: 'svg:google',
    microsoft: 'mdi:microsoft',
    oidc: 'lucide:key-round',
    wechat: 'svg:wechat',
    wecom: 'svg:wechat',
  };
  return builtInIcons[icon] || builtInIcons[provider.providerType] || 'lucide:key-round';
}

function identityDescription(identity?: AuthIdentity) {
  if (!identity) {
    return $t('pages.profile.security.externalAuth.unboundDesc');
  }
  return (
    identity.displayName ||
    identity.email ||
    identity.mobile ||
    identity.subject ||
    $t('pages.profile.security.externalAuth.boundDesc')
  );
}

async function loadExternalAuth() {
  try {
    loading.value = true;
    const [nextProviders, nextIdentities] = await Promise.all([
      getAuthProvidersApi(),
      getAuthIdentities(),
    ]);
    providers.value = nextProviders;
    identities.value = nextIdentities;
  } catch {
    providers.value = [];
    identities.value = [];
  } finally {
    loading.value = false;
  }
}

async function handleBind(provider: AuthApi.AuthProvider) {
  try {
    actionProviderKey.value = provider.providerKey;
    const result = await bindAuthIdentity(provider.providerKey, '/profile');
    if (!result.redirectUrl) {
      message.error($t('pages.profile.security.externalAuth.bindFailed'));
      return;
    }
    window.location.assign(result.redirectUrl);
  } finally {
    actionProviderKey.value = '';
  }
}

async function handleUnbind(provider: AuthApi.AuthProvider) {
  try {
    actionProviderKey.value = provider.providerKey;
    await unbindAuthIdentity(provider.providerKey);
    message.success($t('pages.profile.security.externalAuth.unbindSuccess'));
    await loadExternalAuth();
  } finally {
    actionProviderKey.value = '';
  }
}

onMounted(loadExternalAuth);
</script>

<template>
  <div
    class="mt-[16px] w-full max-w-[42rem]"
    data-testid="profile-auth-identities"
  >
    <div class="mb-4">
      <h3 class="text-base font-medium">
        {{ $t('pages.profile.security.externalAuth.title') }}
      </h3>
      <p class="mt-1 text-sm text-muted-foreground">
        {{ $t('pages.profile.security.externalAuth.description') }}
      </p>
    </div>

    <Spin :spinning="loading">
      <div
        v-if="providerViews.length === 0"
        class="rounded-md border border-dashed px-4 py-6 text-sm text-muted-foreground"
        data-testid="profile-auth-identities-empty"
      >
        {{ $t('pages.profile.security.externalAuth.empty') }}
      </div>
      <div v-else class="divide-y rounded-md border">
        <div
          v-for="provider in providerViews"
          :key="provider.providerKey"
          :data-testid="`profile-auth-provider-${provider.providerKey}`"
          class="flex flex-col gap-3 px-4 py-4 sm:flex-row sm:items-center sm:justify-between"
        >
          <div class="flex min-w-0 items-center gap-3">
            <div
              class="flex size-10 shrink-0 items-center justify-center rounded-md border bg-muted/30"
            >
              <IconifyIcon
                :icon="resolveProviderIcon(provider)"
                class="size-5"
              />
            </div>
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-2">
                <span class="font-medium">{{ provider.name }}</span>
                <Tag v-if="provider.identity" color="success">
                  {{ $t('pages.profile.security.externalAuth.bound') }}
                </Tag>
                <Tag v-else>
                  {{ $t('pages.profile.security.externalAuth.unbound') }}
                </Tag>
              </div>
              <p class="mt-1 break-all text-sm text-muted-foreground">
                {{ identityDescription(provider.identity) }}
              </p>
            </div>
          </div>

          <div class="shrink-0">
            <Popconfirm
              v-if="provider.identity"
              :description="
                $t('pages.profile.security.externalAuth.unbindConfirm')
              "
              :ok-text="$t('pages.common.confirm')"
              :cancel-text="$t('pages.common.cancel')"
              @confirm="handleUnbind(provider)"
            >
              <Button
                :data-testid="`profile-auth-unbind-${provider.providerKey}`"
                :loading="actionProviderKey === provider.providerKey"
              >
                {{ $t('pages.profile.security.externalAuth.unbind') }}
              </Button>
            </Popconfirm>
            <Button
              v-else
              :data-testid="`profile-auth-bind-${provider.providerKey}`"
              :loading="actionProviderKey === provider.providerKey"
              type="primary"
              @click="handleBind(provider)"
            >
              {{ $t('pages.profile.security.externalAuth.bind') }}
            </Button>
          </div>
        </div>
      </div>
    </Spin>
  </div>
</template>
