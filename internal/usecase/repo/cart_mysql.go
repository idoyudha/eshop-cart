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

const queryInsertCart = `INSERT INTO carts (id, user_id, product_id, product_name, product_image_url, product_price, product_quantity, note, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

func (r *CartMySQLRepo) Insert(ctx context.Context, cart *entity.Cart) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryInsertCart)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, insertErr := stmt.ExecContext(ctx, cart.ID, cart.UserID, cart.ProductID, cart.ProductName, cart.ProductImageURL, cart.ProductPrice, cart.ProductQuantity, cart.Note, cart.CreatedAt, cart.UpdatedAt)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

const getCartsQueryByUserID = `SELECT id, user_id, product_id, product_name, product_image_url, product_price, product_quantity, note, created_at, updated_at FROM carts WHERE user_id = ? AND deleted_at IS NULL`

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
		err := rows.Scan(&cart.ID, &cart.UserID, &cart.ProductID, &cart.ProductName, &cart.ProductImageURL, &cart.ProductPrice, &cart.ProductQuantity, &cart.Note, &cart.CreatedAt, &cart.UpdatedAt)
		if err != nil {
			continue
		}
		carts = append(carts, cart)
	}

	return carts, nil
}

const queryUpdateQtyAndNoteCart = `UPDATE carts SET product_quantity = ?, note = ?, updated_at = ? WHERE id = ? AND user_id = ? AND deleted_at IS NOT NULL`
const querySelectUpdatedCart = `SELECT product_id FROM carts WHERE id = ? AND user_id = ? AND deleted_at IS NULL`

func (r *CartMySQLRepo) UpdateQtyAndNote(ctx context.Context, cart *entity.Cart) (*uuid.UUID, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdateQtyAndNoteCart)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx, cart.ProductQuantity, cart.Note, cart.UpdatedAt, cart.ID, cart.UserID)
	if updateErr != nil {
		return nil, updateErr
	}

	// get product id for updating in redis
	stmt, errStmt = r.Conn.PrepareContext(ctx, querySelectUpdatedCart)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, cart.ID, cart.UserID)
	var productID uuid.UUID
	err := row.Scan(&productID)
	if err != nil {
		return nil, err
	}

	return &productID, nil
}

const queryUpdateNameAndPriceCart = `UPDATE carts SET product_name = ?, product_price = ?, updated_at = ? WHERE product_id = ? AND deleted_at IS NOT NULL`

func (r *CartMySQLRepo) UpdateNameAndPrice(ctx context.Context, cart *entity.Cart) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdateNameAndPriceCart)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx, cart.ProductName, cart.ProductPrice, cart.UpdatedAt, cart.ProductID)
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

const queryDeleteCart = `UPDATE carts SET deleted_at = ? WHERE id = ?`
const queryGetProductIDDeletedCart = `SELECT product_id FROM carts WHERE id = ? AND deleted_at IS NOT NULL`

func (r *CartMySQLRepo) DeleteOne(ctx context.Context, cartID uuid.UUID) (*uuid.UUID, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryDeleteCart)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	_, deleteErr := stmt.ExecContext(ctx, time.Now(), cartID)
	if deleteErr != nil {
		return nil, deleteErr
	}

	// get product id for updating in redis
	stmt, errStmt = r.Conn.PrepareContext(ctx, queryGetProductIDDeletedCart)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, cartID)
	var productID uuid.UUID
	err := row.Scan(&productID)
	if err != nil {
		return nil, err
	}

	return &productID, nil
}

const queryUpdateProductQtyCart = `UPDATE carts SET product_quantity = product_quantity + ?, updated_at = ? WHERE id = ? AND user_id = ? AND deleted_at IS NOT NULL`

func (r *CartMySQLRepo) UpdateProductQty(ctx context.Context, cart *entity.Cart) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdateProductQtyCart)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx, cart.ProductQuantity, cart.UpdatedAt, cart.ID, cart.UserID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}
