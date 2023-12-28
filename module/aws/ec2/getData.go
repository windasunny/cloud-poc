package aws

import (
	"fmt"

	awsec2 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gookit/color"
)

func (ec2 *Ec2) SearchService() {

	color.Magenta.Println("S3 Bucket:")
	findS3(ec2.AccessKeyId, ec2.SecretAccessKey, ec2.Token)
	color.Magenta.Println("Ec2 Instances:")
	findEC2(ec2.AccessKeyId, ec2.SecretAccessKey, ec2.Token)

}

func findS3(accessKeyId string, secretAccessKey string, token string) {
	sess := buildSession(accessKeyId, secretAccessKey, token)

	svc := s3.New(sess)
	input := &s3.ListBucketsInput{}

	result, err := svc.ListBuckets(input)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("---------------------------------")

	for _, bucket := range result.Buckets {
		fmt.Println(*bucket.Name)
	}

	fmt.Println("---------------------------------")
}

func findEC2(accessKeyId string, secretAccessKey string, token string) {
	sess := buildSession(accessKeyId, secretAccessKey, token)

	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		fmt.Println("Error describing EC2 instances:", err)
		return
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			// instance id
			fmt.Println("Instance ID:", *instance.InstanceId)

			// instance name
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					fmt.Println("Instance Name:", *tag.Value)
				}
			}

			// instance type
			fmt.Println("Instance Type:", *instance.InstanceType)

			//instance public ip
			if instance.PublicIpAddress != nil {
				fmt.Println("Public IP Address:", *instance.PublicIpAddress)
			}

			// instance private ip
			if instance.PrivateIpAddress != nil {
				fmt.Println("Private IP Address:", *instance.PrivateIpAddress)
			}

			// instance role
			if instance.IamInstanceProfile != nil {
				fmt.Println("IAM Role:", *instance.IamInstanceProfile.Arn)
			}

			// instance security group
			fmt.Println("Security Groups:")
			for _, sg := range instance.SecurityGroups {
				fmt.Println("  -", *sg.GroupName)
			}
		}
		fmt.Println("================================")
	}

	fmt.Println("---------------------------------")
}

func buildSession(accessKeyId string, secretAccessKey string, token string) *session.Session {
	sess := session.Must(session.NewSession(&awsec2.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, token),
		Region:      awsec2.String("us-east-1"),
	}))

	return sess
}

// lamdba
