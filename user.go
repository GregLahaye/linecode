package main

import "fmt"

type User struct {
	LeetCodeSession string
	CSRFToken       string
}

func main() {
	u := User{}
	if err := Login(&u); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Sucessfully logged in", u)
	}
}
