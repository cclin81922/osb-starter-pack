package broker

import (
	"net/http"
	"sync"

	"github.com/golang/glog"
	"github.com/pmorie/osb-broker-lib/pkg/broker"

	"reflect"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

// NewBusinessLogic is a hook that is called with the Options the program is run
// with. NewBusinessLogic is the place where you will initialize your
// BusinessLogic the parameters passed in.
func NewBusinessLogic(o Options) (*BusinessLogic, error) {
	// For example, if your BusinessLogic requires a parameter from the command
	// line, you would unpack it from the Options and set it on the
	// BusinessLogic here.
	return &BusinessLogic{
		async:     o.Async,
		instances: make(map[string]*exampleInstance, 10),
	}, nil
}

// BusinessLogic provides an implementation of the broker.BusinessLogic
// interface.
type BusinessLogic struct {
	// Indicates if the broker should handle the requests asynchronously.
	async bool
	// Synchronize go routines.
	sync.RWMutex
	// Add fields here! These fields are provided purely as an example
	instances map[string]*exampleInstance
}

var _ broker.Interface = &BusinessLogic{}

func truePtr() *bool {
	b := true
	return &b
}

func (b *BusinessLogic) GetCatalog(c *broker.RequestContext) (*broker.CatalogResponse, error) {
	// Your catalog business logic goes here
	response := &broker.CatalogResponse{}
	osbResponse := &osb.CatalogResponse{
		Services: []osb.Service{
			{
				Name:          "example-starter-pack-service",
				ID:            "4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a",
				Description:   "The example service from the osb starter pack!",
				Bindable:      true,
				PlanUpdatable: truePtr(),
				Metadata: map[string]interface{}{
					"displayName": "Example starter pack service",
					"imageUrl":    "https://avatars2.githubusercontent.com/u/19862012?s=200&v=4",
				},
				Plans: []osb.Plan{
					{
						Name:        "default",
						ID:          "86064792-7ea2-467b-af93-ac9694d96d5b",
						Description: "The default plan for the starter pack example service",
						Free:        truePtr(),
						Schemas: &osb.Schemas{
							ServiceInstance: &osb.ServiceInstanceSchema{
								Create: &osb.InputParametersSchema{
									Parameters: map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"color": map[string]interface{}{
												"type":    "string",
												"default": "Clear",
												"enum": []string{
													"Clear",
													"Beige",
													"Grey",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	glog.Infof("catalog response: %#+v", osbResponse)

	response.CatalogResponse = *osbResponse

	return response, nil
}

func (b *BusinessLogic) Provision(request *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
	// Your provision business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	response := broker.ProvisionResponse{}

	exampleInstance := &exampleInstance{
		ID:        request.InstanceID,
		ServiceID: request.ServiceID,
		PlanID:    request.PlanID,
		Params:    request.Parameters,
	}

	// ...
	const (
		baseurl = `https://baas-skeleton-osbapibaas.baas-skeleton.svc.cluster.local`
		ca      = `-----BEGIN CERTIFICATE-----
MIID4zCCAsugAwIBAgIJAJNowMdJfs8UMA0GCSqGSIb3DQEBCwUAMIGHMQswCQYD
VQQGEwJUVzEPMA0GA1UECAwGVGFpd2FuMQ8wDQYDVQQHDAZUYWlwZWkxDjAMBgNV
BAoMBWNjbGluMQ4wDAYDVQQLDAVjY2xpbjERMA8GA1UEAwwIY2EuY2NsaW4xIzAh
BgkqhkiG9w0BCQEWFGNjbGluODE5MjJAZ21haWwuY29tMB4XDTE4MDIyNzA3NTI0
MloXDTI4MDEwNjA3NTI0MlowgYcxCzAJBgNVBAYTAlRXMQ8wDQYDVQQIDAZUYWl3
YW4xDzANBgNVBAcMBlRhaXBlaTEOMAwGA1UECgwFY2NsaW4xDjAMBgNVBAsMBWNj
bGluMREwDwYDVQQDDAhjYS5jY2xpbjEjMCEGCSqGSIb3DQEJARYUY2NsaW44MTky
MkBnbWFpbC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDKveco
eF+ZlJDgiX/Vtlg492DC/wrF4ZFzdZVpIxuD0kJxS9w8i2Xhf5tD72Uu8KAHR28W
u3NjjDtzH/KQWJd9XNeN6lkzTTTQUzFUfP3uYEELnmaTRLzdO5UDJzC/n98YOJwE
HK9K9M7C++EEjbzNYEBdhIaLlbU/p4mZSOqIUhOLjg5EzlaCgvgTuvS9YEh5hXix
5by3j9GCRIj39E8R+gdXkr4XqTVIhg4xUF81iJk1yFYEgoO2gJd0KEXCPW8NqKjZ
lUcYRWw3LvKcuhqjWfOxQVIMZeRw6nM2syJfU6umO0mFVwl19ajcq6Ic2QxZTBEf
3osMjh9966/xefkBAgMBAAGjUDBOMB0GA1UdDgQWBBRiUWcH93pfph/aBKJcsvgq
ObVJpzAfBgNVHSMEGDAWgBRiUWcH93pfph/aBKJcsvgqObVJpzAMBgNVHRMEBTAD
AQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAzMbnRbXLYOFawr1N3i2NOBwqKjyare+TD
JSYnDvIQjnPm1TVLhCRN10XH98nqbT8tx1pWVmmH0XtINVYq9KaWm4089oiopYc5
MG4Ru3a3vdNADJBvh+EUtO3pAYbHExIfBCP0Vo/gp3n1LLcUn49sIkHXKmbbzPcP
scWwlL72mtmtcrbL4HGjX632xpvuyc1ZzIiHcdwKLnxrUtZcl6oQCMlGVGJHLJyH
RWQ8gXZ71sbYOHghHrIcK4XG+ChaZJxYskd3112RIeC/5/QK8U8FKNo2KhcOOv1n
cUupgdgewtbVnl1p09PRnnBkEHhUZeXItoCUAeAC36Lhb981v5h1
-----END CERTIFICATE-----
`
		key = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCvag9C1PCwP/hM
MbGWYzfUwLkFUiQIkshZLRSSIHWYiQuyUxcBbCy8lSY9gS8IgbZoNRPAS8B8YMtv
DHCCLobO8q3ZBr7dpFxbaR3/w2H9TAr4Y0DxTwdLU4xYjt1lu/loQkAeH1AKvC5m
mUX/Y3zdYU1s3kRsjhIbKRqnh07g1sFsRWPYmPQCRxBRYMILuMh50AA6JkOP2v+r
jQzENNu2s6Qt+6KZEv1Q0hbGpHtv1o4J6N5qtOif8jEFNJtsCmeudqcLYq5DmdQp
/OIfF2TGH81q764/3JbzDnBqgWkIC21zOh7dhKRKO68Ni87wmg/2HWC1WtAaQTRx
zCOgfi1ZAgMBAAECggEAUIAf2M/YVUpGLNFxak7GRIDdaC+2EakrAKHLmvQCg6oB
EClJmYGHVlQsZHVwnDrK9y/EjK82+t2A/sl6qIOpojeEyOBrn1PafqjS95k20wOe
1TbXiuZ1tn/1HH8T46hMYShmPGyqUwLhWHxmvzltCDurSJcIV7krXgOTE+bosA2b
fCDDTu0ePwDvz6v9u+5RVIjeyZtEl+anAi3Y+WwLAXfjyO3cEvf8ta9676Qgv24c
9opDPdc9RrN3Yf6Fd6uea1eVB2PjnQKha+P1lCVMvh//CfXvcWcrElnn7mSucuR2
mdJnOeki1s6Igf9Y/eCWkNlxJQaIEoHqL/matL6fcQKBgQDhotcO8mMBS3X4LCwA
xZ6uJnhv2H+CBLga7DkpjqhmbEdZJ0P6CDFnBx50Guh0RO9C6h55s5He3h9MVfN8
0IHg27trxKGEaFajwzYecFbUUrB51ZfY3+AuMeSSd04UFMWWrPXL0Pp1AhnNGOjB
1Ym0kSGl1vjLoq8mhtZ/kvK0ewKBgQDHBRRQ4mTSKgFH8TfF8F74Ju0JVGO7juwE
2MDp2sU+vANzkcBvHxhn1Vx4OLs92djaqT7JBGcSWsHIr3n7lqfDmyZFvNJvEQmp
sJQobVLBPsqy7r4Qn1wCZWX5ho0SoVSruDJoHxek2MlmbaQ3p3Xqafz5BISzu4uu
gHZILIQvOwKBgQDfyKixK1dkRlpnXC/8SAPcR0116Gx2IIYUNatv+wwsIUIWOyph
RlTxEQ90KefYwTHn1NmK7L1FJFo4VJrcdNQLlwLonKlw8CbV3tvDDrofdS+QdnZW
45utVVCUr30hz4Q0r7BMiCSPfhjm4Migzk/4ZWTQ3Uf+d4htlpgRCUZsFwKBgQC2
6bMvV7PD+LkursNM18vhFJ2cioQTGJtRJQnApMHOE6y0ZgvP1Wtv2wfesn1crkCB
TzWWOMamduVNlgFtupw7yfeV9qINVEJmRBUXRsrdMuHHLGdhDaXZyem8OO6lZcNV
A7jIO3NWnawUyMY6JF3acUkAcSeprMAHRKfxU4C1iwKBgQDBkhHb19VsXRhtklR/
bAQiW+AOsPgMM1yWunGoWxUVt5DtVTo+yFbpDz8syZ4PfP5dSVOzyc5h5gXwVqmj
eTSi2eE4mI6aHdG6k63qYAEyTMcRbjaIDI23I1poL1chVjuxVr/GeLRfwcqI81p8
YNGcnPwlp986qk8PCLL3bQVjnQ==
-----END PRIVATE KEY-----
`
		cert = `-----BEGIN CERTIFICATE-----
MIIDijCCAnICCQCm5zS6njopmjANBgkqhkiG9w0BAQUFADCBhzELMAkGA1UEBhMC
VFcxDzANBgNVBAgMBlRhaXdhbjEPMA0GA1UEBwwGVGFpcGVpMQ4wDAYDVQQKDAVj
Y2xpbjEOMAwGA1UECwwFY2NsaW4xETAPBgNVBAMMCGNhLmNjbGluMSMwIQYJKoZI
hvcNAQkBFhRjY2xpbjgxOTIyQGdtYWlsLmNvbTAeFw0xODAyMjcwODAxMDFaFw0x
OTAyMjcwODAxMDFaMIGFMQswCQYDVQQGEwJUVzEPMA0GA1UECAwGVGFpd2FuMQ8w
DQYDVQQHDAZUYWlwZWkxDjAMBgNVBAoMBWNjbGluMQ4wDAYDVQQLDAVjY2xpbjEP
MA0GA1UEAwwGY2xpZW50MSMwIQYJKoZIhvcNAQkBFhRjY2xpbjgxOTIyQGdtYWls
LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAK9qD0LU8LA/+Ewx
sZZjN9TAuQVSJAiSyFktFJIgdZiJC7JTFwFsLLyVJj2BLwiBtmg1E8BLwHxgy28M
cIIuhs7yrdkGvt2kXFtpHf/DYf1MCvhjQPFPB0tTjFiO3WW7+WhCQB4fUAq8LmaZ
Rf9jfN1hTWzeRGyOEhspGqeHTuDWwWxFY9iY9AJHEFFgwgu4yHnQADomQ4/a/6uN
DMQ027azpC37opkS/VDSFsake2/Wjgno3mq06J/yMQU0m2wKZ652pwtirkOZ1Cn8
4h8XZMYfzWrvrj/clvMOcGqBaQgLbXM6Ht2EpEo7rw2LzvCaD/YdYLVa0BpBNHHM
I6B+LVkCAwEAATANBgkqhkiG9w0BAQUFAAOCAQEASJwEJEu6HDJu4+2Ac01AsRuM
praTinX8InVAot+DJLMIqtNm1XKyozNJ/d0zpl43CIpswij8wfCd+3yStPeEQjjX
Hh9bW2yoeCyNDrtXQChDxwwF0mrqM1EqnPWZ/TNQPPZGGLe5fY0EVHXkuz58stIp
MXFj7IiX/bilR9uiSZlToAcziYpsabteOR0FcTeMHIkFOKMvLRrELjgpxHjovCJ1
TolZ/yYHqnM7S1dmgzJdlIO2HqZaVH1CxrJIzIPelBwcKijsasm9VKJHLXeIyAws
iQRlfkYLeXuYwZ5Oqc4HCH7N1076QyiK6n8BN7U1jUwWu/H0EEQvSz19oYyG7Q==
-----END CERTIFICATE-----
`
	)
	exampleInstance.Params["baseurl"] = baseurl
	exampleInstance.Params["ca"] = ca
	exampleInstance.Params["key"] = key
	exampleInstance.Params["cert"] = cert

	// Check to see if this is the same instance
	if i := b.instances[request.InstanceID]; i != nil {
		if i.Match(exampleInstance) {
			response.Exists = true
			return &response, nil
		} else {
			// Instance ID in use, this is a conflict.
			description := "InstanceID in use"
			return nil, osb.HTTPStatusCodeError{
				StatusCode:  http.StatusConflict,
				Description: &description,
			}
		}
	}
	b.instances[request.InstanceID] = exampleInstance

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) Deprovision(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
	// Your deprovision business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	response := broker.DeprovisionResponse{}

	delete(b.instances, request.InstanceID)

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) LastOperation(request *osb.LastOperationRequest, c *broker.RequestContext) (*broker.LastOperationResponse, error) {
	// Your last-operation business logic goes here

	return nil, nil
}

func (b *BusinessLogic) Bind(request *osb.BindRequest, c *broker.RequestContext) (*broker.BindResponse, error) {
	// Your bind business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	instance, ok := b.instances[request.InstanceID]
	if !ok {
		return nil, osb.HTTPStatusCodeError{
			StatusCode: http.StatusNotFound,
		}
	}

	response := broker.BindResponse{
		BindResponse: osb.BindResponse{
			Credentials: instance.Params,
		},
	}
	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) Unbind(request *osb.UnbindRequest, c *broker.RequestContext) (*broker.UnbindResponse, error) {
	// Your unbind business logic goes here
	return &broker.UnbindResponse{}, nil
}

func (b *BusinessLogic) Update(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
	// Your logic for updating a service goes here.
	response := broker.UpdateInstanceResponse{}
	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) ValidateBrokerAPIVersion(version string) error {
	return nil
}

// example types

// exampleInstance is intended as an example of a type that holds information about a service instance
type exampleInstance struct {
	ID        string
	ServiceID string
	PlanID    string
	Params    map[string]interface{}
}

func (i *exampleInstance) Match(other *exampleInstance) bool {
	return reflect.DeepEqual(i, other)
}
