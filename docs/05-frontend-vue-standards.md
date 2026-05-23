# 前端 Vue 3 开发规范

## 1. 基本要求

- 使用 Vue 3 + Vite + TypeScript。
- 样式使用 TailwindCSS，避免大段手写全局 CSS。
- 推荐使用 Vue Router、Pinia、Axios。
- 组件使用 Composition API 和 `<script setup lang="ts">`。
- 页面以实际业务操作为第一屏，不做营销式 landing page。

## 2. 页面范围

旅客端页面：

- 登录页
- 注册页
- 车次查询页
- 车次详情/选择席别页
- 订单确认页
- 支付模拟页
- 订单列表与详情页
- 我的车票页
- 退票确认页
- 改签查询与确认页

管理端页面：

- 管理后台布局
- 站点管理
- 车次管理
- 经停管理
- 票额管理
- 系统设置

## 3. 目录规范

```text
src/
  api/
  components/
  composables/
  layouts/
  pages/
  router/
  stores/
  types/
  utils/
```

- `api` 中只封装请求函数和 Axios 实例。
- `pages` 中组织业务流程。
- `components` 只放可复用 UI，不直接依赖路由跳转。
- `stores` 管理用户登录态、常用站点、查询条件等共享状态。
- `types` 放后端响应 DTO 对应的 TypeScript 类型。

## 4. UI 与交互规范

- 业务工具型系统应清爽、克制、信息密度适中。
- 查询页优先展示表单、结果表格、筛选条件和状态反馈。
- 按钮要有明确状态：默认、hover、disabled、loading。
- 危险操作如退票、取消订单必须有确认对话框。
- 表单错误显示在对应字段附近，同时保留后端错误提示。
- 空状态、加载状态、错误状态必须明确。
- 移动端至少保证主要旅客流程可用。

## 5. TailwindCSS 规范

- 优先使用 Tailwind utility class。
- 重复布局抽象成组件，不复制大段 class。
- 颜色避免单一蓝/紫主题铺满页面；管理端以中性色为主，关键状态使用语义色。
- 不使用过度装饰的渐变背景、漂浮卡片或大型宣传区。
- 卡片圆角不超过 `rounded-lg`，表格和表单保持紧凑可读。

## 6. API 调用规范

- 所有请求通过统一 Axios 实例。
- 开发环境通过 Vite proxy 转发 `/api` 到后端，不依赖后端 CORS。
- token 注入、401 处理、错误消息转换在拦截器中完成。
- 页面不直接拼接后端基础 URL。
- 对写操作按钮做 loading 和防重复点击。
- 支付、退票、改签请求需要携带幂等键或由后端返回并保存业务流水号。

## 7. 路由与权限

- 未登录用户访问订单、车票、管理页面时跳转登录。
- 普通旅客访问管理页面时展示无权限或跳转首页。
- 路由 meta 建议包含 `requiresAuth` 和 `roles`。

## 8. 类型规范

- 后端响应统一建模为：

```ts
export interface ApiResponse<T> {
  code: string
  message: string
  data: T
  traceId?: string
}
```

- 金额字段使用 `number` 表示整数分，展示时通过工具函数格式化。
- 时间字段使用字符串接收，展示时统一格式化。

## 9. 前端测试建议

- 工具函数测试：金额格式化、时间格式化、状态文案。
- 组件测试：查询表单、订单状态标签、确认弹窗。
- 端到端测试：注册登录、查车、下单支付、退票、改签。
