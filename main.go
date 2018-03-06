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
	Position Vector2 `json:"position"`
	TargetPosition Vector2 `json:"-"`
	TargetPositionPtr *Vector2 `json:"-"`
}
func (character Character) String() string {
    return fmt.Sprintf("ID: %s\nName: %s\nPosition: %s\nTargetPosition %s\n", character.ID, character.Name, character.Position, character.TargetPosition)
}

var characters []Character

func ReadDatabaseFile() {
	jsonFile, err := os.Open("database.json")
	if err != nil {
		fmt.Println(err)
	}	
	
	fmt.Println("Successfully Opened users.json")	
	byteValue, _ := ioutil.ReadAll(jsonFile)	
	json.Unmarshal(byteValue, &characters)
	PrintCharacters()	
	fmt.Printf("-- current state --\n")		
}

func MoveCharactersGoroutine() {
	for {
		MoveCharacters()
		PrintCharacters()
		time.Sleep(1000 * time.Millisecond)	
	}	
}

func MoveCharacters() {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for index, character := range characters {
		if character.TargetPositionPtr == nil {			
			var newTargetPosition Vector2
			newTargetPosition.X = r.Float64() * 10
			newTargetPosition.Y = r.Float64() * 10
			characters[index].TargetPosition = newTargetPosition		
			characters[index].TargetPositionPtr = &characters[index].TargetPosition		
//			fmt.Printf("no target for %d\n", index)
		} else {
			direction := Subtract(character.TargetPosition, character.Position)
			fmt.Printf("XXXXXXXXXXXXXXXXXX dir %s\n", direction)
			Normalize(&direction)			
			fmt.Printf("XXXXXXXXXXXXXXXXXX normalized %s\n", direction)
			var newPosition Vector2
			newPosition = Add(character.Position, Multiply(direction, 0.5))	
			fmt.Printf("XXXXXXXXXXXXXXXXXX new pos %s\n", newPosition)
			characters[index].Position = newPosition
			
			distance := Distance(characters[index].Position, characters[index].TargetPosition)
			if distance <= 1 {
				characters[index].TargetPositionPtr = nil
			}			
		}		
	}
}

func PrintCharactersGoroutine() {
	for {
		PrintCharacters()
		time.Sleep(1000 * time.Millisecond)	
	}
}

func PrintCharacters() {
	fmt.Printf("-- current state --\n")
	for index, value := range characters {
		fmt.Printf("%d)\n%s", index, value)
	}
}

func GetCharacters(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(characters)
}


func GetCharacter(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range characters {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
			fmt.Printf("-- GetCharacter -- id: %s\n", params["id"])
            return
        }
    }
	
	fmt.Printf("-- GetCharacter -- id: %s NOT FOUND!\n", params["id"])
    json.NewEncoder(w).Encode(&Character{})
}

func PutCharacter(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, item := range characters {
        if item.ID == params["id"] {		
			var character Character
			_ = json.NewDecoder(r.Body).Decode(&character)		
			fmt.Printf("-- PutCharacter --\n%s", character)
			
			characters[index] = character
        }
    }
	
	PrintCharacters()
}

func PostCharacter(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var character Character
	character.ID = params["id"]
	characters = append(characters, character)
	json.NewEncoder(w).Encode(character)
	
	fmt.Printf("-- PostCharacter -- id: %s\n", params["id"])
	PrintCharacters()
}

// main function to boot up everything
func main() {
	ReadDatabaseFile()	

    router := mux.NewRouter()	
	router.HandleFunc("/characters", GetCharacters).Methods("GET")
	router.HandleFunc("/character/{id}", GetCharacter).Methods("GET")
	router.HandleFunc("/character/{id}", PutCharacter).Methods("PUT")
	router.HandleFunc("/character/{id}", PostCharacter).Methods("POST")
	go MoveCharactersGoroutine()	
    log.Fatal(http.ListenAndServe(":80", router))
}