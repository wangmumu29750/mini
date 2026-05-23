<template>
  <main>
    <PageHeader title="无权限访问" description="当前账号没有访问该页面的角色权限。" />
    <section class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
      <RouterLink class="btn-primary" to="/">返回车次查询</RouterLink>
    </section>
  </main>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

import PageHeader from '@/components/PageHeader.vue'
</script>

