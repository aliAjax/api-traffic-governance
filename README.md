# API Traffic Governance Control Plane

纯 Go 的多租户 API 流量治理控制面，覆盖 Service、Upstream、Route、限流配额、熔断组件、故障实验、请求镜像脱敏、xDS 风格快照和 explain 决策接口。它不是功能发布系统或 RBAC 后台。

```bash
go test ./...
go vet ./...
TRAFFIC_ADDR=:8085 go run ./cmd/control-api
curl http://127.0.0.1:8085/healthz
```

本地模式使用原子 JSON 仓储；生产环境应替换事务数据库、分布式限流存储和 mTLS xDS transport。任何未实现的真实 sidecar 数据面能力不能伪装成成功响应。
