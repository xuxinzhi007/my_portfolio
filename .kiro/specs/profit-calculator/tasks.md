# 实施计划

- [x] 1. 创建项目结构和数据模型
  - 创建 `profit_calculator` 目录
  - 实现 `model.go` 中的核心数据结构：`Investor`、`MonthlyProfit`、`ProfitCalculatorData`
  - 实现投资比例计算函数 `CalculateInvestmentRatio`
  - 实现总投资计算函数 `CalculateTotalInvestment`
  - 实现收益分配函数 `DistributeProfit`
  - _需求: 1.3, 1.4, 3.1, 3.5_

- [x] 2. 实现统计计算功能
  - 实现 `InvestorStats` 和 `OverallStats` 结构体
  - 实现 `CalculateInvestorStats` 函数计算单个投资者统计
  - 实现 `CalculateOverallStats` 函数计算整体统计
  - _需求: 3.2, 3.3, 3.4, 5.3_

- [x] 3. 实现数据持久化层
  - 创建 `storage.go` 文件
  - 定义 `Storage` 接口
  - 实现 `JSONStorage` 结构体
  - 实现 `Load()` 方法从 JSON 文件加载数据
  - 实现 `Save()` 方法保存数据到 JSON 文件
  - 处理文件不存在和空文件的情况
  - _需求: 4.1, 4.2, 4.3, 4.4_

- [x] 4. 创建基础 UI 框架
  - 创建 `ui.go` 文件
  - 实现 `ProfitCalculatorUI` 结构体
  - 实现 `NewProfitCalculatorUI` 构造函数
  - 实现 `MakeUI()` 方法创建主界面布局
  - 实现 `loadData()` 和 `saveData()` 方法
  - _需求: 5.4_

- [x] 5. 实现统计卡片显示
  - 实现 `createStatsCard()` 方法创建统计卡片
  - 添加总投资、总收益、投资者数量的显示组件
  - 实现 `updateStats()` 方法更新统计显示
  - 使用颜色区分不同类型的统计数据
  - _需求: 5.3_

- [x] 6. 实现投资者管理功能
- [x] 6.1 创建投资者列表界面
  - 实现 `createInvestorSection()` 方法
  - 实现 `createInvestorList()` 方法创建投资者列表
  - 显示投资者姓名、投资金额、投资比例、累计收益和最终金额
  - 添加编辑和删除按钮
  - _需求: 5.1_

- [x] 6.2 实现添加投资者功能
  - 实现 `showAddInvestorDialog()` 方法显示添加对话框
  - 添加姓名和金额输入验证（非空、金额 > 0、合理范围）
  - 检查姓名重复
  - 创建新投资者并保存
  - 刷新列表和统计
  - _需求: 1.1, 1.2_

- [x] 6.3 实现编辑投资者功能
  - 实现 `showEditInvestorDialog()` 方法
  - 预填充现有投资者信息
  - 验证修改后的数据
  - 更新投资者信息并保存
  - 重新计算投资比例
  - _需求: 1.5_

- [x] 6.4 实现删除投资者功能
  - 实现 `deleteInvestor()` 方法
  - 显示确认对话框
  - 删除投资者数据
  - 重新计算剩余投资者的比例
  - 保留历史收益记录
  - _需求: 6.1, 6.2, 6.3, 6.4_

- [x] 7. 实现月度收益管理功能
- [x] 7.1 创建收益记录列表界面
  - 实现 `createProfitSection()` 方法
  - 实现 `createProfitList()` 方法创建收益列表
  - 显示日期、总收益金额
  - 添加查看详情和删除按钮
  - _需求: 2.3, 5.2_

- [x] 7.2 实现添加收益记录功能
  - 实现 `showAddProfitDialog()` 方法
  - 添加日期选择器和金额输入
  - 验证输入（日期不能为未来、金额在合理范围）
  - 检查是否有投资者（无投资者时提示）
  - 调用 `DistributeProfit()` 自动计算分配
  - 创建收益记录并保存
  - 刷新列表和统计
  - _需求: 2.1, 2.2, 3.1_

- [x] 7.3 实现收益详情查看功能
  - 实现 `showProfitDetailDialog()` 方法
  - 显示收益日期和总金额
  - 显示每个投资者的分配明细（姓名、金额、比例）
  - 使用表格或列表清晰展示
  - _需求: 3.2, 5.2_

- [x] 7.4 实现删除收益记录功能
  - 实现 `deleteProfitRecord()` 方法
  - 显示确认对话框
  - 删除收益记录
  - 更新统计信息
  - _需求: 2.4_

- [x] 8. 实现数据格式化和显示优化
  - 实现货币金额格式化（保留两位小数）
  - 实现百分比格式化（投资比例）
  - 实现日期格式化
  - 添加空状态提示（无投资者、无收益记录时）
  - 优化列表项布局和间距
  - _需求: 5.5_

- [x] 9. 集成到主应用
  - 修改 `main.go` 导入 `profit_calculator` 包
  - 创建 `ProfitCalculatorUI` 实例
  - 添加"收益计算"标签页到主界面
  - 使用合适的图标（`theme.ConfirmIcon()` 或自定义）
  - 测试标签页切换和功能正常运行
  - _需求: 所有需求的集成验证_

- [ ]* 10. 编写单元测试
  - 为 `model.go` 编写测试：测试投资比例计算、收益分配、统计计算
  - 为 `storage.go` 编写测试：测试保存、加载、错误处理
  - 测试边界情况：零投资、无投资者、负收益等
  - _需求: 3.5, 4.4_
