package profit_calculator

import (
	"errors"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ProfitCalculatorUI 收益计算器UI
type ProfitCalculatorUI struct {
	storage Storage
	data    *ProfitCalculatorData
	window  fyne.Window

	// UI 组件
	mainContent  *fyne.Container
	investorList *widget.List
	profitList   *widget.List

	// 统计显示组件
	totalInvestmentText *canvas.Text
	totalProfitText     *canvas.Text
	investorCountText   *canvas.Text
}

// NewProfitCalculatorUI 创建新的收益计算器UI
func NewProfitCalculatorUI(window fyne.Window) *ProfitCalculatorUI {
	ui := &ProfitCalculatorUI{
		storage: NewJSONStorage("profit_records.json"),
		window:  window,
	}

	// 加载现有数据
	ui.loadData()

	return ui
}

// MakeUI 构建完整的UI界面
func (ui *ProfitCalculatorUI) MakeUI() fyne.CanvasObject {
	// 创建统计卡片
	statsCard := ui.createStatsCard()

	// 创建投资者管理区域
	investorSection := ui.createInvestorSection()

	// 创建收益管理区域
	profitSection := ui.createProfitSection()

	// 组合布局
	ui.mainContent = container.NewVBox(
		statsCard,
		widget.NewSeparator(),
		investorSection,
		widget.NewSeparator(),
		profitSection,
	)

	return container.NewScroll(ui.mainContent)
}

// loadData 从存储加载数据
func (ui *ProfitCalculatorUI) loadData() {
	data, err := ui.storage.Load()
	if err != nil {
		// 如果加载失败，使用空数据
		ui.data = &ProfitCalculatorData{
			Investors:      []Investor{},
			MonthlyProfits: []MonthlyProfit{},
		}
		dialog.ShowError(
			errors.New("加载数据失败: "+err.Error()),
			ui.window,
		)
		return
	}

	ui.data = data
}

// saveData 保存数据到存储
func (ui *ProfitCalculatorUI) saveData() {
	err := ui.storage.Save(ui.data)
	if err != nil {
		dialog.ShowError(
			errors.New("保存失败: "+err.Error()),
			ui.window,
		)
	}
}

// refreshUI 刷新整个UI
func (ui *ProfitCalculatorUI) refreshUI() {
	ui.updateStats()
	
	// 重新创建列表以处理空状态
	if ui.investorList != nil {
		ui.investorList.Refresh()
	}
	if ui.profitList != nil {
		ui.profitList.Refresh()
	}
	
	// 刷新主容器
	if ui.mainContent != nil {
		ui.mainContent.Refresh()
	}
}

// createStatsCard 创建统计卡片
func (ui *ProfitCalculatorUI) createStatsCard() fyne.CanvasObject {
	// 标题
	title := widget.NewLabelWithStyle("💰 收益统计", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// 总投资
	ui.totalInvestmentText = canvas.NewText("¥0.00", nil)
	ui.totalInvestmentText.TextSize = 24
	ui.totalInvestmentText.Alignment = fyne.TextAlignCenter
	totalInvestmentLabel := widget.NewLabel("总投资")
	totalInvestmentLabel.Alignment = fyne.TextAlignCenter

	// 总收益
	ui.totalProfitText = canvas.NewText("¥0.00", nil)
	ui.totalProfitText.TextSize = 24
	ui.totalProfitText.Alignment = fyne.TextAlignCenter
	totalProfitLabel := widget.NewLabel("累计收益")
	totalProfitLabel.Alignment = fyne.TextAlignCenter

	// 投资者数量
	ui.investorCountText = canvas.NewText("0", nil)
	ui.investorCountText.TextSize = 20
	ui.investorCountText.Alignment = fyne.TextAlignCenter
	investorCountLabel := widget.NewLabel("投资者")
	investorCountLabel.Alignment = fyne.TextAlignCenter

	// 更新统计数据
	ui.updateStats()

	// 布局
	statsRow := container.NewHBox(
		container.NewVBox(totalInvestmentLabel, ui.totalInvestmentText),
		widget.NewSeparator(),
		container.NewVBox(totalProfitLabel, ui.totalProfitText),
		widget.NewSeparator(),
		container.NewVBox(investorCountLabel, ui.investorCountText),
	)

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		statsRow,
	)
}

// updateStats 更新统计信息
func (ui *ProfitCalculatorUI) updateStats() {
	stats := CalculateOverallStats(ui.data)

	ui.totalInvestmentText.Text = formatCurrency(stats.TotalInvestment)
	ui.totalProfitText.Text = formatCurrency(stats.TotalProfit)
	ui.investorCountText.Text = formatInt(stats.InvestorCount)

	ui.totalInvestmentText.Refresh()
	ui.totalProfitText.Refresh()
	ui.investorCountText.Refresh()
}

// 辅助函数：格式化货币
func formatCurrency(amount float64) string {
	return fmt.Sprintf("¥%.2f", amount)
}

// 辅助函数：格式化整数
func formatInt(num int) string {
	return fmt.Sprintf("%d", num)
}

// 辅助函数：格式化百分比
func formatPercentage(ratio float64) string {
	return fmt.Sprintf("%.2f%%", ratio*100)
}

// createInvestorSection 创建投资者管理区域
func (ui *ProfitCalculatorUI) createInvestorSection() fyne.CanvasObject {
	// 标题
	title := widget.NewLabelWithStyle("👥 投资者管理", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// 添加按钮
	addButton := widget.NewButton("添加投资者", func() {
		ui.showAddInvestorDialog()
	})

	// 创建投资者列表
	ui.createInvestorList()

	// 空状态提示
	emptyHint := widget.NewLabel("还没有投资者，点击上方按钮添加")
	emptyHint.Alignment = fyne.TextAlignCenter

	var listContainer fyne.CanvasObject
	if len(ui.data.Investors) == 0 {
		listContainer = container.NewCenter(emptyHint)
	} else {
		listContainer = ui.investorList
	}

	return container.NewBorder(
		container.NewVBox(
			container.NewHBox(title, addButton),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		listContainer,
	)
}

// createInvestorList 创建投资者列表
func (ui *ProfitCalculatorUI) createInvestorList() {
	ui.investorList = widget.NewList(
		func() int {
			return len(ui.data.Investors)
		},
		func() fyne.CanvasObject {
			// 列表项模板 - 卡片式布局
			nameLabel := widget.NewLabelWithStyle("姓名", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			
			investmentTitleLabel := widget.NewLabel("投资金额:")
			investmentAmountLabel := widget.NewLabel("¥0.00")
			
			ratioTitleLabel := widget.NewLabel("投资比例:")
			ratioLabel := widget.NewLabel("0%")
			
			profitTitleLabel := widget.NewLabel("累计收益:")
			profitLabel := widget.NewLabel("¥0.00")
			
			finalTitleLabel := widget.NewLabel("最终金额:")
			finalLabel := widget.NewLabelWithStyle("¥0.00", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

			editBtn := widget.NewButton("编辑", nil)
			deleteBtn := widget.NewButton("删除", nil)

			// 第一行：姓名
			row1 := container.NewHBox(nameLabel)
			
			// 第二行：投资金额和比例
			row2 := container.NewHBox(
				investmentTitleLabel,
				investmentAmountLabel,
				widget.NewLabel("  |  "),
				ratioTitleLabel,
				ratioLabel,
			)
			
			// 第三行：累计收益和最终金额
			row3 := container.NewHBox(
				profitTitleLabel,
				profitLabel,
				widget.NewLabel("  |  "),
				finalTitleLabel,
				finalLabel,
			)

			// 第四行：操作按钮
			btnRow := container.NewHBox(
				editBtn,
				deleteBtn,
			)

			return container.NewVBox(
				row1,
				row2,
				row3,
				btnRow,
				widget.NewSeparator(),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(ui.data.Investors) {
				return
			}

			investor := ui.data.Investors[id]
			stats := CalculateInvestorStats(investor.ID, ui.data.Investors, ui.data.MonthlyProfits)

			vbox := obj.(*fyne.Container)
			row1 := vbox.Objects[0].(*fyne.Container)
			row2 := vbox.Objects[1].(*fyne.Container)
			row3 := vbox.Objects[2].(*fyne.Container)
			btnRow := vbox.Objects[3].(*fyne.Container)

			// 更新第一行：姓名
			nameLabel := row1.Objects[0].(*widget.Label)
			nameLabel.SetText("👤 " + investor.Name)

			// 更新第二行：投资金额和比例
			investmentAmountLabel := row2.Objects[1].(*widget.Label)
			ratioLabel := row2.Objects[4].(*widget.Label)
			
			investmentAmountLabel.SetText(formatCurrency(investor.InvestmentAmount))
			ratioLabel.SetText(formatPercentage(stats.InvestmentRatio))

			// 更新第三行：累计收益和最终金额
			profitLabel := row3.Objects[1].(*widget.Label)
			finalLabel := row3.Objects[4].(*widget.Label)

			profitLabel.SetText(formatCurrency(stats.TotalProfit))
			finalLabel.SetText(formatCurrency(stats.FinalAmount))

			// 更新按钮
			editBtn := btnRow.Objects[0].(*widget.Button)
			deleteBtn := btnRow.Objects[1].(*widget.Button)

			editBtn.OnTapped = func() {
				ui.showEditInvestorDialog(&investor)
			}

			deleteBtn.OnTapped = func() {
				ui.deleteInvestor(investor.ID)
			}
		},
	)
}

// showAddInvestorDialog 显示添加投资者对话框
func (ui *ProfitCalculatorUI) showAddInvestorDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("请输入投资者姓名")

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("请输入投资金额")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "姓名", Widget: nameEntry},
			{Text: "投资金额", Widget: amountEntry},
		},
		OnSubmit: func() {
			// 验证姓名
			name := nameEntry.Text
			if name == "" {
				dialog.ShowError(errors.New("姓名不能为空"), ui.window)
				return
			}

			if len(name) > 50 {
				dialog.ShowError(errors.New("姓名长度不能超过50个字符"), ui.window)
				return
			}

			// 检查姓名重复
			for _, investor := range ui.data.Investors {
				if investor.Name == name {
					dialog.ShowError(errors.New("投资者姓名已存在"), ui.window)
					return
				}
			}

			// 验证金额
			amount, err := parseAmount(amountEntry.Text)
			if err != nil {
				dialog.ShowError(errors.New("请输入有效的金额"), ui.window)
				return
			}

			if amount <= 0 {
				dialog.ShowError(errors.New("投资金额必须大于0"), ui.window)
				return
			}

			if amount < 0.01 || amount > 10000000 {
				dialog.ShowError(errors.New("投资金额必须在0.01到10,000,000之间"), ui.window)
				return
			}

			// 创建新投资者
			newInvestor := NewInvestor(name, amount)
			ui.data.Investors = append(ui.data.Investors, *newInvestor)

			// 保存数据
			ui.saveData()

			// 刷新UI
			ui.refreshUI()

			dialog.ShowInformation("成功", fmt.Sprintf("投资者 %s 已添加", name), ui.window)
		},
	}

	d := dialog.NewForm("添加投资者", "添加", "取消", form.Items, func(confirmed bool) {
		if confirmed {
			form.OnSubmit()
		}
	}, ui.window)
	d.Show()
}

// 辅助函数：解析金额
func parseAmount(s string) (float64, error) {
	var amount float64
	_, err := fmt.Sscanf(s, "%f", &amount)
	return amount, err
}

// showEditInvestorDialog 显示编辑投资者对话框
func (ui *ProfitCalculatorUI) showEditInvestorDialog(investor *Investor) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(investor.Name)

	amountEntry := widget.NewEntry()
	amountEntry.SetText(fmt.Sprintf("%.2f", investor.InvestmentAmount))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "姓名", Widget: nameEntry},
			{Text: "投资金额", Widget: amountEntry},
		},
		OnSubmit: func() {
			// 验证姓名
			name := nameEntry.Text
			if name == "" {
				dialog.ShowError(errors.New("姓名不能为空"), ui.window)
				return
			}

			if len(name) > 50 {
				dialog.ShowError(errors.New("姓名长度不能超过50个字符"), ui.window)
				return
			}

			// 检查姓名重复（排除自己）
			for _, inv := range ui.data.Investors {
				if inv.Name == name && inv.ID != investor.ID {
					dialog.ShowError(errors.New("投资者姓名已存在"), ui.window)
					return
				}
			}

			// 验证金额
			amount, err := parseAmount(amountEntry.Text)
			if err != nil {
				dialog.ShowError(errors.New("请输入有效的金额"), ui.window)
				return
			}

			if amount <= 0 {
				dialog.ShowError(errors.New("投资金额必须大于0"), ui.window)
				return
			}

			if amount < 0.01 || amount > 10000000 {
				dialog.ShowError(errors.New("投资金额必须在0.01到10,000,000之间"), ui.window)
				return
			}

			// 更新投资者信息
			for i := range ui.data.Investors {
				if ui.data.Investors[i].ID == investor.ID {
					ui.data.Investors[i].Name = name
					ui.data.Investors[i].InvestmentAmount = amount
					break
				}
			}

			// 保存数据
			ui.saveData()

			// 刷新UI
			ui.refreshUI()

			dialog.ShowInformation("成功", fmt.Sprintf("投资者 %s 已更新", name), ui.window)
		},
	}

	d := dialog.NewForm("编辑投资者", "保存", "取消", form.Items, func(confirmed bool) {
		if confirmed {
			form.OnSubmit()
		}
	}, ui.window)
	d.Show()
}

// deleteInvestor 删除投资者
func (ui *ProfitCalculatorUI) deleteInvestor(investorID string) {
	// 查找投资者姓名
	var investorName string
	for _, investor := range ui.data.Investors {
		if investor.ID == investorID {
			investorName = investor.Name
			break
		}
	}

	// 显示确认对话框
	dialog.ShowConfirm(
		"确认删除",
		fmt.Sprintf("确定要删除投资者 %s 吗？\n\n注意：历史收益记录将被保留用于审计。", investorName),
		func(confirmed bool) {
			if !confirmed {
				return
			}

			// 删除投资者
			newInvestors := []Investor{}
			for _, investor := range ui.data.Investors {
				if investor.ID != investorID {
					newInvestors = append(newInvestors, investor)
				}
			}
			ui.data.Investors = newInvestors

			// 保存数据
			ui.saveData()

			// 刷新UI
			ui.refreshUI()

			dialog.ShowInformation("成功", fmt.Sprintf("投资者 %s 已删除", investorName), ui.window)
		},
		ui.window,
	)
}

// createProfitSection 创建收益管理区域
func (ui *ProfitCalculatorUI) createProfitSection() fyne.CanvasObject {
	// 标题
	title := widget.NewLabelWithStyle("📊 月度收益记录", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// 添加按钮
	addButton := widget.NewButton("添加收益记录", func() {
		ui.showAddProfitDialog()
	})

	// 创建收益列表
	ui.createProfitList()

	// 空状态提示
	emptyHint := widget.NewLabel("还没有收益记录，点击上方按钮添加")
	emptyHint.Alignment = fyne.TextAlignCenter

	var listContainer fyne.CanvasObject
	if len(ui.data.MonthlyProfits) == 0 {
		listContainer = container.NewCenter(emptyHint)
	} else {
		listContainer = ui.profitList
	}

	return container.NewBorder(
		container.NewVBox(
			container.NewHBox(title, addButton),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		listContainer,
	)
}

// createProfitList 创建收益记录列表
func (ui *ProfitCalculatorUI) createProfitList() {
	ui.profitList = widget.NewList(
		func() int {
			return len(ui.data.MonthlyProfits)
		},
		func() fyne.CanvasObject {
			// 列表项模板
			dateLabel := widget.NewLabel("日期")
			amountLabel := widget.NewLabelWithStyle("金额", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

			detailBtn := widget.NewButton("查看详情", nil)
			deleteBtn := widget.NewButton("删除", nil)

			infoRow := container.NewHBox(
				dateLabel,
				amountLabel,
			)

			btnRow := container.NewHBox(
				detailBtn,
				deleteBtn,
			)

			return container.NewVBox(
				infoRow,
				btnRow,
				widget.NewSeparator(),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(ui.data.MonthlyProfits) {
				return
			}

			profit := ui.data.MonthlyProfits[id]

			vbox := obj.(*fyne.Container)
			infoRow := vbox.Objects[0].(*fyne.Container)
			btnRow := vbox.Objects[1].(*fyne.Container)

			// 更新信息行
			dateLabel := infoRow.Objects[0].(*widget.Label)
			amountLabel := infoRow.Objects[1].(*widget.Label)

			dateLabel.SetText(profit.Date.Format("2006-01-02"))
			amountLabel.SetText(formatCurrency(profit.TotalProfit))

			// 更新按钮
			detailBtn := btnRow.Objects[0].(*widget.Button)
			deleteBtn := btnRow.Objects[1].(*widget.Button)

			detailBtn.OnTapped = func() {
				ui.showProfitDetailDialog(&profit)
			}

			deleteBtn.OnTapped = func() {
				ui.deleteProfitRecord(profit.ID)
			}
		},
	)
}

// showAddProfitDialog 显示添加收益记录对话框
func (ui *ProfitCalculatorUI) showAddProfitDialog() {
	// 检查是否有投资者
	if len(ui.data.Investors) == 0 {
		dialog.ShowError(errors.New("请先添加投资者"), ui.window)
		return
	}

	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("YYYY-MM-DD")
	dateEntry.SetText(time.Now().Format("2006-01-02"))

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("请输入总收益金额")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "日期", Widget: dateEntry},
			{Text: "总收益", Widget: amountEntry},
		},
		OnSubmit: func() {
			// 验证日期
			date, err := time.Parse("2006-01-02", dateEntry.Text)
			if err != nil {
				dialog.ShowError(errors.New("日期格式无效，请使用 YYYY-MM-DD 格式"), ui.window)
				return
			}

			// 检查日期不能为未来
			if date.After(time.Now()) {
				dialog.ShowError(errors.New("日期不能为未来"), ui.window)
				return
			}

			// 验证金额
			amount, err := parseAmount(amountEntry.Text)
			if err != nil {
				dialog.ShowError(errors.New("请输入有效的金额"), ui.window)
				return
			}

			if amount < -10000000 || amount > 10000000 {
				dialog.ShowError(errors.New("收益金额必须在-10,000,000到10,000,000之间"), ui.window)
				return
			}

			// 计算收益分配
			distributions := DistributeProfit(amount, ui.data.Investors)

			// 创建新收益记录
			newProfit := NewMonthlyProfit(date, amount, distributions)
			ui.data.MonthlyProfits = append(ui.data.MonthlyProfits, *newProfit)

			// 保存数据
			ui.saveData()

			// 刷新UI
			ui.refreshUI()

			dialog.ShowInformation("成功", fmt.Sprintf("收益记录已添加：%s", formatCurrency(amount)), ui.window)
		},
	}

	d := dialog.NewForm("添加收益记录", "添加", "取消", form.Items, func(confirmed bool) {
		if confirmed {
			form.OnSubmit()
		}
	}, ui.window)
	d.Show()
}

// showProfitDetailDialog 显示收益详情对话框
func (ui *ProfitCalculatorUI) showProfitDetailDialog(profit *MonthlyProfit) {
	// 创建详情内容
	dateLabel := widget.NewLabel(fmt.Sprintf("日期：%s", profit.Date.Format("2006-01-02")))
	totalLabel := widget.NewLabelWithStyle(
		fmt.Sprintf("总收益：%s", formatCurrency(profit.TotalProfit)),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	// 创建分配明细列表
	detailsLabel := widget.NewLabelWithStyle("分配明细：", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	var distributionRows []fyne.CanvasObject
	for _, investor := range ui.data.Investors {
		if amount, exists := profit.Distributions[investor.ID]; exists {
			totalInvestment := CalculateTotalInvestment(ui.data.Investors)
			ratio := CalculateInvestmentRatio(investor, totalInvestment)

			row := widget.NewLabel(fmt.Sprintf(
				"  • %s: %s (%s)",
				investor.Name,
				formatCurrency(amount),
				formatPercentage(ratio),
			))
			distributionRows = append(distributionRows, row)
		}
	}

	// 组合内容
	content := container.NewVBox(
		dateLabel,
		totalLabel,
		widget.NewSeparator(),
		detailsLabel,
	)

	for _, row := range distributionRows {
		content.Add(row)
	}

	// 显示对话框
	dialog.ShowCustom("收益详情", "关闭", content, ui.window)
}

// deleteProfitRecord 删除收益记录
func (ui *ProfitCalculatorUI) deleteProfitRecord(profitID string) {
	// 查找收益记录
	var profitDate string
	var profitAmount float64
	for _, profit := range ui.data.MonthlyProfits {
		if profit.ID == profitID {
			profitDate = profit.Date.Format("2006-01-02")
			profitAmount = profit.TotalProfit
			break
		}
	}

	// 显示确认对话框
	dialog.ShowConfirm(
		"确认删除",
		fmt.Sprintf("确定要删除 %s 的收益记录（%s）吗？", profitDate, formatCurrency(profitAmount)),
		func(confirmed bool) {
			if !confirmed {
				return
			}

			// 删除收益记录
			newProfits := []MonthlyProfit{}
			for _, profit := range ui.data.MonthlyProfits {
				if profit.ID != profitID {
					newProfits = append(newProfits, profit)
				}
			}
			ui.data.MonthlyProfits = newProfits

			// 保存数据
			ui.saveData()

			// 刷新UI
			ui.refreshUI()

			dialog.ShowInformation("成功", "收益记录已删除", ui.window)
		},
		ui.window,
	)
}
