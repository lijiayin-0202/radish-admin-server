-- 网约车 ride 模块（代理 ride-hailing 真实业务库）接口权限（gva 库）
-- 与主干 trip 模块（gva 库自建业务表）并存：ride 模块统一使用 /ride 前缀。
BEGIN;

-- 1. 登记接口
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT now(), now(), v.path, v.description, v.api_group, v.method
FROM (VALUES
    ('/ride/order/list',      '分页获取订单',   '网约车订单', 'GET'),
    ('/ride/order/detail',    '订单详情',       '网约车订单', 'GET'),
    ('/ride/order/dispatch',  '后台派车',       '网约车订单', 'POST'),
    ('/ride/order/advance',   '推进订单状态',   '网约车订单', 'POST'),
    ('/ride/passenger/list',  '分页获取乘客',   '网约车乘客', 'GET'),
    ('/ride/passenger/detail','乘客详情',       '网约车乘客', 'GET'),
    ('/ride/vehicle/list',    '分页获取车辆',   '网约车车辆', 'GET'),
    ('/ride/vehicle/detail',  '车辆详情',       '网约车车辆', 'GET'),
    ('/ride/dispatch/list',   '分页获取派单记录','网约车派单', 'GET')
) AS v(path, description, api_group, method)
WHERE NOT EXISTS (
    SELECT 1 FROM sys_apis a WHERE a.path = v.path AND a.method = v.method AND a.deleted_at IS NULL
);

-- 2. 授予 888 超级管理员与 2222 运营角色
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', role.authority_id, api.path, api.method
FROM (VALUES ('888'), ('2222')) AS role(authority_id)
CROSS JOIN (VALUES
    ('/ride/order/list', 'GET'),
    ('/ride/order/detail', 'GET'),
    ('/ride/order/dispatch', 'POST'),
    ('/ride/order/advance', 'POST'),
    ('/ride/passenger/list', 'GET'),
    ('/ride/passenger/detail', 'GET'),
    ('/ride/vehicle/list', 'GET'),
    ('/ride/vehicle/detail', 'GET'),
    ('/ride/dispatch/list', 'GET')
) AS api(path, method)
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule r
    WHERE r.ptype = 'p' AND r.v0 = role.authority_id AND r.v1 = api.path AND r.v2 = api.method
);

COMMIT;
