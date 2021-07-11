package reporter

import "github.com/jimschubert/docked"

// Reporter defines necessary members for all reporter implementations
type Reporter interface {
	// Write handler for analysis results
	Write(result docked.AnalysisResult) error
}
