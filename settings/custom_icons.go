package settings

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CustomCard 创建带有图标和颜色的卡片
func CustomCard(icon fyne.Resource, title string, content fyne.CanvasObject, accentColor color.Color) fyne.CanvasObject {
	// 创建标题栏
	iconWidget := widget.NewIcon(icon)
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	
	// 创建彩色分隔线
	separator := canvas.NewRectangle(accentColor)
	separator.SetMinSize(fyne.NewSize(0, 2))
	
	// 组合卡片
	card := container.NewVBox(
		container.NewHBox(iconWidget, titleLabel),
		separator,
		content,
	)
	
	// 添加内边距
	return container.NewPadded(card)
}

// InfoRow 创建信息行（图标 + 文本）
func InfoRow(icon fyne.Resource, text string) fyne.CanvasObject {
	iconWidget := widget.NewIcon(icon)
	label := widget.NewLabel(text)
	return container.NewHBox(iconWidget, label)
}

// ColoredButton 创建带颜色的按钮
func ColoredButton(text string, icon fyne.Resource, bgColor color.Color, onTapped func()) *widget.Button {
	btn := widget.NewButtonWithIcon(text, icon, onTapped)
	btn.Importance = widget.LowImportance
	return btn
}

// StatsCard 创建统计卡片
func StatsCard(icon fyne.Resource, title string, value string, valueColor color.Color) fyne.CanvasObject {
	iconWidget := widget.NewIcon(icon)
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle = fyne.TextStyle{Italic: true}
	
	valueText := canvas.NewText(value, valueColor)
	valueText.TextSize = 24
	valueText.TextStyle = fyne.TextStyle{Bold: true}
	valueText.Alignment = fyne.TextAlignCenter
	
	card := container.NewVBox(
		container.NewHBox(iconWidget, titleLabel),
		valueText,
	)
	
	return container.NewPadded(card)
}

// FeatureItem 创建功能列表项
func FeatureItem(icon fyne.Resource, title string, description string) fyne.CanvasObject {
	iconWidget := widget.NewIcon(icon)
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	descLabel := widget.NewLabel(description)
	descLabel.TextStyle = fyne.TextStyle{Italic: true}
	
	content := container.NewVBox(titleLabel, descLabel)
	
	return container.NewBorder(
		nil, nil,
		iconWidget, nil,
		content,
	)
}

// AnimatedIcon 创建可以动画的图标（简单版本）
type AnimatedIcon struct {
	widget.Icon
	rotating bool
}

// NewAnimatedIcon 创建新的动画图标
func NewAnimatedIcon(res fyne.Resource) *AnimatedIcon {
	icon := &AnimatedIcon{}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)
	return icon
}

// StartRotation 开始旋转动画（简化版）
func (a *AnimatedIcon) StartRotation() {
	a.rotating = true
	// 注意：Fyne 不直接支持旋转动画，这里只是示例
	// 实际使用需要自定义渲染
}

// StopRotation 停止旋转动画
func (a *AnimatedIcon) StopRotation() {
	a.rotating = false
}

// GradientBackground 创建渐变背景（使用矩形模拟）
func GradientBackground(topColor, bottomColor color.Color, content fyne.CanvasObject) fyne.CanvasObject {
	// Fyne 不直接支持渐变，使用纯色背景
	bg := canvas.NewRectangle(topColor)
	
	return container.NewStack(bg, content)
}

// IconButton 创建只有图标的按钮
func IconButton(icon fyne.Resource, tooltip string, onTapped func()) *widget.Button {
	btn := widget.NewButtonWithIcon("", icon, onTapped)
	// 可以添加 tooltip（需要额外实现）
	return btn
}

// BadgeLabel 创建带徽章的标签
func BadgeLabel(text string, badgeText string, badgeColor color.Color) fyne.CanvasObject {
	label := widget.NewLabel(text)
	
	badge := canvas.NewText(badgeText, color.White)
	badge.TextSize = 10
	badge.TextStyle = fyne.TextStyle{Bold: true}
	
	badgeBg := canvas.NewRectangle(badgeColor)
	badgeBg.SetMinSize(fyne.NewSize(20, 20))
	
	badgeContainer := container.NewStack(badgeBg, container.NewCenter(badge))
	
	return container.NewHBox(label, badgeContainer)
}

// ProgressCard 创建进度卡片
func ProgressCard(title string, current, total int, icon fyne.Resource) fyne.CanvasObject {
	iconWidget := widget.NewIcon(icon)
	titleLabel := widget.NewLabel(title)
	
	progress := widget.NewProgressBar()
	progress.Min = 0
	progress.Max = float64(total)
	progress.Value = float64(current)
	
	progressText := widget.NewLabel(intToString(current) + " / " + intToString(total))
	progressText.Alignment = fyne.TextAlignCenter
	
	return container.NewVBox(
		container.NewHBox(iconWidget, titleLabel),
		progress,
		progressText,
	)
}

// intToString 辅助函数：整数转字符串
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}
	
	digits := []rune{}
	for n > 0 {
		digits = append([]rune{rune('0' + n%10)}, digits...)
		n /= 10
	}
	
	result := string(digits)
	if negative {
		result = "-" + result
	}
	
	return result
}

// SectionHeader 创建章节标题
func SectionHeader(text string, icon fyne.Resource) fyne.CanvasObject {
	iconWidget := widget.NewIcon(icon)
	label := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	
	separator := widget.NewSeparator()
	
	return container.NewVBox(
		container.NewHBox(iconWidget, label),
		separator,
	)
}

// ToggleCard 创建可切换的卡片
func ToggleCard(title string, icon fyne.Resource, content fyne.CanvasObject) fyne.CanvasObject {
	iconWidget := widget.NewIcon(icon)
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	
	expanded := true
	expandIcon := widget.NewIcon(theme.MenuDropDownIcon())
	
	header := container.NewBorder(
		nil, nil,
		container.NewHBox(iconWidget, titleLabel),
		expandIcon,
	)
	
	card := container.NewVBox(
		header,
		widget.NewSeparator(),
		content,
	)
	
	// 点击标题切换展开/收起（简化版）
	_ = expanded // 避免未使用警告
	
	return card
}
