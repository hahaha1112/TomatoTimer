// package main
//
// import (
//
//	"fyne.io/fyne/v2/app"
//	"fyne.io/fyne/v2/widget"
//
// )
//
//	func main() {
//		a := app.New()
//		w := a.NewWindow("Hello")
//
//		w.SetContent(widget.NewLabel("Hello Fyne!"))
//		w.ShowAndRun()
//
// }
package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/hajimehoshi/oto"
	"io/ioutil"
	"time"
)

func playSound() {
	data, err := ioutil.ReadFile("bell.wav")
	if err != nil {
		fmt.Println("Error reading sound file:", err)
		return
	}

	context, err := oto.NewContext(44100, 2, 2, 8192)
	if err != nil {
		fmt.Println("Error creating audio context:", err)
		return
	}
	defer context.Close()

	player := context.NewPlayer()
	defer player.Close()

	player.Write(data)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("番茄时钟")

	workDuration := 25 * time.Minute
	breakDuration := 5 * time.Minute

	label := widget.NewLabel("点击开始按钮开始工作")
	timeLabel := widget.NewLabel("剩余时间: 25:00")
	modeLabel := widget.NewLabel("当前模式: 工作")

	isPaused := false
	isDarkMode := false
	isWorkMode := true

	startButton := widget.NewButton("开始", func() {
		go func() {
			var remaining time.Duration
			for {
				if isWorkMode {
					remaining = workDuration
					modeLabel.SetText("当前模式: 工作")
				} else {
					remaining = breakDuration
					modeLabel.SetText("当前模式: 休息")
				}

				for remaining > 0 {
					if isPaused {
						time.Sleep(time.Second)
						continue
					}
					timeLabel.SetText(fmt.Sprintf("剩余时间: %02d:%02d", int(remaining.Minutes()), int(remaining.Seconds())%60))
					time.Sleep(time.Second)
					remaining -= time.Second
				}

				playSound() // 播放铃声

				if isWorkMode {
					label.SetText("工作结束，休息中...")
				} else {
					label.SetText("休息结束，可以开始新的工作周期了！")
				}

				isWorkMode = !isWorkMode
			}
		}()
	})

	pauseButton := widget.NewButton("暂停/继续", func() {
		isPaused = !isPaused
		if isPaused {
			label.SetText("已暂停")
		} else {
			label.SetText("继续中...")
		}
	})

	themeButton := widget.NewButton("切换黑夜/白天模式", func() {
		isDarkMode = !isDarkMode
		if isDarkMode {
			myApp.Settings().SetTheme(theme.DarkTheme())
		} else {
			myApp.Settings().SetTheme(theme.LightTheme())
		}
	})

	myWindow.SetContent(container.NewVBox(
		label,
		timeLabel,
		modeLabel,
		startButton,
		pauseButton,
		themeButton,
	))

	myWindow.ShowAndRun()
}
