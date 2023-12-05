package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type RepeatedFieldInRpcMessageAvoidRule struct {
	severity rule.Severity
}

func NewRepeatedFieldInRpcMessageAvoidRule(
	severity rule.Severity,
) RepeatedFieldInRpcMessageAvoidRule {
	return RepeatedFieldInRpcMessageAvoidRule{
		severity: severity,
	}
}

func (r RepeatedFieldInRpcMessageAvoidRule) ID() string {
	return "REPEATED_FIELD_IN_RPC_MESSAGE_AVOID"
}

func (r RepeatedFieldInRpcMessageAvoidRule) Purpose() string {
	return "Verifies that repeated field in rpc message should be avoided."
}

func (r RepeatedFieldInRpcMessageAvoidRule) IsOfficial() bool {
	return false
}

func (r RepeatedFieldInRpcMessageAvoidRule) Severity() rule.Severity {
	return r.severity
}

func (r RepeatedFieldInRpcMessageAvoidRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base := visitor.NewBaseAddVisitor(r.ID(), string(r.Severity()))

	v := &repeatedFieldInRpcMessageAvoidVisitor{
		BaseAddVisitor: base,
		messages:       make(map[string]*parser.Message),
		fields:         make(map[string]*parser.Field),
		groupFields:    make(map[string]*parser.GroupField),
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type repeatedFieldInRpcMessageAvoidVisitor struct {
	*visitor.BaseAddVisitor
	rpcMessages []string
	messages    map[string]*parser.Message
	fields      map[string]*parser.Field
	groupFields map[string]*parser.GroupField
}

func (v *repeatedFieldInRpcMessageAvoidVisitor) VisitRPC(rpc *parser.RPC) bool {
	if rpc.RPCRequest.IsStream {
		v.rpcMessages = append(v.rpcMessages, rpc.RPCResponse.MessageType)
		v.AddFailuref(rpc.Meta.Pos, "%s", rpc.RPCResponse.MessageType)
	}

	return false
}

func (v *repeatedFieldInRpcMessageAvoidVisitor) VisitMessage(m *parser.Message) bool {
	v.messages[m.MessageName] = m

	for i, item := range m.MessageBody {
		// v.AddFailuref(message.Meta.Pos, "[%v]: %s | %T", i, reflect.TypeOf(item).String(), item)
		switch itemType := item.(type) {
		case *parser.Field:
			v.AddFailuref(m.Meta.Pos, "[%v]: Field | FieldName=%s", i, itemType.FieldName)
		default:
			v.AddFailuref(m.Meta.Pos, "[%v]: item is %T", i, item)
		}
	}

	return false
}

func (v *repeatedFieldInRpcMessageAvoidVisitor) VisitField(f *parser.Field) bool {
	if f.IsRepeated {
		v.fields[f.FieldName] = f
	}
	return false
}

func (v *repeatedFieldInRpcMessageAvoidVisitor) VisitGroupField(f *parser.GroupField) bool {
	if f.IsRepeated {
		v.groupFields[f.GroupName] = f
	}
	return false
}

func (v *repeatedFieldInRpcMessageAvoidVisitor) Finally() error {
	// if !v.hasPackage {
	// 	v.AddFailuref(v.pos, "File must have defined package")
	// }

	// for i, s := range v.rpcMessages {
	// 	fmt.Println(i, s)
	// }

	return v.BaseAddVisitor.Finally()
}
