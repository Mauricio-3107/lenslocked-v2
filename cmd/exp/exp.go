package main

import (
	"fmt"

	"github.com/Mauricio-3107/lenslocked-v2/models"
)

func main() {
	gs := models.GalleryService{}
	imgs, err := gs.Images(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(imgs)
	fmt.Println(len(imgs))
}
