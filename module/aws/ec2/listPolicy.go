package aws

import (
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/gookit/color"
)

func (ec2 *Ec2) ListPolicy() {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(ec2.AccessKeyId, ec2.SecretAccessKey, ec2.Token),
		Region:      aws.String(ec2.Region),
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	svc := iam.New(sess)

	roleName := "AmazonSSMRoleForInstancesQuickSetup"

	// Get Policy

	// Inline Policy
	inlineInput := &iam.ListRolePoliciesInput{
		RoleName: aws.String(roleName),
	}

	inlineResult, err := svc.ListRolePolicies(inlineInput)
	if err != nil {
		fmt.Println("Error listing role policies:", err)
		return
	}

	color.Magenta.Println("Inline Policies:")
	for _, policyName := range inlineResult.PolicyNames {

		getPolicyInput := &iam.GetRolePolicyInput{
			RoleName:   aws.String(roleName),
			PolicyName: policyName,
		}

		policyResult, err := svc.GetRolePolicy(getPolicyInput)
		if err != nil {
			fmt.Println("Error getting role policy:", err)
			return
		}

		policyDocument, err := url.QueryUnescape(*policyResult.PolicyDocument)
		if err != nil {
			fmt.Println("Error decoding policy document:", err)
			return
		}

		fmt.Println("Policy Name:", *policyResult.PolicyName)
		fmt.Println("Policy Document:", policyDocument)
	}
	fmt.Println("--------------------------------")
	color.Magenta.Println("Managed Policies:")

	// Managed Policy
	managerInput := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	}

	managerResult, err := svc.ListAttachedRolePolicies(managerInput)
	if err != nil {
		fmt.Println("Error listing role policies:", err)
	}
	for _, policy := range managerResult.AttachedPolicies {
		fmt.Println("Policy ARN:", *policy.PolicyArn)

		input := &iam.GetPolicyInput{
			PolicyArn: aws.String(*policy.PolicyArn),
		}

		policyResult, err := svc.GetPolicy(input)
		if err != nil {
			fmt.Println("Error getting policy:", err)
			return
		}

		fmt.Println("Default Policy Version ID:", *policyResult.Policy.DefaultVersionId)

		versionInput := &iam.GetPolicyVersionInput{
			PolicyArn: aws.String(*policy.PolicyArn),
			VersionId: aws.String(*policyResult.Policy.DefaultVersionId),
		}

		versionResult, err := svc.GetPolicyVersion(versionInput)
		if err != nil {
			fmt.Println("Error getting policy version:", err)
			continue
		}

		fmt.Println("Policy Document:")
		policyDocument, err := url.QueryUnescape(*versionResult.PolicyVersion.Document)
		fmt.Println(policyDocument)
	}
}
