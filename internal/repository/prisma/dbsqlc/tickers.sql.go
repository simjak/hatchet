// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: tickers.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const listActiveTickers = `-- name: ListActiveTickers :many
SELECT
    tickers.id, tickers."createdAt", tickers."updatedAt", tickers."lastHeartbeatAt", tickers."isActive"
FROM "Ticker" as tickers
WHERE
    -- last heartbeat greater than 15 seconds
    "lastHeartbeatAt" > NOW () - INTERVAL '15 seconds'
    -- active
    AND "isActive" = true
`

type ListActiveTickersRow struct {
	Ticker Ticker `json:"ticker"`
}

func (q *Queries) ListActiveTickers(ctx context.Context, db DBTX) ([]*ListActiveTickersRow, error) {
	rows, err := db.Query(ctx, listActiveTickers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListActiveTickersRow
	for rows.Next() {
		var i ListActiveTickersRow
		if err := rows.Scan(
			&i.Ticker.ID,
			&i.Ticker.CreatedAt,
			&i.Ticker.UpdatedAt,
			&i.Ticker.LastHeartbeatAt,
			&i.Ticker.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listStaleTickers = `-- name: ListStaleTickers :many
SELECT
    tickers.id, tickers."createdAt", tickers."updatedAt", tickers."lastHeartbeatAt", tickers."isActive"
FROM "Ticker" as tickers
WHERE
    -- last heartbeat older than 15 seconds
    "lastHeartbeatAt" < NOW () - INTERVAL '15 seconds'
    -- not active
    AND "isActive" = false
`

type ListStaleTickersRow struct {
	Ticker Ticker `json:"ticker"`
}

func (q *Queries) ListStaleTickers(ctx context.Context, db DBTX) ([]*ListStaleTickersRow, error) {
	rows, err := db.Query(ctx, listStaleTickers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListStaleTickersRow
	for rows.Next() {
		var i ListStaleTickersRow
		if err := rows.Scan(
			&i.Ticker.ID,
			&i.Ticker.CreatedAt,
			&i.Ticker.UpdatedAt,
			&i.Ticker.LastHeartbeatAt,
			&i.Ticker.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTickers = `-- name: ListTickers :many
SELECT
    tickers.id, tickers."createdAt", tickers."updatedAt", tickers."lastHeartbeatAt", tickers."isActive"
FROM
    "Ticker" as tickers
`

type ListTickersRow struct {
	Ticker Ticker `json:"ticker"`
}

func (q *Queries) ListTickers(ctx context.Context, db DBTX) ([]*ListTickersRow, error) {
	rows, err := db.Query(ctx, listTickers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListTickersRow
	for rows.Next() {
		var i ListTickersRow
		if err := rows.Scan(
			&i.Ticker.ID,
			&i.Ticker.CreatedAt,
			&i.Ticker.UpdatedAt,
			&i.Ticker.LastHeartbeatAt,
			&i.Ticker.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setTickersInactive = `-- name: SetTickersInactive :many
UPDATE
    "Ticker" as tickers
SET
    "isActive" = false
WHERE
    "id" = ANY ($1::uuid[])
RETURNING
    tickers.id, tickers."createdAt", tickers."updatedAt", tickers."lastHeartbeatAt", tickers."isActive"
`

type SetTickersInactiveRow struct {
	Ticker Ticker `json:"ticker"`
}

func (q *Queries) SetTickersInactive(ctx context.Context, db DBTX, ids []pgtype.UUID) ([]*SetTickersInactiveRow, error) {
	rows, err := db.Query(ctx, setTickersInactive, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*SetTickersInactiveRow
	for rows.Next() {
		var i SetTickersInactiveRow
		if err := rows.Scan(
			&i.Ticker.ID,
			&i.Ticker.CreatedAt,
			&i.Ticker.UpdatedAt,
			&i.Ticker.LastHeartbeatAt,
			&i.Ticker.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
