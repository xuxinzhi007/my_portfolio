package settings

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SettingsUI 设置界面
type SettingsUI struct {
	app    fyne.App
	window fyne.Window
}

// NewSettingsUI 创建设置界面
func NewSettingsUI(app fyne.App, window fyne.Window) *SettingsUI {
	return &SettingsUI{
		app:    app,
		window: window,
	}
}

// MakeUI 构建设置界面
func (s *SettingsUI) MakeUI() fyne.CanvasObject {
	// 标题
	title := widget.NewLabelWithStyle("⚙️ 应用设置", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	titleText := canvas.NewText("自定义您的应用体验", color.RGBA{R: 128, G: 128, B: 128, A: 255})
	titleText.Alignment = fyne.TextAlignCenter
	titleText.TextSize = 12

	// 主题设置
	themeCard := s.createThemeCard()

	// 关于信息
	aboutCard := s.createAboutCard()

	// 组合所有卡片
	content := container.NewVBox(
		title,
		titleText,
		widget.NewSeparator(),
		themeCard,
		widget.NewSeparator(),
		aboutCard,
		layout.NewSpacer(),
	)

	return container.NewScroll(content)
}

// createThemeCard 创建主题设置卡片
func (s *SettingsUI) createThemeCard() fyne.CanvasObject {
	// 卡片标题
	cardTitle := widget.NewLabelWithStyle("🎨 主题设置", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	cardDesc := widget.NewLabel("选择您喜欢的应用主题")
	cardDesc.TextStyle = fyne.TextStyle{Italic: true}

	// 当前主题显示
	currentThemeIcon := widget.NewIcon(theme.ColorPaletteIcon())
	currentThemeLabel := widget.NewLabel("当前主题：")
	currentThemeValue := widget.NewLabel(s.getCurrentThemeName())
	currentThemeValue.TextStyle = fyne.TextStyle{Bold: true}

	// 主题选项 - 使用可用的图标
	lightButton := widget.NewButtonWithIcon("☀️ 浅色主题", theme.ContentClearIcon(), func() {
		s.app.Settings().SetTheme(theme.LightTheme())
		currentThemeValue.SetText("浅色")
		dialog.ShowInformation("✅ 主题已更改", "已切换到浅色主题", s.window)
	})
	lightButton.Importance = widget.LowImportance

	darkButton := widget.NewButtonWithIcon("🌙 深色主题", theme.VisibilityIcon(), func() {
		s.app.Settings().SetTheme(theme.DarkTheme())
		currentThemeValue.SetText("深色")
		dialog.ShowInformation("✅ 主题已更改", "已切换到深色主题", s.window)
	})
	darkButton.Importance = widget.LowImportance

	// 主题预览
	previewLabel := widget.NewLabelWithStyle("预览效果", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	previewBox := s.createThemePreview()

	// 布局
	card := container.NewVBox(
		cardTitle,
		cardDesc,
		widget.NewSeparator(),
		container.NewHBox(currentThemeIcon, currentThemeLabel, currentThemeValue),
		container.NewGridWithColumns(2, lightButton, darkButton),
		widget.NewSeparator(),
		previewLabel,
		previewBox,
	)

	return card
}

// createAboutCard 创建关于信息卡片
func (s *SettingsUI) createAboutCard() fyne.CanvasObject {
	// 卡片标题
	cardTitle := widget.NewLabelWithStyle("ℹ️ 关于应用", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// 应用信息 - 添加图标
	appIcon := widget.NewIcon(theme.HomeIcon())
	appName := widget.NewLabel("我的超级工具箱")
	appName.TextStyle = fyne.TextStyle{Bold: true}
	
	versionIcon := widget.NewIcon(theme.InfoIcon())
	appVersion := widget.NewLabel("版本 1.0.0")
	
	fyneIcon := widget.NewIcon(theme.ComputerIcon())
	fyneVersion := widget.NewLabel("Fyne v2.7.0")

	// 功能列表
	featuresLabel := widget.NewLabelWithStyle("✨ 包含功能", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	
	feature1 := container.NewHBox(widget.NewIcon(theme.MediaRecordIcon()), widget.NewLabel("体重记录追踪（带统计和动画）"))
	feature2 := container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("周报生成器"))
	feature3 := container.NewHBox(widget.NewIcon(theme.FolderIcon()), widget.NewLabel("文件去重工具"))
	feature4 := container.NewHBox(widget.NewIcon(theme.ColorPaletteIcon()), widget.NewLabel("主题切换"))
	feature5 := container.NewHBox(widget.NewIcon(theme.ContentAddIcon()), widget.NewLabel("更多功能开发中..."))

	// 系统信息
	sysInfoLabel := widget.NewLabelWithStyle("💻 系统信息", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	
	dpi := s.window.Canvas().Scale()
	dpiIcon := widget.NewIcon(theme.ZoomInIcon())
	sysInfo := widget.NewLabel("DPI 缩放：" + s.getScaleText(dpi))

	// 重置按钮
	resetButton := widget.NewButtonWithIcon("🔄 重置主题", theme.ViewRefreshIcon(), func() {
		dialog.ShowConfirm("⚠️ 确认重置", "确定要重置主题到默认值吗？", func(confirmed bool) {
			if confirmed {
				s.app.Settings().SetTheme(theme.DefaultTheme())
				dialog.ShowInformation("✅ 重置完成", "主题已恢复默认", s.window)
			}
		}, s.window)
	})
	resetButton.Importance = widget.DangerImportance

	// 布局
	card := container.NewVBox(
		cardTitle,
		widget.NewSeparator(),
		container.NewHBox(appIcon, appName),
		container.NewHBox(versionIcon, appVersion),
		container.NewHBox(fyneIcon, fyneVersion),
		widget.NewSeparator(),
		featuresLabel,
		feature1,
		feature2,
		feature3,
		feature4,
		feature5,
		widget.NewSeparator(),
		sysInfoLabel,
		container.NewHBox(dpiIcon, sysInfo),
		widget.NewSeparator(),
		resetButton,
	)

	return card
}

// getScaleText 获取缩放文本描述
func (s *SettingsUI) getScaleText(scale float32) string {
	percentage := int(scale * 100)
	return string(rune(percentage/100+'0')) + string(rune((percentage/10)%10+'0')) + string(rune(percentage%10+'0')) + "%"
}

// createThemePreview 创建主题预览
func (s *SettingsUI) createThemePreview() fyne.CanvasObject {
	// 创建一些示例组件
	sampleLabel := widget.NewLabel("示例文本")
	sampleButton := widget.NewButton("示例按钮", func() {})
	sampleEntry := widget.NewEntry()
	sampleEntry.SetPlaceHolder("示例输入框")
	sampleCheck := widget.NewCheck("示例复选框", func(bool) {})

	preview := container.NewVBox(
		sampleLabel,
		sampleButton,
		sampleEntry,
		sampleCheck,
	)

	// 添加边框
	bordered := container.NewPadded(preview)

	return bordered
}

// getCurrentThemeName 获取当前主题名称
func (s *SettingsUI) getCurrentThemeName() string {
	// Fyne 没有直接的方法获取当前主题名称
	// 我们可以通过检查背景色来判断
	bg := theme.BackgroundColor()
	
	// 简单判断：深色背景认为是深色主题
	r, g, b, _ := bg.RGBA()
	brightness := (r + g + b) / 3
	
	if brightness < 32768 { // 16位色值的一半
		return "深色"
	}
	return "浅色"
}

// ShowSettingsDialog 显示设置对话框
func ShowSettingsDialog(app fyne.App, window fyne.Window) {
	settingsUI := NewSettingsUI(app, window)
	content := settingsUI.MakeUI()
	
	// 创建自定义对话框
	settingsDialog := dialog.NewCustom("设置", "关闭", content, window)
	settingsDialog.Resize(fyne.NewSize(500, 600))
	settingsDialog.Show()
}
