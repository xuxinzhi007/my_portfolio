package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"my_portfolio/profit_calculator"
	"my_portfolio/settings"
	"my_portfolio/weight_tracker"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("我的超级工具箱")
	myWindow.Resize(fyne.NewSize(400, 600))

	// 创建菜单栏
	settingsItem := fyne.NewMenuItem("设置", func() {
		settings.ShowSettingsDialog(myApp, myWindow)
	})
	settingsItem.Icon = theme.SettingsIcon()

	aboutItem := fyne.NewMenuItem("关于", func() {
		// 可以添加关于对话框
	})
	aboutItem.Icon = theme.InfoIcon()

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("应用", settingsItem, aboutItem),
	)
	myWindow.SetMainMenu(mainMenu)

	// 1. 创建多个"页面"的容器
	homeContent := container.NewVBox(
		widget.NewLabelWithStyle("欢迎来到我的工具箱！", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
		widget.NewLabel("这是一个多功能工具集合应用"),
		widget.NewLabel(""),
		widget.NewButtonWithIcon("打开设置", theme.SettingsIcon(), func() {
			settings.ShowSettingsDialog(myApp, myWindow)
		}),
	)

	tool1Content := container.NewVBox(
		widget.NewLabel("这是工具一：周报生成器"),
		makeTool1Content(),
	)

	tool2Content := container.NewVBox(
		widget.NewLabel("这是工具二：文件去重器"),
		widget.NewLabel("功能开发中..."),
	)

	// 创建体重记录UI
	weightTrackerUI := weight_tracker.NewWeightTrackerUI(myWindow)
	weightTrackerContent := weightTrackerUI.MakeUI()

	// 创建收益计算器UI
	profitCalculatorUI := profit_calculator.NewProfitCalculatorUI(myWindow)
	profitCalculatorContent := profitCalculatorUI.MakeUI()

	// 创建设置UI
	settingsUI := settings.NewSettingsUI(myApp, myWindow)
	settingsContent := settingsUI.MakeUI()

	// 2. 使用 TabContainer 来组织页面
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("首页", theme.HomeIcon(), homeContent),
		container.NewTabItemWithIcon("周报工具", theme.DocumentIcon(), tool1Content),
		container.NewTabItemWithIcon("文件去重", theme.FolderIcon(), tool2Content),
		container.NewTabItemWithIcon("体重记录", theme.MediaRecordIcon(), weightTrackerContent),
		container.NewTabItemWithIcon("收益计算", theme.ConfirmIcon(), profitCalculatorContent),
		container.NewTabItemWithIcon("设置", theme.SettingsIcon(), settingsContent),
	)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

// 在main函数外面，定义我们需要的控件变量，以便在按钮函数里访问
var (
	inputEntry  *widget.Entry
	resultLabel *widget.Label
)

func makeTool1Content() *fyne.Container {
	// 创建输入框
	inputEntry = widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("在这里输入你的工作要点...")

	// 创建显示结果的标签
	resultLabel = widget.NewLabel("生成的周报会显示在这里。")

	// 创建按钮，并绑定点击事件
	generateButton := widget.NewButton("生成周报", func() {
		// 1. 从输入框获取文本
		userInput := inputEntry.Text
		// 2. 调用你的核心业务逻辑（比如之前写的AI调用函数）
		weeklyReport := userInput //generateWeeklyReport(userInput) 这是你需要实现的函数
		// 3. 将结果显示在标签上
		resultLabel.SetText(weeklyReport)
	})

	// 将控件垂直排列
	return container.NewVBox(
		widget.NewLabel("周报生成器"),
		inputEntry,
		generateButton,
		resultLabel,
	)
}
