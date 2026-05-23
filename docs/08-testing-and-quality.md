# 测试与质量规范

## 1. 最低质量门槛

- 后端代码必须通过 `gofmt`。
- 前端代码必须通过项目配置的 lint 和类型检查。
- 核心业务流程必须可通过固定演示数据完整跑通。
- 任何影响接口、数据库或业务状态机的改动都要更新文档。

## 2. 后端测试重点

- 注册：用户名、手机号、身份证唯一性；密码加密。
- 登录：正确登录、错误密码、禁用账号。
- 查询车次：站序合法、日期合法、余票和票价返回正确。
- 创建订单：余票不足、重复提交、库存锁定。
- 支付出票：支付成功、支付失败、重复支付幂等。
- 取消订单：只允许待支付订单取消，库存释放。
- 退票：状态校验、重复退票、库存释放、退款记录。
- 改签：新库存不足、差价、原票关闭、新票生成、事务回滚。
- 权限：普通用户不能访问管理接口，不能访问他人订单。

## 3. 前端测试重点

- 登录态：未登录跳转、登录后恢复用户信息、退出清理状态。
- 查询表单：必填项、日期选择、出发到达站不能相同。
- 订单流程：创建订单后跳转支付，支付成功后显示车票。
- 危险操作：取消订单、退票、改签必须二次确认。
- 错误处理：后端错误消息能正确展示，按钮恢复可点击。
- 响应式：核心旅客流程在常见移动端宽度可用。

## 4. 推荐命令

后端：

```bash
go fmt ./...
go test ./...
go run ./cmd/server
```

前端：

```bash
npm install
npm run dev
npm run lint
npm run typecheck
npm run build
```

具体命令以实际 `package.json` 和 `Makefile` 为准。

## 5. 演示数据

建议准备：

- 管理员账号：`admin / Admin123456`
- 旅客账号：`alice / Password123`
- 至少 6 个站点：北京南、天津南、济南西、南京南、上海虹桥、杭州东。
- 至少 3 个车次：覆盖不同时间、不同席别、不同余票。
- 至少 1 个可退票订单、1 个可改签订单、1 个已完成订单。

演示数据不得依赖真实身份证号或真实银行卡号。

## 6. 提交规范

提交信息建议使用：

- `feat: add train search api`
- `fix: prevent duplicate ticket refund`
- `docs: update api contract`
- `test: cover order payment idempotency`
- `chore: configure lint`

## 7. 文档同步规则

- 新增接口：更新 `docs/06-api-contract.md`。
- 新增表或字段：更新 `docs/07-database-design.md`。
- 改变业务状态：更新 `docs/02-requirements-and-domain-rules.md`。
- 改变目录或技术选型：更新 `docs/03-architecture-and-stack.md`。
- 改变开发命令：更新本文件和 README。
