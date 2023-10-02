package networking

import (
	//"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"strings"
)

func EIPResolver (eipAllocation []string) ([]string, error) {
	var allocationIDs []string
	var err error
        sess, _ := session.NewSession()
	ec2svc := ec2.New(sess)
	for _, nameOrIDs := range eipAllocation {
		if strings.HasPrefix(nameOrIDs, "eipalloc-") {
			allocationIDs = append(allocationIDs, nameOrIDs)
		} else {
			results, err := ec2svc.DescribeAddresses(&ec2.DescribeAddressesInput{
				Filters: []*ec2.Filter{
					{
						Name:   aws.String("tag:Name"),
						Values: aws.StringSlice([]string{nameOrIDs}),
					},
				},
			})
			// if there are no EIPs by the name that is provided, then results.Addresses will be equal to nil so we compare results.Addresses to nil to check for this condition.
			if err != nil {
				return nil, err
			}
			if results.Addresses == nil {
				return nil, errors.Errorf("EIP %s is not found, please provide a valid EIP name",nameOrIDs)
			} else {
				singleallocationID := *results.Addresses[0].AllocationId
				allocationIDs = append(allocationIDs, singleallocationID)
			}
		}
	}
  return allocationIDs, err
}
