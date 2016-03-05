package auth

import (
	"fmt"
	"net"
	"net/url"

	"database/sql"
	_ "github.com/lib/pq"
)

type (
	postgresql string
)

//
func init() {
	authenticators["postgres"] = newPostgres
}

//
func newPostgres(url *url.URL) error {

	//
	host, port, err := net.SplitHostPort(url.Host)
	db := url.Query().Get("db")
	user := url.User.Username()
	// pass, _ := url.User.Password()

	//
	pg := postgresql(fmt.Sprintf("user=%v database=%v sslmode=disable host=%v port=%v", user, db, host, port))

	// create the tables needed to support mist authentication
	if _, err := pg.exec(`
CREATE TABLE IF NOT EXISTS tokens (
	token text NOT NULL,
	token_id SERIAL UNIQUE NOT NULL,
	PRIMARY KEY (token)
)`); err != nil {
		return err
	}

	if _, err = pg.exec(`
CREATE TABLE IF NOT EXISTS tags (
  token_id integer NOT NULL REFERENCES tokens (token_id) ON DELETE CASCADE,
  tag text NOT NULL,
  PRIMARY KEY (token_id, tag)
)`); err != nil {
		return err
	}

	//
	DefaultAuth = pg

	return nil
}

//
func (p postgresql) AddToken(token string) error {
	_, err := p.exec("INSERT INTO tokens (token) VALUES ($1)", token)
	return err
}

//
func (p postgresql) RemoveToken(token string) error {
	_, err := p.exec("DELETE FROM tokens WHERE token = $1", token)
	return err
}

//
func (p postgresql) AddTags(token string, tags []string) error {

	// This could be optimized a LOT
	for _, tag := range tags {
		// errors are ignored, this may not be the best idea.
		p.exec("INSERT INTO tags (token_id,tag) VALUES ((SELECT token_id FROM tokens WHERE token = $1), $2)", token, tag)
	}
	return nil
}

//
func (p postgresql) RemoveTags(token string, tags []string) error {
	for _, tag := range tags {
		p.exec("DELETE FROM tags INNER JOIN tokens ON (tags.token_id = tokens.token_id) WHERE token = $1 AND tag = $2", token, tag)
	}
	return nil
}

//
func (p postgresql) GetTagsForToken(token string) ([]string, error) {

	//
	rows, err := p.query("SELECT tag FROM tags INNER JOIN tokens ON (tags.token_id = tokens.token_id) WHERE token = $1", token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//
	tags := make([]string, 0)

	//
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	err = rows.Err()

	switch {
	case len(tags) == 0:
		return tags, ErrTokenNotFound
	default:
		return tags, err
	}
}

//
func (p postgresql) Clear() error {
	_, err := p.exec("TRUNCATE tokens, tags")
	return err
}

//
func (p postgresql) connect() (*sql.DB, error) {
	return sql.Open("postgres", string(p))
}

// this could really be optimized a lot. instead of opening a new
// conenction for each query, it should reuse connections
func (p postgresql) query(query string, args ...interface{}) (*sql.Rows, error) {
	client, err := p.connect()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.Query(query, args...)
}

// This could also be optimized a lot
func (p postgresql) exec(query string, args ...interface{}) (sql.Result, error) {
	client, err := p.connect()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.Exec(query, args...)
}
