package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type ShortURL struct {
	ID     int    `json:"id"`
	Slug   string `json:"slug"`
	URL    string `json:"url"`
	Clicks int    `json:"clicks"`
}

type ShortURLProvider struct {
	db *sql.DB
}

func NewShortURLProvider(db *sql.DB) (*ShortURLProvider, error) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS urls(
		id integer primary key autoincrement,
		slug text,
		url text,
		clicks integer not null
	);`)
	if err != nil {
		return nil, err
	}

	return &ShortURLProvider{
		db: db,
	}, nil
}

var ErrEmptyURLSlug = errors.New("empty url slug")

func (s *ShortURLProvider) Create(ctx context.Context, url ShortURL) (ShortURL, error) {
	if url.Slug == "" {
		return ShortURL{}, ErrEmptyURLSlug
	}

	query := `
	INSERT INTO urls(slug, url, clicks)
	           VALUES(   ?,   ?,      0);`

	res, err := s.db.ExecContext(ctx, query, url.Slug, url.URL)
	if err != nil {
		return ShortURL{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return ShortURL{}, err
	}

	url.ID = int(id)

	return url, nil
}

func (s *ShortURLProvider) Delete(ctx context.Context, id string) error {
	query := `
	DELETE FROM urls
        WHERE id = ?;`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ShortURLProvider) Get(ctx context.Context, id string) (ShortURL, error) {
	if id == "" {
		return ShortURL{}, nil
	}
	query := `
	SELECT id, slug, url, clicks
	FROM urls
        WHERE id = ?;`

	row := s.db.QueryRowContext(ctx, query, id)

	var l ShortURL

	err := row.Scan(&l.ID, &l.Slug, &l.URL, &l.Clicks)
	if err != nil {
		return ShortURL{}, err
	}

	return l, nil
}

func (s *ShortURLProvider) Update(ctx context.Context, url ShortURL) error {
	log.Println("update:", url)
	if url.Slug == "" {
		return errors.New("empty url ID")
	}

	query := `
	UPDATE urls
	SET
	  id   = ?,
	  slug = ?,
	  url  = ?
	WHERE
	  id   = ?;`

	_, err := s.db.ExecContext(ctx, query, url.ID, url.Slug, url.URL, url.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ShortURLProvider) List(ctx context.Context, offset, limit int) ([]ShortURL, error) {
	query := `
	SELECT id, slug, url, clicks
	FROM urls
	LIMIT ?
	OFFSET ?;`

	var urls []ShortURL
	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// we return an empty slice
			return urls, nil
		}
		return nil, err
	}

	for rows.Next() {
		var l ShortURL
		err = rows.Scan(&l.ID, &l.Slug, &l.URL, &l.Clicks)
		if err != nil {
			return nil, err
		}

		urls = append(urls, l)
	}

	return urls, nil
}

func (s *ShortURLProvider) FindBySlug(ctx context.Context, slug string) (ShortURL, error) {
	if slug == "" {
		return ShortURL{}, nil
	}
	query := `
	SELECT id, slug, url, clicks
	FROM urls
        WHERE slug = ?;`

	row := s.db.QueryRowContext(ctx, query, slug)

	var l ShortURL

	err := row.Scan(&l.ID, &l.Slug, &l.URL, &l.Clicks)
	if err != nil {
		return ShortURL{}, err
	}

	return l, nil
}

func (s *ShortURLProvider) Click(ctx context.Context, slug string) error {
	query := `
	UPDATE urls
	SET clicks = clicks + 1
	WHERE slug = ?`

	_, err := s.db.ExecContext(ctx, query, slug)
	if err != nil {
		return err
	}
	return nil
}
