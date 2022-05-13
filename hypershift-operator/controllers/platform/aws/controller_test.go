package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	configv1 "github.com/openshift/api/config/v1"
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	hyperapi "github.com/openshift/hypershift/support/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type fakeEC2Client struct {
	ec2iface.EC2API
	created   *ec2.CreateVpcEndpointServiceConfigurationInput
	createOut *ec2.CreateVpcEndpointServiceConfigurationOutput
}

func (f *fakeEC2Client) CreateVpcEndpointServiceConfigurationWithContext(ctx aws.Context, in *ec2.CreateVpcEndpointServiceConfigurationInput, o ...request.Option) (*ec2.CreateVpcEndpointServiceConfigurationOutput, error) {
	if f.created != nil {
		return nil, errors.New("already created endpoint service")
	}
	f.created = in
	return f.createOut, nil
}

type fakeElbv2Client struct {
	elbv2iface.ELBV2API
	out *elbv2.DescribeLoadBalancersOutput
}

func (f *fakeElbv2Client) DescribeLoadBalancersWithContext(aws.Context, *elbv2.DescribeLoadBalancersInput, ...request.Option) (*elbv2.DescribeLoadBalancersOutput, error) {
	return f.out, nil
}

func TestReconcileAWSEndpointServiceStatus(t *testing.T) {
	elbClient := &fakeElbv2Client{out: &elbv2.DescribeLoadBalancersOutput{LoadBalancers: []*elbv2.LoadBalancer{{
		LoadBalancerArn: aws.String("lb-arn"),
		State:           &elbv2.LoadBalancerState{Code: aws.String(elbv2.LoadBalancerStateEnumActive)},
	}}}}

	infra := &configv1.Infrastructure{
		ObjectMeta: metav1.ObjectMeta{Name: "cluster"},
		Status:     configv1.InfrastructureStatus{InfrastructureName: "management-cluster-infra-id"},
	}
	client := fake.NewClientBuilder().WithScheme(hyperapi.Scheme).WithObjects(infra).Build()

	ec2Client := &fakeEC2Client{
		createOut: &ec2.CreateVpcEndpointServiceConfigurationOutput{ServiceConfiguration: &ec2.ServiceConfiguration{ServiceName: aws.String("ep-service")}},
	}

	if err := reconcileAWSEndpointServiceStatus(context.Background(), client, &hyperv1.AWSEndpointService{}, ec2Client, elbClient); err != nil {
		t.Fatalf("reconcileAWSEndpointServiceStatus failed: %v", err)
	}

	if actual, expected := *ec2Client.created.TagSpecifications[0].Tags[0].Key, "kubernetes.io/cluster/management-cluster-infra-id"; actual != expected {
		t.Errorf("expected first tag key to be %s, was %s", expected, actual)
	}

	if actual, expected := *ec2Client.created.TagSpecifications[0].Tags[0].Value, "owned"; actual != expected {
		t.Errorf("expected first tags value to be %s, was %s", expected, actual)
	}
}