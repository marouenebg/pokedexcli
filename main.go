package main

import (
	"bufio"
	"fmt"
 	"strings"
	"os"
)

type cliCommand struct {
   name          string
   description   string
   callback      func() error
}
//var commands map[string]cliCommand 
func main() {
commands := make (map[string]cliCommand)

  commands["Usage"] = cliCommand {
    name:             "Usage",
    description:      "",
    callback:         nil,
  }
  
  commands["exit"] = cliCommand {
     name:            "exit",
     description:     "Exit the Pokedex",
     callback:        commandExit,
  }
  commands["help"] = cliCommand {
     name:            "help",
     description:     "Displays a help message",
     callback:        commandHelp,
  }

  commands["map"] = cliCommand {
     name:            "map",
     description:     "map",
     callback:        commandMap,
  }

fmt.Println("Welcome to the Pokedex!")
 scanner := bufio.NewScanner(os.Stdin)
 for {
	 for _,cmd := range commands {
		 fmt.Println(cmd.name+": "+cmd.description )
	 }
	 fmt.Print("Pokedex > ")
	 scanner.Scan()	
	 input := scanner.Text()
	 if cmd, exists := commands[input]; exists {
	   cmd.callback()
	 } else {
		fmt.Println("Unknown command")
	 }
 }

}
func cleanInput(text string) []string {
  cleanOP := []string{}
  cleanOP = strings.Fields(text) 
  for i:=0;i<len(cleanOP);i++ {
 	cleanOP[i] = strings.ToLower(cleanOP[i])
  }

  return cleanOP

}

func commandHelp() error {
  fmt.Println("I'm all the help you can get")
  return nil
}

func commandExit() error {
  fmt.Println("Closing the Pokedex... Goodbye!")  
  os.Exit(0)
  return nil
}

func commandMap() error {
 return nil
}
