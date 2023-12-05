package rules

import (
	"regexp"
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type FieldTypeOutsidePackageRule struct {
	severity            rule.Severity
	allowedTypesRgxpStr string
	allowedTypesRgxp    *regexp.Regexp
}

func NewFieldTypeOutsidePackageRule(
	severity rule.Severity,
	allowedTypesRgxpStr string,
) FieldTypeOutsidePackageRule {
	var defaultAllowedTypesRgxpStr = "^.*"
	if len(allowedTypesRgxpStr) == 0 {
		allowedTypesRgxpStr = defaultAllowedTypesRgxpStr
	}
	allowedTypesRgxp := regexp.MustCompile(allowedTypesRgxpStr)
	return FieldTypeOutsidePackageRule{
		severity:            severity,
		allowedTypesRgxpStr: allowedTypesRgxpStr,
		allowedTypesRgxp:    allowedTypesRgxp,
	}
}

func (r FieldTypeOutsidePackageRule) ID() string {
	return "FIELD_TYPE_OUTSIDE_PACKAGE"
}

func (r FieldTypeOutsidePackageRule) Purpose() string {
	return "Verifies that all used types are from package or from allowed types"
}

func (r FieldTypeOutsidePackageRule) IsOfficial() bool {
	return false
}

func (r FieldTypeOutsidePackageRule) Severity() rule.Severity {
	return r.severity
}

func (r FieldTypeOutsidePackageRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base := visitor.NewBaseAddVisitor(r.ID(), string(r.Severity()))

	v := &fieldTypeOutsidePackageVisitor{
		BaseAddVisitor:      base,
		allowedTypesRgxpStr: r.allowedTypesRgxpStr,
		allowedTypesRgxp:    r.allowedTypesRgxp,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type fieldTypeOutsidePackageVisitor struct {
	*visitor.BaseAddVisitor
	allowedTypesRgxpStr string
	allowedTypesRgxp    *regexp.Regexp
	packageName         string
	fields              []*parser.Field
}

func (v *fieldTypeOutsidePackageVisitor) VisitPackage(p *parser.Package) bool {
	v.packageName = p.Name

	return false
}

func (v *fieldTypeOutsidePackageVisitor) VisitField(f *parser.Field) bool {
	if strings.Contains(f.Type, ".") {
		v.fields = append(v.fields, f)
	}

	return false
}

func (v *fieldTypeOutsidePackageVisitor) Finally() error {
	for _, f := range v.fields {
		if len(v.packageName) > 0 {
			if !strings.HasPrefix(f.Type, v.packageName) {
				if !v.allowedTypesRgxp.MatchString(f.Type) {
					v.AddFailuref(f.Meta.Pos, "Field type %q is outside package %q and not match with allowed types %q", f.Type, v.packageName, v.allowedTypesRgxpStr)
				}
			}
		}
	}

	return v.BaseAddVisitor.Finally()
}
