package database

import (
	"database/sql"
	"errors"

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

	AddPermalink(permalink common.Permalink) (int, error)
	GetPermalinks() ([]common.Permalink, error)
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

func (db SqlDatabase) AddPermalink(permalink common.Permalink) (int, error) {
	res, err := db.Connection.Exec("INSERT INTO post_permalinks(permalink, post_id) VALUES(?, ?)", permalink.Path, permalink.PostId)
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

func (db SqlDatabase) GetPermalinks() ([]common.Permalink, error) {
	rows, err := db.Connection.Query("SELECT permalink, post_id FROM post_permalinks")
	if err != nil {
		return []common.Permalink{}, err
	}
	defer func() {
		err = errors.Join(rows.Close())
	}()

	permalinks := []common.Permalink{}
	for rows.Next() {
		var permalink common.Permalink

		if err = rows.Scan(&permalink.Path, &permalink.PostId); err != nil {
			log.Error().Msgf("could not get permalink from db: %v", err)
			return []common.Permalink{}, err
		}

		permalinks = append(permalinks, permalink)
	}

	return permalinks, nil
}
