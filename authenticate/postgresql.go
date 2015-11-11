// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package authenticate

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net"
)

type (
	postgresql string
)

func NewPostgresqlAuthenticator(user, database, address string) (postgresql, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return postgresql(""), err
	}

	pg := postgresql(fmt.Sprintf("user=%v database=%v sslmode=disable host=%v port=%v", user, database, host, port))
	// create the tables needed to support mist authentication
	_, err = pg.exec(`
CREATE TABLE IF NOT EXISTS tokens (
	token text NOT NULL,
	token_id SERIAL UNIQUE NOT NULL,
	PRIMARY KEY (token)
)`)

	if err != nil {
		return pg, err
	}

	_, err = pg.exec(`
CREATE TABLE IF NOT EXISTS tags (
  token_id integer NOT NULL REFERENCES tokens (token_id) ON DELETE CASCADE,
  tag text NOT NULL,
  PRIMARY KEY (token_id, tag)
)`)

	return pg, err
}

func (p postgresql) Clear() error {
	_, err := p.exec("TRUNCATE tokens, tags")
	return err
}

func (p postgresql) TagsForToken(token string) ([]string, error) {
	rows, err := p.query("SELECT tag FROM tags,tokens WHERE token = $1", token)
	if err != nil {
		return nil, err
	}

	// now to process the result
	defer rows.Close()
	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	err = rows.Err()
	if len(tags) == 0 && err == nil {
		return tags, NotFound
	}
	return tags, err

}

func (p postgresql) AddTags(token string, tags []string) error {
	// This could be optimized a LOT
	for _, tag := range tags {
		// errors are ignored, this may not be the best idea.
		p.exec("INSERT INTO tags (token_id,tag) VALUES ((SELECT token_id FROM tokens WHERE token = $1), $2)", token, tag)
	}
	return nil
}

func (p postgresql) RemoveTags(token string, tags []string) error {
	for _, tag := range tags {
		p.exec("DELETE FROM tags USING tokens WHERE token = $1 AND tag = $2", token, tag)
	}
	return nil
}

func (p postgresql) AddToken(token string) error {
	_, err := p.exec("INSERT INTO tokens (token) VALUES ($1)", token)
	return err
}

func (p postgresql) RemoveToken(token string) error {
	_, err := p.exec("DELETE FROM tokens WHERE token = $1", token)
	return err
}

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
