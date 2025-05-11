package main

import (
	"fmt"
	"os"
	"strconv"

	"text/tabwriter"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

// initializes all the cli commands
func setupCommands(storage *Storage) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tasks",
		Short: "tasks is a simple CLI task manager",
		Long:  `a simple command-line todo application for managing your tasks.`,
		// if no subcommand is provided, list tasks by default
		Run: func(cmd *cobra.Command, args []string) {
			listTasks(storage, false)
		},
	}

	// add command
	addCmd := &cobra.Command{
		Use:   "add [description]",
		Short: "add a new task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			description := args[0]
			err := storage.AddTask(description)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error adding task: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "task added: %s\n", description)
		},
	}

	// list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list all incomplete tasks",
		Run: func(cmd *cobra.Command, args []string) {
			showAll, _ := cmd.Flags().GetBool("all")
			listTasks(storage, showAll)
		},
	}
	listCmd.Flags().BoolP("all", "a", false, "show all tasks, including completed ones")

	// complete command
	completeCmd := &cobra.Command{
		Use:   "complete [taskID]",
		Short: "Mark a task as complete",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid task ID: %s\n", args[0])
				os.Exit(1)
			}

			err = storage.CompleteTask(id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error completing task: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "Task %d marked as complete\n", id)
		},
	}

	// delete command
	deleteCmd := &cobra.Command{
		Use:   "delete [taskID]",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid task ID: %s\n", args[0])
				os.Exit(1)
			}

			err = storage.DeleteTask(id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting task: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "Task %d deleted\n", id)
		},
	}

	rootCmd.AddCommand(addCmd, listCmd, completeCmd, deleteCmd)
	return rootCmd
}

// listTasks displays tasks in a tabular format
func listTasks(storage *Storage, showAll bool) {
	tasks, err := storage.ListTasks(showAll)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing tasks: %v\n", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Fprintln(os.Stdout, "No tasks found!")
		return
	}

	// create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// print header
	if showAll {
		fmt.Fprintln(w, "ID\tTask\tCreated\tDone")
	} else {
		fmt.Fprintln(w, "ID\tTask\tCreated")
	}

	// print tasks
	for _, task := range tasks {
		timeAgo := timediff.TimeDiff(task.CreatedAt)

		if showAll {
			fmt.Fprintf(w, "%d\t%s\t%s\t%t\n", task.ID, task.Description, timeAgo, task.IsComplete)
		} else {
			fmt.Fprintf(w, "%d\t%s\t%s\n", task.ID, task.Description, timeAgo)
		}
	}

	w.Flush()
}
