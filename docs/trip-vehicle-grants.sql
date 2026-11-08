-- 车辆管理页面所需的后台接口权限（gva 库）
-- 菜单 sys_base_menus.id=70 (vehicle 车辆管理) 已存在，这里只补缺失的 API 与 casbin 策略。
BEGIN;

-- 1. 登记接口（list 已存在历史记录，仅补缺）
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT now(), now(), v.path, v.description, v.api_group, v.method
FROM (VALUES
    ('/trip/vehicle/list',   '分页获取车辆',  '车辆管理', 'GET'),
    ('/trip/vehicle/detail', '车辆详情',      '车辆管理', 'GET')
) AS v(path, description, api_group, method)
WHERE NOT EXISTS (
    SELECT 1 FROM sys_apis a WHERE a.path = v.path AND a.method = v.method AND a.deleted_at IS NULL
);

-- 2. 授予已绑定该菜单的角色（888 超级管理员 / 2222 运营）
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', role.authority_id, api.path, api.method
FROM (VALUES ('888'), ('2222')) AS role(authority_id)
CROSS JOIN (VALUES
    ('/trip/vehicle/list', 'GET'),
    ('/trip/vehicle/detail', 'GET')
) AS api(path, method)
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule r
    WHERE r.ptype = 'p' AND r.v0 = role.authority_id AND r.v1 = api.path AND r.v2 = api.method
);

COMMIT;
