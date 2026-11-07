-- 萝卜快跑后台：业务菜单与演示数据 seed
-- 由 source/custom/business_seed.go 在启动时自动执行，版本号 business_data = 1
-- 同一版本只执行一次，之后重启不会覆盖手工修改的数据
-- 也可直接用 Navicat / psql 在 gva 库执行本文件

CREATE TABLE IF NOT EXISTS sys_seed_versions (
    name        varchar(64) PRIMARY KEY,
    version     integer NOT NULL,
    applied_at  timestamptz NOT NULL DEFAULT now()
);

-- ---------- 1. 菜单规整 ----------
-- 隐藏 GVA 原生演示与插件菜单，侧边栏只保留萝卜快跑业务菜单
UPDATE sys_base_menus SET hidden = true, updated_at = now()
WHERE id IN (2, 5, 6, 7, 8, 9, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35,
             36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50)
   OR parent_id IN (5, 6, 9, 41);

-- 业务菜单幂等写入，保证全新库也有一致的菜单树
INSERT INTO sys_base_menus
    (id, menu_level, parent_id, path, name, hidden, component, sort, active_name,
     keep_alive, default_menu, title, icon, close_tab, transition_type, created_at, updated_at)
VALUES
    (53, 0, 0, 'platform',  'platform',  false, 'view/routerHolder.vue', 4, '', false, false, '平台配置', 'setting', false, '', now(), now()),
    (65, 0, 0, 'trip',      'trip',      false, 'view/routerHolder.vue', 2, '', false, false, '出行运营', 'van', false, '', now(), now()),
    (76, 0, 0, 'csCenter',  'csCenter',  false, 'view/routerHolder.vue', 3, '', false, false, '客服中心', 'service', false, '', now(), now()),
    (54, 1, 53, 'role',          'role',          false, 'view/platform/role/role.vue',                 1,  '', false, false, '角色管理',   'user', false, '', now(), now()),
    (55, 1, 53, 'operator',      'operator',      false, 'view/platform/operator/operator.vue',         2,  '', false, false, '账号管理',   'user', false, '', now(), now()),
    (56, 1, 53, 'permission',    'permission',    false, 'view/platform/permission/permission.vue',     3,  '', false, false, '权限管理',   'lock', false, '', now(), now()),
    (57, 1, 53, 'auditLog',      'auditLog',      false, 'view/platform/auditLog/auditLog.vue',         4,  '', false, false, '操作日志',   'document', false, '', now(), now()),
    (58, 1, 53, 'agreement',     'agreement',     false, 'view/platform/agreement/agreement.vue',       5,  '', false, false, '协议管理',   'document', false, '', now(), now()),
    (59, 1, 53, 'announcement',  'announcement',  false, 'view/platform/announcement/announcement.vue', 6,  '', false, false, '公告管理',   'bell', false, '', now(), now()),
    (60, 1, 53, 'setting',       'setting',       false, 'view/platform/setting/setting.vue',           7,  '', false, false, '平台参数',   'tools', false, '', now(), now()),
    (61, 1, 53, 'pricingRule',   'pricingRule',   false, 'view/platform/pricingRule/pricingRule.vue',   8,  '', false, false, '计价规则',   'money', false, '', now(), now()),
    (62, 1, 53, 'mapVersion',    'mapVersion',    false, 'view/platform/mapVersion/mapVersion.vue',     9,  '', false, false, '地图版本',   'map-location', false, '', now(), now()),
    (63, 1, 53, 'pushTemplate',  'pushTemplate',  false, 'view/platform/pushTemplate/pushTemplate.vue', 10, '', false, false, '消息模板',   'message', false, '', now(), now()),
    (64, 1, 53, 'dataDashboard','dataDashboard',false, 'view/platform/dashboard/dashboard.vue',       11, '', false, false, '数据大屏',   'data-analysis', false, '', now(), now()),
    (83, 1, 65, 'operationDispatch','operationDispatch',false,'view/trip/operationDispatch.vue',         0, '', false, false, '运营调度管理','operation', false, '', now(), now()),
    (66, 1, 65, 'passenger',    'passenger',    false, 'view/trip/passenger/passenger.vue',   1,  '', false, false, '乘客管理',   'user', false, '', now(), now()),
    (67, 1, 65, 'tripOrder',    'tripOrder',    false, 'view/trip/order/order.vue',         2,  '', false, false, '订单管理',   'list', false, '', now(), now()),
    (68, 1, 65, 'area',         'area',         false, 'view/trip/area/area.vue',           3,  '', false, false, '运营区域',   'location', false, '', now(), now()),
    (69, 1, 65, 'station',      'station',      false, 'view/trip/station/station.vue',     4,  '', false, false, '站点管理',   'guide', false, '', now(), now()),
    (70, 1, 65, 'vehicle',      'vehicle',      false, 'view/trip/vehicle/vehicle.vue',     5,  '', false, false, '车辆管理',   'van', false, '', now(), now()),
    (71, 1, 65, 'dispatch',     'dispatch',     false, 'view/trip/dispatch/dispatch.vue',   6,  '', false, false, '派单记录',   'sort', false, '', now(), now()),
    (72, 1, 65, 'command',      'command',      false, 'view/trip/command/command.vue',    7,  '', false, false, '车辆指令',   'pointer', false, '', now(), now()),
    (73, 1, 65, 'fault',        'fault',        false, 'view/trip/fault/fault.vue',        8,  '', false, false, '车辆故障',   'warning', false, '', now(), now()),
    (74, 1, 65, 'boarding',     'boarding',     false, 'view/trip/boarding/boarding.vue',  9,  '', false, false, '上车核验',   'ticket', false, '', now(), now()),
    (75, 1, 65, 'incident',     'incident',     false, 'view/trip/incident/incident.vue',  10, '', false, false, '安全事件',   'bell', false, '', now(), now()),
    (77, 1, 76, 'session',      'session',      false, 'view/cs/session/session.vue',             1, '', false, false, '会话管理',  'chat-line-square', false, '', now(), now()),
    (78, 1, 76, 'ticket',       'ticket',       false, 'view/cs/ticket/ticket.vue',               2, '', false, false, '工单管理',  'tickets', false, '', now(), now()),
    (79, 1, 76, 'faqCategory',  'faqCategory',  false, 'view/cs/faqCategory/faqCategory.vue',   3, '', false, false, 'FAQ 分类', 'folder', false, '', now(), now()),
    (80, 1, 76, 'faq',          'faq',          false, 'view/cs/faq/faq.vue',                   4, '', false, false, 'FAQ 管理', 'question-filled', false, '', now(), now()),
    (81, 1, 76, 'notification', 'notification', false, 'view/cs/notification/notification.vue', 5, '', false, false, '消息通知',  'bell', false, '', now(), now())
ON CONFLICT (id) DO UPDATE SET
    menu_level = EXCLUDED.menu_level, parent_id = EXCLUDED.parent_id,
    path = EXCLUDED.path, name = EXCLUDED.name, hidden = EXCLUDED.hidden,
    component = EXCLUDED.component, sort = EXCLUDED.sort,
    title = EXCLUDED.title, icon = EXCLUDED.icon, updated_at = now();

-- 保留的 GVA 原生菜单：统一排序与命名
UPDATE sys_base_menus SET sort = 1, title = '仪表盘',   icon = 'odometer', hidden = false, updated_at = now() WHERE id = 1;
UPDATE sys_base_menus SET sort = 5, title = '系统管理', icon = 'user',     hidden = false, updated_at = now() WHERE id = 3;
UPDATE sys_base_menus SET sort = 6, title = '个人信息', icon = 'message',  hidden = false, updated_at = now() WHERE id = 4;
UPDATE sys_base_menus SET title = 'API 管理', updated_at = now() WHERE id = 12;

SELECT setval(pg_get_serial_sequence('sys_base_menus', 'id'), (SELECT max(id) FROM sys_base_menus));

-- ---------- 2. 角色授权 ----------
-- 超级管理员 888：全部业务菜单
INSERT INTO sys_authority_menus (sys_base_menu_id, sys_authority_authority_id)
SELECT id, 888 FROM sys_base_menus WHERE id BETWEEN 53 AND 83
ON CONFLICT DO NOTHING;
-- 客服角色 1111：仪表盘 + 客服中心
INSERT INTO sys_authority_menus (sys_base_menu_id, sys_authority_authority_id)
SELECT id, 1111 FROM sys_base_menus WHERE id IN (1, 76, 77, 78, 79, 80, 81)
ON CONFLICT DO NOTHING;
-- 运营角色 2222：仪表盘 + 出行运营 + 数据大屏 + 平台配置
INSERT INTO sys_authority_menus (sys_base_menu_id, sys_authority_authority_id)
SELECT id, 2222 FROM sys_base_menus WHERE id IN (1, 53, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 83)
ON CONFLICT DO NOTHING;

-- ---------- 3. 乘客昵称规整，让客服数据能与其他菜单串起来 ----------
UPDATE passengers SET nickname = CASE id
    WHEN 11 THEN '张伟' WHEN 12 THEN '李娜' WHEN 13 THEN '王强'
    WHEN 14 THEN '刘敏' WHEN 15 THEN '陈晨' ELSE nickname END,
    updated_at = now()
WHERE id BETWEEN 11 AND 15;

-- ---------- 4. 客服演示数据重建 ----------
-- 这些表当前都是演示数据，清空后写入一套自洽、贴近真实运营的样例
DELETE FROM cs_ticket_logs;
DELETE FROM cs_tickets;
DELETE FROM cs_sessions;
DELETE FROM cs_faqs;
DELETE FROM cs_faq_categories;
DELETE FROM cs_notifications;

ALTER TABLE cs_tickets ALTER COLUMN priority DROP DEFAULT;
ALTER TABLE cs_tickets ALTER COLUMN priority TYPE varchar(255) USING priority::varchar(255);

-- FAQ 分类
INSERT INTO cs_faq_categories (id, category_name, status, sort, title, created_at, updated_at) VALUES
(1, '账号与登录', 1, 1, '账号与登录', now(), now()),
(2, '订单与派车', 1, 2, '订单与派车', now(), now()),
(3, '支付与发票', 1, 3, '支付与发票', now(), now()),
(4, '安全与行程', 1, 4, '安全与行程', now(), now()),
(5, '优惠与会员', 1, 5, '优惠与会员', now(), now());

-- FAQ 问答，问题、答案、关键词语义一致
INSERT INTO cs_faqs (id, category_id, question, answer, keywords, status, sort, view_count, useful_count, created_at, updated_at) VALUES
(1,  1, '忘记登录密码怎么办？', '在登录页点击「忘记密码」，通过绑定手机号获取验证码后即可重置，重置成功后原密码立即失效。', '忘记密码,重置密码,登录', 1, 1, 328, 41, now(), now()),
(2,  1, '如何修改已绑定的手机号？', '进入「个人信息 → 账号安全 → 更换手机号」，验证原手机号后绑定新号码，24 小时内不可再次变更。', '修改手机号,换绑,账号安全', 1, 2, 215, 28, now(), now()),
(3,  2, '订单一直显示派车中怎么办？', '高峰时段运力紧张时派车会延长，超过 5 分钟未分配车辆可在订单详情点击「催派」，或免费取消后重新下单。', '派车中,催派,等待', 1, 1, 512, 76, now(), now()),
(4,  2, '下单后还可以取消订单吗？', '车辆接单前可免费取消，车辆已出发后取消将按计价规则收取取消费，具体金额以订单详情提示为准。', '取消订单,取消费,退款', 1, 2, 467, 63, now(), now()),
(5,  2, '预约单可以提前多久下单？', '支持提前 7 天预约，系统会在预约时间前 15 分钟开始调度车辆，请保持手机畅通。', '预约单,提前预约', 1, 3, 189, 22, now(), now()),
(6,  3, '如何开具电子发票？', '行程完成后进入「订单详情 → 开发票」，填写抬头与税号，发票将在 24 小时内发送至预留邮箱。', '发票,报销,电子发票', 1, 1, 604, 91, now(), now()),
(7,  3, '被重复扣费了怎么处理？', '可在「订单详情 → 费用疑义」提交申诉，或联系在线客服，核实后 1 至 3 个工作日原路退回。', '重复扣费,退款,支付', 1, 2, 276, 55, now(), now()),
(8,  3, '优惠券为什么没有生效？', '请确认券的有效期、适用城市与最低订单金额，同一订单仅可使用一张券，不满足条件时会自动置灰。', '优惠券,不可用,折扣', 1, 3, 341, 38, now(), now()),
(9,  4, '行程中遇到紧急情况怎么办？', '车内屏幕与 App 均有「一键求助」按钮，点击后将接通 24 小时安全中心并同步行程位置给紧急联系人。', '一键求助,紧急,安全', 1, 1, 388, 82, now(), now()),
(10, 4, '物品遗落在车上如何找回？', '在「订单详情 → 失物招领」登记物品信息，客服会在 2 小时内联系当班司机核查并回电。', '遗失物品,失物招领', 1, 2, 233, 47, now(), now()),
(11, 5, '新人券的使用规则是什么？', '新用户注册后自动发放，有效期 30 天，仅限完成实名认证后的首单使用，不可与其他优惠叠加。', '新人券,首单,有效期', 1, 1, 421, 34, now(), now()),
(12, 5, '会员权益包含哪些内容？', '会员享行程折扣、优先派车、专属客服与每月出行券包，权益有效期以「会员中心」展示为准。', '会员,权益,折扣', 1, 2, 297, 29, now(), now());

-- 客服会话，关联真实乘客与订单
INSERT INTO cs_sessions
    (id, session_no, passenger_id, order_id, channel, scene, issue_category, session_status,
     assigned_to, connected_at, closed_at, user_name, user_phone, clerk_name, status, created_at, updated_at)
VALUES
(1,  'CS202611010001', 11, 10, '在线',  '支付问题',   '支付与发票', '已结束', '客服-小美', '2026-11-01 09:12:00+08', '2026-11-01 09:26:00+08', '张伟', '13800000007', '客服-小美', 0, now(), now()),
(2,  'CS202611010002', 12, 11, '电话',  '订单咨询',   '订单与派车', '已结束', '客服-小王', '2026-11-01 10:05:00+08', '2026-11-01 10:18:00+08', '李娜', '13800000009', '客服-小王', 0, now(), now()),
(3,  'CS202611020003', 13, 12, '在线',  '行程中咨询', '安全与行程', '已结束', '客服-小李', '2026-11-02 11:20:00+08', '2026-11-02 11:35:00+08', '王强', '13800000012', '客服-小李', 0, now(), now()),
(4,  'CS202611020004', 14, 13, 'App内', '退款申请',   '支付与发票', '已结束', '客服-小美', '2026-11-02 14:03:00+08', '2026-11-02 14:22:00+08', '刘敏', '17500009915', '客服-小美', 0, now(), now()),
(5,  'CS202611030005', 11, 14, '在线',  '优惠券咨询', '优惠与会员', '已结束', '客服-小张', '2026-11-03 08:47:00+08', '2026-11-03 08:59:00+08', '张伟', '13800000007', '客服-小张', 0, now(), now()),
(6,  'CS202611030006', 12, 15, '电话',  '派车超时',   '订单与派车', '已结束', '客服-小王', '2026-11-03 16:31:00+08', '2026-11-03 16:52:00+08', '李娜', '13800000009', '客服-小王', 0, now(), now()),
(7,  'CS202611040007', 13, 16, '在线',  '遗失物品',   '安全与行程', '已结束', '客服-小李', '2026-11-04 09:15:00+08', '2026-11-04 09:40:00+08', '王强', '13800000012', '客服-小李', 0, now(), now()),
(8,  'CS202611040008', 15, 17, '在线',  '账号登录',   '账号与登录', '已结束', '客服-小美', '2026-11-04 13:22:00+08', '2026-11-04 13:30:00+08', '陈晨', '17552729253', '客服-小美', 0, now(), now()),
(9,  'CS202611050009', 14, 13, '电话',  '发票开具',   '支付与发票', '已结束', '客服-小张', '2026-11-05 10:08:00+08', '2026-11-05 10:19:00+08', '刘敏', '17500009915', '客服-小张', 0, now(), now()),
(10, 'CS202611060010', 11, 10, '在线',  '投诉建议',   '安全与行程', '服务中', '客服-小美', '2026-11-06 15:44:00+08', NULL, '张伟', '13800000007', '客服-小美', 1, now(), now()),
(11, 'CS202611060011', 12, 11, '在线',  '订单咨询',   '订单与派车', '服务中', '客服-小王', '2026-11-06 16:02:00+08', NULL, '李娜', '13800000009', '客服-小王', 1, now(), now()),
(12, 'CS202611060012', 13, 12, 'App内', '支付问题',   '支付与发票', '排队中', '',          NULL,                     NULL, '王强', '13800000012', '',         2, now(), now());

-- 客服工单，工单号与标题贴近真实业务，关联真实订单号
INSERT INTO cs_tickets
    (id, ticket_no, passenger_id, order_id, order_no, ticket_type, priority, title, description,
     status, assigned_to, resolution, closed_at, user_name, user_phone, category, created_at, updated_at)
VALUES
(1,  'TK202611010001', 11, 10, 'RO20261028000001', '退款',      '高', '订单 RO20261028000001 行程结束后被重复扣费', '乘客反馈同一笔行程被扣款两次，提供了两笔支付流水截图，要求核实并退回多扣金额。', 3, '客服-小美', '核实为支付渠道重复回调，已原路退回 26.80 元并致歉。', '2026-11-02 15:40:00+08', '张伟', '13800000007', '支付与发票', now(), now()),
(2,  'TK202611010002', 12, 11, 'RO20261028000002', '投诉',      '中', '车辆未按预估时间到达上车点', '乘客称车辆比 App 预估时间晚到 12 分钟，导致错过预约，要求说明原因。', 3, '客服-小王', '已调取车辆轨迹核实，补偿 10 元出行券并说明拥堵情况。', '2026-11-01 18:20:00+08', '李娜', '13800000009', '订单与派车', now(), now()),
(3,  'TK202611020003', 13, 12, 'RO20261029000003', '咨询',      '低', '咨询企业发票抬头填写规则', '乘客询问是否可开具公司抬头发票以及税号填写要求。', 3, '客服-小李', '已发送企业开票指引，说明税号与开户行填写规范。', '2026-11-02 12:05:00+08', '王强', '13800000012', '支付与发票', now(), now()),
(4,  'TK202611020004', 14, 13, 'RO20261029000004', '故障反馈',  '高', '行程中车机提示定位异常', '乘客反馈行驶途中车机多次提示定位漂移，担心行程计费不准。', 2, '客服-小张', NULL, NULL, '刘敏', '17500009915', '安全与行程', now(), now()),
(5,  'TK202611030005', 11, 14, 'RO20261029000005', '退款',      '中', '取消订单后优惠券未返还', '乘客取消订单后原优惠券仍显示已使用，希望恢复券的有效期。', 2, '客服-小美', NULL, NULL, '张伟', '13800000007', '优惠与会员', now(), now()),
(6,  'TK202611030006', 12, 15, 'RO20261029000006', '咨询',      '低', '如何修改绑定手机号', '乘客原手机号已停用，咨询换绑流程与所需材料。', 3, '客服-小王', '已指引在账号安全页完成换绑，并提示 24 小时变更限制。', '2026-11-03 17:10:00+08', '李娜', '13800000009', '账号与登录', now(), now()),
(7,  'TK202611040007', 13, 16, 'RO20261028000007', '投诉',      '高', '客服电话等待时间过长', '乘客反馈拨打客服热线等待超过 8 分钟未接通，建议增加高峰坐席。', 1, '客服-小李', NULL, NULL, '王强', '13800000012', '安全与行程', now(), now()),
(8,  'TK202611040008', 15, 17, 'RO20261028000008', '故障反馈',  '中', 'App 行程结束按钮无响应', '乘客升级 App 后点击结束行程无反应，担心持续计费。', 1, '客服-小美', NULL, NULL, '陈晨', '17552729253', '订单与派车', now(), now()),
(9,  'TK202611050009', 14, 13, 'RO20261029000004', '退款',      '中', '系统取消订单后未自动退款', '因运营区域内无可用车辆系统自动取消订单，乘客未收到预付款退回。', 1, '客服-小张', NULL, NULL, '刘敏', '17500009915', '支付与发票', now(), now()),
(10, 'TK202611060010', 11, 10, 'RO20261028000001', '咨询',      '低', '会员权益有效期说明', '乘客咨询年度会员到期后未使用的出行券是否保留。', 1, '客服-小王', NULL, NULL, '张伟', '13800000007', '优惠与会员', now(), now());

-- 工单处理流水
INSERT INTO cs_ticket_logs (id, ticket_id, action, operator, content, created_at, updated_at) VALUES
(1,  1, '创建', '系统',     '乘客提交重复扣费工单，已关联订单 RO20261028000001。', now(), now()),
(2,  1, '回复', '客服-小美', '已收到支付流水截图，正在联系支付渠道核实。', now(), now()),
(3,  1, '关闭', '客服-小美', '重复扣款确认，已原路退回 26.80 元。', now(), now()),
(4,  2, '创建', '系统',     '乘客投诉车辆晚到，已关联订单 RO20261028000002。', now(), now()),
(5,  2, '关闭', '客服-小王', '轨迹核实完毕，补偿出行券并致歉。', now(), now()),
(6,  3, '创建', '系统',     '乘客咨询企业开票规则。', now(), now()),
(7,  3, '关闭', '客服-小李', '已发送开票指引。', now(), now()),
(8,  4, '创建', '系统',     '车机定位异常反馈，已转技术组排查。', now(), now()),
(9,  4, '转派', '客服-小张', '转派至车辆技术组，等待定位日志分析结果。', now(), now()),
(10, 5, '创建', '系统',     '取消订单后优惠券未返还。', now(), now()),
(11, 6, '创建', '系统',     '乘客咨询手机号换绑流程。', now(), now()),
(12, 6, '关闭', '客服-小王', '换绑指引已发送，乘客确认完成。', now(), now());

-- 消息通知，关联真实乘客与订单
INSERT INTO cs_notifications (id, passenger_id, notify_type, title, content, is_read, created_at, updated_at) VALUES
(1, 11, '订单', '订单已创建', '您的订单 RO20261028000001 已创建，正在为您派车。', true,  now(), now()),
(2, 12, '支付', '退款到账通知', '您申请的退款 26.80 元已原路退回，请注意查收。', false, now(), now()),
(3, 13, '安全', '行程安全提示', '行程已开始，请系好安全带，祝您旅途愉快。', true,  now(), now()),
(4, 14, '订单', '订单已取消', '订单 RO20261029000004 因区域内暂无可用车辆已取消，预付款将原路退回。', false, now(), now()),
(5, 15, '活动', '新人券已到账', '欢迎使用萝卜快跑，30 天有效的新人出行券已发放至您的账户。', true,  now(), now()),
(6, 11, '支付', '电子发票已开具', '订单 RO20261028000001 的电子发票已发送至您的预留邮箱。', false, now(), now()),
(7, 12, '系统', '服务时间调整', '亦庄运营区服务时间自 11 月 1 日起调整为 06:30 至 23:00。', true,  now(), now()),
(8, 13, '安全', '遗失物品已找回', '您登记的遗失物品已由当班司机上交，请联系客服预约领取。', false, now(), now());

-- 修正各表自增序列，避免后续新增记录与固定 ID 冲突
SELECT setval(pg_get_serial_sequence('cs_faq_categories', 'id'), (SELECT COALESCE(max(id), 1) FROM cs_faq_categories));
SELECT setval(pg_get_serial_sequence('cs_faqs', 'id'), (SELECT COALESCE(max(id), 1) FROM cs_faqs));
SELECT setval(pg_get_serial_sequence('cs_sessions', 'id'), (SELECT COALESCE(max(id), 1) FROM cs_sessions));
SELECT setval(pg_get_serial_sequence('cs_tickets', 'id'), (SELECT COALESCE(max(id), 1) FROM cs_tickets));
SELECT setval(pg_get_serial_sequence('cs_ticket_logs', 'id'), (SELECT COALESCE(max(id), 1) FROM cs_ticket_logs));
SELECT setval(pg_get_serial_sequence('cs_notifications', 'id'), (SELECT COALESCE(max(id), 1) FROM cs_notifications));

-- ---------- 5. 记录 seed 版本 ----------
INSERT INTO sys_seed_versions (name, version, applied_at) VALUES ('business_data', 1, now())
ON CONFLICT (name) DO UPDATE SET version = EXCLUDED.version, applied_at = now();
