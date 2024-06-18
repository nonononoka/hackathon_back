package dao

import (
	"back/model"
	"firebase.google.com/go/auth"
	"github.com/oklog/ulid/v2"
	"log"
)

func PostLike(token *auth.Token, tweetID string) error {
	var tweet model.Tweet

	tx, err := db.Begin()
	if err != nil {
		log.Printf("fail: db.Begin, %v\n", err)
		return err
	}

	// insert fav into likes table
	likesID := ulid.Make().String()
	_, err = tx.Exec("INSERT INTO likes (id, tweet_id, user_id) VALUES(?, ?, ?)", likesID, tweetID, token.UID)

	if err != nil {
		log.Printf("fail: tx.Exec, %v\n", err)
		tx.Rollback()
		return err
	}

	// add tweet table's likes
	_, err = tx.Exec("UPDATE tweet SET like_count = like_count + 1 WHERE id = ?", tweetID)

	if err != nil {
		log.Printf("fail: tx.Exec, %v\n", err)
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("fail: tx.Commit, %v\n", err)
		return err
	}

	err = db.QueryRow("select id, body, posted_by, posted_at, reply_to, like_count from tweet where id = ?", tweetID).Scan(&tweet.ID, &tweet.Body, &tweet.PostedBy, &tweet.PostedAt, &tweet.ReplyTo, &tweet.LikeCount)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	err = db.QueryRow("SELECT name FROM user WHERE id = ?", tweet.PostedBy).Scan(&tweet.PostedBy)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	return nil
}
