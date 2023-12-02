package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

var defaultRequestSuffix = "Req"
var defaultResponseSuffix = "Res"

type RPCReqResNameSuffixRule struct {
	severity       rule.Severity
	requestSuffix  string
	responseSuffix string
}

func NewRPCReqResNameSuffixRule(
	severity rule.Severity,
	requestSuffix string,
	responseSuffix string,
) RPCReqResNameSuffixRule {
	if len(requestSuffix) == 0 {
		requestSuffix = defaultRequestSuffix
	}
	if len(responseSuffix) == 0 {
		responseSuffix = defaultResponseSuffix
	}
	return RPCReqResNameSuffixRule{
		severity:       severity,
		requestSuffix:  requestSuffix,
		responseSuffix: responseSuffix,
	}
}

func (r RPCReqResNameSuffixRule) ID() string {
	return "RPC_REQ_RES_NAME_SUFFIX"
}

func (r RPCReqResNameSuffixRule) Purpose() string {
	return "Verifies that all rpc request message and response message ends with specified suffixs."
}

func (r RPCReqResNameSuffixRule) IsOfficial() bool {
	return false
}

func (r RPCReqResNameSuffixRule) Severity() rule.Severity {
	return r.severity
}

func (r RPCReqResNameSuffixRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base := visitor.NewBaseAddVisitor(r.ID(), string(r.Severity()))

	v := &rpcReqResNameSuffixVisitor{
		BaseAddVisitor: base,
		requestSuffix:  r.requestSuffix,
		responseSuffix: r.responseSuffix,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type rpcReqResNameSuffixVisitor struct {
	*visitor.BaseAddVisitor
	requestSuffix  string
	responseSuffix string
}

func (v *rpcReqResNameSuffixVisitor) VisitRPC(rpc *parser.RPC) bool {
	if !strings.HasSuffix(rpc.RPCRequest.MessageType, v.requestSuffix) {
		v.AddFailuref(rpc.RPCRequest.Meta.Pos, "RPC Request name %q must ends with %q", rpc.RPCRequest.MessageType, v.requestSuffix)
	}

	if !strings.HasSuffix(rpc.RPCResponse.MessageType, v.responseSuffix) {
		v.AddFailuref(rpc.RPCResponse.Meta.Pos, "RPC Response name %q must ends with %q", rpc.RPCResponse.MessageType, v.responseSuffix)
	}

	return false
}
