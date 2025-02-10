// This example will demonstrate using https://jsonplaceholder.typicode.com/users/%d api.
/* The call returns:
{
  "id": 3,
  "name": "Clementine Bauch",
  "username": "Samantha",
  "email": "Nathan@yesenia.net",
  "address": {
    "street": "Douglas Extension",
    "suite": "Suite 847",
    "city": "McKenziehaven",
    "zipcode": "59590-4157",
    "geo": {
      "lat": "-68.6102",
      "lng": "-47.0653"
    }
  },
  "phone": "1-463-123-4447",
  "website": "ramiro.info",
  "company": {
    "name": "Romaguera-Jacobson",
    "catchPhrase": "Face to face bifurcated interface",
    "bs": "e-enable strategic applications"
  }
}
As an example.*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fayep/go-rest"
)

// User is a struct that represents a user.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo"`
	} `json:"address"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase"`
		BS          string `json:"bs"`
	} `json:"company"`
}

// GetUser returns a user from the API.
const GetUserByID = rest.APIMethod[rest.Get, User]("https://jsonplaceholder.typicode.com/users/%d")

func main() {
	c := http.DefaultClient
	u, err := GetUserByID.Do(c, 3)
	if err != nil {
		panic(err)
	}
	bytes, _ := json.MarshalIndent(u, "", "  ")
	fmt.Printf("User:\n%s", string(bytes))
}
