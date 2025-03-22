# git-hours

A Go implementation of a tool to estimate time spent on a git repository based on commit history.

## 🚀 Installation

```bash
go install github.com/trinhminhtriet/git-hours/cmd/git-hours@latest
```

## 💡 Usage

```bash
git-hours [options]
```

### Options

- `-max-commit-diff`: Maximum difference in minutes between commits counted to one session (default: 120)
- `-first-commit-add`: How many minutes first commit of session should add to total (default: 120)
- `-since`: Analyze data since certain date [always|yesterday|today|lastweek|thisweek|yyyy-mm-dd] (default: always)
- `-until`: Analyze data until certain date [always|yesterday|today|lastweek|thisweek|yyyy-mm-dd] (default: always)
- `-merge-request`: Include merge requests into calculation (default: true)
- `-path`: Git repository to analyze (default: ".")
- `-branch`: Analyze only data on the specified branch (default: all branches)

### Examples

1. Estimate hours of project:

```bash
git-hours
```

2. Estimate hours in repository where developers commit more seldom (4h pause between commits):

```bash
git-hours --max-commit-diff 240
```

3. Estimate hours in repository where developer works 5 hours before first commit in day:

```bash
git-hours --first-commit-add 300
```

4. Estimate hours work in repository since yesterday:

```bash
git-hours --since yesterday
```

5. Estimate hours work in repository since specific date:

```bash
git-hours --since 2015-01-31
```

6. Estimate hours work in repository on the "master" branch:

```bash
git-hours --branch master
```

## Output

The tool outputs JSON with the following structure:

```json
{
  "author@email.com": {
    "name": "Author Name",
    "hours": 42,
    "commits": 100
  },
  "total": {
    "hours": 42,
    "commits": 100
  }
}
```

## Notes

- Cannot analyze shallow copies. Run `git fetch --unshallow` first if needed.
- Time estimation is based on commit dates and configured time windows
- Merge commits can be excluded from the calculation

## 🤝 How to contribute

We welcome contributions!

- Fork this repository;
- Create a branch with your feature: `git checkout -b my-feature`;
- Commit your changes: `git commit -m "feat: my new feature"`;
- Push to your branch: `git push origin my-feature`.

Once your pull request has been merged, you can delete your branch.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
