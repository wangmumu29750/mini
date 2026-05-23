<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<template>
  <div>
    <div class="mx-auto max-w-7xl px-4 pt-8 sm:px-8">
      <h1 class="text-3xl font-black tracking-normal text-slate-950">{{ title }}</h1>
      <p v-if="description" class="mt-2 max-w-4xl text-base font-medium text-slate-500">{{ description }}</p>
    </div>
  </div>
</template>
