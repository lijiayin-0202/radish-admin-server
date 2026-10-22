# radish-admin-server
萝卜快跑后台管理后端服务

## 技术栈
- Go + Gin + GORM(GVA框架)
- PostgreSQL
- Redis（可选）

## 功能
网约车后台管理，司机/乘客/订单管理底座

## 开发规范
本仓库开发约束请看根目录 AGENT.md

## 启动
1. 复制 config.example.yaml 为 config.yaml
2. 修改数据库连接信息
3. go run main.go
