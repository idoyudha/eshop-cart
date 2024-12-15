package repo

import (
	"context"
	"strings"
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

const queryInsertCart = `INSERT INTO carts (id, user_id, product_id, product_name, product_price, product_quantity, note, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

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

const queryUpdateQtyAndNoteCart = `UPDATE carts SET product_quantity = ?, note = ?, updated_at = ? WHERE id = ? AND user_id = ? AND deleted_at IS NOT NULL`

func (r *CartMySQLRepo) UpdateQtyAndNote(ctx context.Context, cart *entity.Cart) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdateQtyAndNoteCart)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx, cart.ProductQuantity, cart.Note, cart.UpdatedAt, cart.ID, cart.UserID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

const querySoftDeleteCart = `UPDATE carts SET deleted_at = ? WHERE id IN`

func (r *CartMySQLRepo) DeleteMany(ctx context.Context, cartIDs uuid.UUIDs) error {
	placeholders := "?" + strings.Repeat(",?", len(cartIDs)-1)
	query := querySoftDeleteCart + " (" + placeholders + ")"

	args := make([]interface{}, len(cartIDs)+1)
	args[0] = time.Now()
	for i, id := range cartIDs {
		args[i+1] = id
	}

	stmt, errStmt := r.Conn.PrepareContext(ctx, query)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, deleteErr := stmt.ExecContext(ctx, args...)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

const queryDeleteCart = `DELETE FROM carts WHERE id = ?`

func (r *CartMySQLRepo) DeleteOne(ctx context.Context, cartID uuid.UUID) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryDeleteCart)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, deleteErr := stmt.ExecContext(ctx, cartID)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
