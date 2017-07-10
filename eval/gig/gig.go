package gig

import (
	"github.com/arahmanhamdy/elvish/eval"
)

func Namespace() eval.Namespace {
	ns := eval.Namespace{}
	eval.AddBuiltinFns(ns, fns...)
	return ns
}

var fns = []*eval.BuiltinFn{
	{"print", gigPrint},
}

func gigPrint(ec *eval.EvalCtx, args []eval.Value, opts map[string]eval.Value) {

	out := ec.OutputChan()
	for _, arg := range args {
		out <- eval.String(arg.Repr(0))
	}
}
