package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"my_portfolio/electronic_fish_tank"
)

func main() {
	electronic_fish_tank.FishTank()
	//myApp := app.New()
	//myWindow := myApp.NewWindow("我的超级工具箱")
	//myWindow.Resize(fyne.NewSize(400, 600)) // 设置一个更适合手机的窗口大小
	//
	//// 1. 创建多个“页面”的容器
	//homeContent := container.NewVBox(
	//	widget.NewLabel("欢迎来到我的工具箱！"),
	//	widget.NewButton("关于我", func() {
	//		// 点击后可以跳转到关于页面
	//	}),
	//)
	//
	//tool1Content := container.NewVBox(
	//	widget.NewLabel("这是工具一：周报生成器"),
	//	makeTool1Content(),
	//)
	//
	//tool2Content := container.NewVBox(
	//	widget.NewLabel("这是工具二：文件去重器"),
	//	// 这里放入文件去重器的所有控件
	//	// 比如：选择文件夹按钮、开始扫描按钮
	//)
	//
	//// 创建体重记录UI
	//weightTrackerUI := weight_tracker.NewWeightTrackerUI(myWindow)
	//weightTrackerContent := weightTrackerUI.MakeUI()
	//
	//// 2. 使用 TabContainer 来组织页面 (就像浏览器标签页)
	//tabs := container.NewAppTabs(
	//	container.NewTabItem("首页", homeContent),
	//	container.NewTabItem("周报工具", tool1Content),
	//	container.NewTabItem("文件去重", tool2Content),
	//	container.NewTabItem("体重记录", weightTrackerContent),
	//	// 未来你可以在这里无限添加新工具...
	//	// container.NewTabItem("新工具", newToolContent),
	//)
	//
	//myWindow.SetContent(tabs)
	//myWindow.ShowAndRun()
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
