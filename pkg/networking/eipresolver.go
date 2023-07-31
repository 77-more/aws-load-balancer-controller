package networking

import (
	//"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"sigs.k8s.io/aws-load-balancer-controller/pkg/aws/services"
	//"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)
type defaultEIPInfoProvider struct {
	ec2Client        services.EC2
}

var availableEIPs []string
var unavailableEIPs []string
func (p *defaultEIPInfoProvider) eipresolver(EIPnameOrIDs []string) ([]string) {
	// Creates session object
	//sess, _ := session.NewSession()
	// opens a new session
	//ec2svc := ec2.New(sess)
	// makes DescribeAddresses api call and stores the output of type DescribeAddressesOutput in results variable. As part of the API call we are looking for Name: test1 and cluster-name:test tags
	// &ec2.DescribeAddressesInput represents the memory address of the
	for _, nameOrIDs := range EIPnameOrIDs {
		results, _ := p.ec2Client.DescribeAddresses(&ec2.DescribeAddressesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: aws.StringSlice([]string{nameOrIDs}),
				},
			},
		})
		
		var allocationIDs []string
		var associationIDs []string

		allocationIDs = append(allocationIDs, *results.Addresses[0].AllocationId)
		associationIDs = append(associationIDs, *results.Addresses[0].AssociationId)
		if len(associationIDs) > 0 {
			unavailableEIPs = append(unavailableEIPs, *results.Addresses[0].AllocationId)
		} else {
			availableEIPs = append(availableEIPs, *results.Addresses[0].AllocationId)
		}
		associationIDs = nil
	
		} 
	//	fmt.Printf("available EIPs:%v, unavailableEIPs:%v", availableEIPs, unavailableEIPs)
        
        return availableEIPs
}
