<script setup lang="ts">
import type { SaveStationPayload } from '@/api/admin'
import type { AdminStation } from '@/types/domain'

defineProps<{
  stations: AdminStation[]
  form: SaveStationPayload & { id?: number }
  saving: boolean
}>()

defineEmits<{
  save: []
  reset: []
  edit: [station: AdminStation]
  disable: [station: AdminStation]
}>()
</script>

<template>
  <section class="mt-6 grid gap-6 xl:grid-cols-[360px_1fr]">
    <form class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm" @submit.prevent="$emit('save')">
      <h2 class="text-lg font-black text-slate-950">{{ form.id ? '编辑站点' : '新增站点' }}</h2>
      <div class="mt-4 space-y-3">
        <input v-model.trim="form.code" class="form-input" placeholder="站点代码，如 BJN" />
        <input v-model.trim="form.name" class="form-input" placeholder="站点名称" />
        <input v-model.trim="form.city" class="form-input" placeholder="城市" />
        <select v-model="form.status" class="form-input">
          <option value="ACTIVE">启用</option>
          <option value="DISABLED">停用</option>
        </select>
        <div class="flex gap-2">
          <button class="btn-primary flex-1" type="submit" :disabled="saving">{{ saving ? '保存中...' : '保存站点' }}</button>
          <button class="btn-secondary" type="button" @click="$emit('reset')">清空</button>
        </div>
      </div>
    </form>

    <div class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
      <table class="min-w-full divide-y divide-slate-100 text-sm">
        <thead class="bg-slate-50 text-left text-xs font-black uppercase text-slate-400">
          <tr><th class="px-4 py-3">代码</th><th class="px-4 py-3">站点</th><th class="px-4 py-3">城市</th><th class="px-4 py-3">状态</th><th class="px-4 py-3">操作</th></tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-for="station in stations" :key="station.id">
            <td class="px-4 py-3 font-black">{{ station.code }}</td>
            <td class="px-4 py-3">{{ station.name }}</td>
            <td class="px-4 py-3">{{ station.city }}</td>
            <td class="px-4 py-3">{{ station.status }}</td>
            <td class="space-x-2 px-4 py-3">
              <button class="text-sm font-black text-teal-700" type="button" @click="$emit('edit', station)">编辑</button>
              <button class="text-sm font-black text-rose-600" type="button" @click="$emit('disable', station)">停用</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
