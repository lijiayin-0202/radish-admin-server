-- 派单记录页面所需的后台接口权限（gva 库）
-- 菜单 sys_base_menus.id=71 (dispatch 派单记录) 已存在；sys_apis 中 /trip/dispatch/list 已登记，
-- 这里只补缺失的 casbin 策略。
BEGIN;

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', role.authority_id, api.path, api.method
FROM (VALUES ('888'), ('2222')) AS role(authority_id)
CROSS JOIN (VALUES
    ('/trip/dispatch/list', 'GET')
) AS api(path, method)
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule r
    WHERE r.ptype = 'p' AND r.v0 = role.authority_id AND r.v1 = api.path AND r.v2 = api.method
);

COMMIT;
