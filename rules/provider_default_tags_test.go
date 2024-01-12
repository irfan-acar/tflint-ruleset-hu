package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_ProviderDefaultTagsTags(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
		RaiseErr error
	}{
		{
			Name: "Provider level tags",
			Content: `variable "aws_region" {
				default = "us-east-1"
				type = string
			  }
			  
			  provider "aws" {
				region = "us-east-1"
				default_tags {
					tags = {
					  "Managed By" = "Owner"
					  "Sample" = "something"
					}
				  }
			  }
			  
			  resource "aws_s3_bucket" "a" {
				provider = aws.foo
				name = "a"
			  }
			  
			  resource "aws_s3_bucket" "b" {
				name = "b"
				tags = var.default_tags
			  }
			  
			  variable "default_tags" {
				default = {
					"Managed By" = "test"
					"Sample" = "test"
				}
				type = map(string)
			  }
			  
			  provider "aws" {
				region = "us-east-1"
				alias = "bar"
				default_tags {
				  tags = var.default_tags
				}
			  }`,
			Config: `rule "provider_default_tags" {
						enabled = true
						tags = ["Managed By", "Sample"]
					}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAProviderDefaultTagsTypeRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"module.tf": tc.Content, ".tflint.hcl": tc.Config})

			err := rule.Check(runner)

			if tc.RaiseErr == nil && err != nil {
				t.Fatalf("Unexpected error occurred in test \"%s\": %s", tc.Name, err)
			}

			assert.Equal(t, tc.RaiseErr, err)

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
