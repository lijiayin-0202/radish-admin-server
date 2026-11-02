-- =====================================================================
-- 模块6：平台系统配置、权限与全局数据统计大屏 —— 菜单 / 接口权限 / 角色授权
-- 目标库：PostgreSQL（gva 库，101.96.206.96:5432，与 GVA 菜单同库）
-- 已在本地执行，保留此文件便于其他环境复现
-- 菜单层级：平台系统配置(platform) → 11 个子菜单
-- =====================================================================

-- 1. 父级菜单
INSERT INTO sys_base_menus
(created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort,
 active_name, keep_alive, default_menu, title, icon, close_tab, transition_type)
VALUES
(now(), now(), 0, 0, 'platform', 'platform', false, 'view/routerHolder.vue', 10,
 '', false, false, '平台系统配置', 'setting', false, '')
RETURNING id;  -- 返回 53

-- 2. 子菜单（parent_id 填父级返回的 id，此处为 53）
INSERT INTO sys_base_menus
(created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort,
 active_name, keep_alive, default_menu, title, icon, close_tab, transition_type)
VALUES
(now(), now(), 1, 53, 'role',         'role',         false, 'view/platform/role/role.vue',                   1, '', true, false, '角色管理', 'user', false, ''),
(now(), now(), 1, 53, 'operator',     'operator',     false, 'view/platform/operator/operator.vue',           2, '', true, false, '账号管理', 'user', false, ''),
(now(), now(), 1, 53, 'permission',   'permission',   false, 'view/platform/permission/permission.vue',       3, '', true, false, '权限管理', 'lock', false, ''),
(now(), now(), 1, 53, 'auditLog',     'auditLog',     false, 'view/platform/auditLog/auditLog.vue',           4, '', true, false, '操作日志', 'document', false, ''),
(now(), now(), 1, 53, 'agreement',    'agreement',    false, 'view/platform/agreement/agreement.vue',         5, '', true, false, '协议管理', 'document', false, ''),
(now(), now(), 1, 53, 'announcement', 'announcement', false, 'view/platform/announcement/announcement.vue',   6, '', true, false, '公告管理', 'bell', false, ''),
(now(), now(), 1, 53, 'setting',      'setting',      false, 'view/platform/setting/setting.vue',             7, '', true, false, '平台参数', 'tools', false, ''),
(now(), now(), 1, 53, 'pricingRule',  'pricingRule',  false, 'view/platform/pricingRule/pricingRule.vue',     8, '', true, false, '计价规则', 'money', false, ''),
(now(), now(), 1, 53, 'mapVersion',   'mapVersion',   false, 'view/platform/mapVersion/mapVersion.vue',       9, '', true, false, '地图版本', 'map-location', false, ''),
(now(), now(), 1, 53, 'pushTemplate', 'pushTemplate', false, 'view/platform/pushTemplate/pushTemplate.vue',  10, '', true, false, '消息模板', 'message', false, ''),
(now(), now(), 1, 53, 'dataDashboard','dataDashboard',false, 'view/platform/dashboard/dashboard.vue',        11, '', true, false, '数据大屏', 'data-analysis', false, '');

-- 3. 角色菜单授权（888 为超级管理员）
INSERT INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id)
SELECT '888', id FROM sys_base_menus WHERE parent_id = 53 OR name = 'platform';

-- 4. 接口权限登记（sys_apis）与 Casbin 授权（casbin_rule）
-- 接口前缀统一为 /platform，共 52 个，节选如下（完整清单见代码 api/v1/platform/platform.go）
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method) VALUES
(now(), now(), '/platform/agreement/list',        '分页获取协议', '平台系统配置', 'GET'),
(now(), now(), '/platform/agreement',             '创建协议',     '平台系统配置', 'POST'),
(now(), now(), '/platform/role/list',             '分页获取角色', '平台系统配置', 'GET'),
(now(), now(), '/platform/role/permissions',      '角色分配权限', '平台系统配置', 'POST'),
(now(), now(), '/platform/operator/list',         '分页获取后台账号', '平台系统配置', 'GET'),
(now(), now(), '/platform/auditLog/list',         '分页获取审计日志', '平台系统配置', 'GET'),
(now(), now(), '/platform/dashboard/overview',    '数据大屏聚合指标', '平台系统配置', 'GET');

INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', '888', '/platform/agreement/list',        'GET'),
('p', '888', '/platform/agreement',             'POST'),
('p', '888', '/platform/role/list',             'GET'),
('p', '888', '/platform/role/permissions',      'POST'),
('p', '888', '/platform/operator/list',         'GET'),
('p', '888', '/platform/auditLog/list',         'GET'),
('p', '888', '/platform/dashboard/overview',    'GET');

-- 5. 补充：数据大屏其余只读接口（cities / areas / station-flow / vehicle-heatmap）
--    说明：middleware.CasbinHandler 已对 GET /platform/dashboard/* 放行；此处登记用于接口管理与文档完整性。
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method) VALUES
(now(), now(), '/platform/dashboard/cities',          '大屏城市下拉', '平台系统配置', 'GET'),
(now(), now(), '/platform/dashboard/areas',           '大屏区域下拉', '平台系统配置', 'GET'),
(now(), now(), '/platform/dashboard/station-flow',    '大屏站点客流', '平台系统配置', 'GET'),
(now(), now(), '/platform/dashboard/vehicle-heatmap', '大屏车辆热力', '平台系统配置', 'GET');

INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', '888', '/platform/dashboard/cities',          'GET'),
('p', '888', '/platform/dashboard/areas',           'GET'),
('p', '888', '/platform/dashboard/station-flow',    'GET'),
('p', '888', '/platform/dashboard/vehicle-heatmap', 'GET');
