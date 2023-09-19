package networking

import (
	//"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	//"sigs.k8s.io/aws-load-balancer-controller/pkg/aws/services"
	//"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"strings"
)
//type eipinfoprovider struct {
//	ec2Client        services.EC2
//}

func EIPResolver (EIPnameOrIDs []string) []string {
	// Creates session object
	sess, _ := session.NewSession()
	// opens a new session
	ec2svc := ec2.New(sess)
	// makes DescribeAddresses api call and stores the output of type DescribeAddressesOutput in results variable.
	//As part of the API call we are looking for Name: test1 and cluster-name:test tags
	// &ec2.DescribeAddressesInput represents the memory address of the
	var availableEIPs []string
	fmt.Println(EIPnameOrIDs)
	if len(EIPnameOrIDs) == 0 {
		fmt.Println("returning from line 27")
		return availableEIPs
	}
	for _, nameOrIDs := range EIPnameOrIDs {
		if strings.HasPrefix(nameOrIDs, "eipalloc-") {
			results, _ := ec2svc.DescribeAddresses(&ec2.DescribeAddressesInput{
				Filters: []*ec2.Filter{
					{
						Name:   aws.String("allocation-id"),
						Values: aws.StringSlice([]string{nameOrIDs}),
					},
				},
			})
			if len(results.Addresses) == 0 {
				continue
			}
			if results.Addresses[0].AssociationId == nil {
				availableEIPs = append(availableEIPs, nameOrIDs)
			}
		}
		results, _ := ec2svc.DescribeAddresses(&ec2.DescribeAddressesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: aws.StringSlice([]string{nameOrIDs}),
				},
			},
		})
		//if the region is set wrong then the results.Addresses len will be 0 so we need to account for that condition as well.
		if len(results.Addresses) == 0 {
			continue
		}
		// *results.Addresses[0].AllocationId is the pointer to the value in AllocationId
		// results.Addresses[0].AllocationId is the address of the AllocationId field
		allocationIDs := *results.Addresses[0].AllocationId
		// the below if condition checks if the AssociationId field exists or no NOTE we are not looking for the address of the field or the value in the address, we are simply looking for the existence of the field itself.

		if results.Addresses[0].AssociationId == nil {
			availableEIPs = append(availableEIPs, allocationIDs)
		}
	}
	fmt.Println("returning from line 67")
	return availableEIPs
}

