package main

import (
	"fmt"

	"github.com/Mauricio-3107/lenslocked-v2/models"
)

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	us := models.UserService{
		DB: db,
	}

	user, err := us.Authenticate("luca@luca.com", "qwerty")
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	// Creating table
	// _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
	// 	id SERIAL PRIMARY KEY,
	// 	name TEXT,
	// 	username TEXT NOT NULL
	// );

	// CREATE TABLE IF NOT EXISTS tweets (
	// 	id SERIAL PRIMARY KEY,
	// 	user_id INT NOT NULL,
	// 	description TEXT,
	// 	retweets INT,
	// 	liked BOOLEAN
	// );`)
	// if err != nil {
	// 	panic(err)
	// }

	// Inser data
	// name := "Mauricio Ramirez"
	// username := "rambito3107"
	// row := db.QueryRow(`
	//   INSERT INTO users (name, username)
	//   VALUES ($1, $2) RETURNING id;`, name, username)
	// var id int
	// row.Scan(&id)
	// fmt.Println("User created,", id)

	//Querying data
	// id := 2
	// row := db.QueryRow(`
	// SELECT name, username
	// FROM users
	// WHERE id = $1;`, id)
	// var name, username string
	// err = row.Scan(&name, &username)
	// if err == sql.ErrNoRows {
	// 	fmt.Println("Error, no rows!")
	// }
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("User info: name > %s, username > %s", name, username)

	// //Creating faking orders
	// userID := 2
	// for i := 1; i <= 5; i++ {
	// 	desc := fmt.Sprintf("tweet #%d", i)
	// 	retweets := i * 1
	// 	liked := true

	// 	_, err := db.Exec(`
	// 	INSERT INTO tweets(user_id, description, retweets, liked)
	// 	VALUES ($1, $2, $3, $4)`, userID, desc, retweets, liked)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("Fake tweets created")
	// }

	// Querying multiple rows
	// userID := 2
	// type Tweet struct {
	// 	ID          int
	// 	UserID      int
	// 	Liked       bool
	// 	Retweets    int
	// 	Description string
	// }
	// var tweets []Tweet

	// rows, err := db.Query(`
	// SELECT id, description, retweets, liked
	// FROM tweets
	// WHERE user_id=$1`, userID)
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var tweet Tweet
	// 	tweet.UserID = userID
	// 	err := rows.Scan(&tweet.ID, &tweet.Description, &tweet.Retweets, &tweet.Liked)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	tweets = append(tweets, tweet)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Mauri's tweets")
	// for _, tweet := range tweets {
	// 	fmt.Println(tweet)
	// }
}
