package edit

import (
	"os"

	"github.com/arahmanhamdy/elvish/edit/highlight"
	"github.com/arahmanhamdy/elvish/eval"
	"github.com/arahmanhamdy/elvish/parse"
	"github.com/arahmanhamdy/elvish/util"
)

func doHighlight(n parse.Node, ed *Editor) {
	s := &highlight.Emitter{
		func(s string) bool { return goodFormHead(s, ed) },
		ed.styling.Add,
	}
	s.EmitAll(n)
}

func goodFormHead(head string, ed *Editor) bool {
	if eval.IsBuiltinSpecial[head] {
		return true
	} else if util.DontSearch(head) {
		// XXX don't stat twice
		return util.IsExecutable(head) || isDir(head)
	} else {
		ev := ed.evaler
		explode, ns, name := eval.ParseVariable(head)
		if !explode {
			switch ns {
			case "":
				if ev.Builtin[eval.FnPrefix+name] != nil || ev.Global[eval.FnPrefix+name] != nil {
					return true
				}
			case "e":
				if ed.isExternal[name] {
					return true
				}
			default:
				if ev.Modules[ns] != nil && ev.Modules[ns][eval.FnPrefix+name] != nil {
					return true
				}
			}
		}
		return ed.isExternal[head]
	}
}

func isDir(fname string) bool {
	stat, err := os.Stat(fname)
	return err == nil && stat.IsDir()
}
