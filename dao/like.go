package dao

import (
	"back/model"
	"firebase.google.com/go/auth"
	"github.com/oklog/ulid/v2"
	"log"
)

// ユーザーがいいねしてるツイート全部持ってくる
func GetLike(token *auth.Token) ([]model.Tweet, error) {
	tweets := make([]model.Tweet, 0)
	// idがtweetIDのtweetを取得する。
	rows, err := db.Query("SELECT tweet.* FROM tweet JOIN likes ON tweet.id = likes.tweet_id WHERE likes.user_id = ?;", token.UID)

	if err != nil {
		log.Printf(err.Error())
		return tweets, err
	}

	tweets, err = returnTweets(rows, token.UID)
	if err := rows.Close(); err != nil {
		log.Printf(err.Error())
		return tweets, err
	}
	log.Println(tweets)
	return tweets, nil
}

func PostLike(token *auth.Token, tweetID string) error {
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

	return nil
}

func DeleteLike(token *auth.Token, tweetID string) error {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("fail: db.Begin, %v\n", err)
		return err
	}

	_, err = tx.Exec("DELETE FROM likes WHERE user_id = ? AND tweet_id = ?", token.UID, tweetID)

	if err != nil {
		log.Printf("fail: tx.Exec, %v\n", err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE tweet SET like_count = like_count - 1 WHERE id = ?", tweetID)

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

	return nil
}
