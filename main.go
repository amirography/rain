package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

var config, _ = os.UserConfigDir()

func main() {
	autorun()
}

func autorun() {
	var autorunWg sync.WaitGroup

	autorunWg.Add(2)

	windowManager(&autorunWg)

	proxyErr := proxy(&autorunWg)
	if proxyErr != nil {
		fmt.Println(proxyErr)
	}

	autorunWg.Wait()

}
func windowManager(autorunWg *sync.WaitGroup) {
	defer autorunWg.Done()

	riverErr := river()
	if riverErr != nil {
		log.Panicln("Err:", "windowManager failed")
	}
	var windowManagerWG sync.WaitGroup

	windowManagerWG.Add(4)
	go swaybg(&windowManagerWG)
	go waybar(&windowManagerWG)
	go mako(&windowManagerWG)
	go setDbus(&windowManagerWG)

	windowManagerWG.Wait()

}

func river() error {
	err := exec.Command("river").Run()

	if err != nil {
		return fmt.Errorf(fmt.Sprintln("Err:", "river failed.", err))
	}

	return nil
}

func waybar(windowManagerWG *sync.WaitGroup) {
	defer windowManagerWG.Done()
	killErr := exec.Command(
		"killall",
		"waybar",
	).Run()
	_ = killErr

	err := exec.Command(
		"waybar",
		"-c",
		config+"/river/waybar/config.json",
		"-s",
		config+"/river/waybar/style.css",
	).Run()
	if err != nil {
		fmt.Println(fmt.Errorf(fmt.Sprintln("Err:", "waybar failed.", err)))
	}

}

func swaybg(windowManagerWG *sync.WaitGroup) {
	defer windowManagerWG.Done()

	killErr := exec.Command(
		"killall",
		"swaybg",
	).Run()
	_ = killErr

	err := exec.Command("swaybg",
		"-m",
		"fill",
		"-i",
		config+"/wallpaper").Run()
	if err != nil {
		log.Println(fmt.Errorf(fmt.Sprintln("Err:", "waybar failed.", err)))
	}

}
func mako(windowManagerWG *sync.WaitGroup) {
	defer windowManagerWG.Done()

	err := exec.Command("mako",
		"--default-timeout",
		"5000",
		"--background-color",
		"#"+"0xE0E0E0",
		"--border-color",
		"#"+"0xE0E0E0",
		"--border-size",
		"2",
		"--font",
		"monospace",
		"--padding",
		"20",
		"--width",
		"350",
	).Run()
	if err != nil {
		log.Println(fmt.Errorf(fmt.Sprintln("Err:", "mako failed.", err)))
	}

}
func proxy(autorunWg *sync.WaitGroup) error {
	defer autorunWg.Done()
	torErr := tor()
	if torErr != nil {
		return fmt.Errorf(fmt.Sprintln("Err:", "proxy failed.", torErr))
	}
	return nil
}
func tor() error {
	err := exec.Command("tor").Run()

	if err != nil {
		return fmt.Errorf(fmt.Sprintln("Err:", "tor failed.", err))
	}

	return nil
}
func setDbus(windowManagerWG *sync.WaitGroup) {
	defer windowManagerWG.Done()
	err := exec.Command(
		"dbus-update-activation-environment",
		"SEATD_SOCK",
		"DISPLAY",
		"WAYLAND_DISPLAY",
		"XDG_SESSION_TYPE",
		"XDG_CURRENT_DESKTOP",
	).Run()

	if err != nil {
		log.Println(fmt.Errorf(fmt.Sprintln("Err:", "setting dbus failed.", err)))
	}

}
