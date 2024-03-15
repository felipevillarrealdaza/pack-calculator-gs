// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const addOrder = `-- name: AddOrder :exec
insert into public.order (order_id, order_quantity) values ($1, $2)
`

type AddOrderParams struct {
	OrderID       uuid.UUID
	OrderQuantity int32
}

func (q *Queries) AddOrder(ctx context.Context, arg AddOrderParams) error {
	_, err := q.db.ExecContext(ctx, addOrder, arg.OrderID, arg.OrderQuantity)
	return err
}

const addOrderPack = `-- name: AddOrderPack :exec
insert into public.order_packs (order_packs_id, order_id, pack_size, pack_quantity) values ($1, $2, $3, $4)
`

type AddOrderPackParams struct {
	OrderPacksID uuid.UUID
	OrderID      uuid.UUID
	PackSize     int32
	PackQuantity int32
}

func (q *Queries) AddOrderPack(ctx context.Context, arg AddOrderPackParams) error {
	_, err := q.db.ExecContext(ctx, addOrderPack,
		arg.OrderPacksID,
		arg.OrderID,
		arg.PackSize,
		arg.PackQuantity,
	)
	return err
}

const addPack = `-- name: AddPack :exec
insert into pack (pack_size) values ($1)
`

func (q *Queries) AddPack(ctx context.Context, packSize int32) error {
	_, err := q.db.ExecContext(ctx, addPack, packSize)
	return err
}

const removePackBySize = `-- name: RemovePackBySize :exec
delete from public.pack where public.pack.pack_size = $1
`

func (q *Queries) RemovePackBySize(ctx context.Context, packSize int32) error {
	_, err := q.db.ExecContext(ctx, removePackBySize, packSize)
	return err
}

const retrieveOrderById = `-- name: RetrieveOrderById :one
select order_id, order_quantity from public.order where public.order.order_id = $1
`

func (q *Queries) RetrieveOrderById(ctx context.Context, orderID uuid.UUID) (Order, error) {
	row := q.db.QueryRowContext(ctx, retrieveOrderById, orderID)
	var i Order
	err := row.Scan(&i.OrderID, &i.OrderQuantity)
	return i, err
}

const retrieveOrderPacksByOrder = `-- name: RetrieveOrderPacksByOrder :many
select
s.order_packs_id,
o.order_quantity,
p.pack_size,
s.pack_quantity
from public.order_packs s
inner join public.order o on s.order_id = o.order_id
inner join public.pack p on s.pack_id = p.pack_id
where o.order_id = $1
`

type RetrieveOrderPacksByOrderRow struct {
	OrderPacksID  uuid.UUID
	OrderQuantity int32
	PackSize      int32
	PackQuantity  int32
}

func (q *Queries) RetrieveOrderPacksByOrder(ctx context.Context, orderID uuid.UUID) ([]RetrieveOrderPacksByOrderRow, error) {
	rows, err := q.db.QueryContext(ctx, retrieveOrderPacksByOrder, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RetrieveOrderPacksByOrderRow
	for rows.Next() {
		var i RetrieveOrderPacksByOrderRow
		if err := rows.Scan(
			&i.OrderPacksID,
			&i.OrderQuantity,
			&i.PackSize,
			&i.PackQuantity,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const retrieveOrders = `-- name: RetrieveOrders :many
select order_id, order_quantity from public.order
`

func (q *Queries) RetrieveOrders(ctx context.Context) ([]Order, error) {
	rows, err := q.db.QueryContext(ctx, retrieveOrders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(&i.OrderID, &i.OrderQuantity); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const retrievePacks = `-- name: RetrievePacks :many
select pack_size from public.pack ORDER BY pack_size DESC
`

func (q *Queries) RetrievePacks(ctx context.Context) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, retrievePacks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var pack_size int32
		if err := rows.Scan(&pack_size); err != nil {
			return nil, err
		}
		items = append(items, pack_size)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
