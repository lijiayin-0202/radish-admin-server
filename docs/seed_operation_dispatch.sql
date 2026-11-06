-- =====================================================================
-- 运营调度管理（operation_dispatch）演示示例数据
-- 目标库：PostgreSQL postgres 库（业务表所在库，与 gva 系统库同主机）
-- 作用：插入覆盖全部 6 种状态的演示数据，供答辩/团队演示使用
-- 注意：菜单/角色/casbin 策略请用同目录 register_operation_dispatch.sql；本文件只管业务数据
-- 可重复执行：先清空再插入（演示数据，非生产数据）
-- =====================================================================
TRUNCATE TABLE operation_dispatch RESTART IDENTITY CASCADE;

INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103001',1001,2001,3001,0,'2026-11-03 08:15:00+00:00','待调度：用户预约出行');
INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103002',1002,2002,3002,0,'2026-11-03 08:40:00+00:00','待调度：企业通勤');
INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103003',1003,2003,3003,1,'2026-11-03 09:05:00+00:00','调度中：就近派车');
INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103004',1004,2004,3004,1,'2026-11-03 09:20:00+00:00','调度中：机场接送');
INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103005',1005,2005,3005,2,'2026-11-03 09:35:00+00:00','已派单：司机已接单');
INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103006',1006,2006,3006,2,'2026-11-03 09:50:00+00:00','已派单：站点循环');
INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103007',1007,2007,3007,3,'2026-11-03 10:10:00+00:00','已完成：行程结束');
INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103008',1008,2008,3008,3,'2026-11-03 10:25:00+00:00','已完成：准时送达');
INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103009',1009,2009,3009,4,'2026-11-03 10:40:00+00:00','已取消：用户改约');
INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103010',1010,2010,3010,4,'2026-11-03 10:55:00+00:00','已取消：无可用车辆');
INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103011',1011,2011,3011,5,'2026-11-03 11:10:00+00:00','异常：车辆故障');
INSERT INTO operation_dispatch (dispatch_no,driver_id,car_id,order_id,status,dispatch_time,remark) VALUES ('DD20261103012',1012,2012,3012,5,'2026-11-03 11:25:00+00:00','异常：超时未接单');
