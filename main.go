package main

import (
	"bufio"
	"fmt"
 	"strings"
	"os"
	"net/http"
	"io"
	"encoding/json"
)

type cliCommand struct {
   name          string
   description   string
   callback      func( *config) error
}

type config struct { 
   Previous      string
   Next          string
}


type LocationResponse struct {
	Count    int `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}


//var commands map[string]cliCommand 
func main() {
commands := make (map[string]cliCommand)

mapping := &config {
   Previous:    "",
   Next:        "https://pokeapi.co/api/v2/location-area",
}

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

  commands["mapb"] = cliCommand {
     name:            "mapb",
     description:     "mapb",
     callback:        commandMapb,
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
	   cmd.callback(mapping)
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

func commandHelp(cfg *config) error {
  fmt.Println("I'm all the help you can get")
  return nil
}

func commandExit(cdg *config) error {
  fmt.Println("Closing the Pokedex... Goodbye!")  
  os.Exit(0)
  return nil
}

func commandMap(cfg * config) error {
  res,err := http.Get(cfg.Next)
  if err != nil {
    return err
  }
  body,err := io.ReadAll(res.Body)
  res.Body.Close()
  if res.StatusCode > 299 {
     fmt.Println("error reading the API %d",res.StatusCode)
     return err
 }
 if err != nil {
   return err	
 }

 var data LocationResponse
 err = json.Unmarshal(body, &data)
 if err != nil {
    return err
 }
cfg.Previous = data.Previous 
cfg.Next = data.Next

// fmt.Printf("Count: %d\n", data.Count)
//fmt.Printf("Next URL: %s\n", data.Next)
//fmt.Println("Locations:")
	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	return nil
}



func commandMapb(cfg * config) error {
if cfg.Previous == "" {
   fmt.Println("you're on the first page")
   return nil
}
  res,err := http.Get(cfg.Previous)
  if err != nil {
    return err
  }
  body,err := io.ReadAll(res.Body)
  res.Body.Close()
  if res.StatusCode > 299 {
     fmt.Println("error reading the API %d",res.StatusCode)
     return err
 }
 if err != nil {
   return err
 }

 var data LocationResponse
 err = json.Unmarshal(body, &data)
 if err != nil {
    return err
 }
cfg.Previous = data.Previous
cfg.Next = data.Next

// fmt.Printf("Count: %d\n", data.Count)
//fmt.Printf("Next URL: %s\n", data.Next)
//fmt.Println("Locations:")
        for _, result := range data.Results {
                fmt.Println(result.Name)
        }

        return nil
}

