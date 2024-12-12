package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-cart/internal/entity"
	mysqlClient "github.com/idoyudha/eshop-cart/pkg/mysql"
)

type CartMySQLRepo struct {
	*mysqlClient.MySQL
}

func NewCartMySQLRepo(client *mysqlClient.MySQL) *CartMySQLRepo {
	return &CartMySQLRepo{
		client,
	}
}

const queryInsertCart = `INSERT INTO carts (id, user_id, product_id, product_name, product_name, product_price, product_quantity, note, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

func (r *CartMySQLRepo) Insert(ctx context.Context, cart *entity.Cart) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryInsertCart)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, insertErr := stmt.ExecContext(ctx, cart.ID, cart.UserID, cart.ProductID, cart.ProductName, cart.ProductPrice, cart.ProductQuantity, cart.Note, cart.CreatedAt, cart.UpdatedAt)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

const getCartsQueryByUserID = `SELECT id, user_id, product_id, product_name, product_price, product_quantity, note, created_at, updated_at FROM carts WHERE user_id = ? AND deleted_at IS NULL`

func (r *CartMySQLRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Cart, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, getCartsQueryByUserID)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carts := make([]*entity.Cart, 0)
	for rows.Next() {
		cart := &entity.Cart{}
		err := rows.Scan(&cart.ID, &cart.UserID, &cart.ProductID, &cart.ProductName, &cart.ProductPrice, &cart.ProductQuantity, &cart.Note, &cart.CreatedAt, &cart.UpdatedAt)
		if err != nil {
			continue
		}
		carts = append(carts, cart)
	}

	return carts, nil
}

const queryUpdateCart = `UPDATE carts SET product_name = ?, product_price = ?, product_quantity = ?, note = ?, updated_at = ? WHERE id = ?`

func (r *CartMySQLRepo) Update(ctx context.Context, cart *entity.Cart) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdateCart)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx, cart.ProductName, cart.ProductPrice, cart.ProductQuantity, cart.Note, cart.UpdatedAt, cart.ID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

const queryDeleteCart = `UPDATE carts SET deleted_at = ? WHERE id IN ?`

func (r *CartMySQLRepo) DeleteMany(ctx context.Context, cartIDs uuid.UUIDs) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryDeleteCart)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, deleteErr := stmt.ExecContext(ctx, time.Now(), cartIDs)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
