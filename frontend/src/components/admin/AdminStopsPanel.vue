<script setup lang="ts">
import type { SaveTrainStopPayload } from '@/api/admin'
import EmptyState from '@/components/EmptyState.vue'
import type { AdminStation, AdminTrain } from '@/types/domain'

defineProps<{
  trains: AdminTrain[]
  stations: AdminStation[]
  stopDrafts: SaveTrainStopPayload[]
  selectedTrainId: number | null
  saving: boolean
}>()

defineEmits<{
  'update:selectedTrainId': [value: number | null]
  load: []
  add: []
  remove: [index: number]
  save: []
}>()
</script>

<template>
  <section class="mt-6 rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <div class="flex items-center gap-3">
        <select :value="selectedTrainId" class="form-input w-44" @change="$emit('update:selectedTrainId', Number(($event.target as HTMLSelectElement).value)); $emit('load')">
          <option v-for="train in trains" :key="train.id" :value="train.id">{{ train.trainNo }}</option>
        </select>
        <span class="text-sm font-bold text-slate-400">完整经停列表参与查询和出票时间计算</span>
      </div>
      <div class="flex gap-2">
        <button class="btn-secondary" type="button" @click="$emit('add')">新增经停</button>
        <button class="btn-primary" type="button" :disabled="saving || !selectedTrainId" @click="$emit('save')">保存经停</button>
      </div>
    </div>

    <div v-if="stopDrafts.length" class="space-y-3">
      <div v-for="(stop, index) in stopDrafts" :key="index" class="grid gap-3 rounded-lg bg-slate-50 p-3 md:grid-cols-[80px_1.2fr_1fr_1fr_80px_90px_60px]">
        <input v-model.number="stop.stopOrder" class="form-input" type="number" min="1" />
        <select v-model.number="stop.stationId" class="form-input">
          <option v-for="station in stations" :key="station.id" :value="station.id">{{ station.name }}</option>
        </select>
        <input v-model.trim="stop.arriveClock" class="form-input" placeholder="到达 HH:mm:ss" />
        <input v-model.trim="stop.departClock" class="form-input" placeholder="发车 HH:mm:ss" />
        <input v-model.number="stop.dayOffset" class="form-input" type="number" min="0" />
        <input v-model.number="stop.mileage" class="form-input" type="number" min="0" />
        <button class="text-sm font-black text-rose-600" type="button" @click="$emit('remove', index)">移除</button>
      </div>
    </div>
    <EmptyState v-else title="暂无经停" description="选择车次后新增至少两个经停站。" />
  </section>
</template>
