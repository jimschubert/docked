package reporter

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/jimschubert/docked"
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/validations"
	"github.com/sirupsen/logrus"
)

// content holds our static web server content.
//go:embed templates/html/fonts/Roboto/Roboto-Bold.ttf
//go:embed templates/html/*
var content embed.FS

var anyCommands = "ADD ARG CMD COPY ENTRYPOINT ENV EXPOSE FROM HEALTHCHECK LABEL MAINTAINER ONBUILD RUN SHELL STOPSIGNAL USER VOLUME WORKDIR "

type htmlRow struct {
	RowNumber int
	Contents  string
	Errors    []validations.ValidationResult
	LineStart int
	LineEnd   int
}
type HtmlReporter struct {
	DockerfilePath string
	OutDirectory   string
}

func (h *HtmlReporter) extractCommand(input string) (command string, found bool) {
	buf := &bytes.Buffer{}
	for _, char := range input {
		if unicode.IsSpace(char) || !unicode.IsLetter(char) {
			break
		} else {
			buf.WriteRune(char)
		}
	}
	if len(buf.Bytes()) < 3 {
		return "", false
	}
	inspect := buf.String()
	if strings.Contains(anyCommands, fmt.Sprintf("%s ", inspect)) {
		return inspect, true
	}

	return "", false
}

func (h *HtmlReporter) Write(result docked.AnalysisResult) error {
	t := template.Must(template.ParseFS(content, "templates/html/index.tmpl"))

	evalCount := len(result.Evaluated)
	notEvaluated := len(result.NotEvaluated)
	total := evalCount + notEvaluated

	dockerfile := path.Join(".", "Dockerfile")
	if h.DockerfilePath != "" {
		dockerfile = h.DockerfilePath
	}

	file, err := os.Open(dockerfile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.WithError(err).Debugf("Failed closing file.")
		}
	}(file)

	rows := h.initializeRows(file)
	errorCount := h.fillErrors(result, rows)

	if h.OutDirectory == "" {
		targetDir := path.Dir(dockerfile)
		h.OutDirectory = path.Join(targetDir, "out")
	}

	err = os.MkdirAll(h.OutDirectory, 0764)
	if err != nil {
		return err
	}

	targetIndexHtml := path.Join(h.OutDirectory, "index.html")
	indexHtml, err := h.file(targetIndexHtml)
	if err != nil {
		return err
	}

	data := struct {
		Filename       string
		EvaluatedCount int
		TotalCount     int
		ErrorCount     int
		Rows           []*htmlRow
	}{
		Filename:       dockerfile,
		EvaluatedCount: evalCount,
		TotalCount:     total,
		ErrorCount:     errorCount,
		Rows:           rows,
	}

	err = t.Execute(indexHtml, data)
	if err == nil {
		return h.syncContents(h.OutDirectory)
	}
	return err
}

func (h *HtmlReporter) fillErrors(result docked.AnalysisResult, rows []*htmlRow) int {
	errorCount := 0
	for _, validation := range result.Evaluated {
		if validation.ValidationResult.Result == model.Failure {
			// This is an error, so add to the errors list for the associated "row"
			// It's important to look fully here for all errors so we report on all offending lines
			for _, ctx := range validation.ValidationResult.Contexts {
				if ctx.CausedFailure {
					errorCount += 1
					line := 1 + ctx.Locations[0].Start.Line
					for _, row := range rows {
						if row.LineStart <= line && line <= row.LineEnd {
							row.Errors = append(row.Errors, validation.ValidationResult)
							break
						}
					}
				}
			}
		}
	}
	return errorCount
}

func (h *HtmlReporter) initializeRows(file *os.File) []*htmlRow {
	rows := make([]*htmlRow, 0)
	// We can't use docker's buildkit parser here because it removes newlines/continuations within commands.
	// We need file formatting fidelity, so we need to work out some naive row parsing here.
	scanner := bufio.NewScanner(file)
	line := 0
	var row *htmlRow
	for scanner.Scan() {
		line += 1
		lineContent := scanner.Text()
		if line == 1 {
			row = &htmlRow{Contents: lineContent, RowNumber: line, LineStart: line, LineEnd: line}
			rows = append(rows, row)
		} else if row != nil {
			if _, ok := h.extractCommand(lineContent); ok {
				row = &htmlRow{Contents: lineContent, RowNumber: line, LineStart: line, LineEnd: line}
				rows = append(rows, row)
			} else {
				suffix := lineContent
				// HACK: prism.js will sometimes collapse a single empty line
				if len(lineContent) == 0 {
					suffix = "\n"
				}
				extendedContents := fmt.Sprintf("%s\n%s", row.Contents, suffix)
				(*row).Contents = extendedContents
			}
			(*row).LineEnd += 1
		}
	}
	return rows
}

func (h *HtmlReporter) ensureParentDir(filename string) error {
	if strings.Contains(filename, string(os.PathSeparator)) {
		parent := path.Dir(filename)
		if _, err := os.Stat(parent); os.IsNotExist(err) {
			err := os.MkdirAll(parent, 0764)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *HtmlReporter) file(dest string) (*os.File, error) {
	err := h.ensureParentDir(dest)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0744)
}

func (h *HtmlReporter) copyFile(src, dest string) (int64, error) {
	source, err := content.Open(src)
	if err != nil {
		return 0, err
	}
	defer func(source fs.File) {
		err := source.Close()
		if err != nil {
			logrus.WithError(err).Debugf("Failed closing embedded source file.")
		}
	}(source)

	destination, err := h.file(dest)
	if err != nil {
		return 0, err
	}
	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {
			logrus.WithError(err).Debugf("Failed closing destination file.")
		}
	}(destination)

	return io.Copy(destination, source)
}

func (h *HtmlReporter) syncContents(targetDir string) error {
	// explicit file list avoids syncing and test/bak/hidden html files or other.
	// for instance: although embed claims to ignore paths starting with '.', we could see .DS_Store
	toSync := []string{
		"templates/html/custom.css",
		"templates/html/normalize.min.css",
		"templates/html/normalize.min.css.map",
		"templates/html/prism.css",
		"templates/html/prism.js",
		"templates/html/fonts/Roboto/Roboto-Bold.ttf",
		"templates/html/fonts/Roboto/LICENSE.txt",
	}
	for _, syncFile := range toSync {
		baseName := strings.TrimPrefix(syncFile, "templates/html/")
		targetFile := path.Join(targetDir, baseName)
		_, err := h.copyFile(syncFile, targetFile)
		if err != nil {
			return err
		}
	}
	return nil
}
