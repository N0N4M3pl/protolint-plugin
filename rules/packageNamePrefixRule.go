package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type PackageNamePrefixRule struct {
	verbose  bool
	severity rule.Severity
}

func NewPackageNamePrefixRule(
	verbose bool,
	severity rule.Severity,
) PackageNamePrefixRule {
	return PackageNamePrefixRule{
		verbose:  verbose,
		severity: severity,
	}
}

func (r PackageNamePrefixRule) ID() string {
	return "PACKAGE_NAME_PREFIX"
}

func (r PackageNamePrefixRule) Purpose() string {
	return "Verifies that the package starts with prefix."
}

func (r PackageNamePrefixRule) IsOfficial() bool {
	return false
}

func (r PackageNamePrefixRule) Severity() rule.Severity {
	return r.severity
}

func (r PackageNamePrefixRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base := visitor.NewBaseAddVisitor(r.ID(), string(r.Severity()))

	v := &packageNamePrefixVisitor{
		BaseAddVisitor: base,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type packageNamePrefixVisitor struct {
	*visitor.BaseAddVisitor
}

func (v *packageNamePrefixVisitor) VisitPackage(p *parser.Package) bool {
	if !strings.HasPrefix(p.Name, "") {
		v.AddFailuref(p.Meta.Pos, "Package name %q must starts with ", p.Name)
	}

	return false
}
