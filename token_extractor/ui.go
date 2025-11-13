package token_extractor

import (
	"context"
	"fmt"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// TokenExtractorUI token提取器UI
type TokenExtractorUI struct {
	window    fyne.Window
	extractor Extractor
	storage   Storage

	// UI组件
	urlEntry      *widget.Entry
	usernameEntry *widget.Entry
	passwordEntry *widget.Entry
	extractButton *widget.Button
	statusLabel   *widget.Label
	resultList    *widget.List
	progressBar   *widget.ProgressBarInfinite

	// 数据
	currentResult *ExtractResult
}

// NewTokenExtractorUI 创建UI实例
func NewTokenExtractorUI(window fyne.Window) *TokenExtractorUI {
	extractor, _ := NewChromeExtractor()

	return &TokenExtractorUI{
		window:    window,
		extractor: extractor,
		storage:   NewJSONStorage("token_history.json"),
	}
}

// MakeUI 构建UI界面
func (ui *TokenExtractorUI) MakeUI() fyne.CanvasObject {
	// 标题
	title := widget.NewLabelWithStyle("🔐 网页Token提取器", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// URL输入框
	ui.urlEntry = widget.NewEntry()
	ui.urlEntry.SetPlaceHolder("https://ankersolix-professional-ci.anker.com/home/systemlist")
	ui.urlEntry.SetText("https://ankersolix-professional-ci.anker.com/home/systemlist")

	// 输入表单 - 使用紧凑的水平布局
	ui.usernameEntry = widget.NewEntry()
	ui.usernameEntry.SetPlaceHolder("账号")

	ui.passwordEntry = widget.NewPasswordEntry()
	ui.passwordEntry.SetPlaceHolder("密码")

	// 按钮
	ui.extractButton = widget.NewButton("开始提取", func() {
		ui.handleExtract()
	})

	// 进度条
	ui.progressBar = widget.NewProgressBarInfinite()
	ui.progressBar.Hide()

	// 状态标签
	ui.statusLabel = widget.NewLabel("等待操作...")
	ui.statusLabel.Alignment = fyne.TextAlignCenter

	// 输入区域 - 紧凑的网格布局
	urlLabel := widget.NewLabel("URL:")
	usernameLabel := widget.NewLabel("账号:")
	passwordLabel := widget.NewLabel("密码:")

	// 使用表格式布局，更紧凑
	inputGrid := container.NewVBox(
		container.NewBorder(nil, nil, urlLabel, nil, ui.urlEntry),
		container.NewHBox(
			container.NewBorder(nil, nil, usernameLabel, nil, ui.usernameEntry),
			container.NewBorder(nil, nil, passwordLabel, nil, ui.passwordEntry),
			ui.extractButton,
		),
	)

	inputSection := container.NewVBox(
		inputGrid,
		ui.progressBar,
		ui.statusLabel,
	)

	// 结果区域
	resultTitle := widget.NewLabelWithStyle("📊 提取结果:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	ui.resultList = widget.NewList(
		func() int {
			if ui.currentResult == nil {
				return 0
			}
			return len(ui.currentResult.Headers)
		},
		func() fyne.CanvasObject {
			// 列表项模板 - 更紧凑的单行布局
			iconLabel := widget.NewLabel("⭐")
			nameLabel := widget.NewLabelWithStyle("Header-Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			valueLabel := widget.NewLabel("value")
			valueLabel.Wrapping = fyne.TextWrapBreak
			copyBtn := widget.NewButton("复制", nil)

			// 使用水平布局，更紧凑
			return container.NewBorder(
				nil, nil, 
				container.NewHBox(iconLabel, nameLabel),
				copyBtn,
				valueLabel,
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if ui.currentResult == nil || id >= len(ui.currentResult.Headers) {
				return
			}

			header := ui.currentResult.Headers[id]
			border := obj.(*fyne.Container)

			// 更新图标和名称（左侧）
			leftBox := border.Objects[2].(*fyne.Container)
			iconLabel := leftBox.Objects[0].(*widget.Label)
			nameLabel := leftBox.Objects[1].(*widget.Label)

			if header.IsKey {
				iconLabel.SetText("⭐")
			} else {
				iconLabel.SetText("📋")
			}
			nameLabel.SetText(header.Name)

			// 更新值（中间）
			valueLabel := border.Objects[4].(*widget.Label)
			// 如果值太长，截断显示
			value := header.Value
			if len(value) > 80 {
				value = value[:77] + "..."
			}
			valueLabel.SetText(value)

			// 更新复制按钮（右侧）
			copyBtn := border.Objects[3].(*widget.Button)
			copyBtn.OnTapped = func() {
				ui.copyToClipboard(header.Value)
			}
		},
	)

	// 操作按钮
	copyAllBtn := widget.NewButton("复制所有关键Token", func() {
		ui.copyAllKeyTokens()
	})

	clearBtn := widget.NewButton("清空结果", func() {
		ui.currentResult = nil
		ui.resultList.Refresh()
		ui.statusLabel.SetText("结果已清空")
	})

	resultSection := container.NewVBox(
		resultTitle,
		widget.NewSeparator(),
		ui.resultList,
		container.NewHBox(copyAllBtn, clearBtn),
	)

	// 主布局
	mainContent := container.NewVBox(
		title,
		widget.NewSeparator(),
		inputSection,
		widget.NewSeparator(),
		resultSection,
	)

	return container.NewScroll(mainContent)
}

// handleExtract 处理提取操作
func (ui *TokenExtractorUI) handleExtract() {
	// 验证输入
	targetURL := ui.urlEntry.Text
	username := ui.usernameEntry.Text
	password := ui.passwordEntry.Text

	if targetURL == "" {
		dialog.ShowError(fmt.Errorf("请输入目标URL"), ui.window)
		return
	}

	if username == "" {
		dialog.ShowError(fmt.Errorf("请输入账号"), ui.window)
		return
	}

	if password == "" {
		dialog.ShowError(fmt.Errorf("请输入密码"), ui.window)
		return
	}

	// 禁用按钮，显示进度
	ui.extractButton.Disable()
	ui.progressBar.Show()
	ui.statusLabel.SetText("正在连接浏览器...")

	// 在goroutine中执行提取
	go func() {
		// 创建请求
		req := LoginRequest{
			Username:  username,
			Password:  password,
			TargetURL: targetURL,
		}

		// 更新状态
		ui.updateStatus("正在登录...")

		// 执行提取
		ctx := context.Background()
		result, err := ui.extractor.Extract(ctx, req)

		// 更新UI（必须在主线程）
		ui.window.Canvas().Content().Refresh()

		if err != nil {
			ui.progressBar.Hide()
			ui.extractButton.Enable()
			ui.statusLabel.SetText(fmt.Sprintf("❌ 提取失败: %v", err))
			dialog.ShowError(err, ui.window)
			return
		}

		// 显示结果
		ui.displayResult(result)

		// 保存历史（可选）
		if result.Success {
			ui.saveHistory(username, result)
		}

		ui.progressBar.Hide()
		ui.extractButton.Enable()
	}()
}

// updateStatus 更新状态（线程安全）
func (ui *TokenExtractorUI) updateStatus(status string) {
	ui.statusLabel.SetText(status)
	ui.statusLabel.Refresh()
}

// displayResult 显示提取结果
func (ui *TokenExtractorUI) displayResult(result *ExtractResult) {
	ui.currentResult = result

	if !result.Success {
		ui.statusLabel.SetText(fmt.Sprintf("❌ 提取失败: %s", result.Error))
		return
	}

	// 排序：关键头部在前
	sort.Slice(result.Headers, func(i, j int) bool {
		if result.Headers[i].IsKey != result.Headers[j].IsKey {
			return result.Headers[i].IsKey
		}
		return result.Headers[i].Name < result.Headers[j].Name
	})

	ui.resultList.Refresh()
	ui.statusLabel.SetText(fmt.Sprintf("✅ 提取成功 (%s) - 共捕获 %d 个请求头",
		result.Timestamp.Format("2006-01-02 15:04:05"),
		len(result.Headers)))
}

// copyToClipboard 复制到剪贴板
func (ui *TokenExtractorUI) copyToClipboard(text string) {
	ui.window.Clipboard().SetContent(text)
	dialog.ShowInformation("成功", "已复制到剪贴板", ui.window)
}

// copyAllKeyTokens 复制所有关键token
func (ui *TokenExtractorUI) copyAllKeyTokens() {
	if ui.currentResult == nil {
		dialog.ShowError(fmt.Errorf("没有可复制的结果"), ui.window)
		return
	}

	var keyTokens string
	for _, header := range ui.currentResult.Headers {
		if header.IsKey {
			keyTokens += fmt.Sprintf("%s: %s\n", header.Name, header.Value)
		}
	}

	if keyTokens == "" {
		dialog.ShowError(fmt.Errorf("未找到关键Token"), ui.window)
		return
	}

	ui.window.Clipboard().SetContent(keyTokens)
	dialog.ShowInformation("成功", "所有关键Token已复制到剪贴板", ui.window)
}

// saveHistory 保存历史记录
func (ui *TokenExtractorUI) saveHistory(username string, result *ExtractResult) {
	if ui.storage == nil {
		return
	}

	// 提取关键头部
	keyHeaders := make(map[string]string)
	for _, header := range result.Headers {
		if header.IsKey {
			// 脱敏处理：只保存前后几位
			value := header.Value
			if len(value) > 20 {
				value = value[:8] + "..." + value[len(value)-8:]
			}
			keyHeaders[header.Name] = value
		}
	}

	record := HistoryRecord{
		ID:         fmt.Sprintf("%d", time.Now().Unix()),
		Timestamp:  result.Timestamp,
		Username:   username,
		Success:    result.Success,
		KeyHeaders: keyHeaders,
	}

	ui.storage.SaveHistory(record)
}
