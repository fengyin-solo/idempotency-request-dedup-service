# 幂等键与请求去重服务

纯 Go 标准库实现的幂等键（Idempotency Key）与请求去重服务。

## 运行

```bash
cd origin
go run ./cmd/server
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/scopes | 创建作用域 |
| GET | /api/scopes | 查询作用域列表（支持 keyword 筛选） |
| GET | /api/scopes/{id} | 获取单个作用域 |
| PUT | /api/scopes/{id} | 更新作用域 |
| DELETE | /api/scopes/{id} | 删除作用域 |
| POST | /api/idemkeys | 创建幂等键 |
| GET | /api/idemkeys | 查询幂等键列表（支持 scope_id、status、keyword 筛选） |
| GET | /api/idemkeys/{id} | 获取单个幂等键 |
| PUT | /api/idemkeys/{id} | 更新幂等键（含状态机流转） |
| DELETE | /api/idemkeys/{id} | 删除幂等键 |
| POST | /api/idemkeys/batch-delete | 批量删除幂等键 |
| POST | /api/idemkeys/expire | 批量过期已到期的 processing 幂等键 |
| POST | /api/requestrecords | 创建请求记录 |
| GET | /api/requestrecords | 查询请求记录列表（支持 idem_key_id、status、request_hash 筛选） |
| GET | /api/requestrecords/{id} | 获取单个请求记录 |
| PUT | /api/requestrecords/{id} | 更新请求记录 |
| DELETE | /api/requestrecords/{id} | 删除请求记录 |
| POST | /api/deduphits | 创建去重命中记录 |
| GET | /api/deduphits | 查询去重命中列表（支持 idem_key_id、keyword 筛选） |
| GET | /api/deduphits/{id} | 获取单个去重命中记录 |
| DELETE | /api/deduphits/{id} | 删除去重命中记录 |
| POST | /api/deduphits/batch-delete | 批量删除去重命中记录 |
| POST | /api/ttlpolicies | 创建 TTL 策略 |
| GET | /api/ttlpolicies | 查询 TTL 策略列表（支持 scope_id、auto_cleanup、keyword 筛选） |
| GET | /api/ttlpolicies/{id} | 获取单个 TTL 策略 |
| PUT | /api/ttlpolicies/{id} | 更新 TTL 策略 |
| DELETE | /api/ttlpolicies/{id} | 删除 TTL 策略 |
| POST | /api/cleanupjobs | 创建清理任务 |
| GET | /api/cleanupjobs | 查询清理任务列表（支持 policy_id、status、keyword 筛选） |
| GET | /api/cleanupjobs/{id} | 获取单个清理任务 |
| PUT | /api/cleanupjobs/{id} | 更新清理任务（含状态机流转） |
| DELETE | /api/cleanupjobs/{id} | 删除清理任务 |
| GET | /api/stats/overview | 聚合统计概览 |
