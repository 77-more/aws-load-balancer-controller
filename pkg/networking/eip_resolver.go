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
	//var unavailableEIPs []string
	//var err error
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
			// I see below error message before the custom message I created below. How to supress it?
                        // The allocation IDs are not available for use status code: 400, request id: e89d6089-dc18-42e6-8400-d620a7845ad4
	
                        if err != nil {
				return nil, err
			}
			if results.Addresses == nil {
				return nil, errors.Errorf("EIP by the name %s not found", nameOrIDs)
			} else if len(results.Addresses) > 1 {
				return nil, errors.Errorf("There are multiple EIPs with the name %s, please use a assign a unique name to the EIP that you want to assign to the load balancer", nameOrIDs)
			} else if results.Addresses[0].AssociationId != nil {
				return nil, errors.Errorf("EIP by the name %s is in use already, please provide a different EIP name or allocation ID", nameOrIDs)
			} else {
				singleallocationID := *results.Addresses[0].AllocationId
				allocationIDs = append(allocationIDs, singleallocationID)
			}
		}
	}
  return allocationIDs, err
}
