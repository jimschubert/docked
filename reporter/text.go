package reporter

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/jimschubert/docked"
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/validations"
	"golang.org/x/sys/unix"
)

const (
	// TIOCGETA will only succeed on a tty
	ioctlReadTermios = unix.TIOCGETA
	validationLine5ColumnFormat = "%s\t%s\t%s\t%s\t%s\t\n"
)

// TextReporter writes formatted output in textual column format to Out.
// Optionally, control whether or not colors are output in supported terminals with DisableColors
type TextReporter struct {
	// Disable colors in supported terminals
	DisableColors bool
	// The output stream
	Out io.Writer
	_isTTY bool
}

// isTerminal returns true if unix-style terminal is supported, which is used as an indicator for color support
func (t *TextReporter) isTerminal() bool {
	if t._isTTY {
		return t._isTTY
	}
	switch w := t.Out.(type) {
	case *os.File:
		// taken from golang.org/x/term
		if _, err := unix.IoctlGetTermios(int(w.Fd()), ioctlReadTermios); err == nil {
			t._isTTY = true
		}
	default:
		t._isTTY = false
	}
	return t._isTTY
}

func (t *TextReporter) formatted(format string, c Color, a ...interface{}) string {
	if !t.DisableColors && t.isTerminal() {
		wrapped := fmt.Sprintf("%s%s%s", c, format, Reset)
		return fmt.Sprintf(wrapped, a...)
	}
	return fmt.Sprintf(format, a...)
}

// writeValidationLine will write the validation in a nice tabular format to the writer.
func (t *TextReporter) writeValidationLine(w io.Writer, v validations.Validation) error {
	indicator := t.formatted("✔", BrightGreenText)
	if v.ValidationResult.Result == model.Failure {
		indicator = t.formatted("⨯", BrightRedText)
	}
	priority := strings.TrimRight((*v.Rule).GetPriority().String(), "Priority")
	lines := make([]string, 0)

	if len(v.Contexts) > 0 {
		for _, context := range v.Contexts {
			for _, location := range context.Locations {
				lines = append(lines, strconv.Itoa(location.Start.Line))
			}
		}
	}
	_, e := fmt.Fprintf(w, validationLine5ColumnFormat, indicator, priority, v.ID, v.Details, strings.Join(lines, ","))
	return e
}

//goland:noinspection GoUnhandledErrorResult
func (t *TextReporter) Write(result docked.AnalysisResult) error {
	evalCount := len(result.Evaluated)
	notEvaluated := len(result.NotEvaluated)
	total := evalCount + notEvaluated
	spacer := strings.Repeat("-", 28)
	errorCount, evalMap := t.prepareLookups(result)

	// all colors, even empty header, have to have equal-with colors. see https://stackoverflow.com/a/46208644/151445
	var status, attention, extra, emptyColor string
	emptyColor = t.formatted(" ", Reset)
	if errorCount > 0 {
		status = t.formatted("Failure", BrightRedText)
		attention = t.formatted("%d errors", BrightRedText, errorCount)
	} else {
		status = t.formatted("Success", BrightGreenText)
		attention = t.formatted("%d errors", BrightGreenText, errorCount)
	}

	if total > evalCount {
		extra = fmt.Sprintf("%d rules were not evaluated", notEvaluated)
	} else {
		extra = "All rules were evaluated"
	}

	w := tabwriter.NewWriter(t.Out, 0, 0, 3, ' ', tabwriter.AlignRight)

	fmt.Fprintf(w, validationLine5ColumnFormat, emptyColor, "Priority", "Rule", "Details", "Line(s)")
	fmt.Fprintf(w, validationLine5ColumnFormat, emptyColor, "--------", "----", "-------", "-------")
	for i := 3; i >= 0; i-- {
		if vs, ok := evalMap[model.Priority(i)]; ok {
			for _, validation := range *vs {
				t.writeValidationLine(w, validation)
			}
		}
	}
	fmt.Fprintf(w, "%s\n", spacer)
	fmt.Fprintf(w, "%s - %s/%d rules\n", status, attention, evalCount)
	fmt.Fprintf(w, "* %s\n", extra)

	return w.Flush()
}

// prepareLookups creates a loop of validations.Validation by priority, returning total error count to avoid iterating the validations elsewhere
func (t *TextReporter) prepareLookups(result docked.AnalysisResult) (errorCount int, errorMap map[model.Priority]*[]validations.Validation) {
	errorCount = 0
	evalMap := make(map[model.Priority]*[]validations.Validation)

	for _, validation := range result.Evaluated {
		if validation.Result == model.Failure {
			errorCount += 1
		}
		if validation.Rule != nil {
			r := *validation.Rule
			priority := r.GetPriority()
			v, ok := evalMap[priority]
			if !ok {
				newSlice := []validations.Validation{validation}
				evalMap[priority] = &newSlice
			} else {
				*v = append(*v, validation)
			}
		}
	}
	return errorCount, evalMap
}
