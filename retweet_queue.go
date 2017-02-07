package curator

import (
	"time"
	"github.com/boltdb/bolt"
	"strconv"
	"log"
	"github.com/ChimeraCoder/anaconda"
)

type RetweetQueue struct{}

type RetweetQueueInterface interface{
	addRetweetToQueue(tweet anaconda.Tweet) error
	saveRetweetsEveryMinute()
	saveRetweet(Id int64, Count int)
}

func (rq *RetweetQueue) addRetweetToQueue(tweet anaconda.Tweet) error {
	var err error

	boltDB.Update(func(tx *bolt.Tx) error {

		// First check if the tweet already exists
		b := tx.Bucket([]byte("Retweets"))
		v := b.Get([]byte(tweet.IdStr)) // use IdStr to to simplify byte conversion

		if v == nil {

			// If not create it with a value of 1
			err := b.Put([]byte(tweet.IdStr), []byte("1"))
			if err != nil { log.Fatal(err) }

		} else {

			// otherwise increment the existing count
			currentValue, err := strconv.Atoi(string(v))
			if err != nil { log.Fatal(err) }

			newValue := strconv.Itoa(currentValue + 1)

			err = b.Put([]byte(tweet.IdStr), []byte(newValue))
			if err != nil { log.Fatal(err) }

		}

		return err
	})

	return err
}

func (rq *RetweetQueue) saveRetweetsEveryMinute() {
	// Save the retweets every minute.
	t := time.NewTicker(time.Minute)
	for {

		// Loop over the entire Retweet bucket.
		boltDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Retweets"))

			b.ForEach(func(k, v []byte) error {

				value, err := strconv.Atoi(string(v))
				if err != nil {
					value = 0
				}

				key, err := strconv.ParseInt(string(k), 10, 64)
				if err != nil {
					log.Fatal(err)
				}

				rq.saveRetweet(key, value)

				// Finally delete the key
				b.Delete(k)
				return nil
			})

			return nil
		})

		<-t.C // blocks the for loop until a minute has passed
	}
}

func (rq *RetweetQueue) saveRetweet(Id int64, Count int) {
	stmt, err := db.Prepare("UPDATE tweets SET retweets = retweets + $1 WHERE id = $2")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(Count, Id)
	if err != nil {
		log.Fatal(err)
	}
}
