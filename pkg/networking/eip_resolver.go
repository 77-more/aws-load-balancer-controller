package networking

import (
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"*/
	"github.com/pkg/errors"
	
)

func EIPResolver (eipAllocation []string) ([]string, error) {
	
	var allocationIDs []string
        sess, _ := session.NewSession()
	ec2svc := ec2.New(sess)
	// Used a for loop with if else, that way a user can use a combination of allocation IDs and EIP names. 
	for _, nameOrIDs := range eipAllocation {
		// Under if condition we check for the allocation IDs and append them to allocationIDs variable. 
                // Under else condition we process EIP names and check if they are unique, if they are already in use, if the EIP name exists in the account at all and, if there are multiple EIPs with the same name. If none of these conditions are true we return the allocation IDs for the particular EIP names.
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
			} 
			if len(results.Addresses) > 1 {
				return nil, errors.Errorf("There are multiple EIPs with the same name %s, please assign a unique name to the EIP that you want to assign to the load balancer", nameOrIDs)
			}
			if results.Addresses[0].AssociationId != nil {
				return nil, errors.Errorf("EIP by the name %s is in use already, please use an EIP that is available to use", nameOrIDs)
			} else {
				allocationIDs = append(allocationIDs, *results.Addresses[0].AllocationId)
			}
		}
	}
  return allocationIDs, nil
}
