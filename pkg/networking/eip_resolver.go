package networking


import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"strings"
)
// testing
var err error

type Resolver interface {
	EIPResolver() ([]string, error)
}

type EIPnameOrIDs struct {
	InputEIPnameOrIDs []string
}

func (e EIPnameOrIDs) EIPResolver() ([]string, error) {

	sess, _ := session.NewSession()
	ec2svc := ec2.New(sess)
	inputEIPSlices := e.InputEIPnameOrIDs
	var allocationIDs []string
	for _, nameOrIDs := range inputEIPSlices {
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
			//if *results.Addresses[0].AssociationId != "" {
			//	return nil, errors.Errorf("EIP by the name %s is in use already, please provide a different EIP name or allocation ID", nameOrIDs)
			//}
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
	return allocationIDs, nil
}


/*func main() {

	var e resolver
	e = EIPnameOrIDs{
		inputEIPnameOrIDs: []string{"test243"},
	}

	var returningEIPs []string

	//	EIPnameOrIDs := []string{"test1"}
	returningEIPs, err = e.eipresolver()
	fmt.Println("available EIPs:", returningEIPs)
	fmt.Println(err)
}*/
