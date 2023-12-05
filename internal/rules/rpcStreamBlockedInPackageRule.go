package rules

import (
	"regexp"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type RPCStreamBlockedInPackageRule struct {
	severity                              rule.Severity
	requestStreamBlockedInPackageRgxpStr  string
	requestStreamBlockedInPackageRgxp     *regexp.Regexp
	responseStreamBlockedInPackageRgxpStr string
	responseStreamBlockedInPackageRgxp    *regexp.Regexp
}

func NewRPCStreamBlockedInPackageRule(
	severity rule.Severity,
	requestStreamBlockedInPackageRgxpStr string,
	responseStreamBlockedInPackageRgxpStr string,
) RPCStreamBlockedInPackageRule {
	var defaultRequestStreamBlockedInPackageRgxpStr = "^$"
	var defaultResponseStreamBlockedInPackageRgxpStr = "^$"
	if len(requestStreamBlockedInPackageRgxpStr) == 0 {
		requestStreamBlockedInPackageRgxpStr = defaultRequestStreamBlockedInPackageRgxpStr
	}
	requestStreamBlockedInPackageRgxp := regexp.MustCompile(requestStreamBlockedInPackageRgxpStr)
	if len(responseStreamBlockedInPackageRgxpStr) == 0 {
		responseStreamBlockedInPackageRgxpStr = defaultResponseStreamBlockedInPackageRgxpStr
	}
	responseStreamBlockedInPackageRgxp := regexp.MustCompile(responseStreamBlockedInPackageRgxpStr)
	return RPCStreamBlockedInPackageRule{
		severity:                              severity,
		requestStreamBlockedInPackageRgxpStr:  requestStreamBlockedInPackageRgxpStr,
		requestStreamBlockedInPackageRgxp:     requestStreamBlockedInPackageRgxp,
		responseStreamBlockedInPackageRgxpStr: responseStreamBlockedInPackageRgxpStr,
		responseStreamBlockedInPackageRgxp:    responseStreamBlockedInPackageRgxp,
	}
}

func (r RPCStreamBlockedInPackageRule) ID() string {
	return "RPC_STREAM_BLOCKED_IN_PACKAGE"
}

func (r RPCStreamBlockedInPackageRule) Purpose() string {
	return "Verifies that request/responce stream is allowed only in specified packages."
}

func (r RPCStreamBlockedInPackageRule) IsOfficial() bool {
	return false
}

func (r RPCStreamBlockedInPackageRule) Severity() rule.Severity {
	return r.severity
}

func (r RPCStreamBlockedInPackageRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base := visitor.NewBaseAddVisitor(r.ID(), string(r.Severity()))

	v := &rpcStreamByPackageVisitor{
		BaseAddVisitor:                        base,
		requestStreamBlockedInPackageRgxpStr:  r.requestStreamBlockedInPackageRgxpStr,
		requestStreamBlockedInPackageRgxp:     r.requestStreamBlockedInPackageRgxp,
		responseStreamBlockedInPackageRgxpStr: r.responseStreamBlockedInPackageRgxpStr,
		responseStreamBlockedInPackageRgxp:    r.responseStreamBlockedInPackageRgxp,
		requestStreamIsBlocked:                false,
		responseStreamIsBlocked:               false,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type rpcStreamByPackageVisitor struct {
	*visitor.BaseAddVisitor
	requestStreamBlockedInPackageRgxpStr  string
	requestStreamBlockedInPackageRgxp     *regexp.Regexp
	responseStreamBlockedInPackageRgxpStr string
	responseStreamBlockedInPackageRgxp    *regexp.Regexp
	packageName                           string
	requestStreamIsBlocked                bool
	responseStreamIsBlocked               bool
}

func (v *rpcStreamByPackageVisitor) VisitPackage(p *parser.Package) bool {
	v.packageName = p.Name
	v.requestStreamIsBlocked = v.requestStreamBlockedInPackageRgxp.MatchString(p.Name)
	v.responseStreamIsBlocked = v.responseStreamBlockedInPackageRgxp.MatchString(p.Name)

	return false
}

func (v *rpcStreamByPackageVisitor) VisitRPC(rpc *parser.RPC) bool {
	if rpc.RPCRequest.IsStream && v.requestStreamIsBlocked {
		v.AddFailuref(rpc.RPCRequest.Meta.Pos, "RPC Request for %q cannot be stream - not allowed in this package %q (match with %q)", rpc.RPCName, v.packageName, v.requestStreamBlockedInPackageRgxpStr)
	}

	if rpc.RPCResponse.IsStream && v.responseStreamIsBlocked {
		v.AddFailuref(rpc.RPCResponse.Meta.Pos, "RPC Response for %q cannot be stream - not allowed in this package %q (match with %q)", rpc.RPCName, v.packageName, v.responseStreamBlockedInPackageRgxpStr)
	}

	return false
}
