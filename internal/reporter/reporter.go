// Package reporter defines the common Issue/Report types produced by every
// diagnostic module and renders them as terminal, JSON, or HTML output.
package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Severity indicates how serious a detected issue is.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityWarning  Severity = "warning"
	SeverityInfo     Severity = "info"
)

// Issue represents a single detected problem, its root cause, and the
// recommended fix.
type Issue struct {
	Category       string   `json:"category"`
	Severity       Severity `json:"severity"`
	Title          string   `json:"title"`
	Reason         string   `json:"reason"`
	Recommendation string   `json:"recommendation"`
	DocsLink       string   `json:"docsLink,omitempty"`
}

// Check represents a passed/skipped diagnostic check surfaced in the report
// summary (e.g. "API Server Healthy").
type Check struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Passed   bool   `json:"passed"`
}

// Report is the aggregate result of a diagnose run.
type Report struct {
	ClusterVersion string  `json:"clusterVersion"`
	Score          int     `json:"score"`
	Health         string  `json:"health"`
	Checks         []Check `json:"checks"`
	Issues         []Issue `json:"issues"`
}

func severityPenalty(s Severity) int {
	switch s {
	case SeverityCritical:
		return 15
	case SeverityWarning:
		return 5
	default:
		return 1
	}
}

func healthLabel(score int) string {
	switch {
	case score >= 90:
		return "Excellent"
	case score >= 75:
		return "Good"
	case score >= 50:
		return "Degraded"
	default:
		return "Critical"
	}
}

// NewReport computes a cluster score from the given checks and issues.
// The score starts at 100 and is reduced by a per-severity penalty for each
// issue found, floored at 0.
func NewReport(clusterVersion string, checks []Check, issues []Issue) Report {
	score := 100
	for _, issue := range issues {
		score -= severityPenalty(issue.Severity)
	}
	if score < 0 {
		score = 0
	}

	return Report{
		ClusterVersion: clusterVersion,
		Score:          score,
		Health:         healthLabel(score),
		Checks:         checks,
		Issues:         issues,
	}
}

func (r Report) counts() (critical, warning int) {
	for _, issue := range r.Issues {
		switch issue.Severity {
		case SeverityCritical:
			critical++
		case SeverityWarning:
			warning++
		}
	}
	return
}

// WriteTable renders a human-readable terminal report.
func (r Report) WriteTable(w io.Writer) {
	fmt.Fprintln(w, "Cluster Status")
	fmt.Fprintln(w)
	for _, check := range r.Checks {
		mark := "✓" // ✓
		if !check.Passed {
			mark = "✗" // ✗
		}
		fmt.Fprintf(w, "%s %s\n", mark, check.Title)
	}

	for _, issue := range r.Issues {
		mark := "⚠" // ⚠
		if issue.Severity == SeverityCritical {
			mark = "✗" // ✗
		}
		fmt.Fprintf(w, "%s %s\n", mark, issue.Title)
		if issue.Reason != "" {
			fmt.Fprintf(w, "  Reason: %s\n", issue.Reason)
		}
		if issue.Recommendation != "" {
			fmt.Fprintf(w, "  Recommendation: %s\n", issue.Recommendation)
		}
		if issue.DocsLink != "" {
			fmt.Fprintf(w, "  Docs: %s\n", issue.DocsLink)
		}
	}

	critical, warning := r.counts()
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Cluster Score: %d/100\n", r.Score)
	fmt.Fprintf(w, "Health: %s\n", r.Health)
	fmt.Fprintf(w, "Warnings: %d\n", warning)
	fmt.Fprintf(w, "Critical: %d\n", critical)
}

// WriteJSON renders the report as indented JSON.
func (r Report) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// WriteHTML renders a minimal single-page HTML report. The full dashboard
// (charts, tables) described in the roadmap for v4.0 is not implemented yet.
func (r Report) WriteHTML(w io.Writer) error {
	var b strings.Builder
	b.WriteString("<!doctype html>\n<html><head><meta charset=\"utf-8\"><title>kubediag report</title></head><body>\n")
	fmt.Fprintf(&b, "<h1>Cluster Score: %d/100 (%s)</h1>\n", r.Score, r.Health)
	b.WriteString("<h2>Checks</h2>\n<ul>\n")
	for _, check := range r.Checks {
		status := "OK"
		if !check.Passed {
			status = "FAIL"
		}
		fmt.Fprintf(&b, "<li>[%s] %s</li>\n", status, check.Title)
	}
	b.WriteString("</ul>\n<h2>Issues</h2>\n<ul>\n")
	for _, issue := range r.Issues {
		fmt.Fprintf(&b, "<li><strong>%s</strong> (%s) — %s. Recommendation: %s</li>\n",
			issue.Title, issue.Severity, issue.Reason, issue.Recommendation)
	}
	b.WriteString("</ul>\n</body></html>\n")
	_, err := w.Write([]byte(b.String()))
	return err
}
