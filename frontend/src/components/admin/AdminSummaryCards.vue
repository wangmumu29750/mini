<script setup lang="ts">
import type { AdminTrain, InventoryQuoteStats, SellableTrainStat } from '@/types/domain'
import { formatMoney } from '@/utils/format'

defineProps<{
  activeStationTotal: number
  stats: SellableTrainStat[]
  quoteStats: InventoryQuoteStats | null
  inventoryCount: number
  totalSellableTrains: number
  lowestPrice: number
  selectedTrain?: AdminTrain
}>()
</script>

<template>
  <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
    <article class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
      <div class="text-sm font-black text-slate-400">启用站点</div>
      <div class="mt-3 text-4xl font-black text-slate-950">{{ activeStationTotal }}</div>
      <div class="mt-2 text-sm font-bold text-slate-500">后台统计</div>
    </article>
    <article class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
      <div class="text-sm font-black text-slate-400">两日可售车次</div>
      <div class="mt-3 text-4xl font-black text-slate-950">{{ totalSellableTrains }}</div>
      <div class="mt-2 text-sm font-bold text-slate-500">{{ stats.map((item) => `${item.date}: ${item.trainCount}`).join(' / ') || '暂无线路' }}</div>
    </article>
    <article class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
      <div class="text-sm font-black text-slate-400">报价项</div>
      <div class="mt-3 text-4xl font-black text-slate-950">{{ quoteStats?.quoteCount || inventoryCount }}</div>
      <div class="mt-2 text-sm font-bold text-slate-500">{{ selectedTrain?.trainNo || '未选车次' }}</div>
    </article>
    <article class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
      <div class="text-sm font-black text-slate-400">最低票价</div>
      <div class="mt-3 text-4xl font-black text-rose-600">{{ lowestPrice ? formatMoney(lowestPrice) : '-' }}</div>
      <div class="mt-2 text-sm font-bold text-slate-500">流转后实时返回</div>
    </article>
  </div>
</template>
