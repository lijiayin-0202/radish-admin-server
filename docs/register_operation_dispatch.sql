-- =====================================================================
-- 运营调度管理（operation_dispatch）菜单 / 按钮权限 / 接口 / casbin 策略 注册
-- 目标库：PostgreSQL gva 库（与 sys_* 系统表、casbin_rule 同库）
-- 幂等：可重复执行，不产生重复数据
-- 作用：让「运营调度管理」出现在后台侧边栏，并给对应角色开通接口访问权限
--       （/operationDispatch/* 不在 casbin 只读白名单内，必须显式注册策略）
-- 执行后请【重启后端服务】使 casbin 缓存与菜单生效（或调用 casbin 刷新接口）
-- =====================================================================

DO $$
DECLARE
    v_parent_id bigint;
    v_menu_id   bigint;
    v_auth      record;
BEGIN
    -- ============ 1. 定位 / 创建「出行运营」父菜单（path = 'trip'） ============
    -- 优先复用已有的 trip 父菜单；否则复用某个 trip 子菜单的父级；都没有再新建
    SELECT id INTO v_parent_id FROM sys_base_menus WHERE path = 'trip' AND deleted_at IS NULL ORDER BY id LIMIT 1;
    IF v_parent_id IS NULL THEN
        SELECT DISTINCT parent_id INTO v_parent_id
        FROM sys_base_menus
        WHERE deleted_at IS NULL AND (path LIKE 'trip/%' OR component LIKE 'view/trip/%')
        ORDER BY parent_id LIMIT 1;
    END IF;
    IF v_parent_id IS NULL THEN
        INSERT INTO sys_base_menus (created_at, updated_at, parent_id, path, name, hidden, component, sort, title, icon, keep_alive)
        VALUES (now(), now(), 0, 'trip', 'trip', false, 'view/routerHolder.vue', 0, '出行运营', 'Van', true)
        RETURNING id INTO v_parent_id;
    END IF;

    -- ============ 2. 创建「运营调度管理」子菜单 ============
    SELECT id INTO v_menu_id FROM sys_base_menus WHERE path = 'operationDispatch' AND deleted_at IS NULL ORDER BY id LIMIT 1;
    IF v_menu_id IS NULL THEN
        INSERT INTO sys_base_menus (created_at, updated_at, parent_id, path, name, hidden, component, sort, title, icon, keep_alive)
        VALUES (now(), now(), v_parent_id, 'operationDispatch', 'operationDispatch', false, 'view/trip/operationDispatch.vue', 50, '运营调度管理', 'Promotion', true)
        RETURNING id INTO v_menu_id;
    ELSE
        -- 菜单已存在：纠正可能被改错的组件/名称/父级（例如曾指向不存在的 view/module3/operationManagement.vue）
        UPDATE sys_base_menus
        SET component = 'view/trip/operationDispatch.vue',
            name = 'operationDispatch',
            parent_id = v_parent_id,
            updated_at = now()
        WHERE id = v_menu_id AND component <> 'view/trip/operationDispatch.vue';
    END IF;

    -- ============ 3. 按钮权限（新增 / 编辑 / 删除 / 撤销派单） ============
    INSERT INTO sys_base_menu_btns (created_at, updated_at, name, desc, sys_base_menu_id)
    SELECT now(), now(), 'add',    '新增', v_menu_id WHERE NOT EXISTS (SELECT 1 FROM sys_base_menu_btns WHERE sys_base_menu_id = v_menu_id AND name = 'add');
    INSERT INTO sys_base_menu_btns (created_at, updated_at, name, desc, sys_base_menu_id)
    SELECT now(), now(), 'edit',   '编辑', v_menu_id WHERE NOT EXISTS (SELECT 1 FROM sys_base_menu_btns WHERE sys_base_menu_id = v_menu_id AND name = 'edit');
    INSERT INTO sys_base_menu_btns (created_at, updated_at, name, desc, sys_base_menu_id)
    SELECT now(), now(), 'delete', '删除', v_menu_id WHERE NOT EXISTS (SELECT 1 FROM sys_base_menu_btns WHERE sys_base_menu_id = v_menu_id AND name = 'delete');
    INSERT INTO sys_base_menu_btns (created_at, updated_at, name, desc, sys_base_menu_id)
    SELECT now(), now(), 'cancel', '撤销派单', v_menu_id WHERE NOT EXISTS (SELECT 1 FROM sys_base_menu_btns WHERE sys_base_menu_id = v_menu_id AND name = 'cancel');

    -- ============ 4. 菜单关联角色：继承父菜单已授权角色 + 超级管理员(888) ============
    FOR v_auth IN
        SELECT DISTINCT am.sys_authority_authority_id AS aid
        FROM sys_authority_menus am
        WHERE am.sys_base_menu_id = v_parent_id
        UNION
        SELECT '888' AS aid
        UNION
        SELECT '8881' AS aid
    LOOP
        INSERT INTO sys_authority_menus (sys_base_menu_id, sys_authority_authority_id)
        SELECT v_menu_id, v_auth.aid
        WHERE NOT EXISTS (SELECT 1 FROM sys_authority_menus WHERE sys_base_menu_id = v_menu_id AND sys_authority_authority_id = v_auth.aid);
    END LOOP;

    -- ============ 5. casbin 接口级策略（白名单之外，需显式授权） ============
    FOR v_auth IN
        SELECT DISTINCT am.sys_authority_authority_id AS aid
        FROM sys_authority_menus am
        WHERE am.sys_base_menu_id = v_menu_id
    LOOP
        INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', v_auth.aid, '/operationDispatch/create', 'POST')
          WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0=v_auth.aid AND v1='/operationDispatch/create' AND v2='POST');
        INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', v_auth.aid, '/operationDispatch/update', 'PUT')
          WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0=v_auth.aid AND v1='/operationDispatch/update' AND v2='PUT');
        INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', v_auth.aid, '/operationDispatch/delete', 'DELETE')
          WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0=v_auth.aid AND v1='/operationDispatch/delete' AND v2='DELETE');
        INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', v_auth.aid, '/operationDispatch/deleteByIds', 'DELETE')
          WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0=v_auth.aid AND v1='/operationDispatch/deleteByIds' AND v2='DELETE');
        INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', v_auth.aid, '/operationDispatch/find', 'GET')
          WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0=v_auth.aid AND v1='/operationDispatch/find' AND v2='GET');
        INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', v_auth.aid, '/operationDispatch/list', 'GET')
          WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0=v_auth.aid AND v1='/operationDispatch/list' AND v2='GET');
        INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', v_auth.aid, '/operationDispatch/cancel', 'POST')
          WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0=v_auth.aid AND v1='/operationDispatch/cancel' AND v2='POST');
    END LOOP;

    -- ============ 6. 自动注册接口（API 管理列表展示） ============
    INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
    SELECT now(), now(), '/operationDispatch/create', '创建运营调度', '运营调度管理', 'POST'
    WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path='/operationDispatch/create' AND method='POST');
    INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
    SELECT now(), now(), '/operationDispatch/update', '更新运营调度', '运营调度管理', 'PUT'
    WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path='/operationDispatch/update' AND method='PUT');
    INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
    SELECT now(), now(), '/operationDispatch/delete', '删除运营调度', '运营调度管理', 'DELETE'
    WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path='/operationDispatch/delete' AND method='DELETE');
    INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
    SELECT now(), now(), '/operationDispatch/deleteByIds', '批量删除运营调度', '运营调度管理', 'DELETE'
    WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path='/operationDispatch/deleteByIds' AND method='DELETE');
    INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
    SELECT now(), now(), '/operationDispatch/find', '根据id获取运营调度', '运营调度管理', 'GET'
    WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path='/operationDispatch/find' AND method='GET');
    INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
    SELECT now(), now(), '/operationDispatch/list', '分页获取运营调度列表', '运营调度管理', 'GET'
    WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path='/operationDispatch/list' AND method='GET');
    INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
    SELECT now(), now(), '/operationDispatch/cancel', '撤销派单', '运营调度管理', 'POST'
    WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path='/operationDispatch/cancel' AND method='POST');
END $$;
