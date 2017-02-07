package curator

type EntityTweet struct {
	TweetId  int64 `db:"tweet_id"`
	EntityId  int `db:"entity_id"`
	KeywordId  int `db:"keyword_id"`
}



