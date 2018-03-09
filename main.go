package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"fmt"
	"os"
	"io/ioutil"
	"math/rand"
)

type Character struct {
	ID string `json:"id"`
	Name string `json:"name"`
	IsAI *bool `json:"isAI"`
	Position *Vector2 `json:"position"`
	TargetPosition *Vector2 `json:"-"`
}
func (character Character) String() string {
	return fmt.Sprintf("ID: %s\nName: %s\nIsAI: %t\nPosition: %s\nTargetPosition %s\n", character.ID, character.Name, *character.IsAI, character.Position, character.TargetPosition)
}

var gCharacters []Character

func ReadDatabaseFile(aCharacters *[]Character) {
	jsonFile, err := os.Open("database.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &aCharacters)
}

func MoveCharactersGoroutine() {
	for {
		MoveCharacters()
		time.Sleep(100 * time.Millisecond)
	}
}

func MoveCharacters() {
	localCharacters := <- charactersChannel
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for index, character := range localCharacters {
		if character.TargetPosition == nil {
			if character.IsAI != nil && *character.IsAI == true {
				var newTargetPosition Vector2
				newTargetPosition.X = r.Float64() * 10
				newTargetPosition.Y = r.Float64() * 10
				localCharacters[index].TargetPosition = &newTargetPosition
	//			fmt.Printf("no target for %d\n", index)
			}
		} else {
			direction := Subtract(*character.TargetPosition, *character.Position)
			Normalize(&direction)
			var newPosition Vector2
			newPosition = Add(*character.Position, Multiply(direction, 0.5))

			*localCharacters[index].Position = newPosition

			distance := Distance(*localCharacters[index].Position, *localCharacters[index].TargetPosition)
			if distance <= 1 {
				localCharacters[index].TargetPosition = nil
			}
		}
	}

	charactersChannel <- localCharacters
}

func PrintCharactersGoroutine() {
	for {
		PrintCharacters()
		time.Sleep(1000 * time.Millisecond)
	}
}

func PrintCharacters() {
	localCharacters := <- charactersChannel
	fmt.Printf("-- current state --\n")
	for index, value := range localCharacters {
		fmt.Printf("%d)\n%s", index, value)
	}
	charactersChannel <- localCharacters
}

func GetCharacters(w http.ResponseWriter, r *http.Request) {
	localCharacters := <- charactersChannel
	json.NewEncoder(w).Encode(localCharacters)
	charactersChannel <- localCharacters
}

func GetCharacter(w http.ResponseWriter, r *http.Request) {
	localCharacters := <- charactersChannel
	params := mux.Vars(r)
	for _, item := range localCharacters {
		if item.ID == params["id"] {
		json.NewEncoder(w).Encode(item)
			fmt.Printf("-- GetCharacter -- id: %s\n", params["id"])
			charactersChannel <- localCharacters
			return
		}
    }

	fmt.Printf("-- GetCharacter -- id: %s NOT FOUND!\n", params["id"])
    json.NewEncoder(w).Encode(&Character{})
	charactersChannel <- localCharacters
}

func PutCharacter(w http.ResponseWriter, r *http.Request) {
	localCharacters := <- charactersChannel
    params := mux.Vars(r)
    for index, item := range localCharacters {
        if item.ID == params["id"] {
			var character Character
			_ = json.NewDecoder(r.Body).Decode(&character)
			fmt.Printf("-- PutCharacter --\n%s", character)
			localCharacters[index] = character
		}
	}

	charactersChannel <- localCharacters
}

func PatchCharacter(w http.ResponseWriter, r *http.Request) {
	localCharacters := <- charactersChannel
    params := mux.Vars(r)
    for index, item := range localCharacters {
        if item.ID == params["id"] {
			var character Character
			_ = json.NewDecoder(r.Body).Decode(&character)
			fmt.Printf("-- PatchCharacter --\n&s", character)

			if character.IsAI != nil {
				localCharacters[index].IsAI = character.IsAI
				if *localCharacters[index].IsAI == false {
					localCharacters[index].TargetPosition = nil
				}
				fmt.Printf("-- Updated IsAi --%t\n", character.IsAI)
			}

			if character.Position != nil {
				localCharacters[index].TargetPosition = character.Position
				fmt.Printf("-- Updated Position --%s\n", character.Position)
			}
		}
	}

	charactersChannel <- localCharacters
}

func PostCharacter(w http.ResponseWriter, r *http.Request) {
	localCharacters := <- charactersChannel
	params := mux.Vars(r)
	var character Character
	character.ID = params["id"]
	localCharacters = append(localCharacters, character)
	json.NewEncoder(w).Encode(character)
	fmt.Printf("-- PostCharacter -- id: %s\n", params["id"])
	charactersChannel <- localCharacters
}

var charactersChannel chan []Character

func main() {
	charactersChannel = make(chan []Character)

	ReadDatabaseFile(&gCharacters)
	go MoveCharactersGoroutine()
	go PrintCharactersGoroutine()
	charactersChannel <- gCharacters

    router := mux.NewRouter()
	router.HandleFunc("/characters", GetCharacters).Methods("GET")
	router.HandleFunc("/character/{id}", GetCharacter).Methods("GET")
	router.HandleFunc("/character/{id}", PutCharacter).Methods("PUT")
	router.HandleFunc("/character/{id}", PatchCharacter).Methods("PATCH")
	router.HandleFunc("/character/{id}", PostCharacter).Methods("POST")

    log.Fatal(http.ListenAndServe(":80", router))
}