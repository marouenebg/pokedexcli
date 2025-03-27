package pokeExplore

import ( "encoding/json"
         "fmt"
 )


type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VersionDetail struct {
	Rate    int     `json:"rate"`
	Version Version `json:"version"`
}

type EncounterMethodRate struct {
	EncounterMethod EncounterMethod `json:"encounter_method"`
	VersionDetails  []VersionDetail `json:"version_details"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Name struct {
	Language Language `json:"language"`
	Name     string   `json:"name"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Method struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EncounterDetail struct {
	Chance          int      `json:"chance"`
	ConditionValues []string `json:"condition_values"` // Assuming it's an array of strings
	MaxLevel        int      `json:"max_level"`
	Method          Method   `json:"method"`
	MinLevel        int      `json:"min_level"`
}

type PokemonVersionDetail struct {
	EncounterDetails []EncounterDetail `json:"encounter_details"`
	MaxChance        int               `json:"max_chance"`
	Version          Version           `json:"version"`
}

type PokemonEncounter struct {
	Pokemon        Pokemon               `json:"pokemon"`
	VersionDetails []PokemonVersionDetail `json:"version_details"`
}

type ResponseData struct {
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	GameIndex            int                   `json:"game_index"`
	ID                   int                   `json:"id"`
	Location             Location              `json:"location"`
	Name                 string                `json:"name"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
} 

func ExplorePokemon(body []byte) {
 var response ResponseData
        err := json.Unmarshal(body, &response)
        if err != nil {
               fmt.Println("error while reading names") 
        }

        for _, result := range response.PokemonEncounters {
                fmt.Println(result.Pokemon.Name)
        }

}

