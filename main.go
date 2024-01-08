package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-template/rules"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "template",
			Version: "1.1.4",
			Rules: []tflint.Rule{
				rules.NewAProviderDefaultTagsTypeRule(),
				rules.NewAwsS3BucketExampleLifecycleRule(),
				rules.NewGoogleComputeSSLPolicyRule(),
				rules.NewTerraformBackendTypeRule(),
			},
		},
	})
}
