# Mini-12306 Project Rules

本仓库是课程作业项目 **Mini-12306 在线车票服务系统**。后续任何会话在回答、设计或编码前，都必须先按需读取本文件和 `docs/` 下的约束文档。

## 固定技术栈

- 前端：Vue 3 + Vite + TypeScript + TailwindCSS，推荐 Vue Router、Pinia、Axios。
- 后端：Go + Gin + GORM + MySQL。
- 第三方能力只做课程级模拟：身份证认证、在线支付、短信/移动服务，不接真实生产服务。

## 必读文档

- 项目计划与范围：`docs/01-project-plan.md`
- 需求与业务规则：`docs/02-requirements-and-domain-rules.md`
- 架构与目录约束：`docs/03-architecture-and-stack.md`
- 后端开发规范：`docs/04-backend-go-standards.md`
- 前端开发规范：`docs/05-frontend-vue-standards.md`
- API 契约：`docs/06-api-contract.md`
- 数据库设计约束：`docs/07-database-design.md`
- 测试与质量规范：`docs/08-testing-and-quality.md`

## 工作原则

- 需求来源以用户提供的《Mini-12306软件系统-案例描述.pdf》和 `docs/` 文档为准，不随意扩大为完整 12306。
- 任何购票、退票、改签、支付、库存相关逻辑必须经过后端服务层，不能只在前端判断。
- 涉及订单、余票、支付、退票、改签的写操作必须考虑事务、幂等和状态流转。
- 金额一律使用整数分 `amount_cents`，时间一律存储明确时区语义，展示按 `Asia/Shanghai`。
- 代码修改如果改变接口、数据结构或业务规则，必须同步更新对应 `docs/` 文档。
- 开发时优先补齐小而可运行的垂直切片：后端模型/接口/测试 + 前端页面/API 调用 + 演示数据。

## 回答要求

- 用户问项目相关问题时，先基于这些规则给出结论，再说明必要假设。
- 用户要求实现时，直接在现有约束内动手，不只停留在方案。
- 用户要求“评审”时，按代码评审方式输出问题、风险和缺失测试。
