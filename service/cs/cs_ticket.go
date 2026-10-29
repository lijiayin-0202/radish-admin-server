package cs

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
	"go.uber.org/zap"
)

type CsTicketService struct{}

// 客服工单状态流转：待受理 -> 处理中 -> 已完结，已完结为终态
// 取值与 cs_tickets.status 注释保持一致：1待处理，2处理中，3已解决，4已关闭
const (
	CsTicketStatusPending    int64 = 1 // 待受理
	CsTicketStatusProcessing int64 = 2 // 处理中
	CsTicketStatusFinished   int64 = 3 // 已完结
	CsTicketStatusClosed     int64 = 4 // 已关闭
)

// csTicketStatusName 工单状态名称，用于错误提示与前端展示对齐
var csTicketStatusName = map[int64]string{
	CsTicketStatusPending:    "待受理",
	CsTicketStatusProcessing: "处理中",
	CsTicketStatusFinished:   "已完结",
	CsTicketStatusClosed:     "已关闭",
}

// csTicketStatusFlow 允许的工单状态流转，已完结为终态
var csTicketStatusFlow = map[int64][]int64{
	CsTicketStatusPending:    {CsTicketStatusProcessing},
	CsTicketStatusProcessing: {CsTicketStatusFinished},
}

// checkCsTicketStatus 校验工单状态取值是否合法：待受理/处理中/已完结
func checkCsTicketStatus(status int64) error {
	if status != CsTicketStatusPending && status != CsTicketStatusProcessing && status != CsTicketStatusFinished {
		return errors.New("工单状态取值不合法，仅支持：1待受理，2处理中，3已完结")
	}
	return nil
}

// checkCsTicketStatusFlow 校验工单状态流转是否合法：待受理 -> 处理中 -> 已完结
func checkCsTicketStatusFlow(oldStatus, newStatus int64) error {
	if oldStatus == newStatus {
		return nil // 状态未变化，不做流转校验
	}
	next, ok := csTicketStatusFlow[oldStatus]
	if !ok {
		return errors.New("工单已" + csTicketStatusName[oldStatus] + "，处于终态不允许再流转")
	}
	for _, s := range next {
		if s == newStatus {
			return nil
		}
	}
	names := make([]string, 0, len(next))
	for _, s := range next {
		names = append(names, csTicketStatusName[s])
	}
	return errors.New("工单状态流转不合法，" + csTicketStatusName[oldStatus] + "只能流转到" + strings.Join(names, "/"))
}

// ==========================================================
// ride-hailing 外部服务（成员C）订单查询
// 说明：本仓库不落地订单/车辆 model，订单信息一律通过 ride-hailing 微服务「订单查询接口」获取，
// 工单表只保存 order_id / order_no 与乘客手机号快照做关联
// ==========================================================

// RideHailingOrderInfo ride-hailing 订单基础信息，字段对齐成员C「订单查询接口」返回结构
type RideHailingOrderInfo struct {
	OrderID        int64  `json:"orderId"`        // 订单ID
	OrderNo        string `json:"orderNo"`        // 订单号
	PassengerName  string `json:"passengerName"`  // 乘客姓名
	PassengerPhone string `json:"passengerPhone"` // 乘客手机号
	StartAddr      string `json:"startAddr"`      // 起点地址
	EndAddr        string `json:"endAddr"`        // 终点地址
	OrderStatus    int64  `json:"orderStatus"`    // 订单状态
}

// fetchRideHailingOrder 查询 ride-hailing（成员C）订单基础信息
// TODO：等 C 接口就绪，替换 mock 为真实 http 调用（替换模板见本文件最下方注释）
// 当前为 mock 版本：orderID 非空时直接返回构造数据，不发起真实 http 请求
func fetchRideHailingOrder(orderID int64) (order *RideHailingOrderInfo, err error) {
	if orderID <= 0 {
		return nil, errors.New("订单ID不合法")
	}
	// 【演示用Mock】对接顺风车订单，暂时注释，避免假手机号入库，联调时打开
	// // mock 数据，仅用于前置联调；接口就绪后本函数整体替换为 http 调用
	// order = &RideHailingOrderInfo{
	// 	OrderID:        orderID,
	// 	OrderNo:        "RH" + strconv.FormatInt(orderID, 10),
	// 	PassengerName:  "mock乘客",
	// 	PassengerPhone: "13800000000",
	// 	StartAddr:      "mock起点地址",
	// 	EndAddr:        "mock终点地址",
	// 	OrderStatus:    2, // TODO(ride-hailing)：订单状态枚举待与成员C确认后对齐
	// }
	// return order, nil
	// Mock 关闭期间不返回任何订单数据，避免假手机号落库；联调时打开上方 Mock 并删除本行
	return nil, errors.New("ride-hailing 订单查询接口未对接，Mock 已关闭")
}

// GetCsTicketList 分页获取客服工单列表，支持工单状态与关键词筛选
func (s *CsTicketService) GetCsTicketList(search csReq.CsTicketSearch) (list []csModel.CsTicket, total int64, err error) {
	db := global.LUOBO_DB.Model(&csModel.CsTicket{})
	if search.TicketNo != "" {
		db = db.Where("ticket_no LIKE ?", "%"+search.TicketNo+"%")
	}
	if search.Keyword != "" {
		kw := "%" + search.Keyword + "%"
		db = db.Where("ticket_no LIKE ? OR title LIKE ? OR description LIKE ?", kw, kw, kw)
	}
	if search.UserID != nil {
		db = db.Where("user_id = ?", *search.UserID)
	}
	if search.UserName != "" {
		db = db.Where("user_name LIKE ?", "%"+search.UserName+"%")
	}
	if search.UserPhone != "" {
		db = db.Where("user_phone = ?", search.UserPhone)
	}
	if search.OrderID != nil {
		db = db.Where("order_id = ?", *search.OrderID)
	}
	if search.OrderNo != "" {
		db = db.Where("order_no = ?", search.OrderNo)
	}
	if search.Status != nil {
		db = db.Where("status = ?", *search.Status)
	}
	if search.Priority != nil {
		db = db.Where("priority = ?", *search.Priority)
	}
	if search.Category != "" {
		db = db.Where("category LIKE ?", "%"+search.Category+"%")
	}
	if search.AssignedClerkID != nil {
		db = db.Where("assigned_clerk_id = ?", *search.AssignedClerkID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(search.PageInfo.Paginate()).Order("id DESC").Find(&list).Error
	return
}

// CreateCsTicket 新增客服工单
// 工单号由 global.Snow 雪花ID自动生成，前端不传；新建工单固定为待受理
// TODO(ride-hailing)：order_id / order_no 只做关联存储，订单与车辆信息调用 ride-hailing 微服务「订单查询接口」获取，
// 本仓库不落地订单/车辆 model
func (s *CsTicketService) CreateCsTicket(operatorID uint, req csReq.CsTicketCreateReq) (entity csModel.CsTicket, err error) {
	if req.Priority < CsTicketStatusPending || req.Priority > 3 { // 优先级：1低，2中，3高
		return entity, errors.New("工单优先级取值不合法，仅支持：1低，2中，3高")
	}
	entity = csModel.CsTicket{
		TicketNo:          strconv.FormatInt(global.Snow.NextID(), 10), // 工单号后端自动生成
		UserID:            req.UserID,
		UserName:          req.UserName,
		UserPhone:         req.UserPhone,
		OrderNo:           req.OrderNo,
		Priority:          req.Priority,
		Category:          req.Category,
		Title:             req.Title,
		Description:       req.Description,
		AttachUrls:        req.AttachUrls,
		AssignedClerkName: req.AssignedClerkName,
		Status:            CsTicketStatusPending, // 新建工单固定为待受理
		CreatedBy:         int64(operatorID),     // 创建者：取当前操作人ID
	}
	if req.AssignedClerkID != nil {
		entity.AssignedClerkID = *req.AssignedClerkID
	}

	// 关联订单：先调用 ride-hailing（成员C）的订单查询接口拿到订单基础信息，再落库
	// TODO(ride-hailing)：等 C 接口就绪，替换 mock 为真实 http 调用（模板见本文件最下方注释）
	if req.OrderID != nil && *req.OrderID > 0 {
		order, orderErr := fetchRideHailingOrder(*req.OrderID)
		if orderErr != nil {
			global.GVA_LOG.Error("调用 ride-hailing 订单查询接口失败!",
				zap.Int64("orderId", *req.OrderID), zap.Error(orderErr))
			return entity, errors.New("获取订单信息失败：" + orderErr.Error())
		}
		entity.OrderID = *req.OrderID
		entity.OrderNo = order.OrderNo          // 订单号快照，来自 ride-hailing
		entity.UserPhone = order.PassengerPhone // 乘客手机号快照，来自 ride-hailing
		if entity.UserName == "" {
			entity.UserName = order.PassengerName // 入参未带乘客姓名时，用订单快照补齐
		}
		// TODO(ride-hailing)：起点/终点（start_addr / end_addr）当前 cs_tickets 无对应列，
		// 需先走 DDL 评审加列后再落库，此处仅记录日志，不做写入
		global.GVA_LOG.Info("已获取 ride-hailing 订单信息",
			zap.Int64("orderId", order.OrderID), zap.String("orderNo", order.OrderNo),
			zap.String("startAddr", order.StartAddr), zap.String("endAddr", order.EndAddr),
			zap.Int64("orderStatus", order.OrderStatus))
	}

	err = global.LUOBO_DB.Create(&entity).Error
	return
}

// UpdateCsTicket 修改客服工单，工单号不允许修改
// 状态变更按流转规则校验：待受理 -> 处理中 -> 已完结；req.Status 为 nil 表示不变更状态
func (s *CsTicketService) UpdateCsTicket(operatorID uint, req csReq.CsTicketUpdateReq) (err error) {
	var entity csModel.CsTicket
	if err = global.LUOBO_DB.Where("id = ?", req.ID).First(&entity).Error; err != nil {
		return errors.New("工单不存在或已删除")
	}
	if req.Priority != 0 && (req.Priority < 1 || req.Priority > 3) {
		return errors.New("工单优先级取值不合法，仅支持：1低，2中，3高")
	}

	// 状态流转校验
	if req.Status != nil {
		if err = checkCsTicketStatus(*req.Status); err != nil {
			return
		}
		if err = checkCsTicketStatusFlow(entity.Status, *req.Status); err != nil {
			return
		}
		// 流转到处理中，必须指定处理客服
		if *req.Status == CsTicketStatusProcessing && req.AssignedClerkID == nil && entity.AssignedClerkID == 0 {
			return errors.New("工单流转到处理中，必须指定处理客服")
		}
		// 流转到已完结，自动补解决时间与关闭时间
		if *req.Status == CsTicketStatusFinished {
			now := time.Now()
			if req.ResolvedAt == nil && entity.ResolvedAt == nil {
				req.ResolvedAt = &now
			}
			if req.ClosedAt == nil && entity.ClosedAt == nil {
				req.ClosedAt = &now
			}
		}
		entity.Status = *req.Status
	}

	entity.UserID = req.UserID
	entity.UserName = req.UserName
	entity.UserPhone = req.UserPhone
	entity.OrderNo = req.OrderNo
	entity.Category = req.Category
	entity.Title = req.Title
	entity.Description = req.Description
	entity.AttachUrls = req.AttachUrls
	entity.AssignedClerkName = req.AssignedClerkName
	entity.Remark = req.Remark
	if req.OrderID != nil {
		entity.OrderID = *req.OrderID
	}
	if req.AssignedClerkID != nil {
		entity.AssignedClerkID = *req.AssignedClerkID
	}
	if req.Priority != 0 {
		entity.Priority = req.Priority // 优先级为0表示不变更
	}
	if req.ResolvedAt != nil {
		entity.ResolvedAt = req.ResolvedAt // 解决时间空值不覆盖已有内容
	}
	if req.ClosedAt != nil {
		entity.ClosedAt = req.ClosedAt // 关闭时间空值不覆盖已有内容
	}
	return global.LUOBO_DB.Save(&entity).Error
}

// DeleteCsTicket 删除客服工单（软删除，依赖 global.GVA_MODEL 自带的 deleted_at）
func (s *CsTicketService) DeleteCsTicket(operatorID uint, id uint) (err error) {
	var entity csModel.CsTicket
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return errors.New("工单不存在或已删除")
	}
	global.GVA_LOG.Info("删除客服工单", zap.Uint("operatorId", operatorID), zap.Uint("ticketId", id))
	return global.LUOBO_DB.Delete(&entity).Error
}

// UpdateCsTicketStatus 修改客服工单状态
// 只变更状态及与之绑定的处理客服、解决/关闭时间、备注；工单号与其它业务字段走 UpdateCsTicket 修改
// 状态变更复用统一流转规则：待受理 -> 处理中 -> 已完结；已完结为终态，状态未变化则不做流转校验
func (s *CsTicketService) UpdateCsTicketStatus(operatorID uint, req csReq.CsTicketStatusReq) (entity csModel.CsTicket, err error) {
	if err = global.LUOBO_DB.Where("id = ?", req.ID).First(&entity).Error; err != nil {
		return entity, errors.New("工单不存在或已删除")
	}
	if err = checkCsTicketStatus(req.Status); err != nil {
		return
	}
	if err = checkCsTicketStatusFlow(entity.Status, req.Status); err != nil {
		return
	}
	// 流转到处理中，必须指定处理客服
	if req.Status == CsTicketStatusProcessing && req.AssignedClerkID == nil && entity.AssignedClerkID == 0 {
		return entity, errors.New("工单流转到处理中，必须指定处理客服")
	}
	// 流转到已完结，自动补解决时间与关闭时间
	if req.Status == CsTicketStatusFinished {
		now := time.Now()
		if req.ResolvedAt == nil && entity.ResolvedAt == nil {
			req.ResolvedAt = &now
		}
		if req.ClosedAt == nil && entity.ClosedAt == nil {
			req.ClosedAt = &now
		}
	}

	entity.Status = req.Status
	if req.AssignedClerkID != nil {
		entity.AssignedClerkID = *req.AssignedClerkID
	}
	if req.AssignedClerkName != "" {
		entity.AssignedClerkName = req.AssignedClerkName
	}
	if req.ResolvedAt != nil {
		entity.ResolvedAt = req.ResolvedAt // 解决时间空值不覆盖已有内容
	}
	if req.ClosedAt != nil {
		entity.ClosedAt = req.ClosedAt // 关闭时间空值不覆盖已有内容
	}
	if req.Remark != "" {
		entity.Remark = req.Remark
	}
	err = global.LUOBO_DB.Save(&entity).Error
	return
}

// GetCsTicket 获取客服工单详情
// TODO(ride-hailing)：详情如需展示订单/车辆信息，调用 ride-hailing 微服务「订单查询接口」按 order_id 获取，
// 本服务只保存 order_id / order_no 做关联，禁止在本地新建订单/车辆 model
func (s *CsTicketService) GetCsTicket(operatorID uint, id uint) (entity csModel.CsTicket, err error) {
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return entity, errors.New("工单不存在或已删除")
	}
	global.GVA_LOG.Info("查询客服工单详情", zap.Uint("operatorId", operatorID), zap.Uint("ticketId", id))
	return
}

// ==========================================================
// TODO(ride-hailing)：等成员C「订单查询接口」就绪后，用下面这段真实 http 实现整段替换 fetchRideHailingOrder 的 mock 版本
// 替换步骤：
//   1. 在 config.yaml 增加 ride-hailing 服务配置（baseUrl / timeout / token），并在 global.GVA_CONFIG
//      对应配置结构体补充字段 —— 属跨模块改动，需先评审，勿直接修改 global 下无关代码
//   2. 保持函数签名 (orderID int64) (*RideHailingOrderInfo, error) 不变，仅替换函数体；
//      CreateCsTicket 其余逻辑、api / router / request / response 均无需改动
//   3. cs_tickets 评审通过增加 start_addr / end_addr 列后，在 model 补字段并打开最下方两行赋值
// ==========================================================
//
// // rideHailingClient ride-hailing 微服务 http 客户端，进程内复用连接，避免频繁建连
// var rideHailingClient = &http.Client{Timeout: 3 * time.Second}
//
// // fetchRideHailingOrder 查询 ride-hailing（成员C）订单基础信息（真实 http 版本）
// func fetchRideHailingOrder(orderID int64) (order *RideHailingOrderInfo, err error) {
// 	if orderID <= 0 {
// 		return nil, errors.New("订单ID不合法")
// 	}
// 	url := fmt.Sprintf("%s/api/ride/order/%d", global.GVA_CONFIG.RideHailing.BaseUrl, orderID)
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Authorization", global.GVA_CONFIG.RideHailing.Token)
//
// 	resp, err := rideHailingClient.Do(req)
// 	if err != nil {
// 		global.GVA_LOG.Error("调用 ride-hailing 订单查询接口失败!",
// 		zap.Int64("orderId", orderID), zap.Error(err))
// 		return nil, errors.New("调用订单服务失败")
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, errors.New("订单服务返回异常状态码：" + resp.Status)
// 	}
//
// 	var result struct {
// 		Code int                  `json:"code"`
// 		Msg  string               `json:"msg"`
// 		Data RideHailingOrderInfo `json:"data"`
// 	}
// 	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, errors.New("解析订单服务响应失败")
// 	}
// 	if result.Code != 0 {
// 		return nil, errors.New(result.Msg)
// 	}
// 	return &result.Data, nil
// }
