<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';
import type { AuthApi } from '#/api/core/auth';

import { computed, onMounted, reactive, ref } from 'vue';

import {
  AuthenticationLogin,
  useVbenForm,
  VbenButton,
  z,
} from '@vben/common-ui';
import { $t } from '@vben/locales';
import { IconifyIcon } from '@vben/icons';

import { message } from 'ant-design-vue';

import {
  authorizeAuthProviderApi,
  getAuthProvidersApi,
} from '#/api/core/auth';
import PluginSlotOutlet from '#/components/plugin/plugin-slot-outlet.vue';
import { pluginSlotKeys } from '#/plugins/plugin-slots';
import { publicFrontendSettings } from '#/runtime/public-frontend';
import { useAuthStore, useTenantStore } from '#/store';

defineOptions({ name: 'Login' });

const authStore = useAuthStore();
const tenantStore = useTenantStore();
const externalProviders = ref<AuthApi.AuthProvider[]>([]);
const externalProviderLoading = ref(false);
const externalProviderAuthorizing = ref('');
const tenantOptions = computed(() =>
  tenantStore.tenants.map((tenant) => ({
    code: tenant.code,
    label: `${tenant.name} (${tenant.code})`,
    name: tenant.name,
    value: String(tenant.id),
  })),
);
const loginSubtitle = computed(
  () =>
    publicFrontendSettings.auth.loginSubtitle ||
    $t('authentication.loginSubtitle'),
);

const tenantSubtitle = computed(() =>
  $t('pages.multiTenant.login.selectTenantSubtitle'),
);

const hasExternalProviders = computed(() => externalProviders.value.length > 0);

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('authentication.usernameTip'),
      },
      fieldName: 'username',
      label: $t('authentication.username'),
      rules: z.string().min(1, { message: $t('authentication.usernameTip') }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        placeholder: $t('authentication.passwordTip'),
      },
      fieldName: 'password',
      label: $t('authentication.password'),
      rules: z.string().min(1, { message: $t('authentication.passwordTip') }),
    },
  ];
});

const tenantFormSchema = computed((): VbenFormSchema[] => [
  {
    component: 'VbenSelect',
    componentProps: {
      class: 'h-11',
      options: tenantOptions.value,
      placeholder: $t('pages.multiTenant.login.selectTenant'),
    },
    fieldName: 'tenantId',
    label: $t('pages.multiTenant.login.selectTenant'),
    rules: 'selectRequired',
  },
]);

const [TenantForm, tenantFormApi] = useVbenForm(
  reactive({
    commonConfig: {
      hideLabel: true,
      hideRequiredMark: true,
    },
    schema: tenantFormSchema,
    showDefaultActions: false,
  }),
);

async function handleSubmit(values: Record<string, any>) {
  const result = await authStore.authLogin(values);
  if (result.requiresTenantSelection && result.tenants?.[0]) {
    await tenantFormApi.setFieldValue(
      'tenantId',
      String(result.tenants[0].id),
    );
  }
}

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

async function loadExternalProviders() {
  try {
    externalProviderLoading.value = true;
    externalProviders.value = await getAuthProvidersApi();
  } catch {
    externalProviders.value = [];
  } finally {
    externalProviderLoading.value = false;
  }
}

async function handleExternalProvider(provider: AuthApi.AuthProvider) {
  try {
    externalProviderAuthorizing.value = provider.providerKey;
    const result = await authorizeAuthProviderApi(provider.providerKey, {
      purpose: 'login',
      redirectUri: '/auth/login',
    });
    if (!result.redirectUrl) {
      message.error($t('pages.auth.externalProviders.authorizeFailed'));
      return;
    }
    window.location.assign(result.redirectUrl);
  } finally {
    externalProviderAuthorizing.value = '';
  }
}

async function handleSelectTenant() {
  const { valid } = await tenantFormApi.validate();
  if (!valid) {
    return;
  }
  const values = await tenantFormApi.getValues<{ tenantId?: string }>();
  const tenantId = Number(values.tenantId);
  if (!Number.isFinite(tenantId) || tenantId <= 0) {
    return;
  }
  await authStore.selectTenant(tenantId);
}

onMounted(loadExternalProviders);
</script>

<template>
  <div>
    <AuthenticationLogin
      v-if="!authStore.pendingPreToken"
      :form-schema="formSchema"
      :loading="authStore.loginLoading"
      :show-code-login="false"
      :show-forget-password="false"
      :show-qrcode-login="false"
      :show-register="false"
      :show-third-party-login="false"
      :sub-title="loginSubtitle"
      @submit="handleSubmit"
    >
      <template #third-party-login>
        <div
          v-if="hasExternalProviders"
          class="mt-5"
          data-testid="external-auth-provider-list"
        >
          <div class="mb-3 flex items-center gap-3">
            <div class="h-px flex-1 bg-border"></div>
            <span class="text-xs text-muted-foreground">
              {{ $t('pages.auth.externalProviders.title') }}
            </span>
            <div class="h-px flex-1 bg-border"></div>
          </div>
          <div class="grid grid-cols-1 gap-2 sm:grid-cols-2">
            <VbenButton
              v-for="provider in externalProviders"
              :key="provider.providerKey"
              :aria-label="
                $t('pages.auth.externalProviders.signInWith', {
                  provider: provider.name,
                })
              "
              :data-testid="`external-auth-provider-${provider.providerKey}`"
              :loading="externalProviderAuthorizing === provider.providerKey"
              class="min-h-10 justify-center gap-2"
              variant="outline"
              @click="handleExternalProvider(provider)"
            >
              <IconifyIcon
                :icon="resolveProviderIcon(provider)"
                class="size-4 shrink-0"
              />
              <span class="truncate">{{ provider.name }}</span>
            </VbenButton>
          </div>
        </div>
        <div
          v-else-if="!externalProviderLoading"
          class="hidden"
          data-testid="external-auth-provider-empty"
        ></div>
      </template>
    </AuthenticationLogin>
    <div
      v-else
      data-testid="login-tenant-selector"
      @keydown.enter.prevent="handleSelectTenant"
    >
      <div class="mb-7 sm:mx-auto sm:w-full sm:max-w-md">
        <h2
          class="mb-3 text-3xl/9 font-bold tracking-tight text-foreground lg:text-4xl"
        >
          {{ $t('pages.multiTenant.login.selectTenant') }}
        </h2>
        <p class="lg:text-md text-sm text-muted-foreground">
          {{ tenantSubtitle }}
        </p>
      </div>
      <TenantForm class="mb-8" data-testid="login-tenant-form" />
      <VbenButton
        :class="{
          'cursor-wait': authStore.loginLoading,
        }"
        :loading="authStore.loginLoading"
        aria-label="select tenant"
        class="w-full"
        data-testid="login-tenant-confirm"
        @click="handleSelectTenant"
      >
        {{ $t('pages.multiTenant.login.enterTenant') }}
      </VbenButton>
    </div>
    <PluginSlotOutlet :slot-key="pluginSlotKeys.authLoginAfter" class="mt-4" />
  </div>
</template>
