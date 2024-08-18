package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type Link struct {
	ID     int    `json:"id"`
	Slug   string `json:"slug"`
	URL    string `json:"url"`
	Clicks int    `json:"clicks"`
}

type LinkProvider struct {
	db *sql.DB
}

func NewLinkProvider(db *sql.DB) (*LinkProvider, error) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS links(
		id integer primary key autoincrement,
		slug text,
		url text,
		clicks integer not null
	);`)
	if err != nil {
		return nil, err
	}

	return &LinkProvider{
		db: db,
	}, nil
}

var ErrEmptyLinkSlug = errors.New("empty link slug")

func (s *LinkProvider) Create(ctx context.Context, link Link) (Link, error) {
	if link.Slug == "" {
		return Link{}, ErrEmptyLinkSlug
	}

	query := `
	INSERT INTO links(slug, url, clicks)
	           VALUES(   ?,   ?,      0);`

	res, err := s.db.ExecContext(ctx, query, link.Slug, link.URL)
	if err != nil {
		return Link{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Link{}, err
	}

	link.ID = int(id)

	return link, nil
}

func (s *LinkProvider) Delete(ctx context.Context, id string) error {
	query := `
	DELETE FROM links
        WHERE id = ?;`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *LinkProvider) Get(ctx context.Context, id string) (Link, error) {
	if id == "" {
		return Link{}, nil
	}
	query := `
	SELECT id, slug, url, clicks
	FROM links
        WHERE id = ?;`

	row := s.db.QueryRowContext(ctx, query, id)

	var l Link

	err := row.Scan(&l.ID, &l.Slug, &l.URL, &l.Clicks)
	if err != nil {
		return Link{}, err
	}

	return l, nil
}

func (s *LinkProvider) Update(ctx context.Context, link Link) error {
	log.Println("update:", link)
	if link.Slug == "" {
		return errors.New("empty link ID")
	}

	query := `
	UPDATE links
	SET
	  id   = ?,
	  slug = ?,
	  url  = ?
	WHERE
	  id   = ?;`

	_, err := s.db.ExecContext(ctx, query, link.ID, link.Slug, link.URL, link.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *LinkProvider) List(ctx context.Context, offset, limit int) ([]Link, error) {
	query := `
	SELECT id, slug, url, clicks
	FROM links
	LIMIT ?
	OFFSET ?;`

	var links []Link
	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// we return an empty slice
			return links, nil
		}
		return nil, err
	}

	for rows.Next() {
		var l Link
		err = rows.Scan(&l.ID, &l.Slug, &l.URL, &l.Clicks)
		if err != nil {
			return nil, err
		}

		links = append(links, l)
	}

	return links, nil
}
