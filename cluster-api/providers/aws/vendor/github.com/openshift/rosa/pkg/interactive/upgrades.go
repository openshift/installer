package interactive

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

func BuildManualUpgradeSchedule(cmd *cobra.Command, scheduleDate string, scheduleTime string) (time.Time, error) {
	// Set the default next run within the next 10 minutes
	now := time.Now().UTC().Add(time.Minute * 10)
	if scheduleDate == "" {
		scheduleDate = now.Format("2006-01-02")
	}
	if scheduleTime == "" {
		scheduleTime = now.Format("15:04")
	}

	if Enabled() {
		// If datetimes are set, use them in the interactive form, otherwise fallback to 'now'
		scheduleParsed, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", scheduleDate, scheduleTime))
		if err != nil {
			return now, fmt.Errorf("schedule date should use the format 'yyyy-mm-dd'\n" +
				"   Schedule time should use the format 'HH:mm'")
		}
		if scheduleParsed.IsZero() {
			scheduleParsed = now
		}
		scheduleDate = scheduleParsed.Format("2006-01-02")
		scheduleTime = scheduleParsed.Format("15:04")

		scheduleDate, err = GetString(Input{
			Question: "Please input desired date in format yyyy-mm-dd",
			Help:     cmd.Flags().Lookup("schedule-date").Usage,
			Default:  scheduleDate,
			Required: true,
		})
		if err != nil {
			return now, fmt.Errorf("expected a valid date: %s", err)
		}
		_, err = time.Parse("2006-01-02", scheduleDate)
		if err != nil {
			return now, fmt.Errorf("date format '%s' invalid", scheduleDate)
		}

		scheduleTime, err = GetString(Input{
			Question: "Please input desired UTC time in format HH:mm",
			Help:     cmd.Flags().Lookup("schedule-time").Usage,
			Default:  scheduleTime,
			Required: true,
		})
		if err != nil {
			return now, fmt.Errorf("expected a valid time: %s", err)
		}
		_, err = time.Parse("15:04", scheduleTime)
		if err != nil {
			return now, fmt.Errorf("time format '%s' invalid", scheduleTime)
		}
	}

	// Parse next run to time.Time
	nextRun, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", scheduleDate, scheduleTime))
	if err != nil {
		return now, fmt.Errorf("schedule date should use the format 'yyyy-mm-dd'\n" +
			"   Schedule time should use the format 'HH:mm'")
	}
	return nextRun, nil
}

func BuildAutomaticUpgradeSchedule(cmd *cobra.Command, schedule string) (string, error) {
	// Check automatic upgrade scheduling
	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	var err error
	if schedule != "" {
		_, err = cronParser.Parse(fmt.Sprintf("CRON_TZ=UTC %s", schedule))
		if err != nil {
			return schedule, fmt.Errorf("Schedule '%s' is not a valid cron expression", schedule)
		}
	}
	if Enabled() {
		schedule, err = GetString(Input{
			Question: "Please input desired automatic schedule with a cron expression",
			Help:     cmd.Flags().Lookup("schedule").Usage,
			Default:  schedule,
			Required: true,
		})
		if err != nil {
			return schedule, fmt.Errorf("Expected a valid automatic schedule: %s", err)
		}
		_, err = cronParser.Parse(fmt.Sprintf("CRON_TZ=UTC %s", schedule))
		if err != nil {
			return schedule, fmt.Errorf("Schedule '%s' is not a valid cron expression", schedule)
		}
	}

	return schedule, nil
}
