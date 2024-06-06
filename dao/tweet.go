package dao

import (
	"back/model"
	"firebase.google.com/go/auth"
	"github.com/oklog/ulid/v2"
	"log"
)

func GetMyTweets(token *auth.Token) ([]model.Tweet, error) {
	tweets := make([]model.Tweet, 0)
	rows, err := db.Query("select tweets.id, tweets.content , tweets.posted_by from tweets "+
		"inner join users on tweets.posted_by = users.id where users.id = ?", token.UID)

	if err != nil {
		log.Printf(err.Error())
		return tweets, err
	}
	for rows.Next() {
		var t model.Tweet
		if err := rows.Scan(&t.ID, &t.Body, &t.PostedBy); err != nil {
			log.Printf("fail: rows.Scan @GetMe, %v\n", err)

			if err := rows.Close(); err != nil {
				return tweets, err
			}
		}
		tweets = append(tweets, t)
	}
	log.Println(tweets)
	return tweets, err
}

func PostTweet(token *auth.Token, content string) (model.Tweet, error) {
	var tweet model.Tweet
	id := ulid.Make().String()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("fail: db.Begin, %v\n", err)
		return tweet, err
	}
	_, err = tx.Exec("INSERT INTO tweets (id, content, posted_by) VALUES(?, ?, ?)", id, content, token.UID)
	if err != nil {
		log.Printf("fail: tx.Exec, %v\n", err)
		tx.Rollback()
		return tweet, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("fail: tx.Commit, %v\n", err)
		return tweet, err
	}
	tweet.ID = token.UID
	tweet.PostedBy = token.UID
	tweet.Body = content
	return tweet, nil
}
