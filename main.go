package main

import (
	"fmt"
	"github.com/caseymrm/menuet"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)
var defaultWorkTime = 14400
var workTime = 14400 // 4 Hours
var isTimerRunning = false
var runEasterEgg = true

func countdownTimer() {
	for {
		if isTimerRunning {
			menuet.App().SetMenuState(&menuet.MenuState{
				Title: "Countdown: " + formatTime(workTime),
			})
			workTime--
			if workTime <= 0 {
				workTime = defaultWorkTime
				isTimerRunning = false
				menuet.App().SetMenuState(&menuet.MenuState{
					Title: "It's the Final Countdown!",
				})
				if runEasterEgg {
					openVictorySong()
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func formatTime(totalSeconds int) string {
	hours := int(math.Floor(float64(totalSeconds / 3600)))
	minutes := int(math.Floor(float64((totalSeconds / 60) - hours * 60)))
	seconds := int(math.Floor(float64(totalSeconds % 60)))
	timeInstance := time.Date(1998, 5, 23, hours, minutes, seconds, 0, time.UTC);
	return timeInstance.Format("15:04:05")
}

func openVictorySong() {
	var err error
	url := "https://www.youtube.com/watch?v=9jK-NcRmVcw"

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func menuItems()[]menuet.MenuItem {
	items := []menuet.MenuItem{}

	menuItemText := "Toggle Timer"

	if isTimerRunning {
		menuItemText = "Stop Timer"
	} else {
		menuItemText = "Start Timer"
	}

	items = append(items, menuet.MenuItem{
		Text: menuItemText,
		Clicked: func() {
			isTimerRunning = !isTimerRunning
		},
	})

	items = append(items, menuet.MenuItem{
		Text: "Set Timer to ...",
		Clicked: func() {
			response := menuet.App().Alert(menuet.Alert{
				MessageText: "How Much Hours Will You Work?",
				Inputs:      []string{"Hours", "Minutes"},
				Buttons:     []string{"Set", "Cancel"},
			})
			if response.Button == 0 && len(response.Inputs) == 2 {
				if response.Inputs[0] == "" {
					response.Inputs[0] = "0"
				}
				if response.Inputs[1] == "" {
					response.Inputs[1] = "0"
				}
				isTimerRunning = false
				workTime = 0
				seconds, err := strconv.Atoi(response.Inputs[1])
				if err != nil {
					fmt.Println(err)
					os.Exit(2)
				}
				hoursInSeconds, err := strconv.Atoi(response.Inputs[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(2)
				}
				workTime += (seconds * 60) + (hoursInSeconds * 3600)
				isTimerRunning = true
			}
		},
	})

	easterEggText := "Toggle Easter Egg"
	if runEasterEgg {
		easterEggText = "No surprise, please!"
	} else {
		easterEggText = "Surprise me broh!"
	}
	items = append(items, menuet.MenuItem{
		Text: easterEggText,
		Clicked: func() {
			runEasterEgg = !runEasterEgg
		},
	})

	return items
}

func main() {
	go countdownTimer()
	menuet.App().Label = "com.github.lahmacun.finalcountdown"
	menuet.App().Children = menuItems
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "Final Countdown",
	})
	menuet.App().RunApplication()
}
