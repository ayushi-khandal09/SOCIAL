package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserId    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetaData struct {
	Post
	CommentCount int `json:"commentts_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts(content, title, user_id, tags)
		VALUES ($1, $2, $3, $4)	RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserId,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetById(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at, tags, version
		FROM posts
		where id = $1 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, postId int64) error {
	query := `DELETE FROM posts WHERE id= $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err == nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts 
		SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetaData, error) {
    // Base query parts
    baseQuery := `
        SELECT 
            p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username,
            COUNT(c.id) AS comments_count 
        FROM posts p
        LEFT JOIN comments c ON c.post_id = p.id
        LEFT JOIN users u ON p.user_id = u.id
        LEFT JOIN followers f ON f.follower_id = p.user_id
    `
    // Dynamic WHERE conditions
    whereClauses := []string{
        "(f.user_id = $1 OR p.user_id = $1)", // always filter by user relation
    }
    args := []interface{}{userID}
    argPos := 2

    // Search
    if fq.Search != "" {
        whereClauses = append(whereClauses, fmt.Sprintf("(p.title ILIKE '%%' || $%d || '%%' OR p.content ILIKE '%%' || $%d || '%%')", argPos, argPos))
        args = append(args, fq.Search)
        argPos++
    }

    // Tags
    if len(fq.Tags) > 0 {
        whereClauses = append(whereClauses, fmt.Sprintf("p.tags @> $%d", argPos))
        args = append(args, pq.Array(fq.Tags))
        argPos++
    }

    // Since
    if !fq.Since.IsZero() {
        whereClauses = append(whereClauses, fmt.Sprintf("p.created_at >= $%d", argPos))
        args = append(args, fq.Since)
        argPos++
    }

    // Until
    if !fq.Until.IsZero() {
        whereClauses = append(whereClauses, fmt.Sprintf("p.created_at <= $%d", argPos))
        args = append(args, fq.Until)
        argPos++
    }

    // Combine WHERE clauses
    query := baseQuery
    if len(whereClauses) > 0 {
        query += " WHERE " + strings.Join(whereClauses, " AND ")
    }

    // Group/order/limit
    query += fmt.Sprintf(`
        GROUP BY p.id, u.username
        ORDER BY p.created_at %s
        LIMIT $%d OFFSET $%d
    `, fq.Sort, argPos, argPos+1)

    args = append(args, fq.Limit, fq.Offset)

    ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
    defer cancel()

    fmt.Println(query, args)
    rows, err := s.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var feed []PostWithMetaData
    for rows.Next() {
        var p PostWithMetaData
        err := rows.Scan(
            &p.ID,
            &p.UserId,
            &p.Title,
            &p.Content,
            &p.CreatedAt,
            &p.Version,
            pq.Array(&p.Tags),
            &p.User.Username,
            &p.CommentCount,
        )
        if err != nil {
            return nil, err
        }
        feed = append(feed, p)
    }
    return feed, nil
}

