package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

var (
	username string
	mailbox  string
)

func User() {
	fmt.Print("Enter Username: ")
	fmt.Scanf("%s", &username)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	fmt.Println("Enter Shared mailbox name: ")
	fmt.Scanf("%s", &mailbox)
}

func main() {
	for {
		User()

		cmd := exec.Command("powershell", "-Command", "Add-RecipientPermission", mailbox, "-AccessRights", "SendAs", "-Trustee", username)
		out, err := cmd.Output()

		if err != nil {
			fmt.Println("Failed to run command:", err)
			return
		}
		fmt.Println(string(out))
		if username == "Done" {
			break
		}
	}
}
