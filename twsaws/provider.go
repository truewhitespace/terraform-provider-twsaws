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
	"time"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"backend": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
				Description: "By default, operate against AWS.  Set to 'localstack' to use against a local localstack instance.  This is primarily intended for debugging.",
			},
			"default_key_expiry": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Default: "30days",
				Description: "Default maximum age of a key before the key is scheduled for deletion",
			},
			"default_key_grace": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Default: "20days",
				Description: "Default age at which we gracefully attempt to rotate out the key",
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
	defaultKeyExpiry time.Duration
	defaultKeyGrace time.Duration
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
	validator := &validationPhase{
		data:     data,
	}

	provider :=  &providerConfig{}
	provider.backend = validator.validateStringOneOf("backend", []string{"default","localstack"})
	provider.defaultKeyGrace = validator.validateExtendedDuration("default_key_grace")
	provider.defaultKeyExpiry = validator.validateExtendedDuration("default_key_expiry")

	if validator.hasProblems() {
		return nil, validator.problems
	}
	return provider, nil
}
