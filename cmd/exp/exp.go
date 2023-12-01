package main

import (
	stdCtx "context"
	"fmt"

	"github.com/Mauricio-3107/lenslocked-v2/context"
	"github.com/Mauricio-3107/lenslocked-v2/models"
)

func main() {
	ctx := stdCtx.Background()

	user := models.User{
		Email: "mau@mauri.com",
	}

	ctx = context.WithUser(ctx, &user)
	retrievedUser := context.User(ctx)
	fmt.Println(retrievedUser.Email)
}
