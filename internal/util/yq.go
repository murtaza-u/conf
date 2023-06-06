package util

import (
	"bytes"
	"io"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

// EvaluateToString evalutes a yaml query expression into the desired
// output.
func EvaluateToString(expr string, path string) (string, error) {
	ev := yqlib.NewAllAtOnceEvaluator()
	buf := new(bytes.Buffer)
	pr := newPrinter(buf, yqlib.YamlOutputFormat)
	dc := yqlib.NewYamlDecoder(yqlib.YamlPreferences{})
	err := ev.EvaluateFiles(expr, []string{path}, pr, dc)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func newPrinter(w io.Writer, f yqlib.PrinterOutputFormat) yqlib.Printer {
	prefs := yqlib.YamlPreferences{
		PrintDocSeparators: true,
		UnwrapScalar:       true,
	}
	enc := yqlib.NewYamlEncoder(2, false, prefs)
	pwr := yqlib.NewSinglePrinterWriter(w)
	return yqlib.NewPrinter(enc, pwr)
}
