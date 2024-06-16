package dao

import (
	"back/model"
	"firebase.google.com/go/auth"
	"github.com/oklog/ulid/v2"
	"log"
)

// tagを含むtweetをgetする。
func GetTweets(tags []string) ([]model.Tweet, error) {
	tweets := make([]model.Tweet, 0)

	if len(tags) == 0 {
		rows, err := db.Query("select * from tweet;")

		if err != nil {
			log.Printf(err.Error())
			return tweets, err
		}
		for rows.Next() {
			var t model.Tweet
			if err := rows.Scan(&t.ID, &t.Body, &t.PostedBy, &t.PostedAt, &t.ReplyTo, &t.LikeCount); err != nil {
				log.Printf(err.Error())
				return tweets, err
			}

			tagRows, err := db.Query("SELECT tag.tag "+
				"FROM tweet  INNER JOIN tweet_tag  ON tweet.id = tweet_tag.tweet_id "+
				"INNER JOIN tag  ON tweet_tag.tag_id = tag.id where tweet.id = ?;", t.ID)

			if err != nil {
				log.Printf(err.Error())
				return tweets, err
			}
			for tagRows.Next() {
				var tag string
				tagRows.Scan(&tag)
				t.Tags = append(t.Tags, tag)
			}
			if err := tagRows.Close(); err != nil {
				log.Printf(err.Error())
				return tweets, err
			}

			err = db.QueryRow("SELECT name FROM user WHERE id = ?", t.PostedBy).Scan(&t.PostedBy)
			if err != nil {
				log.Printf(err.Error())
				return tweets, err
			}
			tweets = append(tweets, t)
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

			for rows.Next() {
				var t model.Tweet
				if err := rows.Scan(&t.ID, &t.Body, &t.PostedBy, &t.PostedAt, &t.ReplyTo, &t.LikeCount); err != nil {
					log.Printf("fail: rows.Scan @GetMe, %v\n", err)
				}

				// 各ツイートに含まれる全部のタグを取得する。
				tagRows, err := db.Query("SELECT tg.tag "+
					"FROM tweet t INNER JOIN tweet_tag tt ON t.id = tt.tweet_id "+
					"INNER JOIN tag tg ON tt.tag_id = tg.id where t.id = ?;", t.ID)

				if err != nil {
					log.Printf(err.Error())
					return tweets, err
				}
				for tagRows.Next() {
					var tag string
					tagRows.Scan(&tag)
					t.Tags = append(t.Tags, tag)
				}
				if err := tagRows.Close(); err != nil {
					return tweets, err
				}
				err = db.QueryRow("SELECT name FROM user WHERE id = ?", t.PostedBy).Scan(&t.PostedBy)
				if err != nil {
					log.Printf(err.Error())
					return tweets, err
				}
				tweets = append(tweets, t)
			}
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
	tweetId := ulid.Make().String()
	log.Println(tweetId, body, token.UID)
	result, err := tx.Exec("INSERT INTO tweet (id, body, posted_by) VALUES(?, ?, ?)", tweetId, body, token.UID)
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

		pairId := ulid.Make().String()
		_, err = tx.Exec("INSERT INTO tweet_tag (id, tweet_id, tag_id) VALUES (?, ?, ?)", pairId, tweetId, tagId)
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

	err = db.QueryRow("select id, body, posted_by, posted_at, reply_to, like_count from tweet where id = ?", tweetId).Scan(&tweet.ID, &tweet.Body, &tweet.PostedBy, &tweet.PostedAt, &tweet.ReplyTo, &tweet.LikeCount)
	if err != nil {
		log.Printf(err.Error())
		return tweet, err
	}

	err = db.QueryRow("SELECT name FROM user WHERE id = ?", tweet.PostedBy).Scan(&tweet.PostedBy)
	if err != nil {
		log.Printf(err.Error())
		return tweet, err
	}
	tweet.Tags = tags

	return tweet, nil
}

// tagをつけてツイートをpostする。
func PostReply(token *auth.Token, body string, repliedTweetID string) (model.Tweet, error) {
	var tweet model.Tweet

	tx, err := db.Begin()
	if err != nil {
		log.Printf("fail: db.Begin, %v\n", err)
		return tweet, err
	}

	// insert tweet into tweet table
	tweetId := ulid.Make().String()
	log.Println(tweetId, body, token.UID)
	result, err := tx.Exec("INSERT INTO tweet (id, body, posted_by, reply_to) VALUES(?, ?, ?, ?)", tweetId, body, token.UID, repliedTweetID)
	log.Println(result)
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

	err = db.QueryRow("select id, body, posted_by, posted_at, reply_to, like_count from tweet where id = ?", tweetId).Scan(&tweet.ID, &tweet.Body, &tweet.PostedBy, &tweet.PostedAt, &tweet.ReplyTo, &tweet.LikeCount)
	if err != nil {
		log.Printf(err.Error())
		return tweet, err
	}

	err = db.QueryRow("SELECT name FROM user WHERE id = ?", tweet.PostedBy).Scan(&tweet.PostedBy)
	if err != nil {
		log.Printf(err.Error())
		return tweet, err
	}
	tweet.Tags = []string{}

	return tweet, nil
}

// tagを含む特定userのツイートを取得する。
func GetUserTweets(userID string, tags []string) ([]model.Tweet, error) {
	tweets := make([]model.Tweet, 0)

	if len(tags) == 0 {
		rows, err := db.Query("select * from tweet where posted_by = ?;", userID)

		if err != nil {
			log.Printf(err.Error())
			return tweets, err
		}
		for rows.Next() {
			var t model.Tweet
			if err := rows.Scan(&t.ID, &t.Body, &t.PostedBy, &t.PostedAt, &t.ReplyTo, &t.LikeCount); err != nil {
				log.Printf(err.Error())
				return tweets, err
			}

			tagRows, err := db.Query("SELECT tag.tag "+
				"FROM tweet  INNER JOIN tweet_tag  ON tweet.id = tweet_tag.tweet_id "+
				"INNER JOIN tag  ON tweet_tag.tag_id = tag.id where tweet.id = ?;", t.ID)

			if err != nil {
				log.Printf(err.Error())
				return tweets, err
			}
			for tagRows.Next() {
				var tag string
				tagRows.Scan(&tag)
				t.Tags = append(t.Tags, tag)
			}
			if err := tagRows.Close(); err != nil {
				log.Printf(err.Error())
				return tweets, err
			}

			err = db.QueryRow("SELECT name FROM user WHERE id = ?", t.PostedBy).Scan(&t.PostedBy)
			if err != nil {
				log.Printf(err.Error())
				return tweets, err
			}
			tweets = append(tweets, t)
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

			for rows.Next() {
				var t model.Tweet
				if err := rows.Scan(&t.ID, &t.Body, &t.PostedBy, &t.PostedAt, &t.ReplyTo, &t.LikeCount); err != nil {
					log.Printf("fail: rows.Scan @GetMe, %v\n", err)
				}

				// 各ツイートに含まれる全部のタグを取得する。
				tagRows, err := db.Query("SELECT tg.tag "+
					"FROM tweet t INNER JOIN tweet_tag tt ON t.id = tt.tweet_id "+
					"INNER JOIN tag tg ON tt.tag_id = tg.id where t.id = ?;", t.ID)

				if err != nil {
					log.Printf(err.Error())
					return tweets, err
				}
				for tagRows.Next() {
					var tag string
					tagRows.Scan(&tag)
					t.Tags = append(t.Tags, tag)
				}
				if err := tagRows.Close(); err != nil {
					return tweets, err
				}
				err = db.QueryRow("SELECT name FROM user WHERE id = ?", t.PostedBy).Scan(&t.PostedBy)
				if err != nil {
					log.Printf(err.Error())
					return tweets, err
				}
				tweets = append(tweets, t)
			}
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
		for rows.Next() {
			var t model.Tweet
			if err := rows.Scan(&t.ID, &t.Body, &t.PostedBy, &t.PostedAt, &t.ReplyTo, &t.LikeCount); err != nil {
				log.Printf("288", err.Error())
				return tweets, err
			}

			tagRows, err := db.Query("SELECT tag.tag "+
				"FROM tweet  INNER JOIN tweet_tag  ON tweet.id = tweet_tag.tweet_id "+
				"INNER JOIN tag  ON tweet_tag.tag_id = tag.id where tweet.id = ?;", t.ID)

			if err != nil {
				log.Printf("297", err.Error())
				return tweets, err
			}
			for tagRows.Next() {
				var tag string
				tagRows.Scan(&tag)
				t.Tags = append(t.Tags, tag)
			}
			if err := tagRows.Close(); err != nil {
				log.Printf("306", err.Error())
				return tweets, err
			}

			err = db.QueryRow("SELECT name FROM user WHERE id = ?", t.PostedBy).Scan(&t.PostedBy)
			if err != nil {
				log.Printf("312", err.Error())
				return tweets, err
			}
			tweets = append(tweets, t)
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

			for rows.Next() {
				var t model.Tweet
				if err := rows.Scan(&t.ID, &t.Body, &t.PostedBy, &t.PostedAt, &t.ReplyTo, &t.LikeCount); err != nil {
					log.Printf("fail: rows.Scan @GetMe, %v\n", err)
				}

				// 各ツイートに含まれる全部のタグを取得する。
				tagRows, err := db.Query("SELECT tg.tag "+
					"FROM tweet t INNER JOIN tweet_tag tt ON t.id = tt.tweet_id "+
					"INNER JOIN tag tg ON tt.tag_id = tg.id where t.id = ?;", t.ID)

				if err != nil {
					log.Printf(err.Error())
					return tweets, err
				}
				for tagRows.Next() {
					var tag string
					tagRows.Scan(&tag)
					t.Tags = append(t.Tags, tag)
				}
				if err := tagRows.Close(); err != nil {
					return tweets, err
				}
				err = db.QueryRow("SELECT name FROM user WHERE id = ?", t.PostedBy).Scan(&t.PostedBy)
				if err != nil {
					log.Printf(err.Error())
					return tweets, err
				}
				tweets = append(tweets, t)
			}
			if err := rows.Close(); err != nil {
				return tweets, err
			}
			log.Println(tweets)
		}
	}
	return tweets, nil
}
