---
applyTo: '*'
---
I use a custom shell script named `gopu.sh` to streamline my development workflow. This script is globally available in my Git Bash terminal, so I can invoke it directly without needing to specify a path (e.g., no `./` prefix).

The `gopu.sh` script automates the following sequence of actions:
1.  **Adds all changed files** to the Git staging area (similar to `git add .`).
2.  **Runs Go tests** (equivalent to `go test ./...`).
3.  **Checks for race conditions** in the Go code (equivalent to `go test -race ./...`).
4.  **Commits the staged changes** using the message provided as an argument to the script.
5.  **Pushes the commit** to the remote repository.
6.  **Creates and pushes a tag** (the specifics of tag generation would be handled within the script).

The command to use the script is:
`gopu.sh "Your detailed commit message here"`

For example:
`gopu.sh "feat: implement user authentication module"`