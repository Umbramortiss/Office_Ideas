package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
)

var (
	records   [][]string
	username  string
	password  string
	groupname string
	ou        string
)

func readCSV() {
	file, err := os.Open("new-hire.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err = reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func main() {
	readCSV()
	for _, record := range records {
		if len(record) < 4 {
			fmt.Println("record has less than 4 fields")
			continue
		}
		username = record[0]
		password = record[1]
		groupname = record[2]
		ou = record[3]

		chanPass()
		chanGroup()
		chanOU()
		updateAzureADUser()

		fmt.Printf("Processed user %s...\n", username)
	}
}

func chanPass() {
	cmd := exec.Command("powershell", "-Command", "Set-ADAccountPassword", "-Identity", username, "-NewPassword", (`"` + password + `"`))
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to change password:", err)
		return
	}
	fmt.Println(string(out))
}

func chanGroup() {
	cmd := exec.Command("powershell", "-Command", "Add-ADPrincipalGroupMembership", "-Identity", username, "-MemberOf", groupname)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to change group:", err)
		return
	}
	fmt.Println(string(out))
}

func chanOU() {
	cmd := exec.Command("powershell", "-Command", "Move-ADObject", "-Identity", username, "-TargetPath", "OU="+ou+", OU=Remote,OU=User Accounts, OU=PS, DC=Pollackpart, DC=Local")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to move OU:", err)
		return
	}
	fmt.Println(string(out))
}

func updateAzureADUser() {
	// Enforce MFA
	mfaCommand := fmt.Sprintf("Set-MsolUser -UserPrincipalName %s -StrongAuthenticationRequirements @(@{State='Enabled'; Detail='Required'})", username)
	runAzureCommand(mfaCommand)

	// Change license (example command, replace with actual command)
	licenseCommand := fmt.Sprintf("Set-MsolUserLicense -UserPrincipalName %s -AddLicenses 'license_sku_id'", username)
	runAzureCommand(licenseCommand)

	// Add user to shared mailbox (example command, replace with actual command)
	mailboxCommand := fmt.Sprintf("Add-MailboxPermission -Identity 'shared@mailbox.com' -User %s -AccessRights FullAccess", username)
	runAzureCommand(mailboxCommand)

	// Any additional Azure AD operations can be added here
}

func runAzureCommand(command string) {
	cmd := exec.Command("powershell", "-Command", command)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to run Azure command (%s): %v\n", command, err)
		return
	}
	fmt.Println(string(out))
}