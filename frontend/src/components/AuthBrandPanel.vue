<script setup lang="ts">
const features = [
  {
    title: '实名资料',
    description: '身份证、手机号、银行卡模拟校验',
  },
  {
    title: '票务闭环',
    description: '下单、支付、出票、库存联动',
  },
  {
    title: '后台维护',
    description: '站点、车次、票额、系统参数',
  },
]
</script>

<template>
  <div class="h-full rounded-l-[28px] bg-gradient-to-br from-[#d7eeff] via-[#e5f7fb] to-[#ddf7f2] px-12 py-16 text-[#164c72] lg:px-16">
    <p class="inline-flex rounded-full bg-[#c3e2f6] px-4 py-2 text-xs font-semibold uppercase tracking-[0.14em] text-[#2d91cf]">
      Mini-12306
    </p>
    <h1 class="mt-8 text-4xl font-bold tracking-normal text-[#154a70]">
      实名登录后完成查询、购票、退票和改签流程
    </h1>
    <p class="mt-6 max-w-xl text-base font-medium leading-8 text-[#5f7f99]">
      以课程级仿真方式串联账号认证、车次查询、订单支付和票务处理，帮助完成 Mini-12306 核心业务演示。
    </p>

    <div class="mt-10 grid gap-5">
      <div v-for="feature in features" :key="feature.title" class="flex items-start gap-3 text-[#2f6d95]">
        <span class="mt-2 h-2.5 w-2.5 rounded-full bg-[#2cc3cb]" />
        <div>
          <p class="text-base font-semibold">{{ feature.title }}</p>
          <p class="mt-1 text-sm leading-6 text-[#6689a1]">{{ feature.description }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
