package edit

import (
	"github.com/arahmanhamdy/elvish/edit/ui"
	"github.com/arahmanhamdy/elvish/eval"
)

func styled(ec *eval.EvalCtx, args []eval.Value, opts map[string]eval.Value) {
	var textv, stylev eval.String
	eval.ScanArgs(args, &textv, &stylev)
	text, style := string(textv), string(stylev)
	eval.TakeNoOpt(opts)

	out := ec.OutputChan()
	out <- &ui.Styled{text, ui.StylesFromString(style)}
}
