package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

var defaultCommonNames = []string{
	"common",
	"general",
	"regular",
	"standard",
	"event",
	"type",
}

type ImportAvoidTypeCommonRule struct {
	severity    rule.Severity
	commonNames []string
}

func NewImportAvoidTypeCommonRule(
	severity rule.Severity,
	commonNames []string,
) ImportAvoidTypeCommonRule {
	if len(commonNames) == 0 {
		commonNames = defaultCommonNames
	}
	return ImportAvoidTypeCommonRule{
		severity:    severity,
		commonNames: commonNames,
	}
}

func (r ImportAvoidTypeCommonRule) ID() string {
	return "IMPORT_AVOID_TYPE_COMMON"
}

func (r ImportAvoidTypeCommonRule) Purpose() string {
	return "Verifies that import is not referring to type 'common'."
}

func (r ImportAvoidTypeCommonRule) IsOfficial() bool {
	return false
}

func (r ImportAvoidTypeCommonRule) Severity() rule.Severity {
	return r.severity
}

func (r ImportAvoidTypeCommonRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base := visitor.NewBaseAddVisitor(r.ID(), string(r.Severity()))

	v := &importAvoidTypeCommonVisitor{
		BaseAddVisitor: base,
		commonNames:    r.commonNames,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type importAvoidTypeCommonVisitor struct {
	*visitor.BaseAddVisitor
	commonNames []string
}

func (v *importAvoidTypeCommonVisitor) VisitImport(i *parser.Import) (next bool) {
	filename := strings.Trim(i.Location, `"`)
	parts := strings.Split(filename, "/")
	filename = parts[len(parts)-1]
	parts = strings.Split(filename, ".")
	parts = strs.SplitSnakeCaseWord(parts[0])
	for _, p := range parts {
		for _, commonName := range v.commonNames {
			if strings.Contains(p, commonName) {
				v.AddFailuref(i.Meta.Pos, "Import %s should not referring to type 'common' (contains %q in filename %q)", i.Location, commonName, filename)
			}
		}
	}

	return false
}
