package dao

import (
	"back/model"
	"firebase.google.com/go/auth"
	"github.com/oklog/ulid/v2"
	"log"
)

func FollowUser(token *auth.Token, userID string) (model.User, error) {
	var followingUser model.User

	tx, err := db.Begin()
	if err != nil {
		log.Printf("fail: db.Begin, %v\n", err)
		return followingUser, err
	}

	// insert friendship into follow table
	friendshipId := ulid.Make().String()

	_, err = tx.Exec("INSERT INTO follow (id, follower_id, followee_id) VALUES(?, ?, ?)", friendshipId, token.UID, userID)

	if err != nil {
		log.Printf("fail: tx.Exec, %v\n", err)
		tx.Rollback()
		return followingUser, err
	}

	err = db.QueryRow("SELECT id, name, email, bio, image from user where id = ?", userID).Scan(&followingUser.ID, &followingUser.Name, &followingUser.Email, &followingUser.Bio, &followingUser.Image)

	if err != nil {
		tx.Rollback()
		log.Printf(err.Error())
		return followingUser, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("fail: tx.Commit, %v\n", err)
		return followingUser, err
	}

	return followingUser, nil
}

func UnfollowUser(token *auth.Token, userID string) error {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("fail: db.Begin, %v\n", err)
		return err
	}

	_, err = tx.Exec("delete from follow where follower_id = ? AND followee_id = ?", token.UID, userID)

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

func GetFollowingUsers(token *auth.Token) ([]model.User, error) {
	var followingUsers []model.User

	rows, err := db.Query("SELECT u.id, u.name, u.email, u.bio, u.image "+
		"FROM user u JOIN follow f ON u.id = f.followee_id "+
		"where f.follower_id = ?;", token.UID)

	if err != nil {
		log.Printf("fail: db.Query, %v\n", err)
		return followingUsers, err
	}

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Bio, &u.Image); err != nil {
			log.Printf(err.Error())
			return followingUsers, err
		}
		following, followed, err := checkFollow(token.UID, u.ID)
		if err != nil {
			return followingUsers, err
		}
		// uにフォローされているか
		u.IsFollowed = followed
		// uをフォローしているか
		u.IsFollowing = following
		followingUsers = append(followingUsers, u)
	}
	if err := rows.Close(); err != nil {
		log.Printf(err.Error())
		return followingUsers, err
	}
	return followingUsers, nil
}

// 自分をフォローしてるユーザーたち
func GetFollowedUsers(token *auth.Token) ([]model.User, error) {
	var followedUsers []model.User

	rows, err := db.Query("SELECT u.id, u.name, u.email, u.bio, u.image "+
		"FROM user u JOIN follow f ON u.id = f.follower_id "+
		"where f.followee_id = ?;", token.UID)

	if err != nil {
		log.Printf("fail: db.Query, %v\n", err)
		return followedUsers, err
	}

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Bio, &u.Image); err != nil {
			log.Printf(err.Error())
			return followedUsers, err
		}
		following, followed, err := checkFollow(token.UID, u.ID)
		if err != nil {
			return followedUsers, err
		}
		// uにフォローされているか
		u.IsFollowed = followed
		// uをフォローしているか
		u.IsFollowing = following
		followedUsers = append(followedUsers, u)
	}
	if err := rows.Close(); err != nil {
		log.Printf(err.Error())
		return followedUsers, err
	}
	return followedUsers, nil
}

// ID1がID2をフォローしている
// ID2がID1をフォローしている
func checkFollow(userID1 string, userID2 string) (bool, bool, error) {
	query := "SELECT COUNT(*) FROM follow WHERE follower_id = ? AND followee_id = ?"
	var user1FollowUser2 bool = false
	var user2FollowUser1 bool = false
	var count int
	// user_id1がuser_id2をフォローしているかを確認するクエリ
	err := db.QueryRow(query, userID1, userID2).Scan(&count)
	if err != nil {
		return false, false, err
		log.Fatal(err)
	}
	if count > 0 {
		user1FollowUser2 = true
	}

	// user_id2がuser_id1をフォローしているかを確認するクエリ
	query = "SELECT COUNT(*) FROM follow WHERE follower_id = ? AND followee_id = ?"
	err = db.QueryRow(query, userID2, userID1).Scan(&count)
	if err != nil {
		return false, false, err
		log.Fatal(err)
	}

	if count > 0 {
		user2FollowUser1 = true
	}
	return user1FollowUser2, user2FollowUser1, nil
}
