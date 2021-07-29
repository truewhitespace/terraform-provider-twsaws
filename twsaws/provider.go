package twsaws

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/truewhitespace/key-rotation/awskeystore"
	"github.com/truewhitespace/key-rotation/rotation"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"backend": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"twsaws_rotating_keys": dataRotatingKeys(),
		},
		DataSourcesMap: map[string]*schema.Resource{
		},
		ConfigureContextFunc: configureProvider,
	}
}

type providerConfig struct {
	backend string
}

func (p *providerConfig) KeyStoreFor(username string) (rotation.KeyStore, error){
	var awsClient *iam.IAM
	var err error
	if p.backend == "default" {
		sess := session.Must(session.NewSession())
		awsClient = iam.New(sess)
	} else if p.backend == "localstack" {
		awsClient, err = awskeystore.NewLocalstackProvider()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("no such backend")
	}
	return awskeystore.NewAWSUserKeyStore(username, awsClient), nil
}

func configureProvider(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return &providerConfig{
		backend: data.Get("backend").(string),
	}, nil
}
