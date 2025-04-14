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

func commandAdd(c *config,args ...string) error{
	if len(args) > 2 || len(args) < 1{
		return fmt.Errorf("wrong usage of add command use help command to learn more")
	}
	if c.user == nil {
		return fmt.Errorf("login to create tasks")
	}
	taskname := args[0]
	var deadline sql.NullTime
	var err error
	if len(args) == 2{
		deadline,err = parseDeadline(args[1])
		if err != nil {
			return fmt.Errorf("enter valid timelimit %v",err)
	}
	}else{
		deadline = sql.NullTime{Time:time.Time{},Valid:false}
	}
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter task description: ")
	reader.Scan()
	description := reader.Text()
	err = c.db.AddTask(c.user.Name,taskname,description,deadline)
	if err != nil{
		return fmt.Errorf("error adding task %v",err)
	}
	fmt.Println("task created successfully")
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
        return sql.NullTime{Time:time.Now().Local().Add(time.Hour * 24 * time.Duration(days)),Valid: true}, nil
    }
    
    // For hours, minutes, seconds, Go's time.ParseDuration works directly
    duration, err := time.ParseDuration(durationStr)
    if err != nil {
        return sql.NullTime{}, err
    }
    
    return sql.NullTime{Time:time.Now().Local().Add(duration),Valid :true}, nil
}