package twsaws

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/truewhitespace/key-rotation/awskeystore"
	"github.com/truewhitespace/key-rotation/rotation"
)

func dataRotatingKeys()  *schema.Resource {
	return &schema.Resource {
		CreateContext: createRotatingKeys,
		ReadContext: dataRotatingKeysRead,
		DeleteContext: deleteRotatingKeys,
		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"active_key_id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"active_key_secret": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Sensitive: true,
			},
		},
	}
}

func planKeyStores(ctx context.Context, d *schema.ResourceData, m interface{}) (*rotation.KeyRotationPlan,rotation.KeyStore,diag.Diagnostics) {
	cfg := m.(*providerConfig)
	user := d.Get("user_name").(string)

	var out diag.Diagnostics
	var err error
	keystore, err := cfg.KeyStoreFor(user)
	if err != nil {
		out = append(out, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "KeyStore failed to create",
			Detail:        err.Error(),
		})
		return nil,nil, out
	}

	rotator, err := rotation.NewKeyRotation(cfg.defaultKeyExpiry, cfg.defaultKeyGrace)
	if err != nil {
		out = append(out, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Bad rotation configuration",
			Detail:        err.Error(),
		})
		return nil,nil, out
	}

	plan, err := rotator.Plan(ctx, keystore)
	if err != nil {
		out = append(out, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Failed to plan",
			Detail:        err.Error(),
		})
		return nil, nil, out
	}
	return plan, keystore, nil
}

func createRotatingKeys(ctx context.Context, d *schema.ResourceData, m interface{}) (out diag.Diagnostics){
	plan, keystore, problems := planKeyStores(ctx,d,m)
	if problems != nil {
		return problems
	}

	keys, err := plan.Apply(ctx,keystore)
	if err != nil {
		out = append(out, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Failed to apply keys",
			Detail:        err.Error(),
		})
		return out
	}
	awsKey := keys[0].(*awskeystore.AWSAccessKey)

	d.Set("active_key_id",awsKey.ID)
	if awsKey.Secret != nil {
		d.Set("active_key_secret",*awsKey.Secret)
	} else {
		d.Set("active_key_secret", d.Get("active_key_secret").(string))
	}
	d.SetId(awsKey.ID)

	return out
}

func deleteRotatingKeys(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

func dataRotatingKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	plan, _, problems := planKeyStores(ctx,d,m)
	if problems != nil {
		return problems
	}

	if plan.CreateKey || len(plan.DestroyKeys) > 0 {
		d.SetId("")
	}

	return nil
}
