package main

import (
	"bufio"
	"database/sql"
	"fmt"
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
		if err.Error() == "UNIQUE constraint failed: tasks.taskname"{
			return fmt.Errorf("task name taken choose another")
		}else{
			return fmt.Errorf("error adding task %v", err)
		}
	}
	fmt.Println("task created successfully")
	return nil

}

func commandRemove(c *config,args ...string) error{
	if len(args) != 1 && (len(args) == 2 && (args[0] != "-a" || args[1] != "tasks")){
		return fmt.Errorf("wrong usage of rem command use help command to learn more")
	}
	if c.user == nil {
		return fmt.Errorf("login to create tasks")
	}
	if len(args) == 1{
		if len(args[0]) == 0{
			return fmt.Errorf("enter valid taskname")
		}
		err := c.db.DeleteTask(c.user.Name,args[0])
		if err != nil{
			return err
		}
		fmt.Println("task successfully deleted")
		return nil
	}else{
		err := c.db.DeleteTasks(c.user.Name)
		if err != nil{
			return err
	}
	fmt.Println("Successfully deleted all tasks for user",c.user.Name)
}
	return nil
}

func commandListTasks(c *config,args ...string)error{
	if c.user == nil {
		return fmt.Errorf("login to see tasks")
	}
	if len(args) > 1 || (len(args) == 1 && args[0] != "-u"){
		return fmt.Errorf("wrong usage of lst command use help command to learn more")
	}
	if len(args) == 0 {
		tasks,err := c.db.GetTasks(c.user.Name)
		if err != nil {
			if err == sql.ErrNoRows{
				return fmt.Errorf("no tasks for current user")
			}
			return fmt.Errorf("error getting tasks %v",err)
		}
		fmt.Printf("%-5s %-10s %-14s %-12s %-30s %-15s %-15s %-10s \n", "No", "TaskName", "CreatedAt", "UpdatedAt", "deadline", "completed", "overDeadline", "descritpion")
		fmt.Println(strings.Repeat("-", 100)) // Add a separator line
		for i,task := range tasks{
			fmt.Printf("%-3v %-15v %-14v %-14v %-30v %-10v %-10v %v \n",
			i+1,
			task.TaskName,
			task.CreatedAt.Format(time.DateOnly),
			task.UpdatedAt.Format(time.DateOnly),
			task.Deadline.Time.Format(time.DateTime),
			task.Completed,
			task.OverDeadline,
			task.Description)
		}

	}
	if len(args) == 1 {
		tasks,err := c.db.GetUncompletedTasks(c.user.Name)
		if err != nil {
			if err == sql.ErrNoRows{
				return fmt.Errorf("no tasks for current user")
			}
			return fmt.Errorf("error getting uncompleted tasks %v",err)
		}
		fmt.Printf("%-5s %-10s %-14s %-12s %-20s %-15s %-15s %-10s \n", "No", "TaskName", "CreatedAt", "UpdatedAt", "deadline", "completed", "overDeadline", "descritpion")
		fmt.Println(strings.Repeat("-", 120)) // Add a separator line
		for i,task := range tasks{
			fmt.Printf("%-3v %-10v %-12v %-12v %-20v %-10v %-10v %v \n",
			i+1,
			task.TaskName,
			task.CreatedAt.Format(time.DateOnly),
			task.UpdatedAt.Format(time.DateOnly),
			task.Deadline.Time.Format(time.DateTime),
			task.Completed,
			task.OverDeadline,
			task.Description)
		}

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
