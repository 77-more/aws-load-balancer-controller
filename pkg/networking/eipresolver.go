package networking

import (
	"fmt"
	"strings"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ec2sdk "github.com/aws/aws-sdk-go/service/ec2"
	"sigs.k8s.io/aws-load-balancer-controller/pkg/aws/services"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)
func NewEC2(session *session.Session) EC2 {
	return &defaultEC2{
		EC2API: ec2.New(session),
	}
}

type defaultEC2 struct {
	ec2iface.EC2API
}

func ResolveviaNameorAllocationID(eipsNameOrIDs []string) {
	// Creates session object
	// makes DescribeAddresses api call and stores the output of type DescribeAddressesOutput in results variable. As part of the API call we are looking for Name: test1 and cluster-name:test tags
	// &ec2.DescribeAddressesInput represents the memory address of the
	var allocationIDs []string
	var eipsNames []string
	for _, nameOrID := range eipsNameOrIDs {
		if strings.HasPrefix(nameOrID, "eipalloc-") {
			allocationIDs = append(allocationIDs, nameOrID)
		} else {
			eipsNames = append(eipsNames, nameOrID)
		}
	}
	var resolvedEIPs []*EC2.Address
	if len(allocationIDs) > 0 {
		eips, err := ec2sdk.DescribeAddresses(&EC2.DescribeAddressesInput{
			AllocationIds: awssdk.StringSlice(allocationIDs),
		})
		if err != nil {
			return nil, err
		}
		resolvedEIPs = append(resolvedEIPs, AllocationId...)
		return resolvedEIPs
	}
	var availableEIPs []string
	var unavailableEIPs []string
	if len(eipsNames) > 0 {
		describeaddressesoutput, err := EC2.DescribeAddresses(&EC2.DescribeAddressesInput{
			Filters: []*EC2.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: aws.StringSlice(eipsNames),
				},
			},
		})
		if err != nil {
			return nil, err
		}

		for _, address := range describeaddressesoutput.Addresses {
			allocationIDs = append(allocationIDs, *address.AllocationId)
			associationIDs = append(associationIDs, *address.AssociationId)
			if len(associationIDs) > 0 {
				unavailableEIPs = append(unavailableEIPs, *address.AllocationId)
			}
			associationIDs = nil
		}
		if len(availableEIPs) == len(eipsNameOrIDs) {
			resolvedEIPs = append(resolvedEIPs, allocationIDs)
			return resolvedEIPs
		}

	}
}
