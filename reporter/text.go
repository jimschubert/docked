package reporter

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/jimschubert/docked"
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/validations"
	"github.com/jimschubert/tabitha"
)

var (
	brightRed   = color.New(color.FgRed, color.Bold)
	brightGreen = color.New(color.FgGreen, color.Bold)
	cyan        = color.New(color.FgCyan, color.Bold)
)

// TextReporter writes formatted output in textual column format to Out.
// Optionally, control whether colors are output in supported terminals with DisableColors
type TextReporter struct {
	DisableColors bool      // Disable colors in supported terminals
	Out           io.Writer // The output stream
}

//goland:noinspection GoUnhandledErrorResult
func (t *TextReporter) Write(result docked.AnalysisResult) (err error) {
	errorCount, recommendations, evalMap := t.prepareLookups(result)

	// all colors, even empty header, have to have equal-with colors. see https://stackoverflow.com/a/46208644/151445
	emptyColor := cyan.Sprint(" ")

	tt := tabitha.NewWriter()
	tt.IgnoreAnsiWidths(!color.NoColor)

	err = tt.Header(emptyColor, "Priority", "Rule", "Details", "Line(s)")
	if err != nil {
		return
	}

	err = tt.SpacerLine()
	if err != nil {
		return
	}

	for i := 3; i >= 0; i-- {
		if vs, ok := evalMap[model.Priority(i)]; ok {
			for _, validation := range *vs {
				err = t.writeValidationLine(tt, validation)
				if err != nil {
					return
				}
			}
		}
	}

	err = tt.SpacerLine()
	if err != nil {
		return
	}

	if _, err := tt.WriteTo(t.Out); err != nil {
		return err
	}

	return t.writeSummary(errorCount, recommendations, &result)
}

// writeValidationLine will write the validation in a nice tabular format to the writer.
func (t *TextReporter) writeValidationLine(tt *tabitha.Writer, v validations.Validation) error {
	indicator := brightGreen.Sprint("✔")
	if v.ValidationResult.Result == model.Failure {
		indicator = brightRed.Sprint("⨯")
	}
	if v.ValidationResult.Result == model.Recommendation {
		indicator = cyan.Sprint("R")
	}
	r := *v.Rule
	priority := strings.TrimSuffix(r.GetPriority().String(), "Priority")

	// int line number to its display text
	lineNumbers := make(map[int]string)
	if len(v.Contexts) > 0 {
		for _, context := range v.Contexts {
			for _, location := range context.Locations {
				currentLine := strconv.Itoa(location.Start.Line)
				if context.CausedFailure {
					lineNumbers[location.Start.Line] = brightRed.Sprint(currentLine)
				} else {
					lineNumbers[location.Start.Line] = currentLine
				}
			}
		}
	}

	// since display text can be wrapped in ansi colors, need to create a sorted key array
	lines := make([]int, 0, len(lineNumbers))
	for k := range lineNumbers {
		lines = append(lines, k)
	}
	sort.Ints(lines)

	// collect in-order display lines, supporting ansi colors.
	// For example, we need \x1b[32;1m14\x1b[0m to come *after* 12.
	displayLines := make([]string, 0)
	for _, key := range lines {
		display := lineNumbers[key]
		displayLines = append(displayLines, display)
	}

	return tt.AddLine(indicator, priority, v.ID, v.Details, strings.Join(displayLines, ","))
}

func (t TextReporter) pluralIf(word string, conditional int) string {
	if conditional != 1 {
		return fmt.Sprintf("%ss", word)
	}
	return word
}

func (t *TextReporter) writeSummary(errorCount int, recommendations int, result *docked.AnalysisResult) (err error) {
	evalCount := len(result.Evaluated)
	notEvaluated := len(result.NotEvaluated)
	total := evalCount + notEvaluated

	if errorCount > 0 {
		if _, err = brightRed.Fprint(t.Out, "Failure\n"); err != nil {
			return
		}
		if _, err = fmt.Fprintf(t.Out, "* %d %s\n", errorCount, t.pluralIf("error", errorCount)); err == nil {
			return
		}
	} else {
		if _, err = fmt.Fprintf(t.Out, "%s\n", brightGreen.Sprint("Success")); err != nil {
			return
		}
	}

	if recommendations > 0 {
		if _, err = fmt.Fprintf(t.Out, "* %d %s\n", recommendations, t.pluralIf("recommendation", recommendations)); err != nil {
			return
		}
	}

	if _, err = fmt.Fprintf(t.Out, "* %d %s evaluated\n", evalCount, t.pluralIf("rule", evalCount)); err != nil {
		return
	}

	if total > evalCount {
		if _, err = fmt.Fprintf(t.Out, "* %d %s not evaluated\n", notEvaluated, t.pluralIf("rule", evalCount)); err != nil {
			return
		}
	} else {
		if _, err = fmt.Fprintf(t.Out, "* All rules were evaluated\n"); err != nil {
			return
		}
	}
	return nil
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
