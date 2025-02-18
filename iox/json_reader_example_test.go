package iox

import (
	"fmt"
	"net/http"
)

func ExampleNewJSONReader() {
	val := NewJSONReader(User{Name: "Tom"})
	_, err := http.NewRequest(http.MethodPost, "/hello", val)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("OK")
}

type User struct {
	Name string `json:"name"`
}
