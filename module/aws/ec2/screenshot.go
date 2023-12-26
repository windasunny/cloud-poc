package aws

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
)

func Screenshot(accessKeyId string, secretAccessKey string) {
	// AWS Session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, ""),
	})

	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	// EC2 Service Client
	svc := awsec2.New(sess)

	// EC2 instance ID
	instanceID := "i-09d20a0b2b49d77cc"

	// Get console screenshot
	input := &awsec2.GetConsoleScreenshotInput{
		InstanceId: aws.String(instanceID),
	}

	// Get screenshot result
	result, err := svc.GetConsoleScreenshot(input)
	if err != nil {
		fmt.Println("Error getting console screenshot:", err)
		return
	}

	// Decode screenshot
	imageData := *result.ImageData
	imgBytes, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		fmt.Println("Error decoding image data:", err)
		return
	}

	// Build image buffer
	imgBuf := bytes.NewReader(imgBytes)

	// Decode image
	img, err := jpeg.Decode(imgBuf)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Build image file
	imgFile, err := os.Create("screenshot.jpg")
	if err != nil {
		fmt.Println("Error creating image file:", err)
		return
	}
	defer imgFile.Close()

	// Encode image to file
	jpeg.Encode(imgFile, img, nil)

	fmt.Println("Screenshot saved as screenshot.png")
}
