package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var parsedHashMap = make(map[string]string)

// CronJob represents a cron job schedule
type CronJob struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Command    string
}

func expandField(field string, min int, max int) []string {
	var result []string
	if field == "*" {
		for i := min; i <= max; i++ {
			result = append(result, strconv.Itoa(i))
		}
	} else if strings.Contains(field, "/") {
		parts := strings.Split(field, "/")
		step, _ := strconv.Atoi(parts[1])
		for i := min; i <= max; i += step {
			result = append(result, strconv.Itoa(i))
		}
	} else if strings.Contains(field, ",") {
		parts := strings.Split(field, ",")
		for _, part := range parts {
			result = append(result, part)
		}
	} else if strings.Contains(field, "-") {
		parts := strings.Split(field, "-")
		start, _ := strconv.Atoi(parts[0])
		end, _ := strconv.Atoi(parts[1])
		for i := start; i <= end; i++ {
			result = append(result, strconv.Itoa(i))
		}
	} else {
		result = append(result, field)
	}
	return result
}

func generateParsedValueFromCronfields(cron *CronJob) {

	parsedHashMap["minuteField"] = strings.Join(expandField(cron.Minute, 0, 59), " ")

	parsedHashMap["hourField"] = strings.Join(expandField(cron.Hour, 0, 23), " ")

	parsedHashMap["dayOfMonthField"] = strings.Join(expandField(cron.DayOfMonth, 1, 31), " ")

	parsedHashMap["monthField"] = strings.Join(expandField(cron.Month, 1, 12), " ")

	parsedHashMap["dayOfWeekField"] = strings.Join(expandField(cron.DayOfWeek, 0, 6), " ")

	parsedHashMap["command"] = cron.Command
}

func ValidateCronField(field string, pattern string) error {
	matched, err := regexp.MatchString(pattern, field)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("invalid cron field: %s", field)
	}
	return nil
}

// ValidateCron validates all fields of the cron job
func validateCron(cron *CronJob) error {
	patterns := map[string]string{
		"minute":     `^(\*|([0-5]?[0-9])|(\*\/[0-5]?[0-9])|([0-5]?[0-9]-[0-5]?[0-9])|([0-5]?[0-9](,[0-5]?[0-9])*))$`,
		"hour":       `^(\*|([0-1]?[0-9]|2[0-3])|(\*\/([0-1]?[0-9]|2[0-3]))|(([0-1]?[0-9]|2[0-3])-([0-1]?[0-9]|2[0-3]))|(([0-1]?[0-9]|2[0-3])(,([0-1]?[0-9]|2[0-3]))*))$`,
		"dayOfMonth": `^(\*|(0?[1-9]|[12][0-9]|3[01])|(\*\/(0?[1-9]|[12][0-9]|3[01]))|((0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|[12][0-9]|3[01]))|((0?[1-9]|[12][0-9]|3[01])(,(0?[1-9]|[12][0-9]|3[01]))*))$`,
		"month":      `^(\*|(0?[1-9]|1[0-2])|(\*\/(0?[1-9]|1[0-2]))|((0?[1-9]|1[0-2])-(0?[1-9]|1[0-2]))|((0?[1-9]|1[0-2])(,(0?[1-9]|1[0-2]))*))$`,
		"dayOfWeek":  `^(\*|([0-6])|(\*\/[0-6])|([0-6]-[0-6])|([0-6](,[0-6])*))$`,
	}

	err := ValidateCronField(cron.Minute, patterns["minute"])
	if err != nil {
		return errors.New("invalid minute field: " + err.Error())
	}

	err = ValidateCronField(cron.Hour, patterns["hour"])
	if err != nil {
		return errors.New("invalid hour field: " + err.Error())
	}

	err = ValidateCronField(cron.DayOfMonth, patterns["dayOfMonth"])
	if err != nil {
		return errors.New("invalid day of month field: " + err.Error())
	}

	err = ValidateCronField(cron.Month, patterns["month"])
	if err != nil {
		return errors.New("invalid month field: " + err.Error())
	}

	err = ValidateCronField(cron.DayOfWeek, patterns["dayOfWeek"])
	if err != nil {
		return errors.New("invalid day of week field: " + err.Error())
	}

	if strings.TrimSpace(cron.Command) == "" {
		return errors.New("command cannot be empty")
	}

	return nil
}

func parseCron(cronStr string) (*CronJob, error) {
	fields := strings.Fields(cronStr)
	if len(fields) < 6 {
		return nil, fmt.Errorf("invalid cron string format")
	}

	return &CronJob{
		Minute:     fields[0],
		Hour:       fields[1],
		DayOfMonth: fields[2],
		Month:      fields[3],
		DayOfWeek:  fields[4],
		Command:    strings.Join(fields[5:], " "),
	}, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run your-program.go '*/15 0 1,15 * 1-5 /usr/bin/find'")
	}

	cronStr := os.Args[1]
	cron, err := parseCron(cronStr)
	if err != nil {
		fmt.Println("Error Validating cron string", err)
	}
	err = validateCron(cron)
	if err != nil {
		fmt.Println("Error Validating cron fields", err)
	} else {
		generateParsedValueFromCronfields(cron)
		for key, value := range parsedHashMap {
			fmt.Printf("key= %s and  value= %s \n", key, value)
		}
	}

}
