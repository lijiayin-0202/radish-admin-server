-- =============================================================
-- 模块5：客服中心与客服运营端 —— PostgreSQL 建表脚本
-- 表前缀：cs_   关联微服务：ride-hailing（订单仅查询转发）
-- 说明：字段注释与 model/cs 下 GORM tag 一一对应，禁止使用 AutoMigrate
-- =============================================================

-- ----------------------------
-- 1. 客服工单表 cs_tickets
-- ----------------------------
CREATE TABLE IF NOT EXISTS cs_tickets (
    id                  BIGSERIAL     PRIMARY KEY,
    created_at          TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    ticket_no           VARCHAR(64)   NOT NULL DEFAULT '',
    user_id             BIGINT        DEFAULT 0,
    user_name           VARCHAR(128)  DEFAULT '',
    user_phone          VARCHAR(32)   DEFAULT '',
    order_id            BIGINT        DEFAULT 0,
    order_no            VARCHAR(64)   DEFAULT '',
    status              INTEGER       DEFAULT 1,
    priority            INTEGER       DEFAULT 1,
    category            VARCHAR(128)  DEFAULT '',
    title               VARCHAR(256)  NOT NULL DEFAULT '',
    description         TEXT,
    attach_urls         TEXT,
    assigned_clerk_id   BIGINT        DEFAULT 0,
    assigned_clerk_name VARCHAR(128)  DEFAULT '',
    resolved_at         TIMESTAMPTZ,
    closed_at           TIMESTAMPTZ,
    created_by          BIGINT        DEFAULT 0,
    remark              TEXT
);

COMMENT ON TABLE cs_tickets IS '客服工单表';
COMMENT ON COLUMN cs_tickets.id                  IS '主键ID';
COMMENT ON COLUMN cs_tickets.created_at          IS '创建时间';
COMMENT ON COLUMN cs_tickets.updated_at          IS '更新时间';
COMMENT ON COLUMN cs_tickets.deleted_at          IS '删除时间（GVA软删除）';
COMMENT ON COLUMN cs_tickets.ticket_no           IS '工单号';
COMMENT ON COLUMN cs_tickets.user_id             IS '用户ID';
COMMENT ON COLUMN cs_tickets.user_name           IS '用户姓名';
COMMENT ON COLUMN cs_tickets.user_phone          IS '用户手机号';
COMMENT ON COLUMN cs_tickets.order_id            IS '关联订单ID';
COMMENT ON COLUMN cs_tickets.order_no            IS '关联订单号';
COMMENT ON COLUMN cs_tickets.status              IS '工单状态：1待处理，2处理中，3已解决，4已关闭';
COMMENT ON COLUMN cs_tickets.priority            IS '优先级：1低，2中，3高';
COMMENT ON COLUMN cs_tickets.category            IS '问题分类';
COMMENT ON COLUMN cs_tickets.title               IS '工单标题';
COMMENT ON COLUMN cs_tickets.description         IS '问题描述';
COMMENT ON COLUMN cs_tickets.attach_urls         IS '附件URL列表';
COMMENT ON COLUMN cs_tickets.assigned_clerk_id   IS '分配客服ID';
COMMENT ON COLUMN cs_tickets.assigned_clerk_name IS '分配客服姓名';
COMMENT ON COLUMN cs_tickets.resolved_at         IS '解决时间';
COMMENT ON COLUMN cs_tickets.closed_at           IS '关闭时间';
COMMENT ON COLUMN cs_tickets.created_by          IS '创建者（操作人ID）';
COMMENT ON COLUMN cs_tickets.remark              IS '备注';

CREATE INDEX IF NOT EXISTS idx_cs_tickets_deleted_at ON cs_tickets (deleted_at);
CREATE INDEX IF NOT EXISTS idx_cs_tickets_ticket_no  ON cs_tickets (ticket_no);
CREATE INDEX IF NOT EXISTS idx_cs_tickets_user_id    ON cs_tickets (user_id);
CREATE INDEX IF NOT EXISTS idx_cs_tickets_order_id   ON cs_tickets (order_id);
CREATE INDEX IF NOT EXISTS idx_cs_tickets_status     ON cs_tickets (status);
CREATE INDEX IF NOT EXISTS idx_cs_tickets_priority   ON cs_tickets (priority);

-- ----------------------------
-- 2. 客服会话表 cs_sessions
-- ----------------------------
CREATE TABLE IF NOT EXISTS cs_sessions (
    id               BIGSERIAL     PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    session_no       VARCHAR(64)   NOT NULL DEFAULT '',
    ticket_id        BIGINT        DEFAULT 0,
    user_id          BIGINT        DEFAULT 0,
    user_name        VARCHAR(128)  DEFAULT '',
    user_phone       VARCHAR(32)   DEFAULT '',
    clerk_id         BIGINT        DEFAULT 0,
    clerk_name       VARCHAR(128)  DEFAULT '',
    status           INTEGER       DEFAULT 1,
    started_at       TIMESTAMPTZ,
    ended_at         TIMESTAMPTZ,
    last_msg_id      BIGINT        DEFAULT 0,
    last_msg_content VARCHAR(512)  DEFAULT '',
    last_msg_time    TIMESTAMPTZ,
    unread_user      INTEGER       DEFAULT 0,
    unread_clerk     INTEGER       DEFAULT 0,
    remark           TEXT
);

COMMENT ON TABLE cs_sessions IS '客服会话表';
COMMENT ON COLUMN cs_sessions.id               IS '主键ID';
COMMENT ON COLUMN cs_sessions.created_at       IS '创建时间';
COMMENT ON COLUMN cs_sessions.updated_at       IS '更新时间';
COMMENT ON COLUMN cs_sessions.deleted_at       IS '删除时间（GVA软删除）';
COMMENT ON COLUMN cs_sessions.session_no       IS '会话编号';
COMMENT ON COLUMN cs_sessions.ticket_id        IS '关联工单ID';
COMMENT ON COLUMN cs_sessions.user_id          IS '用户ID';
COMMENT ON COLUMN cs_sessions.user_name        IS '用户姓名';
COMMENT ON COLUMN cs_sessions.user_phone       IS '用户手机号';
COMMENT ON COLUMN cs_sessions.clerk_id         IS '客服ID';
COMMENT ON COLUMN cs_sessions.clerk_name       IS '客服姓名';
COMMENT ON COLUMN cs_sessions.status           IS '会话状态：1进行中，2已结束';
COMMENT ON COLUMN cs_sessions.started_at       IS '开始时间';
COMMENT ON COLUMN cs_sessions.ended_at         IS '结束时间';
COMMENT ON COLUMN cs_sessions.last_msg_id      IS '最后消息ID';
COMMENT ON COLUMN cs_sessions.last_msg_content IS '最后消息摘要';
COMMENT ON COLUMN cs_sessions.last_msg_time    IS '最后消息时间';
COMMENT ON COLUMN cs_sessions.unread_user      IS '用户未读消息数';
COMMENT ON COLUMN cs_sessions.unread_clerk     IS '客服未读消息数';
COMMENT ON COLUMN cs_sessions.remark           IS '备注';

CREATE INDEX IF NOT EXISTS idx_cs_sessions_deleted_at ON cs_sessions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_cs_sessions_session_no ON cs_sessions (session_no);
CREATE INDEX IF NOT EXISTS idx_cs_sessions_ticket_id  ON cs_sessions (ticket_id);
CREATE INDEX IF NOT EXISTS idx_cs_sessions_user_id    ON cs_sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_cs_sessions_clerk_id   ON cs_sessions (clerk_id);
CREATE INDEX IF NOT EXISTS idx_cs_sessions_status     ON cs_sessions (status);

-- ----------------------------
-- 3. 客服会话消息表 cs_messages
-- ----------------------------
CREATE TABLE IF NOT EXISTS cs_messages (
    id          BIGSERIAL     PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    session_id  BIGINT        DEFAULT 0,
    sender_type INTEGER       DEFAULT 1,
    msg_type    INTEGER       DEFAULT 1,
    content     TEXT          NOT NULL DEFAULT '',
    attach_url  VARCHAR(500)  DEFAULT '',
    seq         BIGINT        DEFAULT 0,
    is_read     INTEGER       DEFAULT 0,
    send_time   TIMESTAMPTZ
);

COMMENT ON TABLE cs_messages IS '客服会话消息表';
COMMENT ON COLUMN cs_messages.id          IS '主键ID';
COMMENT ON COLUMN cs_messages.created_at  IS '创建时间';
COMMENT ON COLUMN cs_messages.updated_at  IS '更新时间';
COMMENT ON COLUMN cs_messages.deleted_at  IS '删除时间（GVA软删除）';
COMMENT ON COLUMN cs_messages.session_id  IS '归属会话ID，关联cs_sessions';
COMMENT ON COLUMN cs_messages.sender_type IS '发送方类型：1乘客，2客服，3系统';
COMMENT ON COLUMN cs_messages.msg_type    IS '消息类型：1文本，2图片，3文件，4系统通知';
COMMENT ON COLUMN cs_messages.content     IS '消息文本内容';
COMMENT ON COLUMN cs_messages.attach_url  IS '附件地址';
COMMENT ON COLUMN cs_messages.seq         IS '消息序号，同一会话内递增';
COMMENT ON COLUMN cs_messages.is_read     IS '已读状态：0未读，1已读';
COMMENT ON COLUMN cs_messages.send_time   IS '消息发送时间';

CREATE INDEX IF NOT EXISTS idx_cs_messages_deleted_at  ON cs_messages (deleted_at);
CREATE INDEX IF NOT EXISTS idx_cs_messages_session_id  ON cs_messages (session_id);
CREATE INDEX IF NOT EXISTS idx_cs_messages_sender_type ON cs_messages (sender_type);
CREATE INDEX IF NOT EXISTS idx_cs_messages_msg_type    ON cs_messages (msg_type);
CREATE INDEX IF NOT EXISTS idx_cs_messages_is_read     ON cs_messages (is_read);

-- ----------------------------
-- 4. 常见问题表 cs_faqs
-- ----------------------------
CREATE TABLE IF NOT EXISTS cs_faqs (
    id           BIGSERIAL     PRIMARY KEY,
    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    question     VARCHAR(512)  NOT NULL DEFAULT '',
    answer       TEXT          NOT NULL,
    category     VARCHAR(128)  DEFAULT '',
    tags         VARCHAR(256)  DEFAULT '',
    sort         INTEGER       DEFAULT 0,
    status       INTEGER       DEFAULT 1,
    view_count   BIGINT        DEFAULT 0,
    useful_count BIGINT        DEFAULT 0,
    remark       TEXT
);

COMMENT ON TABLE cs_faqs IS '常见问题表';
COMMENT ON COLUMN cs_faqs.id           IS '主键ID';
COMMENT ON COLUMN cs_faqs.created_at   IS '创建时间';
COMMENT ON COLUMN cs_faqs.updated_at   IS '更新时间';
COMMENT ON COLUMN cs_faqs.deleted_at   IS '删除时间（GVA软删除）';
COMMENT ON COLUMN cs_faqs.question     IS '问题标题';
COMMENT ON COLUMN cs_faqs.answer       IS '回答内容';
COMMENT ON COLUMN cs_faqs.category     IS '问题分类';
COMMENT ON COLUMN cs_faqs.tags         IS '标签';
COMMENT ON COLUMN cs_faqs.sort         IS '排序号';
COMMENT ON COLUMN cs_faqs.status       IS '状态：0禁用，1启用';
COMMENT ON COLUMN cs_faqs.view_count   IS '查看次数';
COMMENT ON COLUMN cs_faqs.useful_count IS '有用次数';
COMMENT ON COLUMN cs_faqs.remark       IS '备注';

CREATE INDEX IF NOT EXISTS idx_cs_faqs_deleted_at ON cs_faqs (deleted_at);
CREATE INDEX IF NOT EXISTS idx_cs_faqs_category   ON cs_faqs (category);
CREATE INDEX IF NOT EXISTS idx_cs_faqs_status     ON cs_faqs (status);
CREATE INDEX IF NOT EXISTS idx_cs_faqs_sort       ON cs_faqs (sort);

-- ----------------------------
-- 5. 补充升级脚本：cs_tickets 增加 created_by 列
-- 说明：CsTicket 模型新增 CreatedBy 字段（记录操作人ID），第 1 段 CREATE TABLE 已同步加列；
--      已存在的库需手动执行本段。禁止使用 AutoMigrate，请执行后确认再启动服务，
--      否则 cs_tickets 的查询与新增会报 column "created_by" does not exist
-- ----------------------------
ALTER TABLE cs_tickets ADD COLUMN IF NOT EXISTS created_by BIGINT DEFAULT 0;

COMMENT ON COLUMN cs_tickets.created_by IS '创建者（操作人ID）';
