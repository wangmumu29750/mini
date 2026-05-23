<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<template>
  <div class="border-b border-slate-200 bg-white">
    <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
      <h1 class="text-2xl font-bold tracking-normal text-slate-950">{{ title }}</h1>
      <p v-if="description" class="mt-2 max-w-3xl text-sm text-slate-600">{{ description }}</p>
    </div>
  </div>
</template>

