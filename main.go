package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"math/rand"
	"github.com/marouenebg/pokedexcli/internal/pokecache"
	"github.com/marouenebg/pokedexcli/internal/pokeExplore"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type config struct {
	Previous string
	Next     string
	Cache    *pokecache.Cache
}

type LocationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Pokemon struct {
	Name       string `json:"name"`
	Experience int    `json:"base_experience"`
	Height     int    `json:"height"`
	Weight     int    `json:"weight"`

	// Stats is an array of stats
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
		     Name   string `json:"name"`
	     } `json:"type"`
     } `json:"types"`
}

 var Pokedex  = make(map[string]Pokemon)

func main() {
	commands := make(map[string]cliCommand)
	// Create the cache with a 10-second expiration interval
	cache := pokecache.NewCache(10 * time.Second)

	mapping := &config{
		Previous: "",
		Next:     "https://pokeapi.co/api/v2/location-area",
		Cache:    cache,
	}

	commands["Usage"] = cliCommand{
		name:        "Usage",
		description: "",
		callback:    nil,
	}

	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "map",
		callback:    commandMap,
	}

	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "mapb",
		callback:    commandMapb,
	}
	
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Display Pokémon in a location",
		callback:    commandExplore,
	}

	commands["catch"] = cliCommand{
		name:        "catch",
		description: "catching a Pokémon",
		callback:    commandCatch,
	}

	commands["inspect"] = cliCommand{
                name:        "inspect",
                description: "inspecting my Pokémon",
                callback:    commandInspect,
        }

        commands["pokedex"] = cliCommand{
                name:        "pokedex",
                description: "displaying all my pokemons",
                callback:    commandPokedex,
        }
	
	fmt.Println("Welcome to the Pokedex!")
	scanner := bufio.NewScanner(os.Stdin)
	for {
//		for _, cmd := range commands {
//			fmt.Println(cmd.name + ": " + cmd.description)
		//}
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanText := cleanInput(input)
		if cmd, exists := commands[cleanText[0]]; exists {
			if len(cleanText) > 1 {
			cmd.callback(mapping,cleanText[1])
		} else {
                       cmd.callback(mapping,"" )
		}
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

func commandHelp(cfg *config, param string) error {
	fmt.Println("I'm all the help you can get")
	return nil
}

func commandExit(cfg *config, param string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func fetchFromAPI(url string, cache *pokecache.Cache) ([]byte, error) {
	// Check cache before making a request
	if data, found := cache.Get(url); found {
		fmt.Println("Cache hit!")
		return data, nil
	}

	// Make API request if not in cache
	//fmt.Println("Fetching from API...")
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("error reading the API: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Store response in cache
	cache.Add(url, body)

	return body, nil
}

func commandMap(cfg *config, param string) error {
	body, err := fetchFromAPI(cfg.Next, cfg.Cache)
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

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapb(cfg *config, param string) error {
	if cfg.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}

	body, err := fetchFromAPI(cfg.Previous, cfg.Cache)
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

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandExplore(cfg *config, param string) error {
	exploreURL := "https://pokeapi.co/api/v2/location-area/"+param
        body, err := fetchFromAPI(exploreURL, cfg.Cache)
        if err != nil {
                return err
        }
	
	pokeExplore.ExplorePokemon(body)

         
        return nil       

}


func CatchChance(baseExp int) bool {

	rand.Seed(time.Now().UnixNano())
	catchDifficulty := float64(baseExp) / 300.0 
	if catchDifficulty > 0.9 {
		catchDifficulty = 0.9 
	}
	if catchDifficulty < 0.1 {
		catchDifficulty = 0.1 
	}

	catchChance := 100.0 * (1 - catchDifficulty)
	randomRoll := rand.Float64() * 100
	return randomRoll < catchChance
}


func commandCatch(cfg *config, param string) error {
	if _,exist := Pokedex[param];exist {
	  fmt.Println("Pokemon already in the Pokedex")
	  return nil
	}
	fmt.Println("Throwing a Pokeball at "+param+"...")
	catchURL := "https://pokeapi.co/api/v2/pokemon/"+param
        body, err := fetchFromAPI(catchURL, cfg.Cache)
        if err != nil {
                return err
        }

         var data Pokemon
        err = json.Unmarshal(body, &data)
        if err != nil {
                return err
        }

	if CatchChance(data.Experience) {
		fmt.Println(param+" was caught!")
		fmt.Println("You may now inspect it with the inspect command.")
		Pokedex[param] = data
	} else {
		fmt.Println(param+" escaped!")
	}
        

        return nil
}

func commandInspect (cfg *config, param string) error {
        if data,exist := Pokedex[param];exist { 
	fmt.Println("Name:"+data.Name)
        fmt.Println("Height:",data.Height)
        fmt.Println("Weight:",data.Weight)
	fmt.Println("Stats:")
        for _,st := range data.Stats {
                fmt.Print("-"+st.Stat.Name+": ")
                fmt.Println(st.BaseStat)
       		 }
	fmt.Println("Types:")
	for _,tp := range data.Types {
	     fmt.Println("- "+tp.Type.Name)
	}
	} else {
	fmt.Println("you have not caught that pokemon")
	}
	return nil
}

func commandPokedex (cfg *config, param string) error {
	fmt.Println("Your Pokedex:")
	for _,pokemon := range Pokedex {
	 fmt.Println("- "+pokemon.Name)
	}
	return nil
}

