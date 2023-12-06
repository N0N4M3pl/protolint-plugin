package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type FieldTypeOutsidePackageRule struct {
	severity     rule.Severity
	allowedTypes []string
}

func NewFieldTypeOutsidePackageRule(
	severity rule.Severity,
	allowedTypes []string,
) FieldTypeOutsidePackageRule {
	var defaultAllowedTypes = []string{
		// "google.protobuf.",
	}
	if len(allowedTypes) == 0 {
		allowedTypes = defaultAllowedTypes
	}
	return FieldTypeOutsidePackageRule{
		severity:     severity,
		allowedTypes: allowedTypes,
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
		BaseAddVisitor: base,
		allowedTypes:   r.allowedTypes,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type fieldTypeOutsidePackageVisitor struct {
	*visitor.BaseAddVisitor
	allowedTypes []string
	packageName  string
	fields       []*parser.Field
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
	allowedTypesWithPackage := append([]string{v.packageName}, v.allowedTypes...)
	for _, f := range v.fields {
		fieldType := f.Type
		isAllowed := false
		for _, allowedType := range allowedTypesWithPackage {
			if strings.HasPrefix(fieldType, allowedType) {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			v.AddFailuref(f.Meta.Pos, "Field type %q is outside package %q and not match with allowed types %s", fieldType, v.packageName, v.allowedTypes)
		}
	}

	return v.BaseAddVisitor.Finally()
}
