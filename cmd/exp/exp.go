package main

import (
	"html/template"
	"os"
)

type User struct {
	Name  string
	Bio   string
	Roots map[string]string
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name:  "Mauricio Ramirez",
		Bio:   `<script>alert("Haha, you have been h4x0r3d!");</script>`,
		Roots: make(map[string]string),
	}

	user.Roots["Country"] = "Bolivia"
	user.Roots["City"] = "Santa Cruz"

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
