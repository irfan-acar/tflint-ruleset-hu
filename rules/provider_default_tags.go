package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsInstanceExampleTypeRule checks whether ...
type ProviderDeaultTagsRule struct {
	tflint.DefaultRule
}

// NewAwsInstanceExampleTypeRule returns a new rule
func NewAProviderDefaultTagsTypeRule() *ProviderDeaultTagsRule {
	return &ProviderDeaultTagsRule{}
}

// Name returns the rule name
func (r *ProviderDeaultTagsRule) Name() string {
	return "provider_default_tags"
}

// Enabled returns whether the rule is enabled by default
func (r *ProviderDeaultTagsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *ProviderDeaultTagsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *ProviderDeaultTagsRule) Link() string {
	return ""
}

// Check checks whether ...
func (r *ProviderDeaultTagsRule) Check(runner tflint.Runner) error {
	// This rule is an example to get a top-level resource attribute.
	content, err := runner.GetProviderContent("aws", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "default_tags"},
		},
	}, nil)
	if err != nil {
		return err
	}

	if _, exists := content.Attributes["default_tags"]; !exists {
		return runner.EmitIssue(
			r,
			"provider does not have default tags",
			hcl.Range{},
		)
	}

	return nil
}
