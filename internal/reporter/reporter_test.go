package reporter

import "testing"

func TestNewReportScore(t *testing.T) {
	tests := []struct {
		name       string
		issues     []Issue
		wantScore  int
		wantHealth string
	}{
		{name: "no issues", issues: nil, wantScore: 100, wantHealth: "Excellent"},
		{
			name:       "one warning",
			issues:     []Issue{{Severity: SeverityWarning}},
			wantScore:  95,
			wantHealth: "Excellent",
		},
		{
			name:       "one critical",
			issues:     []Issue{{Severity: SeverityCritical}},
			wantScore:  85,
			wantHealth: "Good",
		},
		{
			name: "score floors at zero",
			issues: []Issue{
				{Severity: SeverityCritical}, {Severity: SeverityCritical}, {Severity: SeverityCritical},
				{Severity: SeverityCritical}, {Severity: SeverityCritical}, {Severity: SeverityCritical},
				{Severity: SeverityCritical},
			},
			wantScore:  0,
			wantHealth: "Critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := NewReport("v1.30.0", nil, tt.issues)
			if report.Score != tt.wantScore {
				t.Errorf("Score = %d, want %d", report.Score, tt.wantScore)
			}
			if report.Health != tt.wantHealth {
				t.Errorf("Health = %q, want %q", report.Health, tt.wantHealth)
			}
		})
	}
}
