package dao

import (
	"back/model"
	"database/sql"
	"firebase.google.com/go/auth"
	"github.com/oklog/ulid/v2"
	"log"
)

func formatTweet(rows *sql.Rows, userID string) (model.Tweet, error) {
	var t model.Tweet
	var threadID string
	if err := rows.Scan(&t.ID, &t.Body, &t.PostedBy, &t.PostedAt, &t.ReplyTo, &t.LikeCount, &t.ReplyCount, &threadID); err != nil {
		log.Printf(err.Error())
		return t, err
	}

	tagRows, err := db.Query("SELECT tag.tag "+
		"FROM tweet  INNER JOIN tweet_tag  ON tweet.id = tweet_tag.tweet_id "+
		"INNER JOIN tag  ON tweet_tag.tag_id = tag.id where tweet.id = ?;", t.ID)

	if err != nil {
		log.Printf(err.Error())
		return t, err
	}
	for tagRows.Next() {
		var tag string
		tagRows.Scan(&tag)
		t.Tags = append(t.Tags, tag)
	}
	if err := tagRows.Close(); err != nil {
		log.Printf(err.Error())
		return t, err
	}
	err = db.QueryRow("SELECT name FROM user WHERE id = ?", t.PostedBy).Scan(&t.PostedByName)
	err = db.QueryRow("SELECT image FROM user WHERE id = ?", t.PostedBy).Scan(&t.PostedByImage)
	if err != nil {
		log.Printf(err.Error())
		return t, err
	}
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM favorite WHERE user_id = ? AND tweet_id = ?)", userID, &t.ID).Scan(&t.IsFaved)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(t)
	return t, nil
}

func returnTweets(rows *sql.Rows, userID string) ([]model.Tweet, error) {
	tweets := make([]model.Tweet, 0)
	for rows.Next() {
		t, err := formatTweet(rows, userID)
		if err != nil {
			return tweets, err
		}
		tweets = append(tweets, t)
	}
	return tweets, nil
}

// tagを含むtweetをgetする。
func GetTweets(token *auth.Token, tags []string, tweetID string) ([]model.Tweet, error) {
	tweets := make([]model.Tweet, 0)
	// idがtweetIDのtweetを取得する。
	if tweetID != "" {
		rows, err := db.Query("select * from tweet where id = ?;", tweetID)

		if err != nil {
			log.Printf(err.Error())
			return tweets, err
		}
		tweets, err = returnTweets(rows, token.UID)
		if err != nil {
			return tweets, err
		}
		if err := rows.Close(); err != nil {
			log.Printf(err.Error())
			return tweets, err
		}
		return tweets, nil
	}

	if len(tags) == 0 {
		rows, err := db.Query("select * from tweet;")

		if err != nil {
			log.Printf(err.Error())
			return tweets, err
		}
		tweets, err = returnTweets(rows, token.UID)
		if err != nil {
			return tweets, err
		}
		if err := rows.Close(); err != nil {
			log.Printf(err.Error())
			return tweets, err
		}
	} else {
		for _, tag := range tags {
			// tagを含むツイートたち
			rows, err := db.Query("SELECT t.* from tweet t "+
				"INNER JOIN tweet_tag tt ON t.id = tt.tweet_id "+
				"INNER JOIN tag tg ON tt.tag_id = tg.id "+
				"WHERE tg.tag = ?;", tag)

			if err != nil {
				log.Printf(err.Error())
				return tweets, err
			}

			taggedTweets, err := returnTweets(rows, token.UID)
			if err != nil {
				return tweets, err
			}
			tweets = append(tweets, taggedTweets...)

			if err := rows.Close(); err != nil {
				return tweets, err
			}
			log.Println(tweets)
		}
	}
	return tweets, nil
}

// tagをつけてツイートをpostする。
func PostTweet(token *auth.Token, body string, tags []string) (model.Tweet, error) {
	var tweet model.Tweet

	tx, err := db.Begin()
	if err != nil {
		log.Printf("fail: db.Begin, %v\n", err)
		return tweet, err
	}

	// insert tweet into tweet table
	tweetID := ulid.Make().String()
	result, err := tx.Exec("INSERT INTO tweet (id, body, posted_by, thread_id) VALUES(?, ?, ?, ?)", tweetID, body, token.UID, tweetID)
	log.Println(result)
	if err != nil {
		log.Printf("fail: tx.Exec, %v\n", err)
		tx.Rollback()
		return tweet, err
	}

	// insert tag into tag table and tag&tweet pair into tweet_tag table.
	for _, tag := range tags {
		var tagId string
		err := db.QueryRow("SELECT id FROM tag WHERE tag = ?", tag).Scan(&tagId)
		// if tag does not exist
		if err != nil {
			tagId = ulid.Make().String()
			_, err = tx.Exec("INSERT INTO tag (id, tag) VALUES (?, ?)", tagId, tag)
			if err != nil {
				log.Printf("タグ %s の挿入エラー: %v", tag, err)
				tx.Rollback()
				return tweet, err
			}
		}

		pairID := ulid.Make().String()
		_, err = tx.Exec("INSERT INTO tweet_tag (id, tweet_id, tag_id) VALUES (?, ?, ?)", pairID, tweetID, tagId)
		if err != nil {
			log.Printf(err.Error())
			return tweet, err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("fail: tx.Commit, %v\n", err)
		return tweet, err
	}

	err = db.QueryRow("select id, body, posted_by, posted_at, reply_to, like_count, reply_count from tweet where id = ?", tweetID).Scan(&tweet.ID, &tweet.Body, &tweet.PostedBy, &tweet.PostedAt, &tweet.ReplyTo, &tweet.LikeCount, &tweet.ReplyCount)
	if err != nil {
		log.Printf(err.Error())
		return tweet, err
	}

	err = db.QueryRow("SELECT name FROM user WHERE id = ?", tweet.PostedBy).Scan(&tweet.PostedByName)
	err = db.QueryRow("SELECT image FROM user WHERE id = ?", tweet.PostedBy).Scan(&tweet.PostedByImage)
	if err != nil {
		log.Printf(err.Error())
		return tweet, err
	}
	tweet.Tags = tags
	tweet.IsFaved = false

	return tweet, nil
}

func PostReply(token *auth.Token, body string, tags []string, repliedTweetID string) (model.Tweet, error) {
	var tweet model.Tweet

	tx, err := db.Begin()
	if err != nil {
		log.Printf("fail: db.Begin, %v\n", err)
		return tweet, err
	}

	var threadID string
	err = db.QueryRow("SELECT thread_id FROM tweet WHERE id = ?", repliedTweetID).Scan(&threadID)
	if err != nil {
		return tweet, err
	}

	_, err = tx.Exec("UPDATE tweet SET reply_count = reply_count + 1 WHERE id = ?", repliedTweetID)
	// insert tweet into tweet table
	tweetId := ulid.Make().String()
	log.Println(tweetId, body, token.UID)
	result, err := tx.Exec("INSERT INTO tweet (id, body, posted_by, reply_to, thread_id) VALUES(?, ?, ?, ?, ?)", tweetId, body, token.UID, repliedTweetID, threadID)
	log.Println(result)
	if err != nil {
		log.Printf("fail: tx.Exec, %v\n", err)
		tx.Rollback()
		return tweet, err
	}
	
	for _, tag := range tags {
		var tagId string
		err := db.QueryRow("SELECT id FROM tag WHERE tag = ?", tag).Scan(&tagId)
		// if tag does not exist
		if err != nil {
			tagId = ulid.Make().String()
			_, err = tx.Exec("INSERT INTO tag (id, tag) VALUES (?, ?)", tagId, tag)
			if err != nil {
				log.Printf("タグ %s の挿入エラー: %v", tag, err)
				tx.Rollback()
				return tweet, err
			}
		}

		pairID := ulid.Make().String()
		_, err = tx.Exec("INSERT INTO tweet_tag (id, tweet_id, tag_id) VALUES (?, ?, ?)", pairID, tweetId, tagId)
		if err != nil {
			log.Printf(err.Error())
			return tweet, err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("fail: tx.Commit, %v\n", err)
		return tweet, err
	}

	err = db.QueryRow("select id, body, posted_by, posted_at, reply_to, like_count, reply_count from tweet where id = ?", tweetId).Scan(&tweet.ID, &tweet.Body, &tweet.PostedBy, &tweet.PostedAt, &tweet.ReplyTo, &tweet.LikeCount, &tweet.ReplyCount)
	if err != nil {
		log.Printf(err.Error())
		return tweet, err
	}

	err = db.QueryRow("SELECT name FROM user WHERE id = ?", tweet.PostedBy).Scan(&tweet.PostedByName)
	err = db.QueryRow("SELECT image FROM user WHERE id = ?", tweet.PostedBy).Scan(&tweet.PostedByImage)
	if err != nil {
		log.Printf(err.Error())
		return tweet, err
	}
	tweet.Tags = []string{}
	tweet.IsFaved = false

	return tweet, nil
}

// tagを含む特定userのツイートを取得する。
func GetUserTweets(token *auth.Token, userID string, tags []string) ([]model.Tweet, error) {
	tweets := make([]model.Tweet, 0)

	if len(tags) == 0 {
		rows, err := db.Query("select * from tweet where posted_by = ?;", userID)

		if err != nil {
			log.Printf(err.Error())
			return tweets, err
		}
		tweets, err = returnTweets(rows, token.UID)
		if err != nil {
			return tweets, err
		}
		if err := rows.Close(); err != nil {
			log.Printf(err.Error())
			return tweets, err
		}
	} else {
		for _, tag := range tags {
			// tagを含むツイートたち
			rows, err := db.Query("SELECT t.* from tweet t "+
				"INNER JOIN tweet_tag tt ON t.id = tt.tweet_id "+
				"INNER JOIN tag tg ON tt.tag_id = tg.id "+
				"WHERE tg.tag = ?;", tag)

			if err != nil {
				log.Printf(err.Error())
				return tweets, err
			}

			taggedTweets, err := returnTweets(rows, token.UID)
			if err != nil {
				return tweets, err
			}
			tweets = append(tweets, taggedTweets...)
			if err := rows.Close(); err != nil {
				return tweets, err
			}
			log.Println(tweets)
		}
	}
	return tweets, nil
}

func GetFollowingUserTweets(token *auth.Token, tags []string) ([]model.Tweet, error) {
	tweets := make([]model.Tweet, 0)

	if len(tags) == 0 {
		rows, err := db.Query("select t.* from tweet t join follow f on t.posted_by = f.followee_id where f.follower_id = ?;", token.UID)

		if err != nil {
			log.Printf("282", err.Error())
			return tweets, err
		}
		tweets, err = returnTweets(rows, token.UID)
		if err != nil {
			return tweets, err
		}
		if err := rows.Close(); err != nil {
			log.Printf("318", err.Error())
			return tweets, err
		}
	} else {
		for _, tag := range tags {
			// tagを含むツイートたち
			rows, err := db.Query("SELECT t.* from tweet t "+
				"INNER JOIN tweet_tag tt ON t.id = tt.tweet_id "+
				"INNER JOIN tag tg ON tt.tag_id = tg.id "+
				"WHERE tg.tag = ?;", tag)

			if err != nil {
				log.Printf(err.Error())
				return tweets, err
			}

			taggedTweets, err := returnTweets(rows, token.UID)
			if err != nil {
				return tweets, err
			}
			tweets = append(tweets, taggedTweets...)
			if err := rows.Close(); err != nil {
				return tweets, err
			}
			log.Println(tweets)
		}
	}
	return tweets, nil
}

// 直後のリプライと上にたどって行った時の親を全部返す。
func GetThreadTweets(token *auth.Token, tweetID string) ([]model.Tweet, error) {
	log.Printf("getThreadTweets")
	tweets := make([]model.Tweet, 0)

	// 直後のリプライ
	rows, err := db.Query("select * from tweet where reply_to = ?;", tweetID)
	if err != nil {
		log.Printf(err.Error())
		return tweets, err
	}
	tweets, err = returnTweets(rows, token.UID)

	currentID := tweetID
	for {
		rows, err := db.Query("SELECT * FROM tweet WHERE id = ?", currentID)
		if err != nil {
			return tweets, err
		}
		defer rows.Close()

		if rows.Next() {
			tweet, err := formatTweet(rows, token.UID)
			if err != nil {
				return tweets, err
			}
			tweets = append(tweets, tweet)
			var replyTo sql.NullString
			db.QueryRow("SELECT reply_to FROM tweet WHERE id = ?", currentID).Scan(&replyTo)
			if replyTo.Valid {
				currentID = tweet.ReplyTo.String
			} else {
				break // If reply_to is NULL, we've reached the root tweet
			}
		}
	}
	return tweets, nil
}
