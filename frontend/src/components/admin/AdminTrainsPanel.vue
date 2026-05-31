<script setup lang="ts">
import type { SaveTrainPayload } from '@/api/admin'
import type { AdminTrain } from '@/types/domain'

defineProps<{
  trains: AdminTrain[]
  form: SaveTrainPayload & { id?: number }
  saving: boolean
}>()

defineEmits<{
  save: []
  reset: []
  edit: [train: AdminTrain]
  stops: [train: AdminTrain]
  delete: [train: AdminTrain]
}>()
</script>

<template>
  <section class="mt-6 grid gap-6 xl:grid-cols-[360px_1fr]">
    <form class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm" @submit.prevent="$emit('save')">
      <h2 class="text-lg font-black text-slate-950">{{ form.id ? '编辑车次' : '新增车次' }}</h2>
      <div class="mt-4 space-y-3">
        <input v-model.trim="form.trainNo" class="form-input" placeholder="车次号，如 G101" />
        <select v-model="form.trainType" class="form-input">
          <option value="G">G 高铁</option>
          <option value="C">C 城际</option>
          <option value="D">D 动车</option>
          <option value="Z">Z 直达</option>
          <option value="T">T 特快</option>
          <option value="K">K 快速</option>
        </select>
        <select v-model="form.status" class="form-input">
          <option value="ACTIVE">启用</option>
          <option value="DISABLED">停用</option>
        </select>
        <div class="flex gap-2">
          <button class="btn-primary flex-1" type="submit" :disabled="saving">{{ saving ? '保存中...' : '保存车次' }}</button>
          <button class="btn-secondary" type="button" @click="$emit('reset')">清空</button>
        </div>
      </div>
    </form>

    <div class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
      <table class="min-w-full divide-y divide-slate-100 text-sm">
        <thead class="bg-slate-50 text-left text-xs font-black uppercase text-slate-400">
          <tr><th class="px-4 py-3">车次</th><th class="px-4 py-3">类型</th><th class="px-4 py-3">经停数</th><th class="px-4 py-3">状态</th><th class="px-4 py-3">操作</th></tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-for="train in trains" :key="train.id">
            <td class="px-4 py-3 font-black">{{ train.trainNo }}</td>
            <td class="px-4 py-3">{{ train.trainType }}</td>
            <td class="px-4 py-3">{{ train.stopCount }}</td>
            <td class="px-4 py-3">{{ train.status }}</td>
            <td class="space-x-2 px-4 py-3">
              <button class="text-sm font-black text-teal-700" type="button" @click="$emit('edit', train)">编辑</button>
              <button class="text-sm font-black text-slate-700" type="button" @click="$emit('stops', train)">经停</button>
              <button class="text-sm font-black text-rose-600" type="button" @click="$emit('delete', train)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
