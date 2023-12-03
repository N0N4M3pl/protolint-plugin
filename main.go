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
		// The purpose of this line just illustrates that you can implement the same as internal linter rules.
		// rules.NewEnumsHaveCommentRule(rule.SeverityWarning, *goStyle),

		// A common custom rule example. It's simple.
		// customrules.NewEnumNamesLowerSnakeCaseRule(),

		// Wrapping with RuleGen allows referring to command-line flags.
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
			return rules.NewRPCReqResNameSuffixRule(rule.SeverityError, "Request", "Response")
		}),

		plugin.RuleGen(func(
			verbose bool,
			fixMode bool,
		) rule.Rule {
			return rules.NewImportAvoidTypeCommonRule(rule.SeverityNote, nil)
		}),
	)
}
