package yq

import (
	"bytes"
	"os"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	glog "gopkg.in/op/go-logging.v1"
)

func init() {
	glog.SetLevel(glog.WARNING, "yq-lib")
}

// Evaluate creates a yqlib Evaluator and applies it to one or more
// files with reasonable, predictable defaults for logging, decoding,
// and printing. Only YAML files are supported.  Will read from standard
// input if no arguments are passed.
func Evaluate(expr string, files ...string) error {
	if len(files) == 0 {
		files = append(files, "-")
	}
	prefs := yqlib.YamlPreferences{}
	ev := yqlib.NewAllAtOnceEvaluator()
	enc := yqlib.NewYamlEncoder(prefs)
	pw := yqlib.NewSinglePrinterWriter(os.Stdout)
	pr := yqlib.NewPrinter(enc, pw)
	dc := yqlib.NewYamlDecoder(prefs)
	return ev.EvaluateFiles(expr, files, pr, dc)
}

// EvaluateToString is the same as Evaluate but returns a string with
// the output instead.
func EvaluateToString(expr string, files ...string) (string, error) {
	if len(files) == 0 {
		files = append(files, "-")
	}
	buf := new(bytes.Buffer)
	prefs := yqlib.YamlPreferences{}
	ev := yqlib.NewAllAtOnceEvaluator()
	enc := yqlib.NewYamlEncoder(prefs)
	pw := yqlib.NewSinglePrinterWriter(buf)
	pr := yqlib.NewPrinter(enc, pw)
	dc := yqlib.NewYamlDecoder(prefs)
	if err := ev.EvaluateFiles(expr, files, pr, dc); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}
