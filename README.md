# MicroGit

**MicroGit** is a lightweight, educational version control system designed to help beginners understand how Git works under the hood. It provides basic functionality like tracking file changes, staging, committing, diffing, and reverting.

---

## ✨ Features

- Initialize a new repository
- Track and stage files
- Commit snapshots with messages
- View commit history
- See file differences
- Revert to a previous version
- Simple, linear commit structure (no branching)

---

## 🚀 Getting Started

### 1. Clone or Download

```bash
git clone https://github.com/your-username/microgit.git
cd microgit
```

## Supported commands

### `microgit init`
Initialize a new MicroGit repository in the current directory.
This creates the necessary directory structure and files for version control.
The repository will be initialized in a .microgit directory.

### `microgit add [files...]`
Add files to the staging area for the next commit.

Usage:
- `microgit add <file1> [file2 ...]` - Stage specific files
- `microgit add .` - Stage all files in current directory

The command will:
1. Calculate a SHA-256 hash of the file content
2. Store the file content in the objects directory
3. Update the index with the file path and corresponding hash

### `microgit remove [files...]`
Remove files from the staging area, effectively un-staging them.

Usage:
- `microgit remove <file1> [file2 ...]` - Remove specific files from staging
- `microgit remove .` - Remove all files from staging

The command will:
1. Remove the specified files from the index
2. Keep the files in your working directory
3. Allow you to re-stage them later if needed

### `microgit status`
Show the working tree status.

Displays the state of the working directory and the staging area.
Shows which files have been staged for the next commit and which files
are untracked. This helps you understand what will be included in your
next commit.

### `microgit save "message"`
Save the current state of staged files as a new commit.

This command requires a commit message that describes the changes being saved.
The staged files will be committed and the staging area will be cleared after the save.

### `microgit log`
Show the commit history.

Displays the commit history in chronological order, starting from the most recent commit.
For each commit, it shows:
- The commit hash
- The timestamp
- The commit message
- The list of files that were modified

### `microgit checkout <commit>`
Switch to a specific commit in the repository history.

Usage:
- `microgit checkout <commit-hash>` - Switch to a specific commit
- `microgit checkout latest` - Switch to the most recent commit

This command will:
1. Restore all files to their state at the specified commit
2. Update the HEAD reference to point to the checked out commit
3. Preserve the commit history for future operations

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
