package curator

import (
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	"os"
	"log"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"github.com/boltdb/bolt"
	"strconv"
	"strings"
)

var db *sqlx.DB
var boltDB *bolt.DB

func Run() {
	ConnectDB()
	ConnectBolt()

	LoadAccounts() // sets the accounts global variable
	LoadKeywords() // sets the keywords global variable

	rq  := &RetweetQueue{}
	tag := &Tagger{}
	sv  := &Save{}

	go rq.saveRetweetsEveryMinute()

	s := openTwitterStream(getAccountIdsAsString())

	for {
		item := <-s.C
		switch status := item.(type){
		case anaconda.Tweet:
			go tag.TagTweet(rq, sv, status)
		default:
			fmt.Sprintf("%T", status)
		}
	}

	defer db.Close()
	defer boltDB.Close()
}

func openTwitterStream(idsString string) anaconda.Stream {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	client := anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))

	v := url.Values{}
	v.Set("follow", idsString)
	return *client.PublicStreamFilter(v)
}


func getAccountIdsAsString() string {
	i := 0
	var ids []string
	for _,account := range accounts {
		ids = append(ids, strconv.FormatInt(account.AccountId, 10))
		i++
	}

	// Implode the twitter ids array into a comma separated string
	concatenated := fmt.Sprintf("Following %s twitter writers", strconv.Itoa(i))
	fmt.Println(concatenated)
	return strings.Join(ids, ",")
}

func ConnectDB() {
	var err error

	db, err = sqlx.Connect("postgres", os.Getenv("PG_DSN"))
	if err != nil { log.Fatalln(err) }

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}

	db.SetMaxIdleConns(100)
}

// Bolt connection is used for caching retweets to reduce the number of writes on the main database.
func ConnectBolt() {
	var err error
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	boltDB, err = bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Retweets"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

}
