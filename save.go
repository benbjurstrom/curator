package curator

import (
	"log"
	"github.com/ChimeraCoder/anaconda"
)

type Save struct{}

type SaveInterface interface{
	saveTweet(tweet anaconda.Tweet, entity_tweets []EntityTweet) error
}

func (sv *Save) saveTweet(tweet anaconda.Tweet, entity_tweets []EntityTweet) error {
	var err error
	err = db.Ping()
	if err != nil { ConnectDB() }

	tx := db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
			if err != nil { log.Fatal(err) }
		}
		err = tx.Commit()
	}()

	// Insert our tweet
	 _, err = tx.NamedExec("INSERT INTO " +
		 "tweets " +
		 	"(id, account_id, text, created_at, updated_at) " +
		 "VALUES " +
		 	"(:id, :user.id, :text, NOW(), NOW())", tweet)
	if err != nil { log.Fatal(err) }


	// And then insert its corresponding entities
	for _,entity_tweet := range entity_tweets {
		_, err = tx.NamedExec("INSERT INTO " +
			"entity_tweets " +
				"(tweet_id, keyword_id, entity_id, created_at, updated_at) " +
			"VALUES " +
				"(:tweet_id, :keyword_id, :entity_id, NOW(), NOW())", entity_tweet)
		if err != nil { log.Fatal(err) }
	}

	return err
}