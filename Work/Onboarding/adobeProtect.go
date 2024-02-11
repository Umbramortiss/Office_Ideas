package main

import (
	"fmt"
	"os/exec"
)

func main() {

	path := "\\HKEY_CURRENT_USER\\SOFTWARE\\Adobe\\Adobe Acrobat\\DC\\Privileged"
	name := "bProtectedMode"
	value := "0"
	property := "DWORD"

	cmd := exec.Command("powershell", "-Command", "New-ItemProperty", "-Path", path, "-Name", name, "-Value", value, "-PropertyType", property, "-Force")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println("Failed to run command:", err)
		return
	}
	fmt.Println(string(out))
}
