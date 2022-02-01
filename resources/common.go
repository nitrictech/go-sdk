package resources

import (
	nitricv1 "github.com/nitrictech/apis/go/nitric/v1"
)

func functionResourceDeclareRequest(subject *nitricv1.Resource, actions []nitricv1.Action) *nitricv1.ResourceDeclareRequest {
	return &nitricv1.ResourceDeclareRequest{
		Resource: &nitricv1.Resource{
			Type: nitricv1.ResourceType_Policy,
		},
		Config: &nitricv1.ResourceDeclareRequest_Policy{
			Policy: &nitricv1.PolicyResource{
				Principals: []*nitricv1.Resource{
					{
						Type: nitricv1.ResourceType_Function,
					},
				},
				Actions:   actions,
				Resources: []*nitricv1.Resource{subject},
			},
		},
	}
}
