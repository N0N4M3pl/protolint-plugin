package main

import (
	"flag"

	"github.com/N0N4M3pl/protolint-plugin/internal/rules"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/plugin"
)

var (
	goStyle = flag.Bool("go_style", true, "the comments should follow a golang style")
)

func main() {
	flag.Parse()

	plugin.RegisterCustomRules(
		plugin.RuleGen(func(
			verbose bool,
			fixMode bool,
		) rule.Rule {
			return rules.NewFieldTypeOutsidePackageRule(rule.SeverityError, []string{
				"google.protobuf.",
			})
		}),

		plugin.RuleGen(func(
			verbose bool,
			fixMode bool,
		) rule.Rule {
			return rules.NewFileHasPackageRule(rule.SeverityError)
		}),

		plugin.RuleGen(func(
			verbose bool,
			fixMode bool,
		) rule.Rule {
			return rules.NewImportAvoidCommonRule(rule.SeverityNote, nil)
		}),

		plugin.RuleGen(func(
			verbose bool,
			fixMode bool,
		) rule.Rule {
			return rules.NewOptionJavaPackageUselessRule(rule.SeverityError)
		}),

		plugin.RuleGen(func(
			verbose bool,
			fixMode bool,
		) rule.Rule {
			return rules.NewPackageNamePrefixRule(rule.SeverityError, "^n0n4m3pl\\.protolint_plugin\\..*")
		}),

		plugin.RuleGen(func(
			verbose bool,
			fixMode bool,
		) rule.Rule {
			return rules.NewPackageNameSuffixVersionRule(rule.SeverityError, "")
		}),

		plugin.RuleGen(func(
			verbose bool,
			fixMode bool,
		) rule.Rule {
			return rules.NewRPCReqResNameSuffixRule(rule.SeverityError, "", "")
		}),

		plugin.RuleGen(func(
			verbose bool,
			fixMode bool,
		) rule.Rule {
			return rules.NewRPCStreamBlockedInPackageRule(rule.SeverityError, "\\.public\\.", "")
		}),
	)
}
