<template>
  <main>
    <PageHeader title="管理后台" description="维护站点、车次、经停、票额和系统参数。" />
    <section class="mx-auto grid max-w-7xl gap-4 px-4 py-6 sm:grid-cols-2 sm:px-6 lg:grid-cols-4 lg:px-8">
      <article v-for="item in items" :key="item.title" class="rounded-lg border border-slate-200 bg-white p-4 shadow-subtle">
        <h2 class="text-base font-semibold text-slate-950">{{ item.title }}</h2>
        <p class="mt-2 text-sm text-slate-500">{{ item.description }}</p>
      </article>
    </section>
  </main>
</template>

<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue'

const items = [
  { title: '站点管理', description: '维护车站名称和启用状态。' },
  { title: '车次管理', description: '维护车次基础信息和运行日期。' },
  { title: '经停管理', description: '维护站序、到发时间和里程。' },
  { title: '票额管理', description: '维护席别库存和票价。' },
]
</script>

