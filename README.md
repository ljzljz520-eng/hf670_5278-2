# 合同用印流程页

这是一个单一 Go module 的合同用印流程后端。服务提供合同名称、部门、用印类型和紧急程度的申请接口，并提供法务处理批量文件的接口；根路径直接托管可用的合同用印页面。业务状态与文件资源全部使用确定性的内存 fixture。

## 运行

```bash
CGO_ENABLED=0 go run ./cmd/contract-seal
```

打开 <http://localhost:8080/>。创建申请：

```bash
curl -X POST http://localhost:8080/api/contract-seals \
  -H 'Content-Type: application/json' \
  -d '{"contractName":"年度服务合同","department":"市场部","sealType":"official","urgency":"urgent"}'
```

法务批量处理文件：

```bash
curl -X POST http://localhost:8080/api/contract-seals/seal-001/process \
  -H 'Content-Type: application/json' \
  -d '{"files":[{"name":"contract.docx"},{"name":"appendix.pdf"},{"name":"seal-list.xlsx"}]}'
```

## 构建

```bash
CGO_ENABLED=0 go build ./...
```

## 业务链路测试

```bash
CGO_ENABLED=0 go test -count=1 ./...
```

项目保留了固定验收任务对应的业务链路测试。该验收任务会暴露当前实现中的文件关闭记录问题。
