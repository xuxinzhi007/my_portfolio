# Implementation Plan - Web Token Extractor

## Task List

- [x] 1. 设置项目结构和依赖
  - 创建token_extractor包目录结构
  - 添加chromedp依赖到go.mod
  - 创建基础文件框架（model.go, extractor.go, ui.go, storage.go）
  - _Requirements: 1.1, 1.2, 2.1, 3.1_

- [x] 2. 实现数据模型层
  - [x] 2.1 创建核心数据结构
    - 实现LoginRequest结构体及验证方法
    - 实现HeaderInfo结构体
    - 实现ExtractResult结构体
    - 实现HistoryRecord结构体
    - _Requirements: 1.1, 1.4, 3.1_

  - [x] 2.2 实现数据验证逻辑
    - 实现LoginRequest.Validate()方法
    - 添加URL格式验证
    - 添加字段长度验证
    - _Requirements: 1.2, 1.4_

- [x] 3. 实现浏览器自动化提取服务
  - [x] 3.1 创建ChromeExtractor基础结构
    - 定义Extractor接口
    - 实现ChromeExtractor结构体
    - 实现NewChromeExtractor构造函数
    - 实现Close方法用于资源清理
    - _Requirements: 2.1, 2.2_

  - [x] 3.2 实现浏览器初始化逻辑
    - 配置chromedp上下文
    - 设置headless模式
    - 配置超时参数
    - 添加浏览器启动错误处理
    - _Requirements: 2.1, 2.3, 5.2_

  - [x] 3.3 实现自动登录流程
    - 实现页面导航逻辑
    - 实现表单元素定位和填写
    - 实现登录按钮点击
    - 实现登录成功验证
    - 添加登录失败检测
    - _Requirements: 2.1, 2.2, 2.3, 5.2_

  - [x] 3.4 实现HTTP请求头捕获
    - 使用chromedp网络事件监听
    - 捕获认证后的请求头
    - 过滤和识别关键头部（X-Auth-Token, X-Auth-Ts, Gtoken等）
    - 存储捕获的头部信息
    - _Requirements: 3.1, 3.2, 3.3, 3.5_

  - [x] 3.5 实现错误处理和重试逻辑
    - 定义错误类型常量
    - 实现超时处理
    - 实现网络错误处理
    - 添加重试机制（最多2次）
    - _Requirements: 2.3, 5.3, 5.4_

- [x] 4. 实现存储层（可选功能）
  - [x] 4.1 创建Storage接口和JSONStorage实现
    - 定义Storage接口
    - 实现JSONStorage结构体
    - 实现NewJSONStorage构造函数
    - _Requirements: 3.1_

  - [x] 4.2 实现历史记录保存和读取
    - 实现SaveHistory方法
    - 实现GetHistory方法
    - 实现ClearHistory方法
    - 添加文件操作错误处理
    - 实现token脱敏逻辑（安全考虑）
    - _Requirements: 3.1_

- [x] 5. 实现UI层
  - [x] 5.1 创建TokenExtractorUI基础结构
    - 实现TokenExtractorUI结构体
    - 实现NewTokenExtractorUI构造函数
    - 初始化extractor和storage实例
    - _Requirements: 1.1, 3.1, 5.1_

  - [x] 5.2 实现输入表单UI
    - 创建目标URL显示标签
    - 创建用户名输入框
    - 创建密码输入框（带遮罩）
    - 创建提取按钮
    - 添加输入验证提示
    - _Requirements: 1.1, 1.2, 1.3, 1.4_

  - [x] 5.3 实现结果展示UI
    - 创建结果列表组件
    - 实现关键头部高亮显示（使用⭐标记）
    - 为每个头部添加复制按钮
    - 实现"复制所有"按钮
    - 添加清空结果按钮
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 4.1, 4.2_

  - [x] 5.4 实现状态和进度显示
    - 创建状态标签
    - 添加进度条（无限进度条）
    - 实现状态消息更新逻辑
    - 显示操作时间戳
    - _Requirements: 5.1, 5.2, 5.3_

  - [x] 5.5 实现复制到剪贴板功能
    - 实现copyToClipboard方法
    - 添加复制成功提示
    - 实现单个头部复制
    - 实现批量复制（所有关键token）
    - _Requirements: 4.1, 4.2, 4.3, 4.4_

  - [x] 5.6 实现提取操作处理逻辑
    - 实现handleExtract方法
    - 添加输入验证
    - 在goroutine中执行提取操作（避免UI阻塞）
    - 实现displayResult方法显示结果
    - 添加错误对话框显示
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 5.1, 5.2, 5.3, 5.4_

  - [x] 5.7 实现MakeUI方法组装完整界面
    - 使用container.NewVBox组织布局
    - 添加分隔线
    - 设置合适的间距和对齐
    - 返回可滚动的容器
    - _Requirements: 1.1, 3.1, 5.1_

- [x] 6. 集成到主应用
  - [x] 6.1 在main.go中添加新Tab
    - 导入token_extractor包
    - 创建TokenExtractorUI实例
    - 调用MakeUI方法获取UI内容
    - 在TabContainer中添加新的Tab项（使用合适的图标）
    - _Requirements: 1.1, 3.1, 5.1_

  - [x] 6.2 测试集成效果
    - 验证Tab切换正常
    - 验证UI布局正确
    - 验证与其他Tab不冲突
    - _Requirements: 1.1, 3.1_

- [x] 7. 端到端功能测试
  - [x] 7.1 测试完整提取流程
    - 使用真实账号测试登录
    - 验证token成功提取
    - 验证提取的token有效性
    - 测试复制功能
    - _Requirements: 1.1, 1.2, 2.1, 2.2, 3.1, 3.2, 3.3, 4.1_

  - [x] 7.2 测试错误场景
    - 测试无效凭证处理
    - 测试网络错误处理
    - 测试超时处理
    - 验证错误消息显示正确
    - _Requirements: 1.4, 2.3, 5.3, 5.4_

  - [x] 7.3 测试UI交互
    - 测试所有按钮功能
    - 测试输入验证
    - 测试状态更新
    - 测试不同窗口尺寸下的显示
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 4.1, 4.2, 4.3, 5.1, 5.2_

- [ ]* 8. 优化和完善
  - [ ]* 8.1 性能优化
    - 优化浏览器启动速度
    - 禁用不必要的浏览器功能（图片、CSS加载）
    - 实现浏览器实例复用（如果需要）
    - _Requirements: 2.1, 5.1_

  - [ ]* 8.2 添加历史记录功能
    - 实现查看历史按钮
    - 创建历史记录对话框
    - 实现历史记录列表显示
    - 添加清空历史功能
    - _Requirements: 3.1_

  - [ ]* 8.3 增强安全性
    - 实现密码内存清除
    - 实现token脱敏显示
    - 添加文件权限设置
    - 验证SSL证书
    - _Requirements: 1.3, 3.1_

  - [ ]* 8.4 添加配置选项
    - 实现超时时间配置
    - 实现重试次数配置
    - 添加代理设置（可选）
    - _Requirements: 5.4_
