-- 订单派车页面所需的后台接口权限（gva 库）
-- 菜单 sys_base_menus.id=67 (tripOrder 订单管理) 已存在，这里只补 API 与 casbin 策略。
BEGIN;

-- 1. 登记接口，便于在「系统管理 → API 管理」中看到并授权
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT now(), now(), v.path, v.description, v.api_group, v.method
FROM (VALUES
    ('/trip/order/list',     '分页获取订单',   '订单管理', 'GET'),
    ('/trip/order/detail',   '订单详情',       '订单管理', 'GET'),
    ('/trip/order/dispatch', '后台派车',       '订单管理', 'POST'),
    ('/trip/order/advance',  '推进订单状态',   '订单管理', 'POST')
) AS v(path, description, api_group, method)
WHERE NOT EXISTS (
    SELECT 1 FROM sys_apis a WHERE a.path = v.path AND a.method = v.method AND a.deleted_at IS NULL
);

-- 2. 授予已绑定该菜单的角色（888 超级管理员 / 2222 运营）
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', role.authority_id, api.path, api.method
FROM (VALUES ('888'), ('2222')) AS role(authority_id)
CROSS JOIN (VALUES
    ('/trip/order/list', 'GET'),
    ('/trip/order/detail', 'GET'),
    ('/trip/order/dispatch', 'POST'),
    ('/trip/order/advance', 'POST')
) AS api(path, method)
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule r
    WHERE r.ptype = 'p' AND r.v0 = role.authority_id AND r.v1 = api.path AND r.v2 = api.method
);

COMMIT;
