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

type Person struct {
    ID        string   `json:"id,omitempty"`
    Firstname string   `json:"firstname,omitempty"`
}


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

type Characters struct {
	Characters []Character `json:"characters"`
}

var people []Person
var characters Characters

// Display all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(people)
}

// Display a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range people {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    json.NewEncoder(w).Encode(&Person{})
}

// create a new item
func CreatePerson(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var person Person
    _ = json.NewDecoder(r.Body).Decode(&person)
    person.ID = params["id"]
    people = append(people, person)
    json.NewEncoder(w).Encode(people)
}

// Delete an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, item := range people {
        if item.ID == params["id"] {
            people = append(people[:index], people[index+1:]...)
            break
        }
        json.NewEncoder(w).Encode(people)
    }
}

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

	for index, character := range characters.Characters {
		if character.TargetPositionPtr == nil {			
			var newTargetPosition Vector2
			newTargetPosition.X = r.Float64()		
			characters.Characters[index].TargetPosition = newTargetPosition		
			characters.Characters[index].TargetPositionPtr = &characters.Characters[index].TargetPosition		
//			fmt.Printf("no target for %d\n", index)
		} else {
			direction := Subtract(character.TargetPosition, character.Position)
			fmt.Printf("XXXXXXXXXXXXXXXXXX dir %s\n", direction)
			Normalize(&direction)			
			fmt.Printf("XXXXXXXXXXXXXXXXXX normalized %s\n", direction)
			var newPosition Vector2
			newPosition = Add(character.Position, Multiply(direction, 0.5))	
			fmt.Printf("XXXXXXXXXXXXXXXXXX new pos %s\n", newPosition)
			characters.Characters[index].Position = newPosition
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
	for index, value := range characters.Characters {
		fmt.Printf("%d)\n%s", index, value)
	}
}

func GetCharacters(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(characters)
}


func GetCharacter(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range characters.Characters {
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
    for index, item := range characters.Characters {
        if item.ID == params["id"] {		
			var character Character
			_ = json.NewDecoder(r.Body).Decode(&character)		
			fmt.Printf("-- PutCharacter --\n%s", character)
			
			characters.Characters[index] = character
        }
    }
	
	PrintCharacters()
}

func PostCharacter(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var character Character
	character.ID = params["id"]
	characters.Characters = append(characters.Characters, character)
	json.NewEncoder(w).Encode(character)
	
	fmt.Printf("-- PostCharacter -- id: %s\n", params["id"])
	PrintCharacters()
}


// main function to boot up everything
func main() {

//	c := make(chan int)

    router := mux.NewRouter()
    people = append(people, Person{ID: "1", Firstname: "John"})
    people = append(people, Person{ID: "2", Firstname: "Koko"})
	
	
	router.HandleFunc("/characters", GetCharacters).Methods("GET")
	router.HandleFunc("/character/{id}", GetCharacter).Methods("GET")
	router.HandleFunc("/character/{id}", PutCharacter).Methods("PUT")
	router.HandleFunc("/character/{id}", PostCharacter).Methods("POST")
	
	
    router.HandleFunc("/people", GetPeople).Methods("GET")
    router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
    router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
    router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")

	ReadDatabaseFile()
	
//	go PrintCharactersGoroutine()
	go MoveCharactersGoroutine()

	
    log.Fatal(http.ListenAndServe(":80", router))
}