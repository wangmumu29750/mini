<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { fetchSystemSettings, updateSystemSettings } from '@/api/settings'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import type { ApiErrorPayload } from '@/types/api'
import type { SystemSetting } from '@/types/domain'

const settings = ref<SystemSetting[]>([])
const loading = ref(false)
const saving = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

onMounted(loadSettings)

async function loadSettings() {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    settings.value = await fetchSystemSettings()
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    settings.value = await updateSystemSettings(settings.value.map((item) => ({ key: item.key, value: item.value })))
    successMessage.value = '系统参数已保存。'
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    saving.value = false
  }
}

function labelOf(key: string) {
  const labels: Record<string, string> = {
    order_pay_expire_minutes: '订单支付超时分钟数',
    refund_cutoff_minutes: '退票截止分钟数',
    change_cutoff_minutes: '改签截止分钟数',
    refund_fee_percent: '退票手续费比例',
    mock_payment_enabled: '模拟支付开关',
  }
  return labels[key] || key
}
</script>

<template>
  <main>
    <PageHeader title="系统设置" description="管理员配置系统级参数，影响支付超时、退改签截止时间和课程级模拟能力。" />

    <section class="mx-auto max-w-5xl px-4 py-6 sm:px-8">
      <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
        <div class="text-sm font-bold text-slate-400">{{ loading ? '加载中...' : `共 ${settings.length} 个参数` }}</div>
        <div class="flex gap-2">
          <button class="btn-secondary" type="button" :disabled="loading || saving" @click="loadSettings">刷新</button>
          <button class="btn-primary" type="button" :disabled="loading || saving || !settings.length" @click="handleSave">
            {{ saving ? '保存中...' : '保存设置' }}
          </button>
        </div>
      </div>

      <p v-if="errorMessage" class="mb-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
      <p v-if="successMessage" class="mb-4 rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{{ successMessage }}</p>

      <div v-if="settings.length" class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
        <div v-for="setting in settings" :key="setting.key" class="grid gap-4 border-b border-slate-100 p-5 last:border-b-0 md:grid-cols-[1fr_220px] md:items-center">
          <div>
            <div class="text-base font-black text-slate-950">{{ labelOf(setting.key) }}</div>
            <div class="mt-1 text-sm font-medium text-slate-500">{{ setting.description }}</div>
            <div class="mt-2 text-xs font-bold uppercase tracking-normal text-slate-400">{{ setting.key }} / {{ setting.valueType }}</div>
          </div>
          <select v-if="setting.valueType === 'BOOL'" v-model="setting.value" class="form-input h-11">
            <option value="true">启用</option>
            <option value="false">停用</option>
          </select>
          <input v-else v-model.trim="setting.value" class="form-input h-11" :type="setting.valueType === 'INT' ? 'number' : 'text'" min="0" />
        </div>
      </div>

      <EmptyState v-else class="mt-6" title="暂无系统参数" description="请确认后端 /admin/settings 接口可用。" />
    </section>
  </main>
</template>
