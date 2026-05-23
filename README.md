# Mini-12306 在线车票服务系统

Mini-12306 是一个课程作业项目，目标是实现简化版在线车票服务系统。项目聚焦实名注册登录、车次查询、下单支付、出票、退票、改签和后台基础数据维护，不接入真实铁路、支付、身份证认证或短信服务。

## 技术栈

- 前端：Vue 3 + Vite + TypeScript + TailwindCSS，推荐 Vue Router、Pinia、Axios
- 后端：Go + Gin + GORM + MySQL
- 第三方能力：身份证认证、在线支付、短信通知均使用本地 mock

## 仓库结构

```text
.
├── AGENTS.md
├── docs/
├── backend/
│   ├── cmd/server/
│   ├── internal/
│   │   ├── config/
│   │   ├── dto/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── router/
│   │   ├── service/
│   │   └── validator/
│   ├── migrations/
│   ├── pkg/
│   │   ├── auth/
│   │   ├── errors/
│   │   ├── logger/
│   │   ├── mock/
│   │   └── response/
│   └── tests/
├── frontend/
│   └── src/
│       ├── api/
│       ├── assets/
│       ├── components/
│       ├── composables/
│       ├── layouts/
│       ├── pages/
│       ├── router/
│       ├── stores/
│       ├── types/
│       └── utils/
└── scripts/
```

## 必读文档

- 项目计划与范围：`docs/01-project-plan.md`
- 需求与业务规则：`docs/02-requirements-and-domain-rules.md`
- 架构与目录约束：`docs/03-architecture-and-stack.md`
- 后端开发规范：`docs/04-backend-go-standards.md`
- 前端开发规范：`docs/05-frontend-vue-standards.md`
- API 契约：`docs/06-api-contract.md`
- 数据库设计约束：`docs/07-database-design.md`
- 测试与质量规范：`docs/08-testing-and-quality.md`

## 开发约定

- 订单、余票、支付、退票、改签相关写操作必须放在后端 service 层，并考虑事务、幂等和状态流转。
- 金额使用整数分 `amount_cents` 存储和传输。
- 时间存储必须有明确时区语义，前端展示按 `Asia/Shanghai`。
- 如果修改接口、数据结构或业务规则，需要同步更新 `docs/` 中对应文档。

## 本地运行

### 1. 启动 MySQL

项目提供了 MySQL 8.4 的 Docker Compose 配置：

```bash
docker compose up -d mysql
```

默认数据库连接信息：

- 数据库：`mini_12306`
- 用户名：`mini_user`
- 密码：`mini_password`
- 端口：`3306`

### 2. 启动后端

```bash
cd backend
copy .env.example .env
go run ./cmd/server
```

后端启动时会自动执行 GORM 迁移，并写入演示账号、站点、车次、经停和未来 7 天的票额数据。

健康检查：

```bash
curl http://localhost:8080/api/v1/health
```

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

访问 Vite 输出的地址，默认是 `http://localhost:5173`。开发环境下 `/api` 会代理到 `http://localhost:8080`。

## 演示账号

- 管理员：`admin / Admin123456`
- 旅客：`alice / Password123`

也可以在注册页创建新的旅客账号。注册会做本地 mock 实名校验：手机号需为 11 位数字，身份证号需为 15 位或 18 位格式，银行卡号需为 12 到 24 位数字。

## 已接入的最小闭环

- `POST /api/v1/auth/register`：旅客注册并自动登录。
- `POST /api/v1/auth/login`：账号登录。
- `GET /api/v1/auth/me`：读取当前登录用户。
- `POST /api/v1/auth/logout`：退出登录占位接口。
- `GET /api/v1/stations`：查询启用站点。
- `GET /api/v1/trains/search`：按日期、出发站、到达站查询车次、票价和余票。

当前阶段已经打通“本地 MySQL 环境 + 注册/登录/当前用户 + 车次查询”两个垂直切片。下一步建议接订单创建、模拟支付和出票。
