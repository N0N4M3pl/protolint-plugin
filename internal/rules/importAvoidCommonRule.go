package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type ImportAvoidCommonRule struct {
	severity    rule.Severity
	commonNames []string
}

func NewImportAvoidCommonRule(
	severity rule.Severity,
	commonNames []string,
) ImportAvoidCommonRule {
	var defaultCommonNames = []string{
		"common",
		"general",
		"regular",
		"standard",
		"event",
		"type",
		"data",
		"typical",
	}
	if len(commonNames) == 0 {
		commonNames = defaultCommonNames
	}
	return ImportAvoidCommonRule{
		severity:    severity,
		commonNames: commonNames,
	}
}

func (r ImportAvoidCommonRule) ID() string {
	return "IMPORT_AVOID_COMMON"
}

func (r ImportAvoidCommonRule) Purpose() string {
	return "Verifies that import is not referring to 'common' type."
}

func (r ImportAvoidCommonRule) IsOfficial() bool {
	return false
}

func (r ImportAvoidCommonRule) Severity() rule.Severity {
	return r.severity
}

func (r ImportAvoidCommonRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base := visitor.NewBaseAddVisitor(r.ID(), string(r.Severity()))

	v := &importAvoidCommonVisitor{
		BaseAddVisitor: base,
		commonNames:    r.commonNames,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type importAvoidCommonVisitor struct {
	*visitor.BaseAddVisitor
	commonNames []string
}

func (v *importAvoidCommonVisitor) VisitImport(i *parser.Import) (next bool) {
	filename := strings.Trim(i.Location, "\"")
	parts := strings.Split(filename, "/")
	filename = parts[len(parts)-1]
	parts = strings.Split(filename, ".")
	parts = strs.SplitSnakeCaseWord(parts[0])
	for _, p := range parts {
		for _, commonName := range v.commonNames {
			if strings.Contains(p, commonName) {
				v.AddFailuref(i.Meta.Pos, "Import %s should not referring to 'common' type (contains %q in filename %q)", i.Location, commonName, filename)
			}
		}
	}

	return false
}
