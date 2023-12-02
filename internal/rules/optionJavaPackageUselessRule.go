package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type OptionJavaPackageUselessRule struct {
	severity rule.Severity
}

func NewOptionJavaPackageUselessRule(
	severity rule.Severity,
) OptionJavaPackageUselessRule {
	return OptionJavaPackageUselessRule{
		severity: severity,
	}
}

func (r OptionJavaPackageUselessRule) ID() string {
	return "OPTION_JAVA_PACKAGE_USELESS"
}

func (r OptionJavaPackageUselessRule) Purpose() string {
	return "Verifies that a option 'java_package' is not same as the package."
}

func (r OptionJavaPackageUselessRule) IsOfficial() bool {
	return false
}

func (r OptionJavaPackageUselessRule) Severity() rule.Severity {
	return r.severity
}

func (r OptionJavaPackageUselessRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base := visitor.NewBaseAddVisitor(r.ID(), string(r.Severity()))

	v := &optionJavaPackageUselessVisitor{
		BaseAddVisitor: base,
		packageName:    "",
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type optionJavaPackageUselessVisitor struct {
	*visitor.BaseAddVisitor
	packageName string
}

func (v *optionJavaPackageUselessVisitor) VisitPackage(p *parser.Package) bool {
	v.packageName = p.Name
	return false
}

func (v *optionJavaPackageUselessVisitor) VisitOption(o *parser.Option) bool {
	if o.OptionName == "java_package" {
		optionValue := strings.Trim(o.Constant, "\"")
		if optionValue == v.packageName {
			v.AddFailuref(o.Meta.Pos, "Option %s is useless because has same value as package", o.OptionName)
		}
	}

	return false
}
