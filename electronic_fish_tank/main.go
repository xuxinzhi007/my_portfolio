package electronic_fish_tank

import (
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func FishTank() {
	// 1. 创建应用
	myApp := app.New()
	myWindow := myApp.NewWindow("电子鱼缸 - V0.3")
	myWindow.Resize(fyne.NewSize(800, 500))

	// 2. 创建鱼缸背景
	tank := canvas.NewRectangle(color.NRGBA{R: 20, G: 60, B: 120, A: 255})
	tank.Resize(fyne.NewSize(780, 400))
	tank.Move(fyne.NewPos(10, 50))

	// 3. 创建状态标签
	statusLabel := widget.NewLabel("点击鱼可以标记完成!")
	statusLabel.Move(fyne.NewPos(20, 10))

	// 4. 创建鱼的身体和尾巴
	fishBody := canvas.NewCircle(color.NRGBA{R: 255, G: 100, B: 100, A: 255})
	fishBody.Resize(fyne.NewSize(40, 40))
	fishBody.StrokeWidth = 2
	fishBody.StrokeColor = color.NRGBA{R: 200, G: 80, B: 80, A: 255}

	fishTail := canvas.NewLine(color.NRGBA{R: 255, G: 150, B: 150, A: 255})
	fishTail.StrokeWidth = 3

	// 5. 创建透明按钮来实现点击检测
	fishButton := widget.NewButton("", nil)
	fishButton.Resize(fyne.NewSize(50, 50)) // 比鱼身体稍大，方便点击

	// 设置按钮样式为完全透明
	fishButton.Importance = widget.LowImportance

	// 鱼的点击事件
	fishButton.OnTapped = func() {
		statusLabel.SetText("太棒了! 你完成了这个任务! 🎉")
		fishBody.FillColor = color.NRGBA{R: 100, G: 200, B: 100, A: 255}
		fishBody.Refresh()
	}

	// 6. 创建容器 - 注意层级顺序（后面的元素在上面）
	content := container.NewWithoutLayout(
		tank,
		statusLabel,
		fishBody,
		fishTail,
		fishButton, // 按钮在最上面
	)
	myWindow.SetContent(content)

	// 7. 动画逻辑
	offset := 0.0
	go func() {
		for {
			offset += 0.05

			// 鱼的位置计算
			xPos := float32(100 + math.Sin(offset*0.5)*300)
			yPos := float32(200 + math.Sin(offset)*80)

			// 更新鱼身体位置
			fishBody.Move(fyne.NewPos(xPos, yPos))

			// 更新鱼尾巴位置
			tailDirection := float32(math.Sin(offset * 2))
			fishTail.Position1 = fyne.NewPos(xPos-20, yPos-tailDirection*10)
			fishTail.Position2 = fyne.NewPos(xPos-35, yPos+tailDirection*10)

			// 更新透明按钮位置（跟随鱼）
			fishButton.Move(fyne.NewPos(xPos-5, yPos-5)) // 居中调整

			// 边界检测和重置
			if xPos > 750 {
				offset = math.Pi / 2
				fishBody.FillColor = color.NRGBA{R: 255, G: 100, B: 100, A: 255}
				fishBody.Refresh()
				statusLabel.SetText("点击鱼可以标记完成!")
			}

			fishBody.Refresh()
			fishTail.Refresh()

			time.Sleep(50 * time.Millisecond)
		}
	}()

	myWindow.ShowAndRun()
}
