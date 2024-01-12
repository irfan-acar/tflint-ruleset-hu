package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsInstanceExampleTypeRule checks whether ...
type ProviderDeaultTagsRule struct {
	tflint.DefaultRule
}

type providerTagsRuleConfig struct {
	Tags    []string `hclext:"tags"`
	Exclude []string `hclext:"exclude,optional"`
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
	return false
}

// Severity returns the rule severity
func (r *ProviderDeaultTagsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *ProviderDeaultTagsRule) Link() string {
	return ""
}

const (
	defaultTagsBlockName  = "default_tags"
	tagsAttributeName     = "tags"
	tagBlockName          = "tag"
	providerAttributeName = "provider"
)

func (r *ProviderDeaultTagsRule) getProviderLevelTags(runner tflint.Runner) (map[string][]string, error, map[string]*hclext.Attribute) {
	providerSchema := &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{
				Name:     "alias",
				Required: false,
			},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: defaultTagsBlockName,
				Body: &hclext.BodySchema{Attributes: []hclext.AttributeSchema{{Name: tagsAttributeName}}},
			},
		},
	}

	providerBody, err := runner.GetProviderContent("aws", providerSchema, nil)
	if err != nil {
		return nil, err, nil
	}

	// Get provider default tags
	allProviderTags := make(map[string][]string)
	allProviderAttr := make(map[string]*hclext.Attribute)
	var providerAlias string
	for _, provider := range providerBody.Blocks.OfType(providerAttributeName) {
		// Get the alias attribute, in terraform when there is a single aws provider its called "default"
		providerAttr, ok := provider.Body.Attributes["alias"]
		if !ok {
			providerAlias = "default"
		} else {
			err := runner.EvaluateExpr(providerAttr.Expr, func(alias string) error {
				logger.Debug("Walk `%s` provider", providerAlias)
				providerAlias = alias
				// Init the provider reference even if it doesn't have tags
				allProviderTags[alias] = nil
				return nil
			}, nil)
			if err != nil {
				return nil, err, nil
			}
		}

		for _, block := range provider.Body.Blocks {
			var providerTags []string
			attr, ok := block.Body.Attributes[tagsAttributeName]
			if !ok {
				continue
			}

			err := runner.EvaluateExpr(attr.Expr, func(val cty.Value) error {
				keys, _ := getKeysForValue(val)

				logger.Debug("Walk `%s` provider with tags `%v`", providerAlias, keys)
				providerTags = keys
				return nil
			}, nil)

			if err != nil {
				return nil, err, nil
			}

			allProviderTags[providerAlias] = providerTags
			allProviderAttr[providerAlias] = attr
		}
	}
	return allProviderTags, nil, allProviderAttr
}

// Check checks whether ...
func (r *ProviderDeaultTagsRule) Check(runner tflint.Runner) error {
	// This rule is an example to get a top-level resource attribute.
	config := providerTagsRuleConfig{}

	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return nil
	}
	providerTagsMap, err, providerAttrMap := r.getProviderLevelTags(runner)

	if err != nil {
		return nil
	}

	for _, tag := range config.Tags {
		logger.Warn(fmt.Sprintf("%v", providerTagsMap))
		for provider, v := range providerTagsMap {

			if !stringInSlice(tag, v) {
				return runner.EmitIssue(
					r,
					"Provider "+provider+" does not have "+tag+" in default tags!",
					providerAttrMap[provider].Expr.Range(),
				)
			}
		}

	}
	return nil
}

func getKeysForValue(value cty.Value) (keys []string, known bool) {
	if !value.CanIterateElements() || !value.IsKnown() {
		return nil, false
	}
	if value.IsNull() {
		return keys, true
	}

	return keys, !value.ForEachElement(func(key, _ cty.Value) bool {
		// If any key is unknown or sensitive, return early as any missing tag could be this unknown key.
		if !key.IsKnown() || key.IsNull() || key.IsMarked() {
			return true
		}

		keys = append(keys, key.AsString())

		return false
	})
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
