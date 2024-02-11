package main

import (
	_ "embed"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"math/rand"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"os"

	"image/png"

	"github.com/kbinani/screenshot"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
)

var (
	Green   = color.New(color.FgGreen).Add(color.Bold)
	Red     = color.New(color.FgRed).Add(color.Bold)
	Magenta = color.New(color.FgMagenta).Add(color.Bold)
	Blue    = color.New(color.FgBlue).Add(color.Bold)
	White   = color.New(color.BgWhite).Add(color.Bold)
	//go:embed wordlist.txt
	wordlistData string
	wordlist     = strings.Fields(wordlistData)
)

type MessageType int

const (
	Info MessageType = iota
	Warning
	Error
	Success
	Reg
)

func main() {

	for {
		var option int
		Blue.Println()
		colors(Info, "1: Work Notes \n2: Random Password\n3: Random Alphanumeric Password\n4: Open Applications\n5: Screenshot \n6: Network tools \n7: File Options \n8 Search Options ")
		fmt.Scanln(&option)
		toolControl(option)

	}

}

func toolControl(option int) {
	options := map[int]func(){
		1: func() {
			notes := notes_list()
			titles := listTitle()
			chooseNote(notes, titles)
		},

		2: func() {
			randPass()
		},
		3: func() {
			ranAlpPass()
		},
		4: func() {
			information()
			chooseApp()
		},
		5: func() {
			numScreenshots := 3
			screenCapture(numScreenshots)
		},
		6: func() {
			netTools()
		},
		7: func() {
			FileOptions()
		},
		8: func() {
			fileSearch()
		},
	}
	if action, ok := options[option]; ok {
		action()
	} else {
		colors(Error, "Invalid option, choose another option")
	}

}

func notes_list() []string {

	notesDir := "Work_Notes"
	nList, err := filepath.Glob(filepath.Join(notesDir, "*"))
	if err != nil {
		log.Fatal(err)
	}

	var notes []string
	for _, file := range nList {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		notes = append(notes, string(content))
	}
	return notes
}

func listTitle() []string {

	notesDir := "Work_Notes"
	nList, err := filepath.Glob(filepath.Join(notesDir, "*"))
	if err != nil {
		log.Fatal(err)
	}

	var titles []string
	for _, file := range nList {
		titles = append(titles, file)
	}
	return titles
}

func chooseNote(notes []string, titles []string) {
	colors(Success, "Available Notes: ")
	for i, title := range titles {
		fmt.Printf("%d. %s\n", i+1, extractFileName(title))
	}
	selectedNoteIndex := userInput("Enter the number of the note: ")
	num, err := strconv.Atoi(selectedNoteIndex)
	if err != nil || num < 1 || num > len(titles) {

		colors(Error, "Invalid input. Please enter a valid note number.")
		return
	}

	selectedNote := notes[num-1]
	colors(Info, fmt.Sprintf("You chose note number: %d", num))
	colors(Info, "Note Content: ")
	colors(Reg, selectedNote)

	copyChoice := userInput("Copy this note to clipboard? (y/n): ")
	if strings.ToLower(copyChoice) == "y" {
		err := clipboard.WriteAll(selectedNote)
		if err != nil {
			log.Fatal(err)
		}
		colors(Success, "Note copied to clipboard!")
	} else {
		colors(Error, "Note not copied to clipboard.")
	}

}

func userInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func extractFileName(notes string) string {
	baseName := filepath.Base(notes)

	baseNameWithoutExtension := strings.TrimSuffix(baseName, ".txt")

	parts := strings.Split(baseNameWithoutExtension, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return baseNameWithoutExtension
}

//Networking Tools

func netTools() {
	colorGreen := "\033[32m"
	var app int
	i := 2
	fmt.Print(string(colorGreen), " You have ", i, " options", "\n")
	colors(Info, "Choose which networking tool to use: \n1: Ping a hose \n2: traceroute to host")
	fmt.Scanln(&app)

	var host string
	colors(Success, "Enter the hostname/IP address: ")
	fmt.Scanln(&host)

	var cmd *exec.Cmd
	switch app {
	case 1:
		cmd = exec.Command("ping", host)
	case 2:
		if runtime.GOOS == "windows" {
			cmd = exec.Command("tracert", host)
		} else {
			cmd = exec.Command("traceroute", host)
		}
	default:
		colors(Error, "Invalid option")
		return
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		colors(Error, fmt.Sprintf("Error executing command: %s\n", err))
		return
	}
	fmt.Println(string(output))
}

//Secruity Tools
func generatePassword(length int, includeNumber bool, includeSpecial bool) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var password []byte
	var charSource string

	if includeNumber {
		charSource += "0123456789"
	}
	if includeSpecial {
		charSource += "!@#$%^&*()_+=-"
	}
	charSource += charset

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < length; i++ {
		password = append(password, charSource[b[i]%byte(len(charSource))])
	}
	return string(password)
}
func generateAlphaPass(wordlist []string, length int, includeNumber bool, includeSpecial bool, includeCap bool) string {
	// Combine the character sources
	var charSource string
	if includeNumber {
		charSource += "0123456789"
	}
	if includeSpecial {
		charSource += "!@#$%^&*()_+=-"
	}
	if includeCap {
		charSource += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	// Use crypto/rand for better randomness
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	var passwordBuilder strings.Builder
	for i := 0; i < length; i++ {
		// Randomly choose to append a word or a character from charSource
		if rand.Intn(2) == 0 && len(wordlist) > 0 {
			// Append a random word
			passwordBuilder.WriteString(wordlist[rand.Intn(len(wordlist))])
		} else {
			// Append a random character
			passwordBuilder.WriteByte(charSource[b[i]%byte(len(charSource))])
		}

		// Stop if the password has reached the desired length
		if passwordBuilder.Len() >= length {
			break
		}
	}

	// Truncate the password if it exceeds the desired length
	password := passwordBuilder.String()
	if len(password) > length {
		password = password[:length]
	}

	return password
}
func openCWM() {
	var cwPath string = "C:/Program Files (x86)/ConnectWise/PSA.net/ConnectWiseManage.exe"
	o := exec.Command("cmd", "/C", "start", cwPath)
	if err := o.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
	}

}

//ScreenCapture
func screenCapture(numScreenshots int) (string, error) {

	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		return "", errors.New("no active displays found")
	}

	bounds := screenshot.GetDisplayBounds(0)

	for i := 1; i <= numScreenshots; i++ {
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		filename := fmt.Sprintf("screenshot_%d.png", i)
		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			panic(err)
		}

	}

	return Blue.Sprintf("%d screenshots saved", numScreenshots), nil
}

//Application Launcher
func openChrome(url string) bool {
	url = " "
	var o []string
	switch runtime.GOOS {
	case "darwin":
		o = []string{"open"}
	case "windows":
		o = []string{"cmd", "/c", "start", "Chrome"}
	default:
		o = []string{"xdg-open"}
	}

	cmd := exec.Command(o[0], append(o[1:], url)...)
	if err := cmd.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
	}
	return cmd.Start() == nil

}

func openOutlook() {
	o := exec.Command("cmd", "/C", "start", "Outlook")
	if err := o.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
	}

}

func openRDM() {
	o := exec.Command("cmd", "/C", "start", "RemoteDesktopManager")
	if err := o.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
	}

}

func openCWA() {
	o := exec.Command("cmd", "/C", "start", "LTClient")
	if err := o.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
	}

}

func openSnip() {
	o := exec.Command("cmd", "/C", "start", "SnippingTool")
	if err := o.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
	}
}

func openCode() {
	o := exec.Command("cmd", "/C", "start", "Code")
	if err := o.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
	}
}

func openWord() {
	o := exec.Command("cmd", "/C", "start", "Word")
	if err := o.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
	}
}

func openExcel() {
	o := exec.Command("cmd", "/C", "start", "Excel")
	if err := o.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
	}
}

func openPowershell() {
	o := exec.Command("cmd", "/C", "start", "Powershell")
	if err := o.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
	}
}

func openCMD() {
	o := exec.Command("cmd", "/C", "start", "CMD")
	if err := o.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
	}
}

func openTeams() {
	var tPath string = "C:/Users/cgrissette/AppData/Local/Microsoft/Teams/current/Teams"
	o := exec.Command("cmd", "/C", "start", tPath)
	if err := o.Start(); err != nil {
		colors(Error, "Error:"+err.Error())
		colors(Error, "Note not copied to clipboard.")
	}

}

func chooseApp() {
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	var app int
	i := 11
	message := fmt.Sprintf("You have %d options\n", i)
	colors(Info, message)
	colors(Reg, "Choose which application to open:")
	fmt.Scanln(&app)
	switch app {
	case 1:
		openOutlook()
		fmt.Println(string(colorGreen), "Outlook is open", string(colorReset))
	case 2:
		openCWA()
		fmt.Println(string(colorGreen), "Automate is open", string(colorReset))
	case 3:
		openChrome("https://360smartnet.itglue.com/")
		fmt.Println(string(colorGreen), "Chrome is open", string(colorReset))
	case 4:
		openRDM()
		fmt.Println(string(colorGreen), "Remote Desktop Manager is open", string(colorReset))
	case 5:
		openSnip()
		fmt.Println(string(colorGreen), "Snipping Tool is open", string(colorReset))
	case 6:
		openCode()
		fmt.Println(string(colorGreen), "Visusal Studio is open", string(colorReset))
	case 7:
		openWord()
		fmt.Println(string(colorGreen), "Word is open", string(colorReset))
	case 8:
		openExcel()
		fmt.Println(string(colorGreen), "Excel is open", string(colorReset))
	case 9:
		openPowershell()
		fmt.Println(string(colorGreen), "Powershell is open", string(colorReset))
	case 10:
		openCMD()
		fmt.Println(string(colorGreen), "CMD is open", string(colorReset))
	case 0:
		appStart()
		fmt.Println(string(colorGreen), "All applications are open", string(colorReset))

	}

}

func appStart() {
	openRDM()
	openOutlook()
	//openCWM()
	openChrome("https://360smartnet.itglue.com/")
	openCWA()
	openSnip()
	openCode()
	openWord()
	openExcel()
	openPowershell()
	openCMD()

}

func information() {
	colorRed := "\033[31m"
	colorReset := "\033[0m"

	fmt.Println(string(colorRed), "1: Outlook \n 2: Automate \n 3: Chrome \n 4: Remote Desktop Manager\n 5: Snipping Tool\n 6: Open Visual Studio\n 7: Word\n 8: Excel\n 9: Powershell\n 10: CMD\n 0: All Applications", string(colorReset))
}

//File Managemant Tools
func bulkRename(files map[string]string) []error {
	var errors []error
	successCount := 0

	for src, dst := range files {
		if _, err := os.Stat(src); os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("source file does not exist: %s", src))
			continue
		}
		if _, err := os.Stat(dst); err == nil {
			errors = append(errors, fmt.Errorf("destination file already exists: %s", dst))
			continue
		}
		if err := os.Rename(src, dst); err != nil {
			errors = append(errors, err)
			continue
		}
		successCount++
	}

	if successCount > 0 {
		Blue.Printf("%d files successfully renamed\n", successCount)

	}
	return errors
}

func fileSearch() {
	var dirname, searchQuery, searchType string
	var matches []string
	var err error

	colors(Reg, "Enter the directory to search: ")
	fmt.Scanln(&dirname)

	colors(Reg, "Enter search query (pattern or exact filename): ")
	fmt.Scanln(&searchQuery)

	colors(Reg, "Search by (1) Pattern or (2) Exact Filename: ")
	fmt.Scanln(&searchType)

	switch searchType {
	case "1": // Pattern search
		matches, err = patternSearch(dirname, searchQuery)
	case "2": // Exact filename search
		matches, err = exactSearch(dirname, searchQuery)
	default:
		colors(Error, "Invalid search type selected")
		return
	}

	if err != nil {
		colors(Error, "Error during search:"+err.Error())
		return
	}

	if len(matches) == 0 {
		colors(Error, "No files found")
	} else {
		colors(Success, "Found files:")
		for _, match := range matches {
			colors(Success, match)
		}
	}
}

func patternSearch(dirname, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if matched, _ := filepath.Match(pattern, info.Name()); matched {
			matches = append(matches, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return matches, nil
}

func exactSearch(dirname, filename string) ([]string, error) {
	var matches []string
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == filename {
			matches = append(matches, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return matches, nil
}

// Enhanced DirList Function
func DirList(dirname string) ([]string, error) {
	dir, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, len(fileInfos))
	for i, file := range fileInfos {
		fileNames[i] = file.Name()
	}
	return fileNames, nil
}

// Improved WriteFile Function with Clear Error Messages
func WriteFile(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		colors(Error, fmt.Sprintf("error creating file: %v", err))
	}
	defer file.Close()

	decodedContent, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		colors(Error, fmt.Sprintf("error decoding content: %v", err))
	}

	_, err = file.Write(decodedContent)
	if err != nil {
		colors(Error, fmt.Sprintf("error writing to file: %v", err))
	}

	return nil
}

func ReadFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		colors(Error, fmt.Sprintf("Error opening file: %v", err))
	}
	defer file.Close()

	// For larger files, consider reading in chunks or using ioutil.ReadAll
	content, err := ioutil.ReadAll(file)
	if err != nil {
		colors(Error, fmt.Sprintf("Error reading file: %v", err))
	}

	return base64.StdEncoding.EncodeToString(content), nil
}

func FileSize(path string) (string, error) {
	fInfo, err := os.Stat(path)

	if err != nil {
		log.Fatal(err)

	}
	fsize := fInfo.Size()

	errorMessage := fmt.Sprintf("The file size is %s\n", formatBytes(fsize))
	colors(Error, errorMessage)

	return fmt.Errorf(errorMessage)
}

// Optimized FileOptions Function with Enhanced UI
func FileOptions() {
	var option int
	var file string

	options := []string{
		"Get file size",
		"Read from file",
		"Write to file",
		"Exit",
	}

	colors(Reg, "Choose your file option: ")
	for i, opt := range options {
		Green.Printf("%d. %s\n", i+1, opt)
	}

	fmt.Scanln(&option)
	if option < 1 || option > len(options) {
		colors(Error, "Invalid option")
		return
	}

	switch option {
	case 1:

		colors(Reg, "Enter the name of the file: ")
		fmt.Scanln(&file)
		size, err := FileSize(file)
		if err != nil {
			colors(Error, err.Error())
			return
		}
		Blue.Printf("%s\n", size)
	case 2:
		colors(Reg, "Enter the name of the file to read: ")
		fmt.Scanln(&file)
		read, err := ReadFile(file)
		if err != nil {
			colors(Error, err.Error())
			return
		}
		colors(Success, fmt.Sprintf("contents of the file: %s", read))

	case 3:
		colors(Reg, "Enter the name of the file to write to: ")
		fmt.Scanln(&file)
		colors(Reg, "Enter the content to write: ")
		var content string
		fmt.Scanln(&content)
		err := WriteFile(file, content)
		if err != nil {
			colors(Error, "Error writing to file:"+err.Error())
			return
		}
		colors(Success, "Content written to file successfully.")

	default:

	}

}

func formatBytes(bytes int64) string {
	units := []string{"bytes", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
	const K = 1024.0

	if bytes == 0 {
		return "0 bytes"
	}

	i := math.Floor(logN(float64(bytes), K))
	if i == 0 {
		return strconv.FormatInt(bytes, 10) + " " + units[int(i)]
	}

	//For the units greater than a byte, format to 3 decmial place
	size := float64(bytes) / math.Pow(K, i)
	return fmt.Sprintf("%.3f %s", size, units[int(i)])

}

// logN calculates the logarithm of a number with a specific base, used to determine the unit scale
func logN(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func randPass() {
	var num int
	var pass int
	colors(Reg, "Press 1 for weak password:\nPress 0 for strong password:  ")
	fmt.Scanln(&pass)
	colors(Reg, "How many Chars (20-Max): ")
	fmt.Scanln(&num)
	if pass > 0 {
		password := generatePassword(num, true, false)
		Blue.Println("Weak Random password:", password)
		if num > 20 || num < 0 {
			colors(Error, "invalid input can not be below 0 or above 20")
		}
		copyChoice := userInput("Copy this password to clipboard? (y/n): ")
		if strings.ToLower(copyChoice) == "y" {
			err := clipboard.WriteAll(password)
			if err != nil {
				log.Fatal(err)
			}
			colors(Success, "Note copied to clipboard!")
		} else {
			colors(Error, "Note not copied to clipboard.")
		}
	} else if pass < 1 {
		password := generatePassword(num, true, true)
		Blue.Println("Random password:", password)
		if num > 20 || num < 0 {
			colors(Error, "invalid input can not be below 0 or above 20")
		}
		copyChoice := userInput("Copy this password to clipboard? (y/n): ")
		if strings.ToLower(copyChoice) == "y" {
			err := clipboard.WriteAll(password)
			if err != nil {
				log.Fatal(err)
			}
			colors(Success, "Note copied to clipboard!")
		} else {
			colors(Error, "Note not copied to clipboard.")
		}
	}
}

func ranAlpPass() {
	var alp int
	var pass int
	colors(Reg, "Press 1 for weak password:\nPress 0 for strong password:  ")
	fmt.Scanln(&pass)
	colors(Reg, "How many Chars (20-Max): ")
	fmt.Scanln(&alp)
	if pass > 0 {
		ranPassword := generateAlphaPass(wordlist, alp, true, false, false)

		colors(Success, fmt.Println(" Weak Random Alphanumeric password:", ranPassword))
		if alp > 20 || alp < 0 {
			colors(Error, "invalid input can not be below 0 or above 20")
		}
		copyChoice := userInput("Copy this password to clipboard? (y/n): ")
		if strings.ToLower(copyChoice) == "y" {
			err := clipboard.WriteAll(ranPassword)
			if err != nil {
				log.Fatal(err)
			}
			colors(Success, "Note copied to clipboard!")
		} else {
			colors(Error, "Note not copied to clipboard.")
		}

	} else if pass < 1 {
		ranPassword := generateAlphaPass(wordlist, alp, true, true, true)
		Blue.Println("Strong Random Alphanumeric password:", ranPassword)
		if alp > 20 || alp < 0 {
			colors(Error, "invalid input can not be below 0 or above 20")
		}

		copyChoice := userInput("Copy this password to clipboard? (y/n): ")
		if strings.ToLower(copyChoice) == "y" {
			err := clipboard.WriteAll(ranPassword)
			if err != nil {
				log.Fatal(err)
			}
			colors(Success, "Note copied to clipboard!")
		} else {
			colors(Error, "Note not copied to clipboard.")
		}
	}
}

func colors(MessageType MessageType, message string) {
	switch MessageType {
	case Info:
		Blue.Println(message)
	case Warning:
		Magenta.Println(message)
	case Error:
		Red.Println(message)
	case Success:
		Green.Println(message)
	case Reg:
		White.Println(message)
	default:
		fmt.Println(message)
	}
}
