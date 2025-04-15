package main

type command struct {
	name        string
	usage       string
	description string
	callback    func(c *config, args ...string) error
}

func getcommands() map[string]command {
	return map[string]command{
		"help": {
			name:        "help",
			usage:       "todo help",
			description: "Lists all the existing commands their usage and description",
			callback:    commandHelp,
		},
		"login": {
			name:        "login",
			usage:       "login username",
			description: "logs in the specified user with credentials",
			callback:    commandLogin,
		},
		"logout": {
			name:        "logout",
			usage:       "logout",
			description: "ends the session of current user",
			callback:    commandLogout,
		},
		"reg": {
			name:        "register",
			usage:       "reg [-a] (username)",
			description: "Registers the new user and logs the user in automatically, you can use -a tag to register as admin which is optional",
			callback:    commandRegister,
		},
		"exit": {
			name:        "exit",
			usage:       "exit",
			description: "exits the app",
			callback:    commandExit,
		},
		"del": {
			name:        "delete",
			usage:       "del (username)",
			description: "deletes specified user(yourself) data, only admins can delete other users",
			callback:    commandDelete,
		},
		"add": {
			name:        "add",
			usage:       "add taskname [timelimit]",
			description: "adds a new task to the user user can specify timelimit it defaults with no timelimit",
			callback:    commandAdd,
		},
		"rem": {
			name:        "remove",
			usage:       "rem [-a] [taskname]/[tasks]",
			description: "removes specified task from user's list or all tasks with -a tag(provide tasks at the end with -a tag instead of taskname",
			callback:    commandRemove,
		},
		"mod": {
			name:        "modify",
			usage:       "work in progress",
			description: "modifies specified task or user details",
			callback:    nil,
		},
		"lst": {
			name:        "list tasks",
			usage:       "[-u] [n] tasks",
			description: "lists all the tasks for the user if tag is not provided,\n		-u to list all uncompleted tasks and -c list all completed tasks and\n		provide optional [n(integer)] to list n tasks",
			callback:    nil,
		},
		"ls": {
			name:        "list users",
			usage:       "ls",
			description: "lists all the users available only works if you are admin",
			callback:    commandListUsers,
		},
	}

}
