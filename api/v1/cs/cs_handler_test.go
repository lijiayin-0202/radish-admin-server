package cs_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
	sysReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 客服模块接口冒烟测试（gin-vue-admin httptest 方式）。
//
// 运行前置条件（与本项目一致）：
//   1. 已存在 config.yaml（含 gorm / luoboPgsql / jwt 配置），且资源库与业务库(luobo)可连通
//   2. casbin 已为超管 authorityId=888 授权 /api/cs/* 相关策略（与 ApiPost 调通时一致）
//   3. 执行：go test ./api/v1/cs -v
//
// 测试结构：每个模块按「创建 → 查询 → 更新 → 删除」顺序排列，每个接口至少 1 个正常用例。
// 请求结构体全部复用 model/cs/request 包（Cs 模块），不在此处自定义结构体。
//
// 重要说明：因测试执行顺序问题，先跑更新/删除（如 TestCsXxxUpdate / TestCsXxxDelete / TestCsXxxDetail
// 使用硬编码 id=1）时，若对应模块尚未创建出 id=1 的记录，会报「记录不存在 / record not found」，
// 这属于预期现象，不用处理。建议首次执行时优先跑 Create，再跑其余：
//   go test ./api/v1/cs -run 'TestCs.*Create' -v
//   go test ./api/v1/cs -v

var (
	Router    *gin.Engine
	authToken string

	one int64 = 1
	two int64 = 2
)

// setupEnv 将工作目录切到仓库根并锁定配置文件，确保 core.Viper 能正确加载。
func setupEnv() {
	_, file, _, ok := runtime.Caller(0)
	if ok {
		dir := filepath.Dir(file)
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				_ = os.Chdir(dir)
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	for _, name := range []string{"config.yaml", "resource/config.yaml", "config.test.yaml", "resource/config.test.yaml"} {
		if _, err := os.Stat(name); err == nil {
			if abs, err := filepath.Abs(name); err == nil {
				_ = os.Setenv("GVA_CONFIG", abs)
			}
		}
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setupEnv()

	// 复用 main.go 的初始化流程
	global.GVA_VP = core.Viper()
	initialize.OtherInit()
	global.GVA_LOG = core.Zap()
	global.GVA_DB = initialize.Gorm()
	initialize.InitSnowflake()
	initialize.Timer()
	initialize.DBList()
	initialize.LuoboPgSql()
	initialize.SetupHandlers()
	if global.GVA_DB != nil {
		initialize.RegisterTables()
	}
	Router = initialize.Routers()

	// 直接签发超管 token：JWT 中间件仅校验签名，casbin 使用 authorityId 鉴权
	claims := utils.NewJWT().CreateClaims(sysReq.BaseClaims{
		ID:          1,
		UUID:        uuid.MustParse("88888888-8888-8888-8888-888888888888"),
		Username:    "888888",
		AuthorityId: 888, // 与 CasbinHandler 的超管放行(authorityId==888)对齐，等价于 ApiPost 的真实超管 token
	})
	token, err := utils.NewJWT().CreateToken(claims)
	if err != nil {
		panic("生成测试 token 失败: " + err.Error())
	}
	authToken = token

	os.Exit(m.Run())
}

// doReq 统一发起 HTTP 请求：自动附带 x-token / Authorization；GET 拼 query，POST/PUT 传 json body，DELETE 通过 query 传 id。
func doReq(t *testing.T, method, path string, query url.Values, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body error: %v", err)
		}
		r = httptest.NewRequest(method, path, bytes.NewReader(raw))
		r.Header.Set("Content-Type", "application/json")
	} else {
		u := path
		if len(query) > 0 {
			u += "?" + query.Encode()
		}
		r = httptest.NewRequest(method, u, nil)
	}
	r.Header.Set("x-token", authToken)
	r.Header.Set("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()
	Router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("FAIL %s %s -> status=%d body=%s", method, path, w.Code, w.Body.String())
	}
	return w
}

func q(values map[string]string) url.Values {
	v := url.Values{}
	for k, val := range values {
		v.Set(k, val)
	}
	return v
}

// createSessionAndGetID 先创建一条会话，解析返回体中的真实 sessionId 返回。
// 用于消息等依赖会话的用例，避免硬编码 id=1 导致「找不到会话」。
func createSessionAndGetID(t *testing.T) int64 {
	t.Helper()
	w := doReq(t, "POST", "/api/cs/session", nil, csReq.CsSessionCreateReq{
		UserID:    1001,
		UserName:  "测试用户",
		UserPhone: "13800000000",
		ClerkID:   &one,
		ClerkName: "客服A",
		Remark:    "smoke",
	})
	var resp struct {
		Data struct {
			SessionID uint `json:"sessionId"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析会话创建返回失败: %v body=%s", err, w.Body.String())
	}
	return int64(resp.Data.SessionID)
}

// ===================== csTicketApi（工单） =====================
// 顺序：创建 → 查询(列表/详情) → 更新(改字段/改状态) → 删除

func TestCsTicketCreate(t *testing.T) {
	doReq(t, "POST", "/api/cs/ticket", nil, csReq.CsTicketCreateReq{
		UserID:      1001,
		UserName:    "测试用户",
		UserPhone:   "13800000000",
		Priority:    "2", // 优先级：1低，2中，3高（varchar 字符串，与库列类型一致）
		Category:    "账单",
		Title:       "测试工单标题",
		Description: "测试工单描述",
	})
}

func TestCsTicketList(t *testing.T) {
	doReq(t, "GET", "/api/cs/ticket/list", q(map[string]string{"page": "1", "pageSize": "10"}), nil)
}

func TestCsTicketDetail(t *testing.T) {
	doReq(t, "GET", "/api/cs/ticket/1", nil, nil)
}

func TestCsTicketUpdate(t *testing.T) {
	doReq(t, "PUT", "/api/cs/ticket/1", nil, csReq.CsTicketUpdateReq{
		ID:          1,
		UserID:      1001,
		UserName:    "测试用户",
		UserPhone:   "13800000000",
		Priority:    "2",
		Category:    "账单",
		Title:       "更新工单标题",
		Description: "更新工单描述",
	})
}

func TestCsTicketStatus(t *testing.T) {
	doReq(t, "PUT", "/api/cs/ticket/status", nil, csReq.CsTicketStatusReq{
		ID:                1,
		Status:            2,
		AssignedClerkID:   &one,
		AssignedClerkName: "客服A",
	})
}

func TestCsTicketDelete(t *testing.T) {
	doReq(t, "DELETE", "/api/cs/ticket/1", nil, nil)
}

// ===================== csSessionApi（会话） =====================
// 顺序：创建 → 查询(列表/详情) → 更新 → 删除
// 注：user_name 已允许为空，可不传；此处保留传值仅覆盖字段，不影响。

func TestCsSessionCreate(t *testing.T) {
	doReq(t, "POST", "/api/cs/session", nil, csReq.CsSessionCreateReq{
		UserID:    1001,
		UserName:  "测试用户",
		UserPhone: "13800000000",
		ClerkID:   &one,
		ClerkName: "客服A",
		Remark:    "smoke",
	})
}

// TestCsSessionAccept 验证接单流程：先创建待分配会话（不指定客服，clerk_id 默认 0），再由客服接单
func TestCsSessionAccept(t *testing.T) {
	w := doReq(t, "POST", "/api/cs/session", nil, csReq.CsSessionCreateReq{
		UserID:    1002,
		UserName:  "待接单用户",
		UserPhone: "13800000002",
		Remark:    "accept-smoke",
	})
	var resp struct {
		Data struct {
			SessionID uint `json:"sessionId"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析会话创建返回失败: %v body=%s", err, w.Body.String())
	}
	// 客服接单：乐观锁 clerk_id=0 命中，RowsAffected=1 即成功
	doReq(t, "POST", "/api/cs/session/accept", nil, map[string]interface{}{"id": resp.Data.SessionID})
}

func TestCsSessionList(t *testing.T) {
	doReq(t, "GET", "/api/cs/session/list", q(map[string]string{"page": "1", "pageSize": "10"}), nil)
}

func TestCsSessionDetail(t *testing.T) {
	doReq(t, "GET", "/api/cs/session", q(map[string]string{"id": "1"}), nil)
}

func TestCsSessionUpdate(t *testing.T) {
	doReq(t, "PUT", "/api/cs/session", nil, csReq.CsSessionUpdateReq{
		ID:        1,
		UserID:    1001,
		UserName:  "测试用户",
		UserPhone: "13800000000",
		ClerkID:   &one,
		ClerkName: "客服A",
		Status:    &two,
		Remark:    "smoke",
	})
}

func TestCsSessionDelete(t *testing.T) {
	doReq(t, "DELETE", "/api/cs/session", q(map[string]string{"id": "1"}), nil)
}

// ===================== csMessageApi（会话消息） =====================
// 顺序：创建 → 查询(列表/详情) → 更新 → 删除
// 注：cs_messages.is_read 为独立表，保持 int64（0未读/1已读），不改为 bool。

func TestCsMessageCreate(t *testing.T) {
	// 先创建会话拿到真实 sessionId，避免写死 id 导致找不到会话
	sessionID := createSessionAndGetID(t)
	doReq(t, "POST", "/api/cs/message", nil, csReq.CsMessageCreateReq{
		SessionID:  sessionID,
		SenderType: 1,
		MsgType:    1,
		Content:    "hello",
		IsRead:     0,
	})
}

func TestCsMessageList(t *testing.T) {
	// 使用真实存在的会话 ID 过滤，避免空结果/找不到会话
	sessionID := createSessionAndGetID(t)
	doReq(t, "GET", "/api/cs/message/list", q(map[string]string{"page": "1", "pageSize": "10", "sessionId": strconv.FormatInt(sessionID, 10)}), nil)
}

func TestCsMessageDetail(t *testing.T) {
	doReq(t, "GET", "/api/cs/message", q(map[string]string{"id": "1"}), nil)
}

func TestCsMessageUpdate(t *testing.T) {
	doReq(t, "PUT", "/api/cs/message", nil, csReq.CsMessageUpdateReq{
		ID:         1,
		SenderType: 1,
		MsgType:    1,
		Content:    "updated",
		IsRead:     0,
	})
}

func TestCsMessageDelete(t *testing.T) {
	doReq(t, "DELETE", "/api/cs/message", q(map[string]string{"id": "1"}), nil)
}

// ===================== csFaqApi（常见问题） =====================
// 顺序：创建 → 查询(列表/详情/搜索) → 更新 → 删除

func TestCsFaqCreate(t *testing.T) {
	doReq(t, "POST", "/api/cs/faq", nil, csReq.CsFaqCreateReq{
		Question:   "如何申请退款",
		Answer:     "请联系在线客服",
		CategoryID: 1,
		Tags:       "退款",
		Sort:       1,
		Status:     &one,
	})
}

func TestCsFaqList(t *testing.T) {
	doReq(t, "GET", "/api/cs/faq/list", q(map[string]string{"page": "1", "pageSize": "10"}), nil)
}

func TestCsFaqDetail(t *testing.T) {
	doReq(t, "GET", "/api/cs/faq", q(map[string]string{"id": "1"}), nil)
}

func TestCsFaqSearch(t *testing.T) {
	doReq(t, "GET", "/api/cs/faq/search", q(map[string]string{"keyword": "退款"}), nil)
}

func TestCsFaqUpdate(t *testing.T) {
	doReq(t, "PUT", "/api/cs/faq", nil, csReq.CsFaqUpdateReq{
		ID:         1,
		Question:   "如何申请退款",
		Answer:     "请联系在线客服",
		CategoryID: 1,
		Tags:       "退款",
		Sort:       1,
		Status:     &one,
	})
}

func TestCsFaqDelete(t *testing.T) {
	doReq(t, "DELETE", "/api/cs/faq", q(map[string]string{"id": "1"}), nil)
}

// ===================== csFaqCategoryApi（常见问题分类） =====================
// 顺序：创建 → 查询(列表/详情) → 更新 → 删除
// 注：分类名称统一使用 title 字段，不再使用 name。

func TestCsFaqCategoryCreate(t *testing.T) {
	doReq(t, "POST", "/api/cs/faqCategory", nil, csReq.CsFaqCategoryCreateReq{
		Title:  "常见问题",
		Sort:   1,
		Status: &one,
		Remark: "smoke",
	})
}

func TestCsFaqCategoryList(t *testing.T) {
	doReq(t, "GET", "/api/cs/faqCategory/list", q(map[string]string{"page": "1", "pageSize": "10"}), nil)
}

func TestCsFaqCategoryDetail(t *testing.T) {
	doReq(t, "GET", "/api/cs/faqCategory", q(map[string]string{"id": "1"}), nil)
}

func TestCsFaqCategoryUpdate(t *testing.T) {
	doReq(t, "PUT", "/api/cs/faqCategory", nil, csReq.CsFaqCategoryUpdateReq{
		ID:     1,
		Title:  "常见问题",
		Sort:   1,
		Status: &one,
		Remark: "smoke",
	})
}

func TestCsFaqCategoryDelete(t *testing.T) {
	doReq(t, "DELETE", "/api/cs/faqCategory", q(map[string]string{"id": "1"}), nil)
}

// ===================== csNoticeApi（站内公告） =====================
// 顺序：创建 → 查询(列表/详情) → 更新 → 删除
// 注：已补齐 status / sort 字段，与模型一致。

func TestCsNoticeCreate(t *testing.T) {
	doReq(t, "POST", "/api/cs/notice", nil, csReq.CsNoticeCreateReq{
		Title:   "系统维护公告",
		Content: "本周六凌晨维护",
		Type:    1,
		Status:  1,
		Sort:    1,
	})
}

func TestCsNoticeList(t *testing.T) {
	doReq(t, "GET", "/api/cs/notice/list", q(map[string]string{"page": "1", "pageSize": "10"}), nil)
}

func TestCsNoticeDetail(t *testing.T) {
	doReq(t, "GET", "/api/cs/notice", q(map[string]string{"id": "1"}), nil)
}

func TestCsNoticeUpdate(t *testing.T) {
	doReq(t, "PUT", "/api/cs/notice", nil, csReq.CsNoticeUpdateReq{
		ID:      1,
		Title:   "系统维护公告",
		Content: "本周六凌晨维护",
		Type:    1,
		Status:  1,
		Sort:    1,
	})
}

func TestCsNoticeDelete(t *testing.T) {
	doReq(t, "DELETE", "/api/cs/notice", q(map[string]string{"id": "1"}), nil)
}

// ===================== csNotificationApi（用户通知） =====================
// 仅 List / Detail / Read 三个接口（无 create/update/delete 路由）。
// 注：is_read 已统一为 bool（false未读/true已读），由 /read 接口标记；此处不传 isRead 入参。

func TestCsNotificationList(t *testing.T) {
	doReq(t, "GET", "/api/cs/notification/list", q(map[string]string{"page": "1", "pageSize": "10"}), nil)
}

func TestCsNotificationDetail(t *testing.T) {
	doReq(t, "GET", "/api/cs/notification", q(map[string]string{"id": "1"}), nil)
}

func TestCsNotificationRead(t *testing.T) {
	doReq(t, "PUT", "/api/cs/notification/read", nil, csReq.CsNotificationReadReq{
		ID: 1,
	})
}

// ===================== csTicketLogApi（工单流水） =====================
// 顺序：创建 → 查询(列表/详情)

func TestCsTicketLogCreate(t *testing.T) {
	doReq(t, "POST", "/api/cs/ticketLog", nil, csReq.CsTicketLogCreateReq{
		TicketID:     1,
		StatusBefore: 1,
		StatusAfter:  2,
		Remark:       "smoke",
	})
}

func TestCsTicketLogList(t *testing.T) {
	doReq(t, "GET", "/api/cs/ticketLog/list", q(map[string]string{"page": "1", "pageSize": "10"}), nil)
}

func TestCsTicketLogDetail(t *testing.T) {
	doReq(t, "GET", "/api/cs/ticketLog", q(map[string]string{"id": "1"}), nil)
}
