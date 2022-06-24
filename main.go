package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func main() {
	autorun()
}

func autorun() {
	var wgAR sync.WaitGroup
	wgAR.Add(2)
	go river(&wgAR)
	go proxy(&wgAR)
	wgAR.Wait()
}

func river(wgAR *sync.WaitGroup) {
	defer wgAR.Done()
	killProcess("river")
	err := exec.Command("river").Start()

	if err != nil {
		log.Println(fmt.Errorf(fmt.Sprintln("Err:", "river failed.", err)))
	}
}

func proxy(wgAR *sync.WaitGroup) {
	defer wgAR.Done()

	err := tor()
	if err != nil {
		log.Println(fmt.Errorf(fmt.Sprintln("Err:", "proxy failed.", err)))
	}
	err = clash()
	if err != nil {
		log.Println(fmt.Errorf(fmt.Sprintln("Err:", "proxy failed.", err)))
	}

}
func tor() error {
	killProcess("tor")
	err := exec.Command("tor", "--RunAsDaemon", "1").Start()

	if err != nil {
		return fmt.Errorf(fmt.Sprintln("Err:", "tor failed.", err))
	}

	return nil
}
func clash() error {
	killProcess("clash")
	err := exec.Command("clash").Start()

	if err != nil {
		return fmt.Errorf(fmt.Sprintln("Err:", "Clash failed.", err))
	}

	return nil
}

func killProcess(procName string) {

	p, err := exec.Command("pidof", procName).Output()
	if err != nil {
		return
	}
	pp := strings.Split(string(p), " ")
	if len(pp[0]) == 0 {
		return
	}

	pid, err := strconv.Atoi(pp[0])
	if err != nil {
		log.Panicln("Err:", "trouble converting output of pidof:", "\n", string(p))
	}

	killErr := exec.Command("kill", fmt.Sprint(pid)).Run()
	if killErr != nil {
		log.Println("Err:", "having trouble killing", procName)
	}
}
