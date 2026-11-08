# 模块 3 当前库菜单恢复记录

## 任务状态

测试通过，待页面验收。

## 问题

后台刷新后，“出行运营”菜单下未显示“运营调度管理”。

## 根因

当前运行的后台服务所连接的 PostgreSQL 中不存在对应菜单记录。此前处理的菜单记录属于另一套数据源，未影响当前后台正在读取的库。

## 实施

在用户确认后，只新增一条当前库菜单记录：

- 父菜单：`parent_id=65`（出行运营）
- 路由：`operationDispatch`
- 路由名称：`OperationDispatch`
- 页面组件：`view/module3/operationManagement.vue`
- 标题：运营调度管理
- 显示状态：`hidden=false`
- 排序：`0`

未修改既有菜单、角色权限、接口、业务代码或数据库表结构。

## 验证

使用后台当前配置连接 PostgreSQL，插入后立即回读：

```text
created id=83
parent=65
hidden=false
path=operationDispatch
name=OperationDispatch
component=view/module3/operationManagement.vue
title=运营调度管理
```

## 页面验收

1. 在后台浏览器执行 `Ctrl + F5` 强制刷新。
2. 展开“出行运营”。
3. 确认“运营调度管理”显示在子菜单中。
4. 打开后确认“运营区域”和“站点管理”两个标签页可见。

## 权限补齐

动态菜单接口会按 `sys_authority_menus` 过滤菜单。仅新增菜单记录不会使其出现在已登录用户的侧边栏。

经用户确认后，为菜单 `83` 新增了与“出行运营”父菜单一致的两个角色关联：

```text
menu=83 authority=888
menu=83 authority=2222
```

两条关联均已在写入后立即回读验证。未删除或覆盖任何已有角色菜单权限。
