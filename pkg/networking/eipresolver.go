package networking

import (
	"fmt"
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

func EIPResolver(EIPnameOrIDs []string) []string {
	sess, _ := session.NewSession()
	ec2svc := ec2.New(sess)
	var returningEIPs []string
	if len(EIPnameOrIDs) == 0 {
		return returningEIPs
	}
	for _, nameOrIDs := range EIPnameOrIDs {
		if strings.HasPrefix(nameOrIDs, "eipalloc-") {
				returningEIPs = append(returningEIPs, nameOrIDs)
		} else {
		results, _ := ec2svc.DescribeAddresses(&ec2.DescribeAddressesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: aws.StringSlice([]string{nameOrIDs}),
				},
			},
		})
		allocationIDs := *results.Addresses[0].AllocationId
		if results.Addresses[0].AssociationId == nil {
			returningEIPs = append(returningEIPs, allocationIDs)
		} else {
			returningEIPs = append(returningEIPs, allocationIDs)
		  }
		}
	}
	return returningEIPs
}
