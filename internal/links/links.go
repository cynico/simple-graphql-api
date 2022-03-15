package links

import (
	"log"

	database "github.com/cynico/hackernews/internal/pkg/db/mysql"
	"github.com/cynico/hackernews/internal/users"
)

type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

func (link Link) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title, Address, UserID) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(link.Title, link.Address, link.User.ID)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Row Inserted!")
	return id
}

func GetAll() []Link {
	stmt, err := database.Db.Prepare("SELECT L.id, L.title, L.address, L.UserID, U.Username FROM Links L INNER JOIN Users U ON L.UserID = U.ID")
	if err != nil {
		log.Fatalf("couldn't prepare statement for fetching links from the db: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatalf("couldn't fetch links from the db: %v", err)
	}
	defer rows.Close()

	var links []Link
	var username string
	var id string
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username)
		if err != nil {
			log.Fatalf("couldn't scan row data into link struct: %v", err)
		}
		link.User = &users.User{
			ID:       id,
			Username: username,
		}
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links

}
