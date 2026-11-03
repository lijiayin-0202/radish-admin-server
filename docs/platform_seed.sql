-- =====================================================================
-- 萝卜快跑后台 平台系统配置模块 种子数据（幂等，可重复执行）
-- 目标库：PostgreSQL gva 库（与 sys_* 系统表同库）
-- 覆盖：角色 / 权限树 / 角色权限 / 计价规则 / 平台参数 / 消息模板 / 协议 / 公告 / 地图版本
-- 用法：在 Navicat 打开 gva 库后，整段执行即可；重复执行不会产生重复数据
-- =====================================================================

DO $$
DECLARE
    top_id   bigint;
    m_id     bigint;
    r_super  bigint;
    r_op     bigint;
    r_cs     bigint;
    menu     record;
BEGIN
    -- ============ 1. 角色 ============
    INSERT INTO roles (created_at, updated_at, role_name, role_code, description, status)
    SELECT now(), now(), v.role_name, v.role_code, v.description, v.status
    FROM (VALUES
        ('超级管理员', 'super_admin',      '平台最高权限，管理所有配置与账号', 1),
        ('运营专员',   'operator',         '负责运营区域、计价规则、公告等日常运营', 1),
        ('客服专员',   'customer_service', '负责客服工单、FAQ 维护', 1)
    ) AS v(role_name, role_code, description, status)
    WHERE NOT EXISTS (SELECT 1 FROM roles r WHERE r.role_code = v.role_code);

    -- ============ 2. 顶层权限 ============
    INSERT INTO permissions (created_at, updated_at, parent_id, permission_name, permission_type, route_path, http_method, sort)
    SELECT now(), now(), 0, '平台系统配置', 1, '/platform', '', 1
    WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.parent_id = 0 AND p.route_path = '/platform');

    SELECT id INTO top_id FROM permissions WHERE parent_id = 0 AND route_path = '/platform' ORDER BY id LIMIT 1;

    -- ============ 3. 菜单 + 按钮 ============
    FOR menu IN
        SELECT * FROM (VALUES
            ('角色管理',   '/platform/role',          1,  false),
            ('账号管理',   '/platform/operator',      2,  false),
            ('权限管理',   '/platform/permission',    3,  false),
            ('操作日志',   '/platform/auditLog',      4,  true),
            ('协议管理',   '/platform/agreement',     5,  false),
            ('公告管理',   '/platform/announcement',  6,  false),
            ('平台参数',   '/platform/setting',       7,  false),
            ('计价规则',   '/platform/pricingRule',   8,  false),
            ('地图版本',   '/platform/mapVersion',    9,  false),
            ('消息模板',   '/platform/pushTemplate',  10, false),
            ('数据大屏',   '/platform/dataDashboard', 11, true)
        ) AS t(perm_name, perm_path, perm_sort, perm_readonly)
    LOOP
        INSERT INTO permissions (created_at, updated_at, parent_id, permission_name, permission_type, route_path, http_method, sort)
        SELECT now(), now(), top_id, menu.perm_name, 1, menu.perm_path, '', menu.perm_sort
        WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.parent_id = top_id AND p.route_path = menu.perm_path);

        SELECT id INTO m_id FROM permissions WHERE parent_id = top_id AND route_path = menu.perm_path ORDER BY id LIMIT 1;

        INSERT INTO permissions (created_at, updated_at, parent_id, permission_name, permission_type, route_path, http_method, sort)
        SELECT now(), now(), m_id, '查看', 2, '', '', 1
        WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.parent_id = m_id AND p.permission_name = '查看');

        IF NOT menu.perm_readonly THEN
            INSERT INTO permissions (created_at, updated_at, parent_id, permission_name, permission_type, route_path, http_method, sort)
            SELECT now(), now(), m_id, '新增', 2, '', '', 2
            WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.parent_id = m_id AND p.permission_name = '新增');

            INSERT INTO permissions (created_at, updated_at, parent_id, permission_name, permission_type, route_path, http_method, sort)
            SELECT now(), now(), m_id, '编辑', 2, '', '', 3
            WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.parent_id = m_id AND p.permission_name = '编辑');

            INSERT INTO permissions (created_at, updated_at, parent_id, permission_name, permission_type, route_path, http_method, sort)
            SELECT now(), now(), m_id, '删除', 2, '', '', 4
            WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.parent_id = m_id AND p.permission_name = '删除');
        END IF;
    END LOOP;

    -- ============ 4. 角色-权限关联 ============
    SELECT id INTO r_super FROM roles WHERE role_code = 'super_admin';
    SELECT id INTO r_op    FROM roles WHERE role_code = 'operator';
    SELECT id INTO r_cs    FROM roles WHERE role_code = 'customer_service';

    -- 超级管理员：全部权限
    INSERT INTO role_permissions (created_at, updated_at, role_id, permission_id)
    SELECT now(), now(), r_super, p.id
    FROM permissions p
    WHERE NOT EXISTS (SELECT 1 FROM role_permissions rp WHERE rp.role_id = r_super AND rp.permission_id = p.id);

    -- 运营专员：顶层 + 协议/公告/平台参数/计价规则/地图版本/消息模板/数据大屏（含按钮）
    INSERT INTO role_permissions (created_at, updated_at, role_id, permission_id)
    SELECT now(), now(), r_op, p.id
    FROM permissions p
    WHERE (
        p.id = top_id
        OR p.route_path IN ('/platform/agreement','/platform/announcement','/platform/setting','/platform/pricingRule','/platform/mapVersion','/platform/pushTemplate','/platform/dataDashboard')
        OR p.parent_id IN (SELECT id FROM permissions WHERE route_path IN ('/platform/agreement','/platform/announcement','/platform/setting','/platform/pricingRule','/platform/mapVersion','/platform/pushTemplate','/platform/dataDashboard'))
    )
    AND NOT EXISTS (SELECT 1 FROM role_permissions rp WHERE rp.role_id = r_op AND rp.permission_id = p.id);

    -- 客服专员：顶层 + 公告/消息模板/数据大屏（含按钮）
    INSERT INTO role_permissions (created_at, updated_at, role_id, permission_id)
    SELECT now(), now(), r_cs, p.id
    FROM permissions p
    WHERE (
        p.id = top_id
        OR p.route_path IN ('/platform/announcement','/platform/pushTemplate','/platform/dataDashboard')
        OR p.parent_id IN (SELECT id FROM permissions WHERE route_path IN ('/platform/announcement','/platform/pushTemplate','/platform/dataDashboard'))
    )
    AND NOT EXISTS (SELECT 1 FROM role_permissions rp WHERE rp.role_id = r_cs AND rp.permission_id = p.id);
END $$;

-- ============ 5. 计价规则 ============
INSERT INTO pricing_rules (created_at, updated_at, rule_name, base_fee, mileage_price, time_price, additional_fee, status, remark)
SELECT now(), now(), v.rule_name, v.base_fee, v.mileage_price, v.time_price, v.additional_fee, v.status, v.remark
FROM (VALUES
    ('标准计价', 10, 2.6, 0.4, 0, 1, '工作日标准计价规则'),
    ('高峰计价', 11, 3.0, 0.5, 0, 1, '早晚高峰时段计价规则'),
    ('夜间计价', 12, 3.2, 0.6, 0, 1, '23:00-05:00 夜间计价规则')
) AS v(rule_name, base_fee, mileage_price, time_price, additional_fee, status, remark)
WHERE NOT EXISTS (SELECT 1 FROM pricing_rules p WHERE p.rule_name = v.rule_name);

-- ============ 6. 平台参数 ============
INSERT INTO platform_settings (created_at, updated_at, config_key, config_value, description, updated_by)
SELECT now(), now(), v.config_key, v.config_value, v.description, v.updated_by
FROM (VALUES
    ('platform_name',              '萝卜快跑',     '平台名称',               '超级管理员'),
    ('service_hotline',            '400-800-0000', '客服热线电话',           '超级管理员'),
    ('default_city',               '北京',         '默认运营城市',           '超级管理员'),
    ('max_booking_advance_hours',  '24',           '最大预约提前时长(小时)', '超级管理员'),
    ('order_timeout_minutes',      '30',           '订单超时自动取消时长(分钟)', '超级管理员'),
    ('fee_rounding',               '0.1',          '计费金额舍入精度(元)',   '超级管理员')
) AS v(config_key, config_value, description, updated_by)
WHERE NOT EXISTS (SELECT 1 FROM platform_settings s WHERE s.config_key = v.config_key);

-- ============ 7. 消息模板 ============
INSERT INTO push_templates (created_at, updated_at, template_name, template_code, content, params, status)
SELECT now(), now(), v.template_name, v.template_code, v.content, v.params, v.status
FROM (VALUES
    ('订单创建通知', 'order_created',     '您的订单已创建，订单号：{order_no}', '["order_no"]', 1),
    ('派单成功通知', 'order_dispatched',  '车辆已接单，车牌号：{plate_no}',     '["plate_no"]', 1),
    ('行程完成通知', 'order_completed',   '行程已结束，费用：{amount}元',      '["amount"]',   1),
    ('订单取消通知', 'order_cancelled',   '订单已取消，原因：{reason}',        '["reason"]',   1)
) AS v(template_name, template_code, content, params, status)
WHERE NOT EXISTS (SELECT 1 FROM push_templates t WHERE t.template_code = v.template_code);

-- ============ 8. 协议 ============
INSERT INTO agreements (created_at, updated_at, version, title, content, status, effective_at)
SELECT now(), now(), v.version, v.title, v.content, v.status, now()
FROM (VALUES
    ('v1.0.0', '用户服务协议', '欢迎使用萝卜快跑自动驾驶出行服务。本协议约定您使用平台服务的权利与义务……', 1),
    ('v1.0.0', '隐私政策',     '我们重视您的个人信息保护，本政策说明我们如何收集、使用、存储和保护您的个人信息……', 1)
) AS v(version, title, content, status)
WHERE NOT EXISTS (SELECT 1 FROM agreements a WHERE a.version = v.version AND a.title = v.title);

-- ============ 9. 公告 ============
INSERT INTO announcements (created_at, updated_at, title, content, region, status, publish_at, publisher)
SELECT now(), now(), v.title, v.content, v.region, v.status, now(), v.publisher
FROM (VALUES
    ('萝卜快跑上线公告', '北京部分地区正式开通自动驾驶出行服务，欢迎体验。', '北京', 1, '超级管理员')
) AS v(title, content, region, status, publisher)
WHERE NOT EXISTS (SELECT 1 FROM announcements a WHERE a.title = v.title);

-- ============ 10. 地图版本 ============
INSERT INTO map_versions (created_at, updated_at, version, publisher, status, is_active, remark)
SELECT now(), now(), v.version, v.publisher, v.status, v.is_active, v.remark
FROM (VALUES
    ('v1.0.0', '超级管理员', 1, true,  '首个上线版本（当前生效）'),
    ('v1.1.0', '超级管理员', 0, false, '下一版本迭代中（草稿）'),
    ('v1.0.1', '超级管理员', 2, false, '历史热修复版本（已下线）')
) AS v(version, publisher, status, is_active, remark)
WHERE NOT EXISTS (SELECT 1 FROM map_versions m WHERE m.version = v.version);
