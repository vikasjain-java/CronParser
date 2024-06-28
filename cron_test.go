package main

import (
	"testing"
)

// TestParseCron tests the ParseCron function
func TestParseCron(t *testing.T) {
	tests := []struct {
		cronStr     string
		expected    *CronJob
		expectError bool
	}{
		{
			cronStr: "*/15 0 1,15 * 1-5 /usr/bin/find",
			expected: &CronJob{
				Minute:     "*/15",
				Hour:       "0",
				DayOfMonth: "1,15",
				Month:      "*",
				DayOfWeek:  "1-5",
				Command:    "/usr/bin/find",
			},
			expectError: false,
		},
		{
			cronStr:     "invalid cron string",
			expected:    nil,
			expectError: true,
		},
	}

	for _, test := range tests {
		result, err := parseCron(test.cronStr)
		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input: %s", test.cronStr)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect error for input: %s, got: %v", test.cronStr, err)
			}
			if result.Minute != test.expected.Minute || result.Hour != test.expected.Hour ||
				result.DayOfMonth != test.expected.DayOfMonth || result.Month != test.expected.Month ||
				result.DayOfWeek != test.expected.DayOfWeek || result.Command != test.expected.Command {
				t.Errorf("Expected: %+v, got: %+v", test.expected, result)
			}
		}
	}
}

// TestValidateCron tests the ValidateCron function
func TestValidateCron(t *testing.T) {
	tests := []struct {
		cronJob     *CronJob
		expectError bool
	}{
		{
			cronJob: &CronJob{
				Minute:     "*/15",
				Hour:       "0",
				DayOfMonth: "1,15",
				Month:      "*",
				DayOfWeek:  "1-5",
				Command:    "/usr/bin/find",
			},
			expectError: false,
		},
		{
			cronJob: &CronJob{
				Minute:     "*/60", // Invalid minute field
				Hour:       "0",
				DayOfMonth: "1,15",
				Month:      "*",
				DayOfWeek:  "1-5",
				Command:    "/usr/bin/find",
			},
			expectError: true,
		},
		{
			cronJob: &CronJob{
				Minute:     "*/15",
				Hour:       "0",
				DayOfMonth: "1,15",
				Month:      "*",
				DayOfWeek:  "1-5",
				Command:    "", // Empty command
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		err := validateCron(test.cronJob)
		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input: %+v", test.cronJob)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect error for input: %+v, got: %v", test.cronJob, err)
			}
		}
	}
}
