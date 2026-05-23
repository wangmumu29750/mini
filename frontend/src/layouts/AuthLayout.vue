<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'

import AuthBrandPanel from '@/components/AuthBrandPanel.vue'
</script>

<template>
  <main class="min-h-screen bg-[radial-gradient(circle_at_24%_20%,#bfe2ff_0,#eff9ff_31%,transparent_48%),radial-gradient(circle_at_72%_82%,#97e5cc_0,#dff7f3_26%,transparent_47%),linear-gradient(135deg,#edf7ff_0%,#f8fbf7_52%,#dff5f6_100%)]">
    <div class="mx-auto flex min-h-screen w-full max-w-7xl flex-col px-4 py-6 sm:px-6 lg:px-8">
      <RouterLink to="/" class="inline-flex w-fit items-center gap-2 text-sm font-semibold text-[#238fd1]">
        <span class="flex h-8 w-8 items-center justify-center rounded-md bg-[#2da9ee] text-white">M</span>
        Mini-12306
      </RouterLink>

      <section class="flex flex-1 items-center justify-center py-10">
        <div class="grid min-h-[640px] w-full max-w-6xl overflow-hidden rounded-[28px] border border-white/70 bg-white/90 shadow-[0_30px_90px_rgb(45_137_190_/_0.18)] lg:grid-cols-[1.08fr_0.92fr]">
          <AuthBrandPanel />
          <div class="flex items-center justify-center bg-white px-8 py-12 lg:px-14">
            <RouterView />
          </div>
        </div>
      </section>
    </div>
  </main>
</template>
