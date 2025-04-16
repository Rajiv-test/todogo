package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"os"
	"strconv"
	"strings"
	"time"
)

func commandAdd(c *config, args ...string) error {
	if len(args) > 2 || len(args) < 1 {
		return fmt.Errorf("wrong usage of add command use help command to learn more")
	}
	if c.user == nil {
		return fmt.Errorf("login to create tasks")
	}
	taskname := args[0]
	var deadline sql.NullTime
	var err error
	if len(args) == 2 {
		deadline, err = parseDeadline(args[1])
		if err != nil {
			return fmt.Errorf("enter valid timelimit %v", err)
		}
	} else {
		deadline = sql.NullTime{Time: time.Time{}, Valid: false}
	}
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter task description: ")
	reader.Scan()
	description := reader.Text()
	err = c.db.AddTask(c.user.Name, taskname, description, deadline)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: tasks.taskname" {
			return fmt.Errorf("task name taken choose another")
		} else {
			return fmt.Errorf("error adding task %v", err)
		}
	}
	fmt.Println("task created successfully")
	return nil

}

func commandRemove(c *config, args ...string) error {
	if len(args) != 1 && (len(args) == 2 && (args[0] != "-a" || args[1] != "tasks")) {
		return fmt.Errorf("wrong usage of rem command use help command to learn more")
	}
	if c.user == nil {
		return fmt.Errorf("login to create tasks")
	}
	if len(args) == 1 {
		if len(args[0]) == 0 {
			return fmt.Errorf("enter valid taskname")
		}
		err := c.db.DeleteTask(c.user.Name, args[0])
		if err != nil {
			return err
		}
		fmt.Println("task successfully deleted")
		return nil
	} else {
		err := c.db.DeleteTasks(c.user.Name)
		if err != nil {
			return err
		}
		fmt.Println("Successfully deleted all tasks for user", c.user.Name)
	}
	return nil
}

func commandListTasks(c *config, args ...string) error {
	if c.user == nil {
		return fmt.Errorf("login to see tasks")
	}
	if len(args) > 1 || (len(args) == 1 && args[0] != "-u") {
		return fmt.Errorf("wrong usage of lst command use help command to learn more")
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgRed).SprintfFunc()

	tbl := table.New("No", "TaskName", "CreatedAt", "UpdatedAt", "deadline", "completed", "overDeadline", "descritpion")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	if len(args) == 0 {
		tasks, err := c.db.GetTasks(c.user.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("no tasks for current user")
			}
			return fmt.Errorf("error getting tasks %v", err)
		}

		for i, task := range tasks {
			tbl.AddRow(i+1, task.TaskName, task.CreatedAt.Format(time.DateTime), task.UpdatedAt.Format(time.DateTime), task.Deadline.Time.Format(time.DateTime), task.Completed, task.OverDeadline, task.Description)
		}

		tbl.Print()
		return nil

	}
	if len(args) == 1 {
		tasks, err := c.db.GetUncompletedTasks(c.user.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("no tasks for current user")
			}
			return fmt.Errorf("error getting uncompleted tasks %v", err)
		}
		for i, task := range tasks {
			if !task.Deadline.Valid {
				task.Deadline.Time = time.Time{}
			}
			tbl.AddRow(i+1, task.TaskName, task.CreatedAt.Format(time.DateTime), task.UpdatedAt.Format(time.DateTime), task.Deadline.Time.Format(time.DateTime), task.Completed, task.OverDeadline, task.Description)
		}

		tbl.Print()

	}
	return nil
}

func parseDeadline(durationStr string) (sql.NullTime, error) {
	// Handle days specifically since Go doesn't have a direct "d" unit
	if strings.HasSuffix(durationStr, "d") {
		// Extract the number part
		days, err := strconv.Atoi(durationStr[:len(durationStr)-1])
		if err != nil {
			return sql.NullTime{}, err
		}
		// Calculate deadline
		return sql.NullTime{Time: time.Now().Local().Add(time.Hour * 24 * time.Duration(days)), Valid: true}, nil
	}

	// For hours, minutes, seconds, Go's time.ParseDuration works directly
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return sql.NullTime{}, err
	}

	return sql.NullTime{Time: time.Now().Local().Add(duration), Valid: true}, nil
}

func commandMarkComplete(c *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong usage of tick command use help command to learn more")
	}
	if c.user == nil {
		return fmt.Errorf("login to use commands")
	}
	if args[0] == "" {
		return fmt.Errorf("enter valid taskname")
	}
	task, err := c.db.GetTask(c.user.Name, args[0])
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no task named %v use lst command to view your tasks", args[0])
		}
		return fmt.Errorf("error while fetching task %v", err)
	}
	if task.Completed {
		return fmt.Errorf("task already completed")
	}
	err = c.db.MarkComplete(c.user.Name, task.TaskName)
	if err != nil {
		return fmt.Errorf("error while marking task complete %v", err)
	}
	fmt.Println("Hooray! you completed a task")
	return nil
}
func commandMarkIncomplete(c *config, args ...string) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("wrong usage of tick command use help command to learn more")
	}
	if c.user == nil {
		return fmt.Errorf("login to use commands")
	}
	if args[0] == "" {
		return fmt.Errorf("enter valid taskname")
	}
	var deadline sql.NullTime
	var err error
	if len(args) == 2 {
		deadline, err = parseDeadline(args[1])
		if err != nil {
			return fmt.Errorf("enter valid deadline %v", err)
		}
	}
	task, err := c.db.GetTask(c.user.Name, args[0])
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no task named %v use lst command to view your tasks", deadline)
		}
		return fmt.Errorf("error while fetching task %v", err)
	}
	if !task.Completed {
		return fmt.Errorf("task is still incomplete")
	}
	err = c.db.MarkIncomplete(c.user.Name, task.TaskName, deadline)
	if err != nil {
		return fmt.Errorf("error while marking task incomplete %v", err)
	}
	fmt.Println("successfully changed the status of task to incomplete")
	return nil
}
