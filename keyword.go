package curator

import (
	"log"
	"fmt"
	"strconv"
)

var keywords []Keyword

type Keyword struct {
	KeywordId  int `db:"keyword_id"`
	Keyword  string `db:"keyword"`
	KeywordTypeId  int `db:"keyword_type_id"`
	KeywordTypeName  string `db:"keyword_type_name"`
	EntityId  int `db:"entity_id"`
	EntityName  string `db:"entity_name"`
}

func LoadKeywords() {
	rows, err := db.Queryx(
		"SELECT " +
			"k.id as keyword_id," +
			"k.keyword, " +
			"k.keyword_type_id, " +
			"k.entity_id, " +
			"kt.name as keyword_type_name, " +
			"e.name as entity_name " +
			"FROM entity_keywords k " +
			"JOIN keyword_types kt ON k.keyword_type_id = kt.id " +
			"JOIN entities e ON k.entity_id = e.id")
	if err != nil {
		log.Fatalln(err)
	}

	i := 0
	for rows.Next() {
		keyword := Keyword{}
		err := rows.StructScan(&keyword)
		if err != nil {
			log.Fatalln(err)
		}
		keywords = append(keywords, keyword)
		i++
	}

	fmt.Println(fmt.Sprintf("Loaded %s keywords", strconv.Itoa(i)))
}
