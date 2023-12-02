package rules

import (
	"regexp"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

var defaultRgxpStr = "^pl\\..*"

type PackageNamePrefixRule struct {
	severity rule.Severity
	rgxpStr  string
	rgxp     *regexp.Regexp
}

func NewPackageNamePrefixRule(
	severity rule.Severity,
	rgxpStr string,
) PackageNamePrefixRule {
	if len(rgxpStr) == 0 {
		rgxpStr = defaultRgxpStr
	}
	rgxp := regexp.MustCompile(rgxpStr)
	return PackageNamePrefixRule{
		severity: severity,
		rgxpStr:  rgxpStr,
		rgxp:     rgxp,
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
		rgxp:           r.rgxp,
		rgxpStr:        r.rgxpStr,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type packageNamePrefixVisitor struct {
	*visitor.BaseAddVisitor
	rgxp    *regexp.Regexp
	rgxpStr string
}

func (v *packageNamePrefixVisitor) VisitPackage(p *parser.Package) bool {
	if !(v.rgxp.MatchString(p.Name)) {
		v.AddFailuref(p.Meta.Pos, "Package name %q must match with %q", p.Name, v.rgxpStr)
	}

	return false
}
