<script setup lang="ts">
import { computed } from 'vue'

import type { SaveInventoryPayload } from '@/api/admin'
import type { AdminStation, AdminTrain, Inventory } from '@/types/domain'
import { formatMoney } from '@/utils/format'

const props = defineProps<{
  trains: AdminTrain[]
  stations: AdminStation[]
  inventories: Inventory[]
  form: SaveInventoryPayload
  selectedInventoryId: number | null
  flowAction: string
  flowQuantity: number
  saving: boolean
}>()

defineEmits<{
  'update:selectedInventoryId': [value: number | null]
  'update:flowAction': [value: string]
  'update:flowQuantity': [value: number]
  save: []
  flow: []
  edit: [item: Inventory]
}>()

const seatClassNames: Record<string, string> = {
  BUSINESS: '商务座',
  FIRST: '一等座',
  SECOND: '二等座',
  FIRST_SLEEPER: '一等卧',
  SECOND_SLEEPER: '二等卧',
  DELUXE_SOFT_SLEEPER: '高级软卧',
  SOFT_SLEEPER: '软卧',
  HARD_SLEEPER: '硬卧',
  HARD_SEAT: '硬座',
  NO_SEAT: '无座',
}

const selectedTrain = computed(() => props.trains.find((train) => train.id === props.form.trainId))
const seatOptions = computed(() => seatClassOptions(selectedTrain.value))

function seatClassOptions(train?: AdminTrain) {
  const codes = train?.seatClassCodes?.length ? train.seatClassCodes : ['SECOND', 'FIRST', 'BUSINESS']
  return codes.map((code) => ({ code, name: seatClassNames[code] || code }))
}
</script>

<template>
  <section class="mt-6 grid gap-6 xl:grid-cols-[420px_1fr]">
    <div class="space-y-4">
      <form class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm" @submit.prevent="$emit('save')">
        <h2 class="text-lg font-black text-slate-950">票额报价维护</h2>
        <div class="mt-4 grid gap-3 sm:grid-cols-2">
          <select v-model.number="form.trainId" class="form-input">
            <option v-for="train in trains" :key="train.id" :value="train.id">{{ train.trainNo }}</option>
          </select>
          <input v-model="form.travelDate" class="form-input" type="date" />
          <select v-model.number="form.fromStationId" class="form-input">
            <option v-for="station in stations" :key="station.id" :value="station.id">{{ station.name }}</option>
          </select>
          <select v-model.number="form.toStationId" class="form-input">
            <option v-for="station in stations" :key="station.id" :value="station.id">{{ station.name }}</option>
          </select>
          <select v-model="form.seatClassCode" class="form-input">
            <option
              v-for="seat in seatOptions"
              :key="seat.code"
              :value="seat.code"
            >
              {{ seat.name }}
            </option>
          </select>
          <input v-model.number="form.priceCents" class="form-input" type="number" min="0" placeholder="票价分，填 0 自动计算" />
          <input v-model.number="form.totalCount" class="form-input" type="number" min="0" placeholder="总票额" />
          <input v-model.number="form.availableCount" class="form-input" type="number" min="0" placeholder="可售" />
          <input v-model.number="form.lockedCount" class="form-input" type="number" min="0" placeholder="锁定" />
          <input v-model.number="form.soldCount" class="form-input" type="number" min="0" placeholder="已售" />
        </div>
        <button class="btn-primary mt-4 w-full" type="submit" :disabled="saving">保存票额</button>
      </form>

      <form class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm" @submit.prevent="$emit('flow')">
        <h2 class="text-lg font-black text-slate-950">票额流转</h2>
        <div class="mt-4 space-y-3">
          <select :value="selectedInventoryId" class="form-input" @change="$emit('update:selectedInventoryId', Number(($event.target as HTMLSelectElement).value))">
            <option v-for="item in inventories" :key="item.id" :value="item.id">{{ item.trainNo }} {{ item.travelDate }} {{ item.seatClassName }}</option>
          </select>
          <select :value="flowAction" class="form-input" @change="$emit('update:flowAction', ($event.target as HTMLSelectElement).value)">
            <option value="LOCK">购票锁定</option>
            <option value="PAY">支付出票</option>
            <option value="RELEASE">取消释放</option>
            <option value="REFUND">退票释放</option>
            <option value="CHANGE_OUT">改签转出</option>
            <option value="CHANGE_IN">改签转入</option>
          </select>
          <input :value="flowQuantity" class="form-input" type="number" min="1" @input="$emit('update:flowQuantity', Number(($event.target as HTMLInputElement).value))" />
          <button class="btn-primary w-full" type="submit" :disabled="saving || !selectedInventoryId">执行流转</button>
        </div>
      </form>
    </div>

    <div class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
      <table class="min-w-full divide-y divide-slate-100 text-sm">
        <thead class="bg-slate-50 text-left text-xs font-black uppercase text-slate-400">
          <tr><th class="px-4 py-3">车次</th><th class="px-4 py-3">日期</th><th class="px-4 py-3">区间</th><th class="px-4 py-3">席别</th><th class="px-4 py-3">库存</th><th class="px-4 py-3">操作</th></tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-for="item in inventories" :key="item.id">
            <td class="px-4 py-3 font-black">{{ item.trainNo }}</td>
            <td class="px-4 py-3">{{ item.travelDate }}</td>
            <td class="px-4 py-3">{{ item.fromStation.name }} → {{ item.toStation.name }}</td>
            <td class="px-4 py-3">{{ item.seatClassName }} {{ formatMoney(item.priceCents) }}</td>
            <td class="px-4 py-3">可售 {{ item.availableCount }} / 锁定 {{ item.lockedCount }} / 已售 {{ item.soldCount }}</td>
            <td class="px-4 py-3"><button class="text-sm font-black text-teal-700" type="button" @click="$emit('edit', item)">编辑</button></td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
