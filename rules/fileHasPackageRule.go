package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type FileHasPackageRule struct {
	verbose  bool
	severity rule.Severity
}

func NewFileHasPackageRule(
	verbose bool,
	severity rule.Severity,
) FileHasPackageRule {
	return FileHasPackageRule{
		verbose:  verbose,
		severity: severity,
	}
}

func (r FileHasPackageRule) ID() string {
	return "FILE_HAS_PACKAGE"
}

func (r FileHasPackageRule) Purpose() string {
	return "Verifies that a fiel has package."
}

func (r FileHasPackageRule) IsOfficial() bool {
	return false
}

func (r FileHasPackageRule) Severity() rule.Severity {
	return r.severity
}

func (r FileHasPackageRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base := visitor.NewBaseAddVisitor(r.ID(), string(r.Severity()))

	v := &fileHasPackageVisitor{
		BaseAddVisitor: base,
		hasPackage:     false,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type fileHasPackageVisitor struct {
	*visitor.BaseAddVisitor
	hasPackage bool
	pos        meta.Position
}

func (v *fileHasPackageVisitor) VisitSyntax(s *parser.Syntax) bool {
	v.pos = s.Meta.Pos
	return false
}

func (v *fileHasPackageVisitor) VisitPackage(p *parser.Package) bool {
	v.hasPackage = true
	return false
}

func (v *fileHasPackageVisitor) Finally() error {
	if !v.hasPackage {
		v.AddFailuref(v.pos, "File must have defined package")
	}

	return v.BaseAddVisitor.Finally()
}
