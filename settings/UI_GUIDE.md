# UI 组件使用指南

## 🎨 可用的 Fyne 主题图标

Fyne 提供了丰富的内置图标，以下是常用的图标列表：

### 📁 文件和文档
- `theme.DocumentIcon()` - 文档
- `theme.DocumentCreateIcon()` - 创建文档
- `theme.DocumentPrintIcon()` - 打印
- `theme.DocumentSaveIcon()` - 保存
- `theme.FileIcon()` - 文件
- `theme.FolderIcon()` - 文件夹
- `theme.FolderNewIcon()` - 新建文件夹
- `theme.FolderOpenIcon()` - 打开文件夹

### 🎬 媒体控制
- `theme.MediaPlayIcon()` - 播放
- `theme.MediaPauseIcon()` - 暂停
- `theme.MediaStopIcon()` - 停止
- `theme.MediaRecordIcon()` - 录制
- `theme.MediaReplayIcon()` - 重播
- `theme.MediaSkipNextIcon()` - 下一个
- `theme.MediaSkipPreviousIcon()` - 上一个
- `theme.MediaFastForwardIcon()` - 快进
- `theme.MediaFastRewindIcon()` - 快退
- `theme.MediaVideoIcon()` - 视频
- `theme.MediaMusicIcon()` - 音乐
- `theme.MediaPhotoIcon()` - 照片

### ✏️ 编辑操作
- `theme.ContentAddIcon()` - 添加
- `theme.ContentRemoveIcon()` - 移除
- `theme.ContentClearIcon()` - 清除
- `theme.ContentCopyIcon()` - 复制
- `theme.ContentCutIcon()` - 剪切
- `theme.ContentPasteIcon()` - 粘贴
- `theme.ContentRedoIcon()` - 重做
- `theme.ContentUndoIcon()` - 撤销
- `theme.DeleteIcon()` - 删除

### 🔍 导航和视图
- `theme.NavigateBackIcon()` - 返回
- `theme.NavigateNextIcon()` - 前进
- `theme.ZoomInIcon()` - 放大
- `theme.ZoomOutIcon()` - 缩小
- `theme.ZoomFitIcon()` - 适应
- `theme.ViewFullScreenIcon()` - 全屏
- `theme.ViewRestoreIcon()` - 还原
- `theme.ViewRefreshIcon()` - 刷新
- `theme.VisibilityIcon()` - 可见
- `theme.VisibilityOffIcon()` - 不可见

### ⚙️ 系统和设置
- `theme.SettingsIcon()` - 设置
- `theme.InfoIcon()` - 信息
- `theme.QuestionIcon()` - 问题
- `theme.WarningIcon()` - 警告
- `theme.ErrorIcon()` - 错误
- `theme.ConfirmIcon()` - 确认
- `theme.CancelIcon()` - 取消
- `theme.CheckButtonIcon()` - 选中
- `theme.CheckButtonCheckedIcon()` - 已选中
- `theme.RadioButtonIcon()` - 单选
- `theme.RadioButtonCheckedIcon()` - 已选单选

### 🏠 应用和界面
- `theme.HomeIcon()` - 主页
- `theme.ComputerIcon()` - 计算机
- `theme.StorageIcon()` - 存储
- `theme.DownloadIcon()` - 下载
- `theme.UploadIcon()` - 上传
- `theme.SearchIcon()` - 搜索
- `theme.SearchReplaceIcon()` - 搜索替换
- `theme.MenuIcon()` - 菜单
- `theme.MenuExpandIcon()` - 展开菜单
- `theme.MenuDropDownIcon()` - 下拉菜单
- `theme.MenuDropUpIcon()` - 上拉菜单

### 🎨 其他
- `theme.ColorPaletteIcon()` - 调色板
- `theme.ColorChromaticIcon()` - 色彩
- `theme.ColorAchromaticIcon()` - 灰度
- `theme.HistoryIcon()` - 历史
- `theme.MailAttachmentIcon()` - 附件
- `theme.MailComposeIcon()` - 撰写邮件
- `theme.MailForwardIcon()` - 转发
- `theme.MailReplyIcon()` - 回复
- `theme.MailReplyAllIcon()` - 全部回复
- `theme.MailSendIcon()` - 发送

## 🎯 自定义 UI 组件

我们创建了一些自定义组件来增强 UI：

### CustomCard - 自定义卡片
```go
card := settings.CustomCard(
    theme.InfoIcon(),
    "标题",
    content,
    color.RGBA{R: 33, G: 150, B: 243, A: 255},
)
```

### InfoRow - 信息行
```go
row := settings.InfoRow(theme.ComputerIcon(), "这是一条信息")
```

### StatsCard - 统计卡片
```go
stats := settings.StatsCard(
    theme.MediaRecordIcon(),
    "记录数",
    "42",
    color.RGBA{R: 76, G: 175, B: 80, A: 255},
)
```

### FeatureItem - 功能列表项
```go
feature := settings.FeatureItem(
    theme.DocumentIcon(),
    "功能标题",
    "功能描述文本",
)
```

### SectionHeader - 章节标题
```go
header := settings.SectionHeader("章节标题", theme.InfoIcon())
```

## 🌈 推荐的颜色方案

### Material Design 颜色
```go
// 蓝色系
color.RGBA{R: 33, G: 150, B: 243, A: 255}   // 主蓝色
color.RGBA{R: 25, G: 118, B: 210, A: 255}   // 深蓝色
color.RGBA{R: 100, G: 181, B: 246, A: 255}  // 浅蓝色

// 绿色系
color.RGBA{R: 76, G: 175, B: 80, A: 255}    // 主绿色
color.RGBA{R: 67, G: 160, B: 71, A: 255}    // 深绿色
color.RGBA{R: 129, G: 199, B: 132, A: 255}  // 浅绿色

// 红色系
color.RGBA{R: 244, G: 67, B: 54, A: 255}    // 主红色
color.RGBA{R: 229, G: 57, B: 53, A: 255}    // 深红色
color.RGBA{R: 239, G: 154, B: 154, A: 255}  // 浅红色

// 橙色系
color.RGBA{R: 255, G: 152, B: 0, A: 255}    // 主橙色
color.RGBA{R: 251, G: 140, B: 0, A: 255}    // 深橙色
color.RGBA{R: 255, G: 183, B: 77, A: 255}   // 浅橙色

// 紫色系
color.RGBA{R: 156, G: 39, B: 176, A: 255}   // 主紫色
color.RGBA{R: 142, G: 36, B: 170, A: 255}   // 深紫色
color.RGBA{R: 186, G: 104, B: 200, A: 255}  // 浅紫色

// 灰色系
color.RGBA{R: 158, G: 158, B: 158, A: 255}  // 中灰色
color.RGBA{R: 97, G: 97, B: 97, A: 255}     // 深灰色
color.RGBA{R: 224, G: 224, B: 224, A: 255}  // 浅灰色
```

## 🎨 UI 设计最佳实践

### 1. 使用一致的间距
```go
container.NewPadded(content)  // 标准内边距
widget.NewSeparator()         // 分隔线
layout.NewSpacer()            // 弹性空间
```

### 2. 合理使用按钮重要性
```go
btn.Importance = widget.HighImportance    // 主要操作（蓝色）
btn.Importance = widget.MediumImportance  // 次要操作（默认）
btn.Importance = widget.LowImportance     // 辅助操作（灰色）
btn.Importance = widget.DangerImportance  // 危险操作（红色）
btn.Importance = widget.WarningImportance // 警告操作（橙色）
btn.Importance = widget.SuccessImportance // 成功操作（绿色）
```

### 3. 使用图标增强识别
- 每个功能都应该有对应的图标
- 图标应该与功能语义相关
- 保持图标使用的一致性

### 4. 文本样式
```go
// 标题
widget.NewLabelWithStyle("标题", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

// 副标题
label := widget.NewLabel("副标题")
label.TextStyle = fyne.TextStyle{Italic: true}

// 强调文本
label.TextStyle = fyne.TextStyle{Bold: true}
```

### 5. 布局选择
```go
// 垂直布局
container.NewVBox(widget1, widget2, widget3)

// 水平布局
container.NewHBox(widget1, widget2, widget3)

// 网格布局
container.NewGridWithColumns(2, widget1, widget2, widget3, widget4)

// 边框布局
container.NewBorder(top, bottom, left, right, center)

// 堆叠布局
container.NewStack(background, foreground)

// 滚动容器
container.NewScroll(content)
```

## 🚀 性能优化建议

1. **避免频繁刷新** - 只在必要时调用 `Refresh()`
2. **使用虚拟列表** - 对于大量数据使用 `widget.List`
3. **延迟加载** - 复杂内容可以延迟创建
4. **减少嵌套** - 避免过深的容器嵌套
5. **复用组件** - 相同的组件可以复用

## 📱 响应式设计

```go
// 根据窗口大小调整布局
size := window.Canvas().Size()
if size.Width < 600 {
    // 小屏幕布局
    return container.NewVBox(...)
} else {
    // 大屏幕布局
    return container.NewHBox(...)
}
```

## 🎭 动画效果

Fyne 的动画支持有限，但可以通过以下方式实现：

1. **渐变效果** - 修改透明度
2. **位置动画** - 修改组件位置
3. **大小动画** - 修改组件大小
4. **颜色动画** - 修改颜色值

示例：
```go
go func() {
    for i := 0; i < 10; i++ {
        alpha := uint8(255 * i / 10)
        text.Color = color.RGBA{R: 255, G: 0, B: 0, A: alpha}
        text.Refresh()
        time.Sleep(50 * time.Millisecond)
    }
}()
```
