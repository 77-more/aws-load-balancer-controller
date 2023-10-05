package networking

import (
    "testing"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/pkg/errors"
)

func initAWSSession() *session.Session {
    // Initialize an AWS session with a specific region
    awsSession, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"), 
    })

    if err != nil {
        panic("Failed to create AWS session: " + err.Error())
    }

    return awsSession
}

func DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
    // Initialize an AWS session
    awsSession := initAWSSession()

    // Simulate the behavior of DescribeAddresses based on the test case
    if len(input.Filters) == 1 && *input.Filters[0].Values[0] == "existing-eip" {
        return &ec2.DescribeAddressesOutput{
            Addresses: []*ec2.Address{
                {
                    AllocationId: aws.String("allocation-id"),
                    AssociationId: aws.String("association-id"),
                },
            },
        }, nil
    }
    return nil, errors.New("EIP not found")
}




func TestEIPResolver(t *testing.T) {
	tests := []struct {
		name           string
		input          []string
		expectedResult []string
		expectedError  error
	}{
		{
			name:           "Resolve EIP by Allocation ID",
			input:          []string{"eipalloc-1"},
			expectedResult: []string{"eipalloc-1"},
			expectedError:  nil,
		},
		{
			name:           "Resolve EIP by Name",
			input:          []string{"existing-eip"},
			expectedResult: []string{"allocation-id"},
			expectedError:  nil,
		},
		{
			name:           "EIP not found",
			input:          []string{"non-existent-eip"},
			expectedResult: nil,
			expectedError:  errors.New("EIP not found"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resultIDs, err := EIPResolver(tc.input)

			if err != nil && tc.expectedError == nil {
				t.Errorf("Expected no error, but got error: %v", err)
			}
			if err == nil && tc.expectedError != nil {
				t.Errorf("Expected error: %v, but got no error", tc.expectedError)
			}
			if err == nil && tc.expectedError == nil && !stringSliceEqual(resultIDs, tc.expectedResult) {
				t.Errorf("Expected result: %v, but got result: %v", tc.expectedResult, resultIDs)
			}
		})
	}
}

func stringSliceEqual(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
