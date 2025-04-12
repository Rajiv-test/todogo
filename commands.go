package main


type command struct{
	name string 
	usage string 
	description string
	callback func(c *config,args... string)error
}

func getcommands()map[string]command{
	return map[string]command{
		"help":{
			name: "help",
			usage: "todo help",
			description: "Lists all the existing commands their usage and description",
			callback: commandHelp,
		},
		"login":{
			name: "login",
			usage: "todo login username",
			description: "logs in the specified user with credentials",
			callback: nil,
		},
		"reg":{
			name: "register",
			usage: "todo reg (username)",
			description: "Registers the new user",
			callback: nil,
		},
		"exit":{
			name: "exit",
			usage: "todo exit",
			description: "exits the currect session",
			callback: commandExit,
		},
		"del":{
			name: "delete",
			usage: "todo delete",
			description: "deletes current user data",
			callback: nil,
		},
		"add":{
			name: "add",
			usage: "todo add taskname [timelimit]",
			description: "adds a new task to the user user can specify timelimit it defaults with no timelimit",
			callback: nil,
		},
		"rem":{
			name: "remove",
			usage: "todo rem [-a] [taskname]",
			description: "removes specified task from user's list or all tasks with -a tag ignores taskname even if it's provided",
			callback: nil,
		},
		"mod":{
			name: "modify",
			usage: "work in progress",
			description: "modifies specified task or user details",
			callback: nil,
		},
		"tasks":{
			name: "tasks",
			usage: "todo [-u][-c] [n] tasks",
			description: "lists all the tasks for the user if tag is not provided,\n		-u to list all uncompleted tasks and -c list all completed tasks and\n		provide optional [n(integer)] to list n tasks",
			callback: nil,
		},
	}

}
