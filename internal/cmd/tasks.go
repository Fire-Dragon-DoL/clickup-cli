package cmd

import (
	"fmt"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/api"
	"github.com/Fire-Dragon-DoL/clickup-cli/internal/resolver"
	"github.com/spf13/cobra"
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Manage tasks",
}

var tasksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks in a list",
	RunE: func(cmd *cobra.Command, args []string) error {
		listArg, err := cmd.Flags().GetString("list")
		if listArg == "" {
			PrintError(fmt.Errorf("--list flag is required"))
			return fmt.Errorf("--list flag is required")
		}

		recursive, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			PrintError(err)
			return err
		}

		kr := GetKeyring()
		apiKey, err := kr.GetAPIKey()
		if err != nil {
			PrintError(err)
			return err
		}

		client := api.NewClient(apiKey, "")
		cfg := GetConfig()
		res := resolver.New(client, cfg.StrictResolve)

		listID, err := res.ResolveList(listArg)
		if err != nil {
			PrintError(err)
			return err
		}

		resp, err := api.GetTasks(client, listID, recursive)
		if err != nil {
			PrintError(err)
			return err
		}

		formatted, err := formatTasksListView(resp.Tasks)
		if err != nil {
			PrintError(err)
			return err
		}

		fmt.Println(formatted)
		return nil
	},
}

var tasksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, err := cmd.Flags().GetString("title")
		if title == "" {
			PrintError(fmt.Errorf("--title flag is required"))
			return fmt.Errorf("--title flag is required")
		}

		listArg, err := cmd.Flags().GetString("list")
		if listArg == "" {
			PrintError(fmt.Errorf("--list flag is required"))
			return fmt.Errorf("--list flag is required")
		}

		kr := GetKeyring()
		apiKey, err := kr.GetAPIKey()
		if err != nil {
			PrintError(err)
			return err
		}

		client := api.NewClient(apiKey, "")
		cfg := GetConfig()
		res := resolver.New(client, cfg.StrictResolve)

		listID, err := res.ResolveList(listArg)
		if err != nil {
			PrintError(err)
			return err
		}

		payload := map[string]any{
			"name":    title,
			"list_id": listID,
		}

		description, _ := cmd.Flags().GetString("description")
		if description != "" {
			payload["description"] = description
		}

		priority, _ := cmd.Flags().GetInt("priority")
		if priority != 0 {
			payload["priority"] = priority
		}

		status, _ := cmd.Flags().GetString("status")
		if status != "" {
			payload["status"] = status
		}

		dueDate, _ := cmd.Flags().GetString("due")
		if dueDate != "" {
			payload["due_date"] = dueDate
		}

		assignee, _ := cmd.Flags().GetString("assignee")
		if assignee != "" {
			assigneeID, err := res.ResolveUser(assignee)
			if err != nil {
				PrintError(err)
				return err
			}
			payload["assignee"] = assigneeID
		}

		parent, _ := cmd.Flags().GetString("parent")
		if parent != "" {
			parentID, err := res.ResolveTask(parent)
			if err != nil {
				PrintError(err)
				return err
			}
			payload["parent"] = parentID
		}

		task, err := api.CreateTask(client, payload)
		if err != nil {
			PrintError(err)
			return err
		}

		formatted, err := formatTaskDetailsView(task)
		if err != nil {
			PrintError(err)
			return err
		}

		fmt.Println(formatted)
		return nil
	},
}

func formatTasksListView(tasks []api.Task) (string, error) {
	formatter := GetFormatter()

	type TaskListView struct {
		ID       string
		Title    string
		Assignee string
		Status   string
		Priority string
	}

	var views []TaskListView
	for _, task := range tasks {
		view := TaskListView{
			ID:    task.ID,
			Title: task.Name,
		}

		if task.Assignee != nil {
			view.Assignee = task.Assignee.Username
		}

		if task.Status != nil {
			view.Status = task.Status.Status
		}

		if task.Priority != nil {
			view.Priority = task.Priority.Priority
		}

		views = append(views, view)
	}

	return formatter.Format(views)
}

func formatTaskDetailsView(task api.Task, comments ...api.Comment) (string, error) {
	formatter := GetFormatter()

	type CommentView struct {
		Author  string
		Content string
		Date    string
	}

	type TaskDetailsView struct {
		ID          string
		Title       string
		Description string
		Assignee    string
		Status      string
		Priority    string
		DueDate     string
		Comments    []CommentView
	}

	view := TaskDetailsView{
		ID:          task.ID,
		Title:       task.Name,
		Description: task.Description,
		DueDate:     task.DueDate,
	}

	if task.Assignee != nil {
		view.Assignee = task.Assignee.Username
	}

	if task.Status != nil {
		view.Status = task.Status.Status
	}

	if task.Priority != nil {
		view.Priority = task.Priority.Priority
	}

	if len(comments) > 0 {
		for _, comment := range comments {
			view.Comments = append(view.Comments, CommentView{
				Author:  comment.User.Username,
				Content: comment.TextContent,
				Date:    comment.DateCreated,
			})
		}
	}

	return formatter.Format(view)
}

var tasksShowCmd = &cobra.Command{
	Use:   "show <task-id|name|url>",
	Short: "Show task details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskArg := args[0]

		kr := GetKeyring()
		apiKey, err := kr.GetAPIKey()
		if err != nil {
			PrintError(err)
			return err
		}

		client := api.NewClient(apiKey, "")
		cfg := GetConfig()
		res := resolver.New(client, cfg.StrictResolve)

		taskID, err := res.ResolveTask(taskArg)
		if err != nil {
			PrintError(err)
			return err
		}

		task, err := api.GetTask(client, taskID)
		if err != nil {
			PrintError(err)
			return err
		}

		comments, err := api.GetTaskComments(client, task.ID)
		if err != nil {
			PrintError(err)
			return err
		}

		formatted, err := formatTaskDetailsView(task, comments...)
		if err != nil {
			PrintError(err)
			return err
		}

		fmt.Println(formatted)
		return nil
	},
}

var tasksUpdateCmd = &cobra.Command{
	Use:   "update <task-id|name|url>",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskArg := args[0]

		kr := GetKeyring()
		apiKey, err := kr.GetAPIKey()
		if err != nil {
			PrintError(err)
			return err
		}

		client := api.NewClient(apiKey, "")
		cfg := GetConfig()
		res := resolver.New(client, cfg.StrictResolve)

		taskID, err := res.ResolveTask(taskArg)
		if err != nil {
			PrintError(err)
			return err
		}

		// Build update payload from flags
		payload := make(map[string]any)

		title, _ := cmd.Flags().GetString("title")
		if title != "" {
			payload["name"] = title
		}

		assignee, _ := cmd.Flags().GetString("assignee")
		if assignee != "" {
			assigneeID, err := res.ResolveUser(assignee)
			if err != nil {
				PrintError(fmt.Errorf("failed to resolve assignee: %w", err))
				return err
			}
			payload["assignee"] = assigneeID
		}

		status, _ := cmd.Flags().GetString("status")
		if status != "" {
			payload["status"] = status
		}

		priority, _ := cmd.Flags().GetString("priority")
		if priority != "" {
			payload["priority"] = priority
		}

		description, _ := cmd.Flags().GetString("description")
		if description != "" {
			payload["description"] = description
		}

		dueDate, _ := cmd.Flags().GetString("due")
		if dueDate != "" {
			payload["due_date"] = dueDate
		}

		parent, _ := cmd.Flags().GetString("parent")
		if parent != "" {
			parentID, err := res.ResolveTask(parent)
			if err != nil {
				PrintError(fmt.Errorf("failed to resolve parent task: %w", err))
				return err
			}
			payload["parent"] = parentID
		}

		// Perform update
		updated, err := api.UpdateTask(client, taskID, payload)
		if err != nil {
			PrintError(err)
			return err
		}

		// Format as Details View
		formatted, err := formatTaskDetailsView(updated)
		if err != nil {
			PrintError(err)
			return err
		}

		fmt.Println(formatted)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tasksCmd)
	tasksCmd.AddCommand(tasksListCmd)
	tasksCmd.AddCommand(tasksShowCmd)
	tasksCmd.AddCommand(tasksCreateCmd)
	tasksCmd.AddCommand(tasksUpdateCmd)
	tasksListCmd.Flags().StringP("list", "l", "", "list name, ID, or URL")
	tasksListCmd.Flags().BoolP("recursive", "r", false, "include subtasks")
	tasksCreateCmd.Flags().StringP("title", "t", "", "task title")
	tasksCreateCmd.Flags().StringP("list", "l", "", "list name, ID, or URL")
	tasksCreateCmd.Flags().StringP("description", "d", "", "task description")
	tasksCreateCmd.Flags().IntP("priority", "p", 0, "task priority (1-5, 0=none)")
	tasksCreateCmd.Flags().String("status", "", "task status")
	tasksCreateCmd.Flags().String("due", "", "due date")
	tasksCreateCmd.Flags().String("assignee", "", "assignee name, ID, or username")
	tasksCreateCmd.Flags().String("parent", "", "parent task ID or name")
	tasksUpdateCmd.Flags().StringP("title", "t", "", "task title")
	tasksUpdateCmd.Flags().StringP("assignee", "a", "", "assignee name, ID, or username")
	tasksUpdateCmd.Flags().StringP("status", "s", "", "task status")
	tasksUpdateCmd.Flags().StringP("priority", "p", "", "task priority")
	tasksUpdateCmd.Flags().StringP("description", "d", "", "task description (markdown)")
	tasksUpdateCmd.Flags().String("due", "", "due date")
	tasksUpdateCmd.Flags().String("parent", "", "parent task name, ID, or URL")

	tasksCmd.AddCommand(tasksDeleteCmd)
	tasksCmd.AddCommand(tasksArchiveCmd)
}

var tasksDeleteCmd = &cobra.Command{
	Use:   "delete <task-id|name|url>",
	Short: "Delete a task permanently",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskArg := args[0]

		kr := GetKeyring()
		apiKey, err := kr.GetAPIKey()
		if err != nil {
			PrintError(err)
			return err
		}

		client := api.NewClient(apiKey, "")
		cfg := GetConfig()
		res := resolver.New(client, cfg.StrictResolve)

		taskID, err := res.ResolveTask(taskArg)
		if err != nil {
			PrintError(err)
			return err
		}

		err = api.DeleteTask(client, taskID)
		if err != nil {
			PrintError(err)
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Task %s deleted\n", taskID)
		return nil
	},
}

var tasksArchiveCmd = &cobra.Command{
	Use:   "archive <task-id|name|url>",
	Short: "Archive a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskArg := args[0]

		kr := GetKeyring()
		apiKey, err := kr.GetAPIKey()
		if err != nil {
			PrintError(err)
			return err
		}

		client := api.NewClient(apiKey, "")
		cfg := GetConfig()
		res := resolver.New(client, cfg.StrictResolve)

		taskID, err := res.ResolveTask(taskArg)
		if err != nil {
			PrintError(err)
			return err
		}

		err = api.ArchiveTask(client, taskID)
		if err != nil {
			PrintError(err)
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Task %s archived\n", taskID)
		return nil
	},
}
