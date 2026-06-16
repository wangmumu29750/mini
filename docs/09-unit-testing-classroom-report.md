# Mini-12306 单元测试课堂报告

## 1. 单元测试工具选择

- 测试工具：Go 标准库 `testing`
- 运行命令：`go test ./internal/service -v`
- 选择原因：项目后端使用 Go，`testing` 是 Go 官方内置单元测试框架，不需要额外安装依赖；VS Code 安装 Go 扩展后可以在测试函数上方直接点击 `run test`，也可以在测试面板查看通过/失败结果。

## 2. 被测核心类或核心方法

本次选择后端服务层票价计算方法作为核心单元：

- 文件：`backend/internal/service/ticket_pricing.go`
- 方法：`CalculateTicketPrice(basePriceCents int64, trainType, seatType, ticketType string) (int64, error)`
- 业务意义：创建订单时需要根据基础票价、车次类型、席别、票种计算实际支付金额。金额使用整数分，属于购票核心流程中的关键业务规则。

## 3. 测试用例设计

| 编号 | 测试点 | 输入 | 期望结果 |
| --- | --- | --- | --- |
| TC-01 | 成人票原价 | 基础价 10000 分，G 车，二等座，成人票 | 返回 10000 分，无错误 |
| TC-02 | 高铁/动车学生票折扣 | 基础价 10000 分，G 车，二等座，学生票 | 返回 7500 分，无错误 |
| TC-03 | 普速车学生票折扣 | 基础价 10000 分，K 车，硬座，学生票 | 返回 6000 分，无错误 |
| TC-04 | 儿童票半价和四舍五入 | 基础价 10001 分，D 车，二等座，儿童票 | 返回 5001 分，无错误 |
| TC-05 | 非法票种 | 票种为 `SENIOR` | 返回错误 |
| TC-06 | 必填参数缺失 | 席别为空 | 返回错误 |
| TC-07 | 输入规范化 | 车次类型、席别、票种带空格且为小写 | 能正常识别并返回 6666 分 |

## 4. 具体测试代码

测试代码已添加在：

- `backend/internal/service/ticket_pricing_test.go`
- `backend/internal/service/fare_rules_test.go`

核心测试代码采用表驱动方式，把多组输入和期望输出放在同一个测试表中，便于扩展更多票种和车次类型。

## 5. 测试结果及分析

执行命令：

```bash
cd backend
go test ./internal/service -v
```

当前结果：服务层单元测试通过。

结果分析：

- 成人票、高铁学生票、普速学生票、儿童票均能按照业务规则返回正确金额。
- 非法票种和缺失参数能够返回错误，说明方法具备基本输入校验能力。
- 小写和带空格输入可以被规范化处理，提升了接口对前端输入的容错性。
- 本次测试属于纯业务单元测试，不依赖 MySQL，因此适合课堂演示和在 VS Code 测试面板中快速运行。

## 6. 在 VS Code 中查看测试结果

1. 打开项目目录 `D:\mini`。
2. 安装 VS Code 的 Go 扩展。
3. 打开 `backend/internal/service/ticket_pricing_test.go`。
4. 点击测试函数上方的 `run test`，或打开左侧 Testing 面板运行测试。
5. 也可以在 VS Code 终端执行：

```bash
cd backend
go test ./internal/service -v
```
