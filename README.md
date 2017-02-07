# Curator

[![Build Status](https://api.travis-ci.org/benbjurstrom/curator.png)](https://travis-ci.org/benbjurstrom/curator)

A go application that follows twitter accounts and scans their tweets for keywords.

## What does it do?
Curator loads a list of twitter accounts from a postgres database and passes them to a [twitter streaming api](https://dev.twitter.com/streaming/overview/request-parameters#follow) connection. Curator then scans each tweet for the presence of a list of keywords. Tweets containing a keyword are saved to the database and tagged with the _keyword_id_ and the keyword's corresponding  _entity_id_. Curator also scores saved tweets by tracking the number of retweets and quotes they receive.

For more details on the data structure see the [curator model](https://github.com/benbjurstrom/curator-model) package.

## Installation

1. Install the curator package by running 
  ```bash 
    go get github.com/benbjurstrom/curator
  ```

2. Setup your main.go file as follows
```Go
  package main
  
  import (
  	"github.com/benbjurstrom/curator"
  )
  
  func main() {
  	curator.Run()
  }
```
3. Load the [curator model](https://github.com/benbjurstrom/curator-model) data structure into a postgres database and add your database connection information to the following environmental variable:
````bash
export PG_DSN="user=postgres password=postgres dbname=postgres host=localhost sslmode=disable"
```` 

4. Use Twitter's [application management console](https://apps.twitter.com) to create a new twitter app and copy its credentials to the following environmental variables:
  ```bash
  export CONSUMER_KEY=your_consumer_key
  export CONSUMER_SECRET=your_consumer_secret
  export ACCESS_TOKEN=your_access_token
  export ACCESS_TOKEN_SECRET=your_access_token_secret
  ``` 

5. Build an executable with 
  ```bash
  go build main.go
  ```
  
6. Then to begin curating tweets simply run the executable.
  ```bash
    ./main
  ```
  
  if everything is working you should see something similar to this.
  
  ```bash
    Loaded 100 twitter accounts
    Loaded 6 keywords
    Following 100 twitter writers
  ```

## License
MIT