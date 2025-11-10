# 实现计划

- [x] 1. 创建项目结构和数据模型
  - 创建 `weight_tracker` 目录
  - 实现 `WeightRecord` 结构体，包含 ID、Weight、Date、Change 和 ChangeType 字段
  - 实现 `CalculateChange` 函数，根据当前和上一条记录计算变化值和类型
  - 实现 `FormatChange` 方法，返回格式化的变化显示文本（如 "↑ +0.5 kg"）
  - 实现 `FormatDate` 方法，返回格式化的日期字符串（如 "2025-11-10 10:30"）
  - _需求: 1.5, 2.1, 2.2, 2.3, 2.4, 2.5_

- [x] 2. 实现数据持久化层
  - 定义 `Storage` 接口，包含 Load 和 Save 方法
  - 实现 `JSONStorage` 结构体，使用 `weight_records.json` 作为存储文件
  - 实现 `Load` 方法，从 JSON 文件读取记录列表，处理文件不存在的情况
  - 实现 `Save` 方法，将记录列表写入 JSON 文件
  - 添加错误处理，确保文件读写失败时有适当的错误返回
  - _需求: 4.1, 4.2, 4.3, 4.4_

- [x] 3. 构建用户界面组件
  - 创建 `WeightTrackerUI` 结构体，包含 storage、records、weightEntry 和 recordList 字段
  - 实现 `NewWeightTrackerUI` 构造函数，初始化存储和加载现有记录
  - 实现输入区域：创建 Entry 组件（带占位符）和"添加记录"按钮
  - 实现记录列表：使用 `widget.List` 显示所有记录，包含日期、体重和变化标识
  - 实现列表项渲染，根据 ChangeType 显示不同的视觉标识（↑/↓/●）和颜色
  - 使用 `container.NewBorder` 布局，将输入区域放在顶部，列表放在中间
  - _需求: 1.1, 1.2, 3.1, 3.2, 3.3, 3.4, 5.3_

- [x] 4. 实现添加记录功能
  - 实现 `addRecord` 方法，从输入框获取体重值
  - 添加输入验证：检查空输入、非数字、负数和超出范围的值
  - 使用 `dialog.ShowError` 显示验证错误信息
  - 创建新的 `WeightRecord`，生成 UUID 作为 ID，使用当前时间作为日期
  - 如果存在上一条记录，调用 `CalculateChange` 计算变化
  - 将新记录插入到列表开头（保持倒序）
  - 调用 `saveRecords` 保存到文件
  - 刷新 UI 列表显示
  - 清空输入框
  - _需求: 1.3, 1.4, 1.5, 2.1, 2.2, 2.3, 2.4, 2.5, 4.3_

- [x] 5. 实现数据加载和保存
  - 实现 `loadRecords` 方法，在 UI 初始化时调用 storage.Load()
  - 处理加载错误，如果文件不存在则使用空列表
  - 实现 `saveRecords` 方法，调用 storage.Save() 保存当前记录列表
  - 处理保存错误，使用 `dialog.ShowError` 提示用户
  - 确保记录按时间倒序排列
  - _需求: 3.2, 4.1, 4.2, 4.3_

- [x] 6. 集成到主应用
  - 在 `main.go` 中导入 `weight_tracker` 包
  - 创建 `WeightTrackerUI` 实例
  - 调用 `MakeUI()` 获取体重记录界面组件
  - 在 `TabContainer` 中添加新的 `TabItem`，标签名为"体重记录"
  - 验证新标签页与其他工具独立运行
  - _需求: 5.1, 5.2, 5.4_

- [x] 7. 处理空状态和边界情况
  - 当记录列表为空时，显示友好的提示信息（如"还没有记录，添加第一条吧！"）
  - 确保第一条记录正确标记为 "首次记录"
  - 测试快速连续添加多条记录的情况
  - 验证列表滚动功能正常工作
  - _需求: 2.5, 3.5_

- [ ]* 8. 编写单元测试
  - 为 `CalculateChange` 函数编写测试用例，覆盖增加、减少、持平场景
  - 为 `FormatChange` 和 `FormatDate` 编写测试用例
  - 为 `JSONStorage` 的 Load 和 Save 方法编写测试
  - 测试空文件和损坏文件的处理
  - _需求: 所有需求的验证_
