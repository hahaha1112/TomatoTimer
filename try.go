package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/hajimehoshi/oto"
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
	var stopChan chan bool
	isRunning := false

	startButton := widget.NewButton("开始", func() {
		if isRunning {
			close(stopChan)
			isRunning = false
			time.Sleep(time.Millisecond * 100)
			label.SetText("点击开始按钮开始工作")
			timeLabel.SetText("剩余时间: 25:00")
			modeLabel.SetText("当前模式: 工作")
			isWorkMode = true
			isPaused = false
			return
		}

		isPaused = false
		label.SetText("开始计时...")
		isRunning = true
		stopChan = make(chan bool)
		go func() {
			defer func() {
				isRunning = false
			}()

			for {
				select {
				case <-stopChan:
					return
				default:
					duration := workDuration
					if !isWorkMode {
						duration = breakDuration
					}

					remaining := duration
					for remaining > 0 {
						select {
						case <-stopChan:
							return
						default:
							if isPaused {
								time.Sleep(time.Second)
								continue
							}
							timeLabel.SetText(fmt.Sprintf("剩余时间: %02d:%02d",
								int(remaining.Minutes()),
								int(remaining.Seconds())%60))
							time.Sleep(time.Second)
							remaining -= time.Second
						}
					}

					select {
					case <-stopChan:
						return
					default:
						playSound()
						if isWorkMode {
							label.SetText("工作结束，休息中...")
							isWorkMode = false
							modeLabel.SetText("当前模式: 休息")
						} else {
							label.SetText("休息结束，可以开始新的工作周期了！")
							isWorkMode = true
							modeLabel.SetText("当前模式: 工作")
						}
					}
				}
			}
		}()
	})

	pauseButton := widget.NewButton("暂停/继续", func() {
		if !isRunning {
			return
		}
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
