package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

var defaultPrefix = "pl."

type PackageNamePrefixRule struct {
	severity rule.Severity
	prefix   string
}

func NewPackageNamePrefixRule(
	severity rule.Severity,
	prefix string,
) PackageNamePrefixRule {
	if len(prefix) == 0 {
		prefix = defaultPrefix
	}
	return PackageNamePrefixRule{
		severity: severity,
		prefix:   prefix,
	}
}

func (r PackageNamePrefixRule) ID() string {
	return "PACKAGE_NAME_PREFIX"
}

func (r PackageNamePrefixRule) Purpose() string {
	return "Verifies that the package starts with specified prefix."
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
		prefix:         r.prefix,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type packageNamePrefixVisitor struct {
	*visitor.BaseAddVisitor
	prefix string
}

func (v *packageNamePrefixVisitor) VisitPackage(p *parser.Package) bool {
	if !strings.HasPrefix(p.Name, v.prefix) {
		v.AddFailuref(p.Meta.Pos, "Package name %q must starts with %q", p.Name, v.prefix)
	}

	return false
}
