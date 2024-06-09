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
		followingUsers = append(followingUsers, u)
	}
	if err := rows.Close(); err != nil {
		log.Printf(err.Error())
		return followingUsers, err
	}
	return followingUsers, nil
}
