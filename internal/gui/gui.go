package gui

import (
	"fmt"
	"os/exec"
)

func ShowWord(title, text string) error {
	notif := exec.Command("zenity", "--info", fmt.Sprint("--title=", title), fmt.Sprint("--text=", text))
	return notif.Start()
}

func Fatal(title string, err error) error {
	notif := exec.Command("zenity", "--error", fmt.Sprint("--title=", title), fmt.Sprint("--text=", err))
	return notif.Start()
}

func Input(title, text string) error {
	notif := exec.Command("zenity", "--error", fmt.Sprint("--title=", title), fmt.Sprint("--text=", text))
	return notif.Start()
}
