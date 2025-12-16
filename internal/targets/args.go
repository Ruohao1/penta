package targets

import (
	"strings"

	"github.com/spf13/pflag"
)

func SplitDashDash(argv []string) (pre, post []string) {
	for i, a := range argv {
		if a == "--" {
			return argv[:i], argv[i+1:]
		}
	}
	return argv, nil
}

func ExtractTargets(pre []string, fs *pflag.FlagSet) string {
	i := 0
	if len(pre) > 0 {
		i++
	} // binary
	if len(pre) > 1 {
		i++
	} // "scan" (and maybe more subcmds if you want)

	for ; i < len(pre); i++ {
		tok := pre[i]

		// stop parsing flags after `--` (shouldn't be here, but safe)
		if tok == "--" {
			break
		}

		// non-flag => target
		if !strings.HasPrefix(tok, "-") {
			return tok
		}

		// flag: --name=value
		if strings.HasPrefix(tok, "--") && strings.Contains(tok, "=") {
			continue
		}

		// flag: --name  value  OR  -c value
		name := strings.TrimLeft(tok, "-")
		f := fs.Lookup(name)
		if f == nil && strings.HasPrefix(tok, "--") {
			// long flag might be like --min-rate; Lookup expects exact
			f = fs.Lookup(strings.TrimPrefix(tok, "--"))
		}

		// if unknown flag, treat as target-ish? (I recommend: return error instead)
		if f == nil {
			continue
		}

		// bool flag => no value to consume
		if f.Value.Type() == "bool" {
			continue
		}

		// consume the next token as value if present
		if i+1 < len(pre) && pre[i+1] != "--" {
			i++
		}
	}

	return ""
}
