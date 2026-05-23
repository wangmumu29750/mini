<template>
  <main class="min-h-screen bg-rail-surface">
    <div class="mx-auto flex min-h-screen w-full max-w-6xl flex-col px-4 py-6 sm:px-6 lg:px-8">
      <RouterLink to="/" class="inline-flex w-fit items-center gap-2 text-sm font-semibold text-teal-800">
        <span class="flex h-8 w-8 items-center justify-center rounded-md bg-teal-700 text-white">M</span>
        Mini-12306
      </RouterLink>

      <section class="grid flex-1 items-center gap-8 py-10 lg:grid-cols-[1fr_420px]">
        <div class="max-w-2xl">
          <p class="text-sm font-semibold text-orange-700">课程演示系统</p>
          <h1 class="mt-3 text-3xl font-bold tracking-normal text-slate-950 sm:text-4xl">
            实名登录后完成查询、购票、退票和改签流程
          </h1>
          <div class="mt-8 grid gap-3 sm:grid-cols-3">
            <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-subtle">
              <p class="text-sm font-semibold text-slate-900">实名资料</p>
              <p class="mt-1 text-sm text-slate-500">身份证、手机号、银行卡模拟校验</p>
            </div>
            <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-subtle">
              <p class="text-sm font-semibold text-slate-900">票务闭环</p>
              <p class="mt-1 text-sm text-slate-500">下单、支付、出票、库存联动</p>
            </div>
            <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-subtle">
              <p class="text-sm font-semibold text-slate-900">后台维护</p>
              <p class="mt-1 text-sm text-slate-500">站点、车次、票额、系统参数</p>
            </div>
          </div>
        </div>

        <RouterView />
      </section>
    </div>
  </main>
</template>

