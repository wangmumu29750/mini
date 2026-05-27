# API 契约草案

## 1. 通用约定

- 基础路径：`/api/v1`
- Content-Type：`application/json`
- 鉴权头：`Authorization: Bearer <token>`
- 时间格式：ISO 8601 字符串
- 金额单位：整数分

统一响应：

```json
{
  "code": "OK",
  "message": "success",
  "data": {},
  "traceId": "req_202605230001"
}
```

分页响应：

```json
{
  "items": [],
  "page": 1,
  "pageSize": 10,
  "total": 0
}
```

## 2. 系统接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/health` | 健康检查，返回服务状态、运行模式、数据库状态和服务器时间 |

响应数据示例：

```json
{
  "status": "ok",
  "env": "debug",
  "database": "ok",
  "time": "2026-05-23T17:20:00+08:00"
}
```

## 3. 认证接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/auth/register` | 旅客实名注册 |
| POST | `/auth/login` | 登录 |
| GET | `/auth/me` | 当前用户信息 |
| POST | `/auth/logout` | 退出登录，MVP 可仅前端清 token |

注册请求示例：

```json
{
  "username": "alice",
  "password": "Password123",
  "realName": "张三",
  "idCardNo": "110101199001011234",
  "phone": "13800138000",
  "bankCardNo": "6222020202020202020"
}
```

登录响应数据示例：

```json
{
  "accessToken": "jwt-token",
  "user": {
    "id": 1,
    "username": "alice",
    "role": "PASSENGER"
  }
}
```

## 4. 车站与车次查询

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/stations` | 车站列表 |
| GET | `/trains/search` | 按日期、出发站、到达站查询车次 |
| GET | `/trains/{trainId}` | 车次详情 |

查询参数：

- `date=2026-06-01`
- `fromStationId=1`
- `toStationId=5`

查询响应单项：

```json
{
  "trainId": 1,
  "trainNo": "G101",
  "travelDate": "2026-06-01",
  "fromStation": {"id": 1, "name": "北京南"},
  "toStation": {"id": 5, "name": "上海虹桥"},
  "departTime": "2026-06-01T08:00:00+08:00",
  "arriveTime": "2026-06-01T13:30:00+08:00",
  "durationMinutes": 330,
  "seatOptions": [
    {"seatClassCode": "SECOND", "seatClassName": "二等座", "priceCents": 55300, "availableCount": 20}
  ]
}
```

## 5. 订单与支付

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/orders` | 创建购票订单并锁定库存 |
| GET | `/orders` | 我的订单列表 |
| GET | `/orders/{orderId}` | 订单详情 |
| POST | `/orders/{orderId}/cancel` | 取消待支付订单 |
| POST | `/orders/{orderId}/payments` | 模拟支付 |

创建订单请求：

```json
{
  "trainId": 1,
  "travelDate": "2026-06-01",
  "fromStationId": 1,
  "toStationId": 5,
  "seatClassCode": "SECOND",
  "idempotencyKey": "uuid-from-client"
}
```

支付请求：

```json
{
  "payMethod": "MOCK_BANK_CARD",
  "idempotencyKey": "uuid-from-client"
}
```

## 6. 车票、退票、改签

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/tickets` | 我的车票列表 |
| GET | `/tickets/{ticketId}` | 车票详情 |
| POST | `/tickets/{ticketId}/refund` | 退票 |
| GET | `/tickets/{ticketId}/change-options` | 查询可改签车次 |
| POST | `/tickets/{ticketId}/change` | 改签 |

退票请求：

```json
{
  "reason": "行程变更",
  "idempotencyKey": "uuid-from-client"
}
```

改签请求：

```json
{
  "newTrainId": 2,
  "newTravelDate": "2026-06-01",
  "newSeatClassCode": "FIRST",
  "idempotencyKey": "uuid-from-client"
}
```

## 7. 管理接口

管理接口均要求 `ADMIN` 角色。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/admin/stations` | 站点列表 |
| POST | `/admin/stations` | 创建站点 |
| PUT | `/admin/stations/{stationId}` | 修改站点 |
| DELETE | `/admin/stations/{stationId}` | 删除站点 |
| GET | `/admin/trains` | 车次列表 |
| POST | `/admin/trains` | 创建车次 |
| PUT | `/admin/trains/{trainId}` | 修改车次 |
| DELETE | `/admin/trains/{trainId}` | 删除车次 |
| GET | `/admin/trains/{trainId}/stops` | 经停列表 |
| PUT | `/admin/trains/{trainId}/stops` | 覆盖保存经停 |
| GET | `/admin/inventories` | 票额列表 |
| PUT | `/admin/inventories/{inventoryId}` | 修改票额 |
| GET | `/admin/settings` | 系统设置 |
| PUT | `/admin/settings` | 更新系统设置 |

## 8. 错误码

| code | HTTP | 说明 |
| --- | --- | --- |
| `OK` | 200 | 成功 |
| `VALIDATION_ERROR` | 400 | 参数错误 |
| `UNAUTHORIZED` | 401 | 未登录或 token 无效 |
| `FORBIDDEN` | 403 | 无权限 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `CONFLICT` | 409 | 唯一性或状态冲突 |
| `INSUFFICIENT_INVENTORY` | 409 | 余票不足 |
| `INVALID_ORDER_STATE` | 409 | 订单状态不允许 |
| `PAYMENT_FAILED` | 402 | 模拟支付失败 |
| `TICKET_NOT_REFUNDABLE` | 409 | 车票不可退 |
| `TICKET_NOT_CHANGEABLE` | 409 | 车票不可改签 |
| `INTERNAL_ERROR` | 500 | 服务端错误 |
## 9. Current Implemented Order API Notes

Implemented authenticated endpoints:

| Method | Path | Description |
| --- | --- | --- |
| POST | `/orders` | Create a one-ticket order and lock inventory |
| GET | `/orders` | List current user's orders |
| GET | `/orders/{orderId}` | Get current user's order detail |
| POST | `/orders/{orderId}/cancel` | Cancel pending-payment order and release locked inventory |
| POST | `/orders/{orderId}/payments` | Mock payment and ticket issuing |

Create order request:

```json
{
  "trainId": 1,
  "travelDate": "2026-05-24",
  "fromStationId": 1,
  "toStationId": 5,
  "seatClassCode": "SECOND",
  "idempotencyKey": "uuid-from-client"
}
```

Order response:

```json
{
  "id": 1,
  "orderNo": "O20260523153000123456",
  "trainId": 1,
  "trainNo": "G101",
  "travelDate": "2026-05-24",
  "fromStation": {"id": 1, "name": "北京南"},
  "toStation": {"id": 5, "name": "上海虹桥"},
  "seatClassCode": "SECOND",
  "seatClassName": "二等座",
  "passengerName": "张三",
  "amountCents": 55300,
  "status": "PENDING_PAYMENT",
  "payExpiresAt": "2026-05-23T15:45:00+08:00"
}
```

Payment request:

```json
{
  "channel": "MOCK_BANK"
}
```

Payment response data:

```json
{
  "paymentNo": "P20260523153100123456",
  "order": {
    "id": 1,
    "orderNo": "O20260523153000123456",
    "status": "PAID",
    "ticketNo": "T20260523153100123456",
    "ticketStatus": "ISSUED"
  }
}
```

## 10. Current Implemented Ticket API Notes

Implemented authenticated endpoints:

| Method | Path | Description |
| --- | --- | --- |
| GET | `/tickets` | List current user's issued tickets |
| GET | `/tickets/{ticketId}` | Get current user's ticket detail |
| POST | `/tickets/{ticketId}/refund` | Refund an issued ticket before departure |
| GET | `/tickets/{ticketId}/change-options` | Query changeable train and seat options for an issued ticket |
| POST | `/tickets/{ticketId}/change` | Change an issued ticket to another train/date/seat class |

Ticket response:

```json
{
  "id": 1,
  "ticketNo": "T20260523153100123456",
  "orderId": 1,
  "trainId": 1,
  "trainNo": "G101",
  "travelDate": "2026-05-24",
  "fromStation": {"id": 1, "name": "北京南"},
  "toStation": {"id": 5, "name": "上海虹桥"},
  "departTime": "2026-05-24T08:00:00+08:00",
  "arriveTime": "2026-05-24T13:30:00+08:00",
  "seatClassCode": "SECOND",
  "seatClassName": "二等座",
  "coachNo": "04",
  "seatNo": "08A",
  "passengerName": "张三",
  "idCardNoMasked": "1101**********1234",
  "status": "ISSUED",
  "issuedAt": "2026-05-23T15:31:00+08:00"
}
```

## 11. Current Implemented Clerk and Settings API Notes

Implemented clerk/admin endpoints:

| Method | Path | Role | Description |
| --- | --- | --- | --- |
| POST | `/clerk/orders` | `CLERK` / `ADMIN` | Create an order for a walk-up passenger and lock inventory |
| GET | `/admin/settings` | `ADMIN` | List system settings |
| PUT | `/admin/settings` | `ADMIN` | Update supported system settings |

Clerk order request:

```json
{
  "trainId": 1,
  "travelDate": "2026-05-25",
  "fromStationId": 1,
  "toStationId": 5,
  "seatClassCode": "SECOND",
  "passengerName": "张三",
  "idCardNo": "110101199001011234",
  "phone": "13800138000",
  "bankCardNo": "6222020202020202020",
  "idempotencyKey": "uuid-from-client"
}
```

System setting item:

```json
{
  "key": "order_pay_expire_minutes",
  "value": "15",
  "valueType": "INT",
  "description": "Order payment timeout in minutes"
}
```

Update settings request:

```json
{
  "settings": [
    {"key": "order_pay_expire_minutes", "value": "20"},
    {"key": "mock_payment_enabled", "value": "true"}
  ]
}
```

Refund response data:

```json
{
  "refundNo": "R20260523154000123456",
  "ticket": {
    "id": 1,
    "ticketNo": "T20260523153100123456",
    "status": "REFUNDED",
    "refundedAt": "2026-05-23T15:40:00+08:00"
  }
}
```

Change options query:

- `date=2026-05-25`

Change options response data:

```json
{
  "originalTicket": {
    "id": 1,
    "ticketNo": "T20260523153100123456",
    "status": "ISSUED"
  },
  "options": [
    {
      "trainId": 2,
      "trainNo": "G102",
      "travelDate": "2026-05-25",
      "fromStation": {"id": 1, "name": "北京南"},
      "toStation": {"id": 5, "name": "上海虹桥"},
      "departTime": "2026-05-25T09:00:00+08:00",
      "arriveTime": "2026-05-25T14:30:00+08:00",
      "durationMinutes": 330,
      "seatOptions": [
        {"seatClassCode": "SECOND", "seatClassName": "二等座", "priceCents": 55300, "availableCount": 18}
      ]
    }
  ]
}
```

Change request:

```json
{
  "newTrainId": 2,
  "newTravelDate": "2026-05-25",
  "newSeatClassCode": "FIRST",
  "idempotencyKey": "uuid-from-client"
}
```

Change response data:

```json
{
  "changeNo": "C20260523154500123456",
  "priceDiffCents": 38000,
  "oldTicket": {"id": 1, "status": "CHANGED_OUT"},
  "newTicket": {"id": 2, "ticketNo": "T20260523154500123456", "status": "ISSUED"}
}
```
