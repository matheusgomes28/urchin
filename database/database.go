package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/matheusgomes28/urchin/common"
	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

type Database interface {
<<<<<<< HEAD
	GetPosts(int, int) ([]common.Post, error)
=======

	// Post related stuff
	GetPosts() ([]common.Post, error)
>>>>>>> 446f776 (Adding basic support for flexible cards)
	GetPost(post_id int) (common.Post, error)
	AddPost(title string, excerpt string, content string) (int, error)
	ChangePost(id int, title string, excerpt string, content string) error
	DeletePost(id int) error

	// Image related stuff
	AddImage(uuid string, name string, alt string) error

	// Card related stuff
	AddCard(uuid string, image_location string, json_data string, schema_name string) error
	GetCard(uuid int) (common.Card, error)
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
	rows, err := db.Connection.Query("SELECT title, content FROM posts WHERE id=?;", post_id)
	if err != nil {
		return common.Post{}, err
	}
	defer func() {
		err = errors.Join(rows.Close())
	}()

	rows.Next()
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
func (db SqlDatabase) DeletePost(id int) error {
	if _, err := db.Connection.Exec("DELETE FROM posts WHERE id=?;", id); err != nil {
		return err
	}

	return nil
}

// AddImage will add the image metadata to the
// database.
// name - file name saved to the disk
// alt - the alternative text
// returns (uuid, nil) if succeeded, ("", err) otherwise
func (db SqlDatabase) AddImage(uuid string, name string, alt string) (err error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if commit_err := tx.Commit(); commit_err != nil {
			err = errors.Join(err, tx.Rollback(), commit_err)
		}
	}()

	log.Info().Msgf("adding stuff to the DB")
	if name == "" {
		return fmt.Errorf("cannot have empty name")
	}

	if alt == "" {
		return fmt.Errorf("cannot have empty alt text")
	}

	query := "INSERT INTO images(uuid, name, alt) VALUES(?, ?, ?);"
	_, err = tx.Exec(query, uuid, name, alt)
	if err != nil {
		return err
	}

	return nil
}

func (db SqlDatabase) AddCard(uuid string, image_location string, json_data string, schema_name string) (err error) {
	// Check that the file exists and is a file
	// not a directory. Ideally, check the ext
	if image_location == "" {
		return fmt.Errorf("cannot have image")
	}
	image_stat, err := os.Stat(image_location)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file does not exist: %s", image_location)
	}
	if err != nil {
		return err
	}
	if image_stat.IsDir() {
		return fmt.Errorf("given path is a directory: %s", image_location)
	}

	// Load the schema
	// TODO : probably pass the schema data instead
	if json_data == "" {
		return fmt.Errorf("cannot have empty data")
	}

	if schema_name == "" {
		return fmt.Errorf("cannot have an empty schema name")
	}

	_, err = validateJson(json_data, schema_name)
	if err != nil {
		return err
	}

	// Insert everything to the DB
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if commit_err := tx.Commit(); commit_err != nil {
			err = errors.Join(err, tx.Rollback(), commit_err)
		}
	}()

	query := "INSERT INTO cards(uuid, image_location, json_data, json_schema) VALUES(?, ?, ?, ?);"
	_, err = tx.Exec(query, uuid, image_location, json_data, schema_name)
	if err != nil {
		return err
	}

	return nil
}

// / This function gets a post from the database
// / with the given ID.
func (db SqlDatabase) GetCard(uuid int) (card common.Card, err error) {
	rows, err := db.Connection.Query("SELECT image_location, json_data, json_schema FROM cards WHERE uuid=?;", uuid)
	if err != nil {
		return common.Card{}, err
	}
	defer func() {
		err = errors.Join(rows.Close())
	}()

	rows.Next()
	if err = rows.Scan(&card.ImageLocation, &card.JsonData, &card.SchemaName); err != nil {
		return common.Card{}, err
	}

	// Validate the json
	validateJson(card.JsonData, card.SchemaName)

	return card, nil
}

func validateJson(json_data string, schema_name string) (bool, error) { 
	schema_data, err := os.ReadFile(filepath.Join("schemas", fmt.Sprintf("%s.json", schema_name)))
	if err != nil {
		return false, err
	}

	schemaLoader := gojsonschema.NewBytesLoader(schema_data)
	documentLoader := gojsonschema.NewStringLoader(json_data)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, fmt.Errorf("could not read json_data: %v", err)
	}
	if !result.Valid() {
		json_err := fmt.Errorf("invalid card json: ")
		for _, e := range result.Errors() {
			json_err = fmt.Errorf("%v %s", json_err, e)
		}
		return false, json_err
	}

	return true, nil
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
