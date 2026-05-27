# 需求与业务规则

## 1. 需求来源

需求来源为《Mini-12306软件系统-案例描述.pdf》。原始需求强调：

- 旅客实名注册，登记身份证、手机号和银行卡等信息。
- 合法用户通过账号和密码登录系统。
- 根据时间、出发点、目的地查询车次。
- 通过在线支付购票。
- 在线退票。
- 通过在线支付改签车票。
- 管理员通过系统设置配置参数。
- 系统可集成移动服务、身份证认证服务、在线支付服务等第三方服务。

## 2. 核心领域对象

- User：登录账号，包含角色、密码摘要、账号状态。
- PassengerProfile：旅客实名资料，包含姓名、身份证号、手机号、银行卡号、认证状态。
- Station：车站。
- Train：车次基础信息，如 G101。
- TrainStop：车次经停站，包含站序、到达时间、出发时间、里程。
- SeatClass：席别，如二等座、一等座、商务座、硬座、硬卧。
- Inventory：某日期、某车次、某区间或席别的可售库存。
- Order：购票订单，承载支付和出票流程。
- Ticket：车票，表示旅客某次乘车凭证。
- Payment：支付记录，课程内为模拟支付。
- Refund：退款记录，课程内为模拟退款。
- ChangeRecord：改签记录，关联原票、新票和差价。
- SystemSetting：系统参数，如订单支付超时时间、退票截止时间。

## 3. 用户与权限规则

- 身份证号必须唯一绑定一个旅客实名资料。
- 手机号必须唯一绑定一个账号或旅客资料，除非明确支持换绑流程。
- 密码不得明文存储，后端必须使用 bcrypt 或同等级方案保存密码摘要。
- 普通旅客只能访问自己的订单、车票、支付和退款记录。
- 管理员可以维护基础数据和系统参数，但管理员不能直接绕过业务服务修改订单状态。
- 所有需要登录的接口必须经过 JWT 鉴权中间件。

## 4. 车次查询规则

- 查询条件至少包含乘车日期、出发站、到达站。
- 出发站和到达站必须在同一车次经停表中，且出发站站序小于到达站站序。
- 查询结果应包含车次号、出发/到达站、出发/到达时间、历时、席别、票价、余票。
- 只展示乘车日期仍可售的车次；已发车车次不允许购买。
- 管理端修改车次或经停后，应保证站序、时间和价格仍可计算。

## 5. 购票规则

- 创建订单前必须校验旅客实名资料已通过模拟认证。
- 同一用户可以购买多张票，但 MVP 可限制每个订单仅包含一名乘车人和一张票，降低复杂度。
- 下单时必须锁定或扣减库存，支付超时或取消时释放库存。
- 支付成功后订单状态变为 `PAID`，随后生成车票并将订单标记为 `TICKETED`。
- 重复支付请求必须幂等，不能重复扣款、重复出票或重复扣库存。
- 金额以整数分存储和计算，不使用浮点数。

## 6. 订单状态机

订单建议状态：

- `PENDING_PAYMENT`：待支付，库存已锁定。
- `CANCELLED`：已取消，库存已释放。
- `PAID`：已支付，等待出票或出票中。
- `TICKETED`：已出票。
- `CLOSED`：全部车票已退或已改签关闭。

允许流转：

- `PENDING_PAYMENT -> PAID`
- `PENDING_PAYMENT -> CANCELLED`
- `PAID -> TICKETED`
- `TICKETED -> CLOSED`

禁止跳过支付直接出票，禁止已取消订单再次支付。

## 7. 车票状态机

车票建议状态：

- `ISSUED`：已出票，可乘车。
- `REFUNDED`：已退票。
- `CHANGED_OUT`：已改签为其他车票。
- `USED`：已使用，课程 MVP 可由管理员或定时任务模拟。

退票只允许 `ISSUED -> REFUNDED`。

改签只允许 `ISSUED -> CHANGED_OUT`，同时生成新票 `ISSUED`。

## 8. 退票规则

- 已发车或超过系统退票截止时间的车票不可退。
- 退票必须在一个数据库事务中完成：校验状态、更新车票、释放库存、创建退款记录、写审计日志。
- 重复退票请求必须返回已处理结果或明确错误，不能重复退款或重复释放库存。
- MVP 可不计算手续费；如启用手续费，必须在系统设置中配置并在退款记录中保存。

## 9. 改签规则

- 改签前必须校验原车票状态为 `ISSUED`。
- 新车次必须满足乘车日期、出发站、到达站和发车时间规则。
- 改签必须在一个事务中完成：锁定新库存、关闭原票、释放原库存、创建新票、创建改签记录。
- 新票价格高于原票时，需要模拟补差价支付。
- 新票价格低于原票时，需要创建模拟退款记录。
- 改签请求必须幂等，避免重复生成新票。

## 10. 第三方服务模拟规则

- 身份证认证：提供本地模拟校验，至少校验姓名、身份证号格式和唯一性，不调用真实接口。
- 在线支付：生成模拟支付流水号，支持成功、失败、超时三类结果。
- 短信/移动服务：记录通知日志即可，不发送真实短信。
- 第三方模拟服务必须封装在后端服务层或 `pkg/mock`，前端不能伪造关键业务结果。

## 11. 时间与金额规则

- 后端保存时间使用 `time.Time`，接口输出 ISO 8601 字符串。
- 所有业务展示时间按 `Asia/Shanghai`。
- 不使用浮点数保存票价、退款、支付金额；统一使用 `amount_cents`。
- 数据库中金额字段使用 `BIGINT`。

## 12. 审计规则

以下操作应记录审计日志：

- 用户注册、登录失败次数过多、修改实名资料。
- 管理员增删改站点、车次、经停、票额、系统设置。
- 创建订单、支付、取消订单、出票、退票、改签。

审计日志至少包含操作人、操作类型、业务对象类型、业务对象 ID、结果、时间、备注。
## 13. Current MVP Implementation Notes

- Order creation currently supports one passenger and one ticket per order.
- `POST /orders` checks the logged-in user's verified profile, locks one inventory item, and creates a `PENDING_PAYMENT` order in one database transaction.
- Repeated create requests with the same `(user_id, idempotency_key)` return the existing order.
- `POST /orders/{orderId}/payments` is a local mock payment. It changes the order to `PAID`, moves one ticket from `locked_count` to `sold_count`, creates a payment record, and creates one `ISSUED` ticket in one transaction.
- If a pending order is already past `pay_expires_at` when payment is attempted, the service releases the locked inventory and changes the order to `CLOSED`.
- `POST /orders/{orderId}/cancel` cancels only `PENDING_PAYMENT` orders, releases the locked inventory, and changes the order to `CANCELLED`.
- The current code treats `PAID` as "paid and ticket issued" for the first runnable vertical slice. Later slices may split this into `PAID -> TICKETED` if asynchronous ticketing is needed.
- The ticket slice now supports transactional refund and change operations.
- `POST /tickets/{ticketId}/refund` only accepts `ISSUED` tickets before departure. It marks the ticket `REFUNDED`, releases one sold inventory back to available inventory, creates a successful mock refund record, and closes the related order in one database transaction.
- `POST /tickets/{ticketId}/change` only accepts `ISSUED` tickets before departure. It marks the old ticket `CHANGED_OUT`, releases old sold inventory, consumes new available inventory, creates a new `ISSUED` ticket, updates the related order snapshot, and writes a successful change record with `price_diff_cents` in one database transaction.
- Refund fees and real difference-payment collection are not enabled in this coursework slice; difference amounts are recorded for display/audit as simulated change settlement.
- User roles now include `PASSENGER`, `CLERK`, and `ADMIN`.
- `CLERK` users can create ticket orders for walk-up passengers through the backend service layer. The clerk flow must still verify passenger identity, lock inventory, use an idempotency key, and leave payment/order state transitions to the existing order service.
- Administrators can maintain system parameters through `SystemSetting`. The current implemented settings are `order_pay_expire_minutes`, `refund_cutoff_minutes`, `change_cutoff_minutes`, `refund_fee_percent`, and `mock_payment_enabled`.
