# 后端 Go 开发规范

## 1. 基本要求

- 使用 Go modules 管理依赖。
- HTTP 框架固定使用 Gin。
- ORM 固定使用 GORM。
- 数据库固定使用 MySQL。
- 所有代码必须通过 `gofmt`。
- 包名使用小写单词，避免下划线和复数堆叠。

## 2. Handler 规范

- Handler 只负责解析请求、调用 service、返回响应。
- 请求体使用 DTO，不直接绑定到 GORM model。
- 使用 Gin binding 和自定义 validator 做基础参数校验。
- 不在 handler 中开启事务，不在 handler 中拼 SQL。
- 从上下文读取当前用户：如 `auth.CurrentUser(c)`。

## 3. Service 规范

- 核心业务规则必须写在 service 层。
- 订单、支付、退票、改签必须显式使用事务。
- 每个状态流转都必须校验当前状态，禁止无条件覆盖状态。
- 金额计算使用整数分。
- 对重复提交敏感的接口要支持幂等键或唯一业务流水号。

## 4. Repository 规范

- Repository 负责数据访问，不负责业务状态决策。
- 查询方法命名表达意图，如 `FindAvailableTrains`、`LockInventoryForUpdate`。
- 复杂查询可以使用 GORM chain 或原生 SQL，但必须封装在 repository 中。
- 涉及库存扣减时必须使用事务对象 `tx *gorm.DB`。

## 5. Model 规范

- 公共字段建议：

```go
type BaseModel struct {
    ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

- 数据库列名使用 `snake_case`，JSON 字段使用 `camelCase`。
- 枚举字段使用字符串，提高可读性和调试体验。
- 密码摘要、内部备注、删除字段等敏感字段不得直接 JSON 输出。

## 6. 错误处理

- 使用统一业务错误类型，包含错误码、HTTP 状态码、用户可读消息。
- Handler 最终只调用统一响应函数。
- 日志中可以记录内部错误，响应中不能暴露 SQL、堆栈或密钥。
- 常见错误码建议：
  - `VALIDATION_ERROR`
  - `UNAUTHORIZED`
  - `FORBIDDEN`
  - `NOT_FOUND`
  - `CONFLICT`
  - `INSUFFICIENT_INVENTORY`
  - `INVALID_ORDER_STATE`
  - `PAYMENT_FAILED`
  - `TICKET_NOT_CHANGEABLE`
  - `TICKET_NOT_REFUNDABLE`

## 7. 响应规范

成功：

```json
{
  "code": "OK",
  "message": "success",
  "data": {}
}
```

失败：

```json
{
  "code": "VALIDATION_ERROR",
  "message": "出发站不能为空",
  "data": null
}
```

## 8. 安全规范

- 密码使用 bcrypt。
- JWT secret 从环境变量读取。
- 登录失败不提示账号是否存在。
- 管理接口必须校验角色。
- 用户 ID 以 token 为准，不能完全相信前端传入的 userId。
- 开发环境跨域由 Vite proxy 处理，后端不启用 CORS 中间件。

## 9. GORM 与 MySQL 规范

- 后端启动时必须读取并校验 `MYSQL_DSN`，缺失或连接失败时直接退出，不进入无数据库运行模式。
- 查询用户、订单、车票时必须带归属条件，避免越权。
- 库存扣减推荐使用条件更新：

```sql
UPDATE inventories
SET available_count = available_count - 1, locked_count = locked_count + 1
WHERE id = ? AND available_count > 0;
```

- 支付、退票、改签流水号加唯一索引。
- 对软删除表，唯一索引要考虑删除字段或避免软删除造成唯一冲突。

## 10. 测试规范

- Service 层状态机和金额计算必须有单元测试。
- Repository 可使用测试数据库或事务回滚测试。
- Handler 至少覆盖认证失败、参数错误、成功响应。
- 库存扣减、支付幂等、退票幂等、改签事务是重点测试对象。
