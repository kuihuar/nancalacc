package biz

// 优化全量同步任务，internal/biz/full_sync.go的CreateSyncAccount方法，
// 1. 验证是否提交过uc.localCache，否则返回错误
// 2. 获取token并从第三方获取部门和用户数据
// 3. 数据入库，包括SaveCompanyCfg公司配置入库，SaveDepartments部门入库，SaveUsers用户入库，SaveDepartmentUserRelations关系入库
// 4. 调用wps接口，通知调用 PostEcisaccountsyncAll 开始同步
// 5. 更新任务状态uc.localCache ，包括任务状态，进度，开始时间，结束时间，实际时间
// 6. 返回结果，包括任务id

// 其中3里面的的三项可以是并发，但需要保证数据一致性，
// 2里面已经有并发，但需要保证数据一致性，
// 确保可以入库完成，再有通知调用
