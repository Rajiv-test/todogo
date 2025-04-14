package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	Id           int
	Username     string
	TaskName     string
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Completed    bool
	Deadline     sql.NullTime
	OverDeadline bool
}

func (c *Client) AddTask(username string, taskname string, description string, deadline sql.NullTime) error {
	query := `INSERT INTO tasks (username,taskname,description,deadline)
				VALUES (?,?,?,?);`
	_, err := c.db.Exec(query, username, taskname, description, deadline)
	if err != nil {
		return fmt.Errorf("error while creating task %v", err)
	}
	fmt.Print("task successfully created")
	return nil
}

func (c *Client) GetTask(username, taskname string) (Task, error) {
	query := `SELECT * FROM tasks WHERE username = ? AND taskname = ?;`
	taskRow, err := c.db.Query(query, username, taskname)
	if err != nil {
		return Task{}, fmt.Errorf("error while fetching task %v", err)
	}
	var task Task
	err = taskRow.Scan(&task.Id, &task.Username, &task.TaskName, &task.Description, &task.CreatedAt, &task.UpdatedAt, &task.Completed, &task.Deadline, &task.OverDeadline)
	if err != nil {
		return Task{}, fmt.Errorf("error while scanning task %v", err)
	}
	return task, nil
}

func (c *Client) GetTasks(username string) ([]Task, error) {
	query := `SELECT * FROM tasks WHERE username = ?;`
	taskRows, err := c.db.Query(query, username)
	if err != nil {
		return []Task{}, fmt.Errorf("error while getting tasks %v", err)
	}
	var tasks []Task
	for taskRows.Next() {
		var task Task // or whatever your user type is
		err := taskRows.Scan(&task.Id, &task.Username, &task.TaskName, &task.Description, &task.CreatedAt, &task.UpdatedAt, &task.Completed, &task.Deadline, &task.OverDeadline)
		if err != nil {
			return []Task{}, err
		}
		// Add task to your slice
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (c *Client) MarkComplete(username, taskname string) error {
	query := `UPDATE tasks SET completed = 1,updated_at = ? WHERE username = ? AND taskname = ?;`
	_, err := c.db.Exec(query, time.Now(), username, taskname)
	if err != nil {
		return fmt.Errorf("error while updating task completion %v", err)
	}
	return nil
}

func (c *Client) ExtendDeadline(username, taskname string, deadline sql.NullTime) error {
	query := `UPDATE tasks SET deadline = ? WHERE username = ? AND taskname = ?;`
	_, err := c.db.Exec(query,deadline, username, taskname)
	if err != nil {
		return fmt.Errorf("error while extending deadline %v", err)
	}
	return nil
}
