<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { LOGIN_PATH } from '@vben/constants';

import { message, Spin } from 'ant-design-vue';

import { callbackAuthProviderApi } from '#/api/core/auth';
import { $t } from '#/locales';
import { useAuthStore } from '#/store';

defineOptions({ name: 'AuthProviderCallback' });

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const failed = ref(false);

const providerKey = computed(() => readCallbackValue('providerKey') || '');

function readQueryValue(key: string) {
  const value = route.query[key];
  const firstValue = Array.isArray(value) ? value[0] : value;
  return firstValue ?? undefined;
}

function readHashValue(key: string) {
  const rawHash = window.location.hash.startsWith('#')
    ? window.location.hash.slice(1)
    : window.location.hash;
  const queryStart = rawHash.indexOf('?');
  const search = queryStart >= 0 ? rawHash.slice(queryStart + 1) : rawHash;
  return new URLSearchParams(search).get(key) ?? undefined;
}

function readCallbackValue(key: string) {
  return readHashValue(key) || readQueryValue(key);
}

async function handleCallback() {
  const payload = readCallbackValue('payload');
  if (payload) {
    await handleHostCallbackPayload(payload);
    return;
  }

  const state = readCallbackValue('state');
  if (!providerKey.value || !state) {
    failed.value = true;
    message.error($t('pages.auth.externalProviders.callbackFailed'));
    await router.replace(LOGIN_PATH);
    return;
  }

  try {
    const result = await callbackAuthProviderApi(providerKey.value, {
      code: readCallbackValue('code'),
      error: readCallbackValue('error'),
      state,
    });
    const redirectUri = result.redirectUri || undefined;
    if (redirectUri && !result.accessToken && !result.preToken) {
      message.success($t('pages.profile.security.externalAuth.bindSuccess'));
      await router.replace(redirectUri);
      return;
    }
    await authStore.completeExternalLogin(result);
  } catch {
    failed.value = true;
    message.error($t('pages.auth.externalProviders.callbackFailed'));
    await router.replace(LOGIN_PATH);
  }
}

async function handleHostCallbackPayload(payload: string) {
  try {
    const result = JSON.parse(window.atob(toBase64(payload)));
    const redirectUri =
      result.redirectUri || readCallbackValue('target') || undefined;
    if (redirectUri && !result.accessToken && !result.preToken) {
      message.success($t('pages.profile.security.externalAuth.bindSuccess'));
      await router.replace(redirectUri);
      return;
    }
    await authStore.completeExternalLogin(result);
  } catch {
    failed.value = true;
    message.error($t('pages.auth.externalProviders.callbackFailed'));
    await router.replace(LOGIN_PATH);
  }
}

function toBase64(raw: string) {
  const normalized = raw.replaceAll('-', '+').replaceAll('_', '/');
  const padLength = (4 - (normalized.length % 4)) % 4;
  return normalized + '='.repeat(padLength);
}

onMounted(handleCallback);
</script>

<template>
  <div
    class="flex min-h-[220px] flex-col items-center justify-center gap-4 text-center"
    data-testid="external-auth-callback"
  >
    <Spin :spinning="!failed" />
    <p class="text-sm text-muted-foreground">
      {{
        failed
          ? $t('pages.auth.externalProviders.callbackFailed')
          : $t('pages.auth.externalProviders.callbackProcessing')
      }}
    </p>
  </div>
</template>
