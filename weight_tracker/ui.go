package weight_tracker

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// WeightTrackerUI 体重记录UI
type WeightTrackerUI struct {
	storage        Storage
	records        []WeightRecord
	weightEntry    *widget.Entry
	recordList     *widget.List
	window         fyne.Window
	mainContent    *fyne.Container
	listContainer  fyne.CanvasObject
	statsContainer *fyne.Container
	currentWeight  *canvas.Text
	totalChange    *canvas.Text
	recordCount    *canvas.Text
	highestWeight  *canvas.Text
	lowestWeight   *canvas.Text
}

// NewWeightTrackerUI 创建新的体重记录UI
func NewWeightTrackerUI(window fyne.Window) *WeightTrackerUI {
	ui := &WeightTrackerUI{
		storage: NewJSONStorage("weight_records.json"),
		window:  window,
	}

	// 加载现有记录
	ui.loadRecords()

	return ui
}

// MakeUI 构建完整的UI界面
func (ui *WeightTrackerUI) MakeUI() fyne.CanvasObject {
	// 创建统计卡片
	statsCard := ui.createStatsCard()

	// 创建输入区域
	inputCard := ui.createInputCard()

	// 创建记录列表
	ui.createRecordList()

	// 历史记录标题
	historyTitle := widget.NewLabelWithStyle("📊 历史记录", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// 创建列表容器
	ui.updateListContainer()

	// 组合布局
	ui.mainContent = container.NewBorder(
		container.NewVBox(
			statsCard,
			widget.NewSeparator(),
			inputCard,
			widget.NewSeparator(),
			historyTitle,
		),
		nil, nil, nil,
		ui.listContainer,
	)

	return ui.mainContent
}

// createStatsCard 创建统计卡片
func (ui *WeightTrackerUI) createStatsCard() fyne.CanvasObject {
	// 标题
	title := widget.NewLabelWithStyle("📈 体重统计", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// 当前体重
	ui.currentWeight = canvas.NewText("--", color.RGBA{R: 33, G: 150, B: 243, A: 255})
	ui.currentWeight.TextSize = 32
	ui.currentWeight.Alignment = fyne.TextAlignCenter
	currentLabel := widget.NewLabel("当前体重")
	currentLabel.Alignment = fyne.TextAlignCenter

	// 总变化
	ui.totalChange = canvas.NewText("--", color.RGBA{R: 76, G: 175, B: 80, A: 255})
	ui.totalChange.TextSize = 24
	ui.totalChange.Alignment = fyne.TextAlignCenter
	totalChangeLabel := widget.NewLabel("总变化")
	totalChangeLabel.Alignment = fyne.TextAlignCenter

	// 记录数量
	ui.recordCount = canvas.NewText("0", color.RGBA{R: 156, G: 39, B: 176, A: 255})
	ui.recordCount.TextSize = 20
	ui.recordCount.Alignment = fyne.TextAlignCenter
	recordCountLabel := widget.NewLabel("记录数")
	recordCountLabel.Alignment = fyne.TextAlignCenter

	// 最高体重
	ui.highestWeight = canvas.NewText("--", color.RGBA{R: 255, G: 87, B: 34, A: 255})
	ui.highestWeight.TextSize = 16
	ui.highestWeight.Alignment = fyne.TextAlignCenter
	highestLabel := widget.NewLabel("最高")
	highestLabel.Alignment = fyne.TextAlignCenter

	// 最低体重
	ui.lowestWeight = canvas.NewText("--", color.RGBA{R: 0, G: 150, B: 136, A: 255})
	ui.lowestWeight.TextSize = 16
	ui.lowestWeight.Alignment = fyne.TextAlignCenter
	lowestLabel := widget.NewLabel("最低")
	lowestLabel.Alignment = fyne.TextAlignCenter

	// 主要统计区域
	mainStats := container.NewVBox(
		currentLabel,
		ui.currentWeight,
		totalChangeLabel,
		ui.totalChange,
	)

	// 次要统计区域
	secondaryStats := container.NewHBox(
		layout.NewSpacer(),
		container.NewVBox(recordCountLabel, ui.recordCount),
		layout.NewSpacer(),
		container.NewVBox(highestLabel, ui.highestWeight),
		layout.NewSpacer(),
		container.NewVBox(lowestLabel, ui.lowestWeight),
		layout.NewSpacer(),
	)

	// 更新统计数据
	ui.updateStats()

	// 创建带背景的卡片
	card := container.NewVBox(
		title,
		widget.NewSeparator(),
		mainStats,
		widget.NewSeparator(),
		secondaryStats,
	)

	return card
}

// createInputCard 创建输入卡片
func (ui *WeightTrackerUI) createInputCard() fyne.CanvasObject {
	// 创建输入框
	ui.weightEntry = widget.NewEntry()
	ui.weightEntry.SetPlaceHolder("例如: 70.5")
	ui.weightEntry.OnSubmitted = func(s string) {
		ui.addRecord()
	}

	// 创建添加按钮（使用图标）
	addButton := widget.NewButtonWithIcon("添加记录", theme.ContentAddIcon(), func() {
		ui.addRecord()
	})
	addButton.Importance = widget.HighImportance

	// 输入标签
	inputLabel := widget.NewLabelWithStyle("⚖️  输入体重 (kg)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// 布局
	inputContainer := container.NewVBox(
		inputLabel,
		container.NewBorder(nil, nil, nil, addButton, ui.weightEntry),
	)

	return inputContainer
}

// createRecordList 创建记录列表
func (ui *WeightTrackerUI) createRecordList() {
	ui.recordList = widget.NewList(
		func() int {
			return len(ui.records)
		},
		func() fyne.CanvasObject {
			// 创建更美观的列表项模板
			dateIcon := widget.NewIcon(theme.HistoryIcon())
			dateLabel := widget.NewLabel("日期")
			dateLabel.TextStyle = fyne.TextStyle{Italic: true}

			weightIcon := widget.NewIcon(theme.InfoIcon())
			weightLabel := widget.NewLabelWithStyle("体重", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

			changeLabel := canvas.NewText("变化", color.Black)
			changeLabel.TextStyle = fyne.TextStyle{Bold: true}

			// 卡片式布局
			card := container.NewVBox(
				container.NewHBox(dateIcon, dateLabel),
				container.NewHBox(
					weightIcon,
					weightLabel,
					layout.NewSpacer(),
					changeLabel,
				),
				widget.NewSeparator(),
			)

			return card
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			// 更新列表项内容
			if id >= len(ui.records) {
				return
			}

			record := ui.records[id]
			vbox := obj.(*fyne.Container)

			// 日期行
			dateRow := vbox.Objects[0].(*fyne.Container)
			dateLabel := dateRow.Objects[1].(*widget.Label)
			dateLabel.SetText(record.FormatDate())

			// 体重和变化行
			weightRow := vbox.Objects[1].(*fyne.Container)
			weightLabel := weightRow.Objects[1].(*widget.Label)
			changeText := weightRow.Objects[3].(*canvas.Text)

			weightLabel.SetText(fmt.Sprintf("%.1f kg", record.Weight))
			changeText.Text = record.FormatChange()

			// 根据变化类型设置颜色和样式
			switch record.ChangeType {
			case "increase":
				changeText.Color = color.RGBA{R: 244, G: 67, B: 54, A: 255} // 红色
			case "decrease":
				changeText.Color = color.RGBA{R: 76, G: 175, B: 80, A: 255} // 绿色
			case "stable":
				changeText.Color = color.RGBA{R: 158, G: 158, B: 158, A: 255} // 灰色
			case "first":
				changeText.Color = color.RGBA{R: 33, G: 150, B: 243, A: 255} // 蓝色
			}

			changeText.Refresh()
		},
	)
}

// updateListContainer 更新列表容器
func (ui *WeightTrackerUI) updateListContainer() {
	if len(ui.records) == 0 {
		// 空状态
		emptyIcon := widget.NewIcon(theme.DocumentCreateIcon())
		emptyLabel := widget.NewLabel("还没有记录")
		emptyLabel.Alignment = fyne.TextAlignCenter
		emptyLabel.TextStyle = fyne.TextStyle{Bold: true}

		emptyHint := widget.NewLabel("添加第一条体重记录开始追踪吧！")
		emptyHint.Alignment = fyne.TextAlignCenter

		ui.listContainer = container.NewCenter(
			container.NewVBox(
				emptyIcon,
				emptyLabel,
				emptyHint,
			),
		)
	} else {
		ui.listContainer = ui.recordList
	}
}

// updateStats 更新统计信息
func (ui *WeightTrackerUI) updateStats() {
	stats := CalculateStats(ui.records)

	if stats.TotalRecords == 0 {
		ui.currentWeight.Text = "--"
		ui.totalChange.Text = "--"
		ui.recordCount.Text = "0"
		ui.highestWeight.Text = "--"
		ui.lowestWeight.Text = "--"
	} else {
		ui.currentWeight.Text = fmt.Sprintf("%.1f kg", stats.CurrentWeight)

		// 设置总变化的颜色和文本
		if stats.TotalChange > 0 {
			ui.totalChange.Text = fmt.Sprintf("↑ +%.1f kg", stats.TotalChange)
			ui.totalChange.Color = color.RGBA{R: 244, G: 67, B: 54, A: 255} // 红色
		} else if stats.TotalChange < 0 {
			ui.totalChange.Text = fmt.Sprintf("↓ %.1f kg", stats.TotalChange)
			ui.totalChange.Color = color.RGBA{R: 76, G: 175, B: 80, A: 255} // 绿色
		} else {
			ui.totalChange.Text = "● 持平"
			ui.totalChange.Color = color.RGBA{R: 158, G: 158, B: 158, A: 255} // 灰色
		}

		ui.recordCount.Text = fmt.Sprintf("%d 条", stats.TotalRecords)
		ui.highestWeight.Text = fmt.Sprintf("%.1f kg", stats.HighestWeight)
		ui.lowestWeight.Text = fmt.Sprintf("%.1f kg", stats.LowestWeight)
	}

	// 刷新所有文本
	ui.currentWeight.Refresh()
	ui.totalChange.Refresh()
	ui.recordCount.Refresh()
	ui.highestWeight.Refresh()
	ui.lowestWeight.Refresh()
}

// addRecord 添加新记录（带动画效果）
func (ui *WeightTrackerUI) addRecord() {
	// 获取输入值
	weightStr := ui.weightEntry.Text

	// 验证：检查空输入
	if weightStr == "" {
		dialog.ShowError(
			errors.New("请输入体重值"),
			ui.window,
		)
		return
	}

	// 验证：检查是否为有效数字
	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		dialog.ShowError(
			errors.New("请输入有效的数字"),
			ui.window,
		)
		return
	}

	// 验证：检查是否为正数
	if weight <= 0 {
		dialog.ShowError(
			errors.New("体重必须大于 0"),
			ui.window,
		)
		return
	}

	// 验证：检查合理范围
	if weight < 20 || weight > 300 {
		dialog.ShowError(
			errors.New("请输入合理的体重值 (20-300 kg)"),
			ui.window,
		)
		return
	}

	// 获取上一条记录（如果存在）
	var previousRecord *WeightRecord
	if len(ui.records) > 0 {
		previousRecord = &ui.records[0]
	}

	// 创建新记录
	newRecord := NewWeightRecord(weight, previousRecord)

	// 插入到列表开头（保持倒序）
	ui.records = append([]WeightRecord{*newRecord}, ui.records...)

	// 保存到文件
	ui.saveRecords()

	// 更新统计信息（带动画效果）
	ui.animateStatsUpdate()

	// 更新列表容器（处理从空到有记录的情况）
	ui.updateListContainer()
	if ui.mainContent != nil {
		ui.mainContent.Objects[0] = ui.listContainer
		ui.mainContent.Refresh()
	}

	// 刷新列表
	if ui.recordList != nil {
		ui.recordList.Refresh()
	}

	// 清空输入框
	ui.weightEntry.SetText("")

	// 显示成功提示
	dialog.ShowInformation(
		"✅ 成功",
		fmt.Sprintf("体重记录已添加：%.1f kg", weight),
		ui.window,
	)
}

// animateStatsUpdate 动画更新统计信息
func (ui *WeightTrackerUI) animateStatsUpdate() {
	// 简单的淡入效果
	oldAlpha := uint8(255)

	// 淡出
	for i := 0; i < 5; i++ {
		oldAlpha -= 50
		ui.currentWeight.Color = color.RGBA{R: 33, G: 150, B: 243, A: oldAlpha}
		ui.currentWeight.Refresh()
		time.Sleep(20 * time.Millisecond)
	}

	// 更新数据
	ui.updateStats()

	// 淡入
	newAlpha := uint8(0)
	for i := 0; i < 5; i++ {
		newAlpha += 50
		ui.currentWeight.Color = color.RGBA{R: 33, G: 150, B: 243, A: newAlpha}
		ui.currentWeight.Refresh()
		time.Sleep(20 * time.Millisecond)
	}

	// 确保最终完全显示
	ui.currentWeight.Color = color.RGBA{R: 33, G: 150, B: 243, A: 255}
	ui.currentWeight.Refresh()
}

// loadRecords 从存储加载记录
func (ui *WeightTrackerUI) loadRecords() {
	records, err := ui.storage.Load()
	if err != nil {
		// 如果加载失败，使用空列表
		ui.records = []WeightRecord{}
		return
	}

	ui.records = records
}

// saveRecords 保存记录到存储
func (ui *WeightTrackerUI) saveRecords() {
	err := ui.storage.Save(ui.records)
	if err != nil {
		dialog.ShowError(
			errors.New("保存失败: "+err.Error()),
			ui.window,
		)
	}
}
