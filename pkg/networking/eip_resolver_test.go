package networking

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/ec2"
)

// MockEC2API is a mock implementation of the EC2 API
type MockEC2API struct {
    mock.Mock
}

func (m *MockEC2API) DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
    args := m.Called(input)
    return args.Get(0).(*ec2.DescribeAddressesOutput), args.Error(1)
}

func TestEIPResolver(t *testing.T) {
    // Create a new instance of the mock EC2 API
    mockEC2 := new(MockEC2API)

    // Initialize your EIPResolver function with the mockEC2 instance
    resolver := EIPResolver{}

    // Define the test data (input EIP names)
    inputNames := []string{"eip-1", "eip-2"}

    // Set up mock behavior for DescribeAddresses
    mockEC2.On("DescribeAddresses", mock.Anything).Return(
        &ec2.DescribeAddressesOutput{
            Addresses: []*ec2.Address{
                {AllocationId: aws.String("alloc-1")},
                {AllocationId: aws.String("alloc-2")},
            },
        },
        nil,
    )

    // Call your function under test
    resultIDs, err := resolver.EIPResolver(inputNames)

    // Assertions
    assert.Nil(t, err)  // Check that there's no error
    assert.Equal(t, []string{"alloc-1", "alloc-2"}, resultIDs) // Check that the result matches the expected IDs

    // Ensure that the DescribeAddresses method was called with the expected input
    mockEC2.AssertCalled(t, "DescribeAddresses", &ec2.DescribeAddressesInput{
        Filters: []*ec2.Filter{
            {
                Name:   aws.String("tag:Name"),
                Values: aws.StringSlice(inputNames),
            },
        },
    })
}
