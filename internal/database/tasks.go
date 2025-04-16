package database

import (
	"database/sql"
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
	query := `INSERT INTO tasks (username,taskname,description,created_at,updated_at,deadline)
				VALUES (?,?,?,?,?,?);`
	_, err := c.db.Exec(query, username, taskname, description,time.Now(),time.Now(), deadline)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetTask(username, taskname string) (Task, error) {
	query := `SELECT * FROM tasks WHERE username = ? AND taskname = ?;`
	taskRow, err := c.db.Query(query, username, taskname)
	if err != nil {
		return Task{},err
	}
	var task Task
	err = taskRow.Scan(&task.Id, &task.Username, &task.TaskName, &task.Description, &task.CreatedAt, &task.UpdatedAt, &task.Deadline, &task.Completed, &task.OverDeadline)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (c *Client) GetTasks(username string) ([]Task, error) {
	updateQuery := `
        UPDATE tasks 
        SET over_deadline = 1 
        WHERE username = ? 
        AND completed = 0
        AND deadline IS NOT NULL 
        AND deadline < datetime('now')
    `
    _, err := c.db.Exec(updateQuery, username)
    if err != nil {
        return nil, err
    }
	query := `SELECT * FROM tasks
		WHERE username = ? 
		ORDER BY 
			completed ASC,
			CASE 
				WHEN deadline IS NULL THEN 1
				ELSE 0
			END,
			deadline ASC,
			id ASC;`
	taskRows, err := c.db.Query(query, username)
	if err != nil {
		return []Task{}, err
	}
	var tasks []Task
	for taskRows.Next() {
		var task Task // or whatever your user type is
		err := taskRows.Scan(&task.Id, &task.Username, &task.TaskName, &task.Description, &task.CreatedAt, &task.UpdatedAt, &task.Deadline, &task.Completed, &task.OverDeadline)
		if err != nil {
			return []Task{}, err
		}
		// Add task to your slice
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (c *Client) GetUncompletedTasks(username string) ([]Task,error){
	updateQuery := `
        UPDATE tasks 
        SET over_deadline = 1 
        WHERE username = ? 
        AND completed = 0
        AND deadline IS NOT NULL 
        AND deadline < datetime('now')
    `
    _, err := c.db.Exec(updateQuery, username)
    if err != nil {
        return nil, err
    }
	query := `SELECT * FROM tasks
		WHERE username = ? AND completed = 0
		ORDER BY 
			CASE WHEN deadline IS NULL THEN 1 ELSE 0 END,
			deadline ASC;`
	taskRows, err := c.db.Query(query, username)
	if err != nil {
		return []Task{}, err
	}
	var tasks []Task
	for taskRows.Next() {
		var task Task // or whatever your user type is
		err := taskRows.Scan(&task.Id, &task.Username, &task.TaskName, &task.Description, &task.CreatedAt, &task.UpdatedAt,&task.Deadline, &task.Completed, &task.OverDeadline)
		if err != nil {
			return []Task{}, err
		}
		// Add task to your slice
		tasks = append(tasks, task)
	}
	return tasks, nil

}





func (c *Client) DeleteTask(username, taskname string) error{
	query := "DELETE FROM tasks WHERE username = ? AND taskname = ?;"
	_,err := c.db.Exec(query,username,taskname) 
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteTasks(username string) error{
	query := "DELETE FROM tasks WHERE username = ?;"
	_,err := c.db.Exec(query,username) 
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) MarkComplete(username, taskname string) error {
	deadline := sql.NullTime{Time: time.Time{},Valid: false}
	query := `UPDATE tasks SET completed = 1,updated_at = ?, deadline = ? WHERE username = ? AND taskname = ?;`
	_, err := c.db.Exec(query, time.Now(),deadline, username, taskname)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ExtendDeadline(username, taskname string, deadline sql.NullTime) error {
	query := `UPDATE tasks SET updated_at = ?, deadline = ? WHERE username = ? AND taskname = ?;`
	_, err := c.db.Exec(query,time.Now(),deadline, username, taskname)
	if err != nil {
		return err
	}
	return nil
}
