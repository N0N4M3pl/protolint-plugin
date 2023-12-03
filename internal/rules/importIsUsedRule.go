package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type ImportIsUsedRule struct {
	severity rule.Severity
}

func NewImportIsUsedRule(
	severity rule.Severity,
) ImportIsUsedRule {
	return ImportIsUsedRule{
		severity: severity,
	}
}

func (r ImportIsUsedRule) ID() string {
	return "IMPORT_IS_USED"
}

func (r ImportIsUsedRule) Purpose() string {
	return "Verifies that import is used."
}

func (r ImportIsUsedRule) IsOfficial() bool {
	return false
}

func (r ImportIsUsedRule) Severity() rule.Severity {
	return r.severity
}

func (r ImportIsUsedRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base := visitor.NewBaseAddVisitor(r.ID(), string(r.Severity()))

	v := &importIsUsedVisitor{
		BaseAddVisitor: base,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type importIsUsedVisitor struct {
	*visitor.BaseAddVisitor
}

func (v *importIsUsedVisitor) VisitImport(i *parser.Import) (next bool) {
	// parts := strings.Split(i.Location, "/")
	// filename := parts[len(parts)-1]
	// parts = strings.Split(filename, ".")
	// parts = strs.SplitSnakeCaseWord(parts[0])
	// for _, p := range parts {
	// 	for _, commonName := range v.commonNames {
	// 		if strings.Contains(p, commonName) {
	// 			v.AddFailuref(i.Meta.Pos, "Import %s should not referring to type 'common' (contains %q in filename %q)", i.Location, commonName, filename)
	// 		}
	// 	}
	// }

	return false
}
