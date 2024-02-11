package main


import (
	"fmt"
	"os/exec"
)

var(
	username string
	groupname string
)
func User() {
    fmt.Print("Enter Username: ")
    fmt.Scanf("%s", &username)
	fmt.Println("Enter Groupname: ")
	fmt.Scanf("%s", &groupname)
}




func main(){
	User()


	cmd := exec.Command("powershell", "-Command", "Add-ADPrincipalGroupMembership", "-Identity", username , "-MemberOf" groupname)
	out, err := cmd.Output()

	if err != nil{
		fmt.Println("Failed to run command:", err)
		return
	}
	fmt.Println(string(out))
}