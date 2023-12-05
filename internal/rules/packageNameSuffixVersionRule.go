package rules

import (
	"regexp"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type PackageNameSuffixVersionRule struct {
	severity rule.Severity
	rgxpStr  string
	rgxp     *regexp.Regexp
}

func NewPackageNameSuffixVersionRule(
	severity rule.Severity,
	rgxpStr string,
) PackageNameSuffixVersionRule {
	var defaultRgxpStr = ".*\\.v((\\d+)|(\\d+test.*)|(\\d+(alpha|beta)))$"
	if len(rgxpStr) == 0 {
		rgxpStr = defaultRgxpStr
	}
	rgxp := regexp.MustCompile(rgxpStr)
	return PackageNameSuffixVersionRule{
		severity: severity,
		rgxpStr:  rgxpStr,
		rgxp:     rgxp,
	}
}

func (r PackageNameSuffixVersionRule) ID() string {
	return "PACKAGE_NAME_SUFFIX_VERSION"
}

func (r PackageNameSuffixVersionRule) Purpose() string {
	return "Verifies that the package ends with specified suffix equal version."
}

func (r PackageNameSuffixVersionRule) IsOfficial() bool {
	return false
}

func (r PackageNameSuffixVersionRule) Severity() rule.Severity {
	return r.severity
}

func (r PackageNameSuffixVersionRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base := visitor.NewBaseAddVisitor(r.ID(), string(r.Severity()))

	v := &packageNameSuffixVersionVisitor{
		BaseAddVisitor: base,
		rgxp:           r.rgxp,
		rgxpStr:        r.rgxpStr,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type packageNameSuffixVersionVisitor struct {
	*visitor.BaseAddVisitor
	rgxp    *regexp.Regexp
	rgxpStr string
}

func (v *packageNameSuffixVersionVisitor) VisitPackage(p *parser.Package) bool {
	if !(v.rgxp.MatchString(p.Name)) {
		v.AddFailuref(p.Meta.Pos, "Package name %q must ends with specified suffix equal version (match with %q)", p.Name, v.rgxpStr)
	}

	return false
}
