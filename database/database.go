package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/matheusgomes28/urchin/common"
	"github.com/rs/zerolog/log"
)

type Database interface {
	GetPosts() ([]common.Post, error)
	GetPost(post_id int) (common.Post, error)
	AddPost(title string, excerpt string, content string) (int, error)
	ChangePost(id int, title string, excerpt string, content string) error
	DeletePost(id int) error
}

type SqlDatabase struct {
	Address    string
	Port       int
	User       string
	Connection *sql.DB
}

// / This function gets all the posts from the current
// / database connection.
func (db SqlDatabase) GetPosts() ([]common.Post, error) {
	rows, err := db.Connection.Query("SELECT title, excerpt, id FROM posts;")
	if err != nil {
		return make([]common.Post, 0), err
	}
	defer rows.Close()

	all_posts := make([]common.Post, 0)
	for rows.Next() {
		var post common.Post
		if err = rows.Scan(&post.Title, &post.Excerpt, &post.Id); err != nil {
			return make([]common.Post, 0), err
		}
		all_posts = append(all_posts, post)
	}

	return all_posts, nil
}

// / This function gets a post from the database
// / with the given ID.
func (db SqlDatabase) GetPost(post_id int) (common.Post, error) {
	rows, err := db.Connection.Query("SELECT title, content FROM posts WHERE id=?;", post_id)
	if err != nil {
		return common.Post{}, err
	}
	defer rows.Close()

	rows.Next()
	var post common.Post
	if err = rows.Scan(&post.Title, &post.Content); err != nil {
		return common.Post{}, err
	}

	return post, nil
}

// / This function adds a post to the database
func (db SqlDatabase) AddPost(title string, excerpt string, content string) (int, error) {
	res, err := db.Connection.Exec("INSERT INTO posts(content, title, excerpt) VALUES(?, ?, ?)", content, title, excerpt)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Warn().Msgf("could not get last ID: %v", err)
		return -1, nil
	}

	// TODO : possibly unsafe int conv,
	// make sure all IDs are i64 in the
	// future
	return int(id), nil
}

// / This function changes a post based on the values
// / provided. Note that empty strings will mean that
// / the value will not be updated.
func (db SqlDatabase) ChangePost(id int, title string, excerpt string, content string) error {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(tx.Rollback())
	}()

	if len(title) > 0 {
		_, err := tx.Exec("UPDATE posts SET title = ? WHERE id = ?;", title, id)
		if err != nil {
			return err
		}
	}

	if len(excerpt) > 0 {
		_, err := tx.Exec("UPDATE posts SET excerpt = ? WHERE id = ?;", excerpt, id)
		if err != nil {
			return err
		}
	}

	if len(content) > 0 {
		_, err := tx.Exec("UPDATE posts SET content = ? WHERE id = ?;", content, id)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// / This function changes a post based on the values
// / provided. Note that empty strings will mean that
// / the value will not be updated.
func (db SqlDatabase) DeletePost(id int) error {
	if _, err := db.Connection.Exec("DELETE FROM posts WHERE id=?;", id); err != nil {
		return err
	}

	return nil
}

func MakeSqlConnection(user string, password string, address string, port int, database string) (SqlDatabase, error) {
	/// Checking the DB connection
	/// TODO : let user specify the DB
	connection_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, address, port, database)
	db, err := sql.Open("mysql", connection_str)
	if err != nil {
		return SqlDatabase{}, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return SqlDatabase{
		Address:    address,
		Port:       port,
		User:       user,
		Connection: db,
	}, nil
}
