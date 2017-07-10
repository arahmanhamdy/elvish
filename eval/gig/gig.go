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
	{"gigtest", gigTest},
}

func gigTest(ec *eval.EvalCtx, args []eval.Value, opts map[string]eval.Value) {

	print(ec, args, opts)
	//ec.ports[1].File.WriteString("\n")
}
