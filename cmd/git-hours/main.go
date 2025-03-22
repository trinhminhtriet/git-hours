package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/trinhminhtriet/git-hours/pkg/analyzer"
)

const (
	defaultMaxCommitDiff  = 120 // 2 hours in minutes
	defaultFirstCommitAdd = 120 // 2 hours in minutes
	defaultMergeRequest   = true
	defaultGitPath        = "."
	dateFormat            = "2006-01-02"
)

func main() {
	config := analyzer.Config{
		MaxCommitDiffInMinutes:       defaultMaxCommitDiff,
		FirstCommitAdditionInMinutes: defaultFirstCommitAdd,
		MergeRequest:                 defaultMergeRequest,
		GitPath:                      defaultGitPath,
		EmailAliases:                 make(map[string]string),
	}

	// Parse command line flags
	flag.IntVar(&config.MaxCommitDiffInMinutes, "max-commit-diff", defaultMaxCommitDiff, "maximum difference in minutes between commits counted to one session")
	flag.IntVar(&config.FirstCommitAdditionInMinutes, "first-commit-add", defaultFirstCommitAdd, "how many minutes first commit of session should add to total")
	since := flag.String("since", "always", "analyze data since certain date [always|yesterday|today|lastweek|thisweek|yyyy-mm-dd]")
	until := flag.String("until", "always", "analyze data until certain date [always|yesterday|today|lastweek|thisweek|yyyy-mm-dd]")
	flag.BoolVar(&config.MergeRequest, "merge-request", defaultMergeRequest, "include merge requests into calculation")
	flag.StringVar(&config.GitPath, "path", defaultGitPath, "git repository to analyze")
	flag.StringVar(&config.Branch, "branch", "", "analyze only data on the specified branch")

	flag.Parse()

	// Parse since/until dates
	config.Since = parseDateInput(*since)
	config.Until = parseDateInput(*until)

	// Check for shallow repository
	if _, err := os.Stat(".git/shallow"); err == nil {
		fmt.Println("Cannot analyze shallow copies!")
		fmt.Println("Please run git fetch --unshallow before continuing!")
		os.Exit(1)
	}

	// Run the analysis
	result, err := analyzer.AnalyzeRepository(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output results as JSON
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding results: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func parseDateInput(input string) time.Time {
	switch input {
	case "always":
		return time.Time{}
	case "today":
		return time.Now().Truncate(24 * time.Hour)
	case "yesterday":
		return time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	case "thisweek":
		now := time.Now()
		return now.AddDate(0, 0, -int(now.Weekday())).Truncate(24 * time.Hour)
	case "lastweek":
		now := time.Now()
		return now.AddDate(0, 0, -int(now.Weekday())-7).Truncate(24 * time.Hour)
	default:
		t, err := time.Parse(dateFormat, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid date format: %s. Expected format: YYYY-MM-DD\n", input)
			os.Exit(1)
		}
		return t
	}
}
