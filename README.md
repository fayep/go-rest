# Go-Rest

Writing a type safe API client isn't hard any more.

A package for making requests to a Rest API. This package avoids making you use particular underlying clients. As long as it implements the Doer interface, you can use it with this package.

The package also allows you to define the base URL for the API, and then define methods that can be formatted with input to create a URL to be added to the base URL.

The package is eminently testable, as it uses interfaces for the Doer and BaseURLer, and the client struct is easily mockable.

The package is also easily extensible, as you can define new methods by defining new types that implement the Requester interface, and then defining new APIMethods with the new Requester type.

## Usage

Import: 

```go
import "github.com/fayep/go-rest"
```

Define a type to receive your response:

```go
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
...
```

Call the API

```go
	client := http.DefaultClient
  GetUserByID := rest.APIMethod[rest.Get, User]("https://jsonplaceholder.typicode.com/users/%d")
	user, err := GetUserByID.Do(client, 3)
```

That's it (3 will be substituted into the request)!  Add any headers/body _after_ any format string parameters.

