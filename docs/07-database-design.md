# 数据库设计约束

## 1. 命名与通用字段

- 表名使用复数或业务名，统一 `snake_case`。
- 字段名使用 `snake_case`。
- 主键建议使用 `BIGINT UNSIGNED AUTO_INCREMENT`。
- 通用字段：`id`、`created_at`、`updated_at`、`deleted_at`。
- 金额字段使用 `BIGINT`，单位为分。
- 枚举字段使用 `VARCHAR` 保存明确字符串。

## 2. 核心表

### users

账号表。

| 字段 | 说明 |
| --- | --- |
| id | 主键 |
| username | 登录名，唯一 |
| password_hash | 密码摘要 |
| role | `PASSENGER` / `ADMIN` |
| status | `ACTIVE` / `DISABLED` |
| last_login_at | 最近登录时间 |

### passenger_profiles

旅客实名资料。

| 字段 | 说明 |
| --- | --- |
| id | 主键 |
| user_id | 关联 users |
| real_name | 姓名 |
| id_card_no | 身份证号，唯一 |
| phone | 手机号，唯一 |
| bank_card_no | 银行卡号，可脱敏展示 |
| verified_status | `PENDING` / `VERIFIED` / `FAILED` |

### stations

车站表。

| 字段 | 说明 |
| --- | --- |
| id | 主键 |
| name | 站名，唯一 |
| city | 城市 |
| code | 站点代码，唯一 |
| status | `ACTIVE` / `DISABLED` |

### trains

车次表。

| 字段 | 说明 |
| --- | --- |
| id | 主键 |
| train_no | 车次号，唯一 |
| train_type | 如 `G`、`D`、`K` |
| status | `ACTIVE` / `DISABLED` |

### train_stops

经停表。

| 字段 | 说明 |
| --- | --- |
| id | 主键 |
| train_id | 车次 ID |
| station_id | 车站 ID |
| stop_order | 站序 |
| arrive_time | 到达时间，可按一天内时间保存 |
| depart_time | 出发时间 |
| day_offset | 跨天偏移 |
| mileage | 里程 |

约束：

- `(train_id, stop_order)` 唯一。
- `(train_id, station_id)` 唯一。

### seat_classes

席别表。

| 字段 | 说明 |
| --- | --- |
| code | 席别代码，如 `SECOND` |
| name | 席别名称 |
| sort_order | 展示顺序 |

### train_seat_prices

票价表，MVP 可按出发站、到达站、席别配置价格。

| 字段 | 说明 |
| --- | --- |
| train_id | 车次 ID |
| from_station_id | 出发站 |
| to_station_id | 到达站 |
| seat_class_code | 席别 |
| price_cents | 票价，分 |

### inventories

票额表。课程 MVP 可按车次、日期、出发站、到达站、席别维护库存，避免实现真实区间占座算法。

| 字段 | 说明 |
| --- | --- |
| train_id | 车次 ID |
| travel_date | 乘车日期 |
| from_station_id | 出发站 |
| to_station_id | 到达站 |
| seat_class_code | 席别 |
| total_count | 总票额 |
| available_count | 可售数量 |
| locked_count | 锁定数量 |
| sold_count | 已售数量 |

唯一索引：

- `(train_id, travel_date, from_station_id, to_station_id, seat_class_code)`

### orders

订单表。

| 字段 | 说明 |
| --- | --- |
| order_no | 订单号，唯一 |
| user_id | 用户 ID |
| passenger_profile_id | 旅客资料 ID |
| status | 订单状态 |
| total_amount_cents | 总金额 |
| pay_deadline_at | 支付截止时间 |
| paid_at | 支付时间 |
| idempotency_key | 创建订单幂等键 |

### tickets

车票表。

| 字段 | 说明 |
| --- | --- |
| ticket_no | 票号，唯一 |
| order_id | 订单 ID |
| user_id | 用户 ID |
| passenger_profile_id | 旅客资料 ID |
| train_id | 车次 ID |
| train_no | 车次号快照 |
| travel_date | 乘车日期 |
| from_station_id | 出发站 |
| to_station_id | 到达站 |
| seat_class_code | 席别 |
| coach_no | 车厢号快照，课程级模拟 |
| seat_no | 座位号快照，课程级模拟 |
| price_cents | 票价 |
| status | 车票状态 |
| refunded_at | 退票完成时间，可空 |

### payments

支付流水。

| 字段 | 说明 |
| --- | --- |
| payment_no | 支付流水号，唯一 |
| order_id | 订单 ID |
| user_id | 用户 ID |
| amount_cents | 支付金额 |
| status | `PENDING` / `SUCCESS` / `FAILED` |
| pay_method | 支付方式 |
| idempotency_key | 支付幂等键 |
| paid_at | 支付完成时间 |

### refunds

退款流水。

| 字段 | 说明 |
| --- | --- |
| refund_no | 退款流水号，唯一 |
| ticket_id | 车票 ID |
| payment_id | 原支付 ID |
| amount_cents | 退款金额 |
| status | `PENDING` / `SUCCESS` / `FAILED` |
| reason | 原因 |
| idempotency_key | 幂等键 |

### change_records

改签记录。

| 字段 | 说明 |
| --- | --- |
| change_no | 改签流水号，唯一 |
| old_ticket_id | 原票 ID |
| new_ticket_id | 新票 ID |
| user_id | 用户 ID |
| price_diff_cents | 差价，正数补款，负数退款 |
| status | `SUCCESS` / `FAILED` |
| idempotency_key | 幂等键 |

### system_settings

系统设置。

| 字段 | 说明 |
| --- | --- |
| setting_key | 配置键，唯一 |
| setting_value | 配置值 |
| description | 描述 |

建议配置：

- `order_pay_expire_minutes`
- `refund_deadline_minutes_before_departure`
- `change_deadline_minutes_before_departure`
- `refund_fee_rate`

### audit_logs

审计日志。

| 字段 | 说明 |
| --- | --- |
| actor_user_id | 操作人 |
| action | 操作类型 |
| resource_type | 资源类型 |
| resource_id | 资源 ID |
| result | `SUCCESS` / `FAILED` |
| detail | JSON 或文本详情 |

## 3. 关键索引

- `users.username`
- `passenger_profiles.id_card_no`
- `passenger_profiles.phone`
- `stations.name`
- `stations.code`
- `trains.train_no`
- `train_stops(train_id, stop_order)`
- `inventories(train_id, travel_date, from_station_id, to_station_id, seat_class_code)`
- `orders(user_id, created_at)`
- `orders.order_no`
- `tickets(user_id, travel_date)`

## 4. Current Implementation Notes

- `users.role` now supports `PASSENGER`, `CLERK`, and `ADMIN`.
- The `system_settings` table is implemented by the `SystemSetting` model with columns `key`, `value`, `value_type`, and `description`.
- Default system settings are seeded lazily by the settings service: `order_pay_expire_minutes`, `refund_cutoff_minutes`, `change_cutoff_minutes`, `refund_fee_percent`, and `mock_payment_enabled`.
- Clerk-created orders are stored in the same `orders` table and reuse the order/inventory transaction rules. The `user_id` is the clerk account that created the order, while passenger name and ID card snapshots are stored on the order/ticket records.
- `tickets.ticket_no`
- `payments.payment_no`
- `payments(order_id, idempotency_key)`
- `refunds.refund_no`
- `change_records.change_no`

## 4. 库存简化说明

真实铁路售票需要区间库存算法。课程 MVP 为降低实现风险，可以按完整出发站到到达站区间维护库存。这样查询、扣减、释放都更直观，适合答辩演示。

如果后续要增强为真实区间库存，应新增区间段库存表，并在购票时扣减经过的每个区间段。
## 5. Current MVP Schema Notes

Current `orders` fields in code: `order_no`, `user_id`, `train_id`, `train_no`, `travel_date`, `from_station_id`, `from_station_name`, `to_station_id`, `to_station_name`, `seat_class_code`, `seat_class_name`, `passenger_name`, `id_card_no`, `amount_cents`, `status`, `pay_expires_at`, `paid_at`, `idempotency_key`.

Current `orders` unique indexes: `order_no`; `(user_id, idempotency_key)`.

Current `tickets` fields in code: `ticket_no`, `order_id`, `user_id`, `train_id`, `train_no`, `travel_date`, `from_station_id`, `from_station_name`, `to_station_id`, `to_station_name`, `seat_class_code`, `seat_class_name`, `coach_no`, `seat_no`, `passenger_name`, `id_card_no`, `status`, `issued_at`, `refunded_at`.

Current `payments` fields in code: `payment_no`, `order_id`, `user_id`, `amount_cents`, `channel`, `status`, `paid_at`.

Current `refunds` fields in code: `refund_no`, `ticket_id`, `payment_id`, `user_id`, `amount_cents`, `status`, `reason`, `idempotency_key`, `refunded_at`.

Current `change_records` fields in code: `change_no`, `old_ticket_id`, `new_ticket_id`, `user_id`, `price_diff_cents`, `status`, `idempotency_key`, `changed_at`.
