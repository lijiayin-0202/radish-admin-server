# 分册 4｜Git 分支与提交规范

> 归属总纲：`../1-总纲/AGENT.md`
> 强制：所有 commit、分支、PR 行为以本分册为准

---

## 1. 分支模型

| 分支 | 命名 | 用途 | 合入方式 |
|------|------|------|----------|
| 主分支 | `main` | **只承载框架/基础设施类变更** | PR + 评审 |
| 需求分支 | `feature/xxx` | 新功能开发 | PR → `main` |
| 缺陷分支 | `fix/xxx` | 已提 Bug 单的修复 | PR → `main` |
| 紧急分支 | `hotfix/xxx` | 线上紧急修复 | PR → `main`（可走快速评审） |
| 优化分支 | `opt/xxx` | 性能/重构/代码质量 | PR → `main` |
| 文档分支 | `docs/xxx` | 文档变更 | PR → `main` |
| 测试分支 | `test/xxx` | 测试代码/用例 | PR → `main` |

**硬规则**
1. **禁止直接向 `main` 推送业务代码**；业务代码一律在分支完成并经 PR 合入。
2. 一个分支只做一件事（一个需求/一个 Bug/一次优化），禁止"顺手改别的"。
3. 分支基于最新 `main` 创建：`git switch main && git pull && git switch -c feature/xxx`。
4. 分支生命周期随需求结束关闭，长期分支必须定期 `rebase main`。

---

## 2. Commit 格式

```
[类型] [版本号] 描述
```

> 本项目不启用工单编号：提交信息不写 `#数字`，也不校验编号。

### 类型

| 类型 | 用途 |
|------|------|
| feat | 新增功能需求 |
| fix | 修复 Bug |
| hotfix | 线上紧急修复 |
| opt | 项目优化（性能、代码质量、重构等） |
| docs | 文档变更 |
| test | 测试相关 |

### 必填规则

| 类型 | 版本号 | 说明 |
|------|--------|------|
| feat | 必填 | 新功能 |
| fix | 必填 | Bug 修复 |
| hotfix | 必填 | 线上紧急修复 |
| opt | 必填 | 优化 |
| docs | 可选 | 文档变更可省略 |
| test | 可选 | 测试相关可省略 |

### 示例

```
[feat] [v2.1.0] 新增用户注册接口
[fix] [v2.1.0] 修复详情接口缺少 appId 过滤
[hotfix] [v2.1.0] 修复线上用户无法登录问题
[hotfix] [v2.1.0] 修复线上支付回调超时
[opt] [v2.1.0] 优化首页列表查询性能，增加复合索引
[opt] [v2.1.0] 重构支付回调逻辑
[docs] 更新 API 接口文档
[test] [v2.1.0] 补充订单模块单元测试
```

### 规则细节

- 描述使用**中文**，简洁明了，动词开头（新增/修复/优化/删除/调整）。
- 版本号格式：`v{主版本}.{次版本}.{修订号}`，如 `v2.1.0`，**必须保留**。
- **禁止出现 `#数字`**：不写工单/Bug/需求编号，描述里也不带任何 `#xxxx`。
- **一次 commit 只包含一个类型**；混合变更必须拆分为多次 commit。
- **严禁整体提交**：禁止 `git add .` / `git add -A`；提交前先输出【变更文件清单】并逐文件确认。
- 禁止提交：密钥文件、`.env`、日志、编译产物、`vendor/`、本地临时文件。

---

## 3. 提交前自检（每次 commit 必做）

```bash
git status                     # 只看到本次相关文件
git diff --stat                # 确认无夹带
go build ./...                 # 编译通过
go vet ./...                   # 静态检查
gofmt -l ./api ./service ./model ./router
git add <明确的文件路径>        # 逐文件暂存，禁止 git add .
git commit -m "[feat] [v2.1.0] 描述"   # 不带 #编号
```

---

## 4. 冲突解决方案（标准流程）

> 原则：**以最新 `main` 为基线**，用 rebase 保持线性历史；冲突解决后必须重新构建 + 自测，禁止"解决完就推"。

### 4.1 拉取阶段冲突（本地未推送）

```bash
git fetch origin
git rebase origin/main
# 出现冲突：
#   1) git status 查看冲突文件
#   2) 逐文件打开，按业务语义保留双方正确逻辑（禁止整文件覆盖对方代码）
#   3) 冲突标记 <<<<<<< ======= >>>>>>> 必须全部清除
git add <解决后的文件>
git rebase --continue
# 若判断本次 rebase 方向错误
git rebase --abort          # 回退到 rebase 前状态，重新评估
```

### 4.2 推送被拒（远端有新提交）

```bash
git push origin feature/xxx     # 被拒：! [rejected] (fetch first)
git fetch origin
git pull --rebase origin feature/xxx   # 或 git rebase origin/feature/xxx
# 解决冲突 → add → rebase --continue
go build ./... && go test ./... -count=1
git push origin feature/xxx
```

### 4.3 已推送又需改写历史（谨慎）

```bash
# 仅限自己的功能分支，且已通知协作人
git fetch origin
git rebase -i origin/main        # squash / reword 整理提交
git push --force-with-lease origin feature/xxx   # 禁止 --force，禁止对 main 使用
```

### 4.4 冲突解决硬规则

1. **禁止对 `main` 强推**（`git push --force` / `--force-with-lease` 均禁止）。
2. 禁止"用我的版本覆盖全部"（`git checkout --ours .` / `--theirs .`），必须逐文件判断。
3. 框架类文件（`core/`、`middleware/`、`utils/`、`config/`）冲突：**停下来问**，由框架负责人裁决。
4. 冲突解决后必须：`go build ./...` 无 error + 关键接口能调通，才可 push。
5. 解决不了的冲突保留现场（`git rebase --abort` 或新建备份分支），禁止硬合并留 `<<<<<<<` 标记。
6. 提交信息中如涉及他人代码合并，描述注明"合并 main，解决 xxx 冲突"。

---

## 5. PR 规范

PR 标题与 commit 同格式：`[feat] [v2.1.0] 新增客服工单模块`（不带 `#编号`）

PR 描述模板：

```markdown
## 变更类型
feat / fix / hotfix / opt / docs / test

## 变更文件清单
- api/v1/cs/cs_ticket.go
- service/cs/cs_ticket.go
- model/cs/cs_tickets.go
- ...

## 实现说明
（方案要点、关键决策、是否使用代码生成器）

## 自测情况
- [ ] `go build ./...` 无 error
- [ ] 服务能启动，新增/变更接口 Apipost 调通（附截图）
- [ ] 前端页面浏览器打开无 Console 报错（附截图）
- [ ] 已知遗留问题 / 本次放宽项（如有，列明）

## 影响面与回滚
（影响模块、是否需要 DDL、回滚方式）
```

**红线**：PR 未通过评审禁止合入 `main`；禁止把多个需求打包成一个 PR。
