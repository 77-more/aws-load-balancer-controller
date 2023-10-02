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
			// if the EIP that customer wants to use is already associated with another component, then the EIP will have an association ID, it that case we can error out. 
			if *results.Addresses[0].AssociationId != "" {
				return nil, errors.Errorf("EIP by the name %s is in use already, please provide a different EIP name or allocation ID", nameOrIDs)
			}
			if err != nil {
				return nil, err
			}
			if results.Addresses == nil {
				return nil, errors.Errorf("EIP with the name %s is not found, please provide a valid EIP name",nameOrIDs)
			} else {
				singleallocationID := *results.Addresses[0].AllocationId
				allocationIDs = append(allocationIDs, singleallocationID)
			}
		}
	}
  return allocationIDs, err
}
