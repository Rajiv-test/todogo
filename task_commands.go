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
