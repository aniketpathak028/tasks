# Todo App

A simple CLI application for managing tasks in the terminal.

## Usage

```bash
$ tasks
```

## Requirements

The application supports full CRUD operations via a command-line interface on a data file of tasks. The available commands are:

### Add

Create a new task by providing a description.

```bash
$ tasks add <description>
```

**Example:**

```bash
$ tasks add "Tidy my desk"
```

This will add a new task with the description "Tidy my desk".

### List

Lists all **uncompleted** tasks by default.

```bash
$ tasks list
```

**Example Output:**

```
ID    Task                                                Created
1     Tidy up my desk                                     a minute ago
3     Change my keyboard mapping to use escape/control    a few seconds ago
```

To show **all tasks** (including completed ones), use the `-a` or `--all` flag:

```bash
$ tasks list -a
```

**Example Output:**

```
ID    Task                                                Created          Done
1     Tidy up my desk                                     2 minutes ago    false
2     Write up documentation for new project feature      a minute ago     true
3     Change my keyboard mapping to use escape/control    a minute ago     false
```

### Complete

Mark a task as completed using its ID.

```bash
$ tasks complete <taskid>
```

### Delete

Delete a task from the data store using its ID.

```bash
$ tasks delete <taskid>
```

## Features

- file locking for secure access
- stdout and stderr for diagnostics or errors
