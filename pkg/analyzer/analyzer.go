package analyzer

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Config holds the configuration for the git analysis
type Config struct {
	MaxCommitDiffInMinutes       int
	FirstCommitAdditionInMinutes int
	Since                        time.Time
	Until                        time.Time
	MergeRequest                 bool
	GitPath                      string
	Branch                       string
	EmailAliases                 map[string]string
}

// AuthorWork represents the work done by a single author
type AuthorWork struct {
	Name    string `json:"name"`
	Hours   int    `json:"hours"`
	Commits int    `json:"commits"`
}

// Result represents the final analysis result
type Result map[string]AuthorWork

// AnalyzeRepository performs the git analysis and returns the results
func AnalyzeRepository(config *Config) (Result, error) {
	// Open the repository
	repo, err := git.PlainOpen(config.GitPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the appropriate reference (branch or HEAD)
	var ref *plumbing.Reference
	if config.Branch != "" {
		ref, err = repo.Reference(plumbing.NewBranchReferenceName(config.Branch), true)
	} else {
		ref, err = repo.Head()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get reference: %w", err)
	}

	// Get the commit history
	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit history: %w", err)
	}

	// Group commits by author
	authorCommits := make(map[string][]time.Time)
	var totalCommits int

	err = commitIter.ForEach(func(c *object.Commit) error {
		// Skip merge commits if configured
		if !config.MergeRequest && strings.HasPrefix(c.Message, "Merge ") {
			return nil
		}

		// Check date range
		if !config.Since.IsZero() && c.Author.When.Before(config.Since) {
			return nil
		}
		if !config.Until.IsZero() && c.Author.When.After(config.Until) {
			return nil
		}

		email := c.Author.Email
		if alias, ok := config.EmailAliases[email]; ok {
			email = alias
		}

		authorCommits[email] = append(authorCommits[email], c.Author.When)
		totalCommits++
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process commits: %w", err)
	}

	// Calculate hours for each author
	result := make(Result)
	var totalHours int

	for email, dates := range authorCommits {
		hours := estimateHours(dates, config)
		name := getAuthorName(email)

		result[email] = AuthorWork{
			Name:    name,
			Hours:   hours,
			Commits: len(dates),
		}
		totalHours += hours
	}

	// Add total
	result["total"] = AuthorWork{
		Hours:   totalHours,
		Commits: totalCommits,
	}

	return result, nil
}

// estimateHours calculates the estimated hours worked based on commit dates
func estimateHours(dates []time.Time, config *Config) int {
	if len(dates) < 2 {
		return 0
	}

	// Sort dates (oldest first)
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	var totalMinutes float64

	// Calculate time between consecutive commits
	for i := 0; i < len(dates)-1; i++ {
		diffMinutes := dates[i+1].Sub(dates[i]).Minutes()

		if diffMinutes < float64(config.MaxCommitDiffInMinutes) {
			totalMinutes += diffMinutes
		} else {
			// Add the first commit addition for this new session
			totalMinutes += float64(config.FirstCommitAdditionInMinutes)
		}
	}

	return int(math.Round(totalMinutes / 60))
}

// getAuthorName extracts the author name from the first commit
func getAuthorName(email string) string {
	// In a real implementation, we would store the author name from the first commit
	// For now, just return the email as the name
	return email
}
