package networking

import (
	"testing"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

type MockEC2API struct{}

func (m *MockEC2API) DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
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
	mockEC2 := &MockEC2API{}
	resolver := EIPResolver{}
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
			resultIDs, err := resolver.EIPResolver(tc.input)

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
