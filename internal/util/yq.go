package util

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	logging "gopkg.in/op/go-logging.v1"
)

// EvaluateToString evalutes a yaml query expression into the desired
// output.
func EvaluateToString(expr string, path string) (string, error) {
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`,
	)
	b1 := logging.NewLogBackend(os.Stderr, "", 0)
	b2 := logging.AddModuleLevel(logging.NewBackendFormatter(b1, format))
	b2.SetLevel(logging.ERROR, "")
	logging.SetBackend(b2)

	ev := yqlib.NewAllAtOnceEvaluator()
	buf := new(bytes.Buffer)
	pr := newPrinter(buf, yqlib.YamlFormat)
	dc := yqlib.NewYamlDecoder(yqlib.YamlPreferences{})
	err := ev.EvaluateFiles(expr, []string{path}, pr, dc)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func newPrinter(w io.Writer, f *yqlib.Format) yqlib.Printer {
	prefs := yqlib.YamlPreferences{
		Indent:                      2,
		ColorsEnabled:               false,
		LeadingContentPreProcessing: false,
		PrintDocSeparators:          true,
		UnwrapScalar:                true,
	}
	enc := yqlib.NewYamlEncoder(prefs)
	pwr := yqlib.NewSinglePrinterWriter(w)
	return yqlib.NewPrinter(enc, pwr)
}
