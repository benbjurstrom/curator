package curator

import (
	"log"
	"fmt"
	"strconv"
	"database/sql"
)

var accounts map[int64]Account

type Account struct {
	AccountId  int64 `db:"id"`
	Username string `db:"username"`
	EntityId  sql.NullString `db:"entity_id"`
}

func LoadAccounts() {
	rows, err := db.Queryx(
		"SELECT " +
			"a.id," +
			"null as entity_id " + // TODO: update query to match model
			"FROM accounts a")
	if err != nil {
		log.Fatalln(err)
	}

	i := 0
	accounts = make(map[int64]Account)
	for rows.Next() {
		account := Account{}
		err := rows.StructScan(&account)
		if err != nil {
			log.Fatalln(err)
		}
		accounts[account.AccountId] = account
		i++
	}


	fmt.Println(fmt.Sprintf("Loaded %s twitter accounts", strconv.Itoa(i)))
}
