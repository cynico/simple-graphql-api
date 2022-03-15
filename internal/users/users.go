package users

import (
	"database/sql"
	"log"

	database "github.com/cynico/hackernews/internal/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

/*
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDczNjQ1MzksInVzZXJuYW1lIjoieXVzdWZzIn0.5U0zwNLfxdx1zFxUSwICfzg6mok17xZBhKOm8QQPkEA
*/

func (user *User) Create() error {
	statement, err := database.Db.Prepare("INSERT INTO Users(Username, Password) VALUES(?,?)")
	if err != nil {
		return err
	}
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	_, err = statement.Exec(user.Username, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

//GetUserIdByUsername check if a user exists in database by given username
func (user *User) Authenticate() bool {
	statement, err := database.Db.Prepare("SELECT Password FROM Users WHERE Username = ?")
	if err != nil {
		log.Fatalf("error preparing statement for getting user's password: %v", err)
	}

	row := statement.QueryRow(user.Username)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return CheckPasswordHash(user.Password, hashedPassword)
}

//GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(username string) (int, error) {
	statement, err := database.Db.Prepare("select ID from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}

//GetUserByID check if a user exists in database and return the user object.
func GetUserByID(userId string) (User, error) {
	statement, err := database.Db.Prepare("select Username from Users WHERE ID = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(userId)

	var username string
	err = row.Scan(&username)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return User{}, err
	}

	return User{ID: userId, Username: username}, nil
}

//HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
