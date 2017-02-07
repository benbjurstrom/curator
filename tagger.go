package curator

import (
	"github.com/ChimeraCoder/anaconda"
	"fmt"
	"regexp"
)

type Tagger struct{}

func (tag *Tagger) TagTweet(rq RetweetQueueInterface, sv SaveInterface, tweet anaconda.Tweet) error {
	var err error

	// return early if this is a retweet
	if tag.isRetweet(tweet) {
		rq.addRetweetToQueue(tweet)
		return err
	}

	// return early if this is a quoted status
	if tag.isQuote(tweet) {
		rq.addRetweetToQueue(tweet)
		return err
	}

	/**
	 * Return early if this is from an account we don't follow. These tweets are included in the stream due to
	 * retweets, mentions ect. At this poit since retweets and quotes are already filtered out, this should
	 * be just mentions.
	 */
	if tag.isNonFollowedAccount(tweet, accounts) {
		return err
	}

	/**
	 * Return early if this is a mention. Since accounts we don't follow have already been filtered out these
	 * should just be mentions contained in tweets written by our followed accounts.
	 */
	if tag.isMention(tweet) {
		return err
	}

	// Check global keywords
	entity_tweets := tag.globalKeywordMatch(tweet, keywords)
	if len(entity_tweets) == 0 {
		return err
	}

	sv.saveTweet(tweet, entity_tweets)

	// Filter out any entities that also match the tweet for a negative word.
	//entity_tweets := negativeKeywordMatch(tweet, keywords, entity_tweets)


	//TODO: check for entity keywords and append to entity_tweets

	return nil
}

// Check if this is a retweet. IF so send to the retweetQueue
func (tag *Tagger) isRetweet(tweet anaconda.Tweet) bool {

	if tweet.RetweetedStatus != nil {
		// TODO: look for RT as string
		return true
	}

	return false
}

// Check if this quotes another status. If so send to the retweetQueue
func (tag *Tagger) isQuote(tweet anaconda.Tweet) bool {
	if tweet.QuotedStatusID != 0 {
		return true
	}

	return false
}

//
/**
 * Check whether this was tweeted by an account we don't follow. If so the tweet will be dropped.
 */
func (tag *Tagger) isNonFollowedAccount(tweet anaconda.Tweet, accounts map[int64]Account) bool {

	if _, ok := accounts[tweet.User.Id]; ok {
		return false
	}

	return true
}

func (tag *Tagger) isMention(tweet anaconda.Tweet) bool {

	if tweet.InReplyToUserID != 0 {
		//SaveTweetToJsonFile(tweet, "mention_in_reply_to_user")
		return true
	}

	if tweet.InReplyToStatusID != 0 {
		//SaveTweetToJsonFile(tweet, "mention_in_reply_to_status")
		return true
	}

	// TODO: check for @
	// TODO: check whether mention is an entity. For example @kingjames

	return false
}

func (tag *Tagger) globalKeywordMatch(tweet anaconda.Tweet, keywords []Keyword) []EntityTweet {
	var entity_tweets []EntityTweet
	entity_tweets = nil

	for _,keyword := range keywords {
		var pattern = fmt.Sprintf(`(?i)(^|[\W\d])%s([\W\d]|$)`, keyword.Keyword)
		match, _ := regexp.MatchString(pattern, tweet.Text)
		if match {

			// create a new entity with the tweet and keyword information
			entity := EntityTweet{}
			entity.TweetId = tweet.Id
			entity.KeywordId = keyword.KeywordId
			entity.EntityId = keyword.EntityId

			// and append it to the entities array
			entity_tweets = append(entity_tweets, entity)

			fmt.Println(tweet.Text)
		}
	}

	return entity_tweets
}






