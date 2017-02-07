package curator

import (
	"testing"
	"github.com/ChimeraCoder/anaconda"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
|--------------------------------------------------------------------------
| Setup mocked interfaces
|--------------------------------------------------------------------------
|
|
*/

type MockedRetweetQueue struct{
	mock.Mock
}
var _ RetweetQueueInterface = (*MockedRetweetQueue)(nil) // make sure it satisfies the interface
func (m *MockedRetweetQueue) addRetweetToQueue(tweet anaconda.Tweet) error {
	_ = m.Called()
	return nil
}
func (m *MockedRetweetQueue) saveRetweetsEveryMinute()  {}
func (m *MockedRetweetQueue) saveRetweet(Id int64, Count int)  {}


type MockedSave struct{
	mock.Mock
}
var _ SaveInterface = (*MockedSave)(nil)
func (sv *MockedSave) saveTweet(tweet anaconda.Tweet, entityTweets []EntityTweet) error {
	_ = sv.Called()
	return nil
}

/*
|--------------------------------------------------------------------------
| Begin Mocked Tests
|--------------------------------------------------------------------------
|
|
*/


func TestTagTweetCallsKeywordMatch(t *testing.T) {

	tweet, _, _ := getStructs()

	tag := &Tagger{}
	rq := new(MockedRetweetQueue)
	sv := new(MockedSave)

	sv.On("saveTweet").Once()

	tag.TagTweet(rq, sv, tweet)
	sv.AssertExpectations(t)
}

func TestTagTweetCallsRetweetQueue(t *testing.T) {

	tweet, _, _ := getStructs()
	tweet.RetweetedStatus = &anaconda.Tweet{}

	tag := &Tagger{}
	rq := new(MockedRetweetQueue)
	sv := new(MockedSave)

	rq.On("addRetweetToQueue").Return(nil)
	tag.TagTweet(rq, sv, tweet)
	rq.AssertExpectations(t)
}

func TestTagTweetDoesNotCallAnyMethodsOnNonFollowedAccount(t *testing.T) {

	tweet, _, _ := getStructs()
	tweet.User.Id = 123

	tag := &Tagger{}
	rq := new(MockedRetweetQueue)
	sv := new(MockedSave)

	tag.TagTweet(rq, sv, tweet)
	rq.AssertExpectations(t)
}

func TestTagTweetDoesNotCallAnyMethodsOnMention(t *testing.T) {

	tweet, _, _ := getStructs()
	tweet.InReplyToUserID = 123

	tag := &Tagger{}
	rq := new(MockedRetweetQueue)
	sv := new(MockedSave)

	tag.TagTweet(rq, sv, tweet)
	rq.AssertExpectations(t)
}

/*
|--------------------------------------------------------------------------
| Begin Unmocked Tests
|--------------------------------------------------------------------------
|
|
*/


func TestTagIsRetweet(t *testing.T){
	var result bool
	tweet, _, _ := getStructs()

	tag := &Tagger{}
	result = tag.isRetweet(tweet)
	assert.False(t, result)

	tweet.RetweetedStatus = &anaconda.Tweet{}
	result = tag.isRetweet(tweet)
	assert.True(t, result)
}

func TestTagIsQuotes(t *testing.T){
	var result bool
	tweet, _, _ := getStructs()

	tag := &Tagger{}
	result = tag.isQuote(tweet)
	assert.False(t, result)

	tweet.QuotedStatusID = 2531941425319414
	result = tag.isQuote(tweet)
	assert.True(t, result)
}

func TestTagIsNonFollowedAccount(t *testing.T){
	var result bool
	tweet, accounts, _ := getStructs()

	tag := &Tagger{}
	result = tag.isNonFollowedAccount(tweet, accounts)
	assert.False(t, result) // Default account struct matches

	tweet.User.Id = 2531941425319413
	result = tag.isNonFollowedAccount(tweet, accounts)
	assert.True(t, result)
}

func TestTagIsMention(t *testing.T){
	var result bool
	tweet, _, _ := getStructs()

	tag := &Tagger{}
	result = tag.isMention(tweet)
	assert.False(t, result)

	tweet.InReplyToUserID = 2531941425319414
	result = tag.isMention(tweet)
	assert.True(t, result)
}

func TestTagGlobalKeywordMatch(t *testing.T){
	var result = []EntityTweet{}
	tweet, _, keywords := getStructs()

	tag := &Tagger{}
	result = tag.globalKeywordMatch(tweet, keywords)
	assert.Equal(t, keywords[0].EntityId,result[0].EntityId)
	assert.Equal(t, keywords[0].KeywordId,result[0].KeywordId)
	assert.Equal(t, tweet.Id,result[0].TweetId)

	result = []EntityTweet{}
	tweet.Text = "Nba Some other Jamz."
	result = tag.globalKeywordMatch(tweet, keywords)
	assert.Equal(t, keywords[0].EntityId,result[0].EntityId)
	assert.Equal(t, keywords[0].KeywordId,result[0].KeywordId)
	assert.Equal(t, tweet.Id,result[0].TweetId)

	assert.Equal(t, keywords[1].EntityId,result[1].EntityId)
	assert.Equal(t, keywords[1].KeywordId,result[1].KeywordId)
	assert.Equal(t, tweet.Id,result[1].TweetId)

	result = []EntityTweet{}
	tweet.Text = "Some other text"
	result = tag.globalKeywordMatch(tweet, keywords)
	assert.Empty(t, result)
}

func getStructs() (anaconda.Tweet, map[int64]Account, []Keyword) {
	var keyword = Keyword{}
	keyword.EntityId = 1
	keyword.Keyword = "NBA"
	keyword.KeywordId = 1
	keywords = append(keywords, keyword)

	keyword.EntityId = 2
	keyword.Keyword = "jamz"
	keyword.KeywordId = 2
	keywords = append(keywords, keyword)

	accounts = make(map[int64]Account)
	var account = Account{}
	account.AccountId = 2557521 // @espn
	accounts[account.AccountId] = account

	tweet := anaconda.Tweet{}
	tweet.Text = "this is a test NBA tweet"
	tweet.User.Id = 2557521 // @espn
	tweet.Id = 826471981805662209

	return tweet, accounts, keywords
}