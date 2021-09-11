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
	"golang.org/x/term"
)

const (
	validationLine5ColumnFormat = "%s\t%s\t%s\t%s\t%s\t\n"
)

// TextReporter writes formatted output in textual column format to Out.
// Optionally, control whether colors are output in supported terminals with DisableColors
type TextReporter struct {
	DisableColors bool      // Disable colors in supported terminals
	Out           io.Writer // The output stream
	_isTTY        *bool
}

// isTerminal returns true if unix-style terminal is supported, which is used as an indicator for color support.
// Note that this function caches the value, and ensures colors are only supported when output is a supporting file descriptor.
func (t *TextReporter) isTerminal() bool {
	if t._isTTY != nil {
		return *t._isTTY
	}
	var isTTY bool
	switch w := t.Out.(type) {
	case *os.File:
		// taken from golang.org/x/term
		isTTY := term.IsTerminal(int(w.Fd()))
		t._isTTY = &isTTY
	default:
	}
	t._isTTY = &isTTY
	return isTTY
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
	if v.ValidationResult.Result == model.Recommendation {
		indicator = t.formatted("R", BrightGreenText)
	}
	r := *v.Rule
	priority := strings.TrimSuffix(r.GetPriority().String(), "Priority")
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

func (t TextReporter) pluralIf(word string, conditional int) string {
	if conditional != 1 {
		return fmt.Sprintf("%ss", word)
	}
	return word
}

//goland:noinspection GoUnhandledErrorResult
func (t *TextReporter) Write(result docked.AnalysisResult) error {
	evalCount := len(result.Evaluated)
	notEvaluated := len(result.NotEvaluated)
	total := evalCount + notEvaluated
	spacer := strings.Repeat("-", 28)
	errorCount, recommendations, evalMap := t.prepareLookups(result)

	// all colors, even empty header, have to have equal-with colors. see https://stackoverflow.com/a/46208644/151445
	emptyColor := t.formatted(" ", Reset)

	w := tabwriter.NewWriter(t.Out, 0, 0, 3, ' ', tabwriter.AlignRight)

	fmt.Fprintf(w, validationLine5ColumnFormat, emptyColor, "Priority", "Rule", "Details", "Line(s)")
	fmt.Fprintf(w, validationLine5ColumnFormat, emptyColor, "--------", "----", "-------", "-------")
	for i := 3; i >= 0; i-- {
		if vs, ok := evalMap[model.Priority(i)]; ok {
			for _, validation := range *vs {
				if err := t.writeValidationLine(w, validation); err != nil {
					return err
				}

			}
		}
	}
	fmt.Fprintf(w, "%s\n", spacer)

	if errorCount > 0 {
		status := t.formatted("Failure", BrightRedText)
		attention := t.formatted("%d %s", BrightRedText, errorCount, t.pluralIf("error", errorCount))
		fmt.Fprintf(w, "%s - %s\n", status, attention)
	} else {
		fmt.Fprintf(w, "%s\n", t.formatted("Success", BrightGreenText))
	}

	if recommendations > 0 {
		fmt.Fprintf(w, "* %d %s\n", recommendations, t.pluralIf("recommendation", recommendations))
	}

	fmt.Fprintf(w, "* %d %s evaluated\n", evalCount, t.pluralIf("rule", evalCount))

	if total > evalCount {
		fmt.Fprintf(w, "* %d %s not evaluated\n", notEvaluated, t.pluralIf("rule", evalCount))
	} else {
		fmt.Fprintf(w, "* All rules were evaluated\n")
	}

	return w.Flush()
}

// prepareLookups creates a loop of validations.Validation by priority, returning total error count to avoid iterating the validations elsewhere
func (t *TextReporter) prepareLookups(result docked.AnalysisResult) (errorCount int, recommendations int, errorMap map[model.Priority]*[]validations.Validation) {
	errorCount = 0
	recommendations = 0
	evalMap := make(map[model.Priority]*[]validations.Validation)

	for _, validation := range result.Evaluated {
		if validation.Result == model.Failure {
			errorCount++
		}
		if validation.Result == model.Recommendation {
			recommendations++
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
	return errorCount, recommendations, evalMap
}
