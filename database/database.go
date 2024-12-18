package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/matheusgomes28/urchin/common"
	"github.com/rs/zerolog/log"
)

type Database interface {
	GetPosts(int, int) ([]common.Post, error)
	GetPost(post_id int) (common.Post, error)
	AddPost(title string, excerpt string, content string) (int, error)
	ChangePost(id int, title string, excerpt string, content string) error
	DeletePost(id int) (int, error)
	AddPage(title string, content string, link string) (int, error)
	GetPage(link string) (common.Page, error)
	AddCard(title string, image string, schema string, content string) (string, error)
	AddCardSchema(json_id string, json_schema string, json_title string, schema string) (string, error)
}

type SqlDatabase struct {
	Address    string
	Port       int
	User       string
	Connection *sql.DB
}

// / This function gets all the posts from the current
// / database connection.
func (db SqlDatabase) GetPosts(limit int, offset int) (all_posts []common.Post, err error) {
	query := "SELECT title, excerpt, id FROM posts LIMIT ? OFFSET ?;"
	rows, err := db.Connection.Query(query, limit, offset)
	if err != nil {
		return make([]common.Post, 0), err
	}
	defer func() {
		err = errors.Join(rows.Close())
	}()

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
func (db SqlDatabase) GetPost(post_id int) (post common.Post, err error) {
	rows, err := db.Connection.Query("SELECT id, title, content, excerpt FROM posts WHERE id=?;", post_id)
	if err != nil {
		return common.Post{}, err
	}
	defer func() {
		err = errors.Join(rows.Close())
	}()

	rows.Next()
	if err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.Excerpt); err != nil {
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
		if commit_err := tx.Commit(); commit_err != nil {
			err = errors.Join(err, tx.Rollback(), commit_err)
		}
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

	return nil
}

// / This function changes a post based on the values
// / provided. Note that empty strings will mean that
// / the value will not be updated.
func (db SqlDatabase) DeletePost(id int) (int, error) {
	var res, err = db.Connection.Exec("DELETE FROM posts where id=?", id)
	if err != nil {
		return 0, err
	}

	rows_affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rows_affected), nil
}

func (db SqlDatabase) AddPage(title string, content string, link string) (int, error) {
	res, err := db.Connection.Exec("INSERT INTO pages(content, title, link) VALUES(?, ?, ?)", content, title, link)
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

func (db SqlDatabase) GetPage(link string) (common.Page, error) {
	rows, err := db.Connection.Query("SELECT id, title, content, link FROM pages WHERE link=?;", link)
	if err != nil {
		return common.Page{}, err
	}
	defer func() {
		err = errors.Join(rows.Close())
	}()

	page := common.Page{}
	rows.Next()
	if err = rows.Scan(&page.Id, &page.Title, &page.Content, &page.Link); err != nil {
		return common.Page{}, err
	}

	return page, nil
}

// / This function adds the card metadata to the cards table.
// / Returns the uuid as a string if successful, otherwise error
// / won't be null
func (db SqlDatabase) AddCard(title string, image string, schema string, content string) (string, error) {

	uuid := uuid.New().String()

	_, err := db.Connection.Exec("INSERT INTO cards(uuid, image_location, json_data, json_schema) VALUES(?, ?, ?, ?)", uuid, image, content, schema)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (db SqlDatabase) AddCardSchema(json_id string, json_schema string, json_title string, schema string) (string, error) {
	uuid := uuid.New().String()

	_, err := db.Connection.Exec(
		"INSERT INTO card_schemas(uuid, json_id, json_schema, json_title, schema, card_ids) VALUES(?, ?, ?, ?)",
		uuid,
		json_id,
		json_schema,
		json_title,
		schema,
		"[]")

	if err != nil {
		return "", err
	}

	return uuid, nil
}

func MakeSqlConnection(user string, password string, address string, port int, database string) (SqlDatabase, error) {

	/// TODO : let user specify the DB
	connection_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, address, port, database)
	db, err := sql.Open("mysql", connection_str)
	if err != nil {
		return SqlDatabase{}, err
	}

	if err := db.Ping(); err != nil {
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
