package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jithesh-wq/oolio-food-cart/logger"
	"github.com/jithesh-wq/oolio-food-cart/models"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Con *sql.DB
}

func CreatePostgres() (*Postgres, error) {
	logger.Log.Infoln("Entering : create db connection")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db_directLink := os.Getenv("DATABASE_PUBLIC_URL")
	conString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", dbHost, dbPort, dbUser, dbPassword, dbName)
	log.Println(conString)
	log.Println(db_directLink)
	var db, err = sql.Open("postgres", db_directLink)
	if err != nil {
		logger.Log.Infoln("Error connecting to db" + err.Error())
		return nil, err
	} else {
		logger.Log.Infoln("Database Connected")
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("DAtabase ping failed")
		return nil, err
	}
	logger.Log.Infoln("Exiting : create db connection")
	return &Postgres{Con: db}, nil
}

func (p *Postgres) InsertProducts(products []models.Product) error {
	logger.Log.Infoln("Inside store products data")
	query := `
		INSERT INTO products (
			name, category, price,
			image_thumbnail, image_mobile,image_tablet,image_desktop
		) VALUES `
	var args []interface{}
	var placeholders []string
	for i, product := range products {
		n := i * 7
		placeholders = append(placeholders,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d,$%d)",
				n+1, n+2, n+3, n+4, n+5, n+6, n+7))
		args = append(args,
			product.Name, product.Category, product.Price,
			product.Image.Thumbnail, product.Image.Mobile, product.Image.Tablet,
			product.Image.Desktop)

	}
	query += strings.Join(placeholders, ", ")
	if _, err := p.Con.Exec(query, args...); err != nil {
		logger.Log.Infoln("Failed to insert product:", err)
		return fmt.Errorf("failed to insert products")
	}
	logger.Log.Infoln("Exiting store products data")
	return nil
}

func (p *Postgres) GetProducts(id string) ([]models.Product, error) {
	var (
		rows *sql.Rows
		err  error
	)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("notfound")

	}
	if idInt != 0 {
		query := `
			SELECT id, name, category, price,
			       image_thumbnail, image_mobile,image_tablet, image_desktop
			FROM products
			WHERE id = $1;
		`
		rows, err = p.Con.Query(query, idInt)
	} else {
		query := `
			SELECT id, name, category, price,
			       image_thumbnail, image_mobile,image_tablet,image_desktop
			FROM products
			ORDER BY id;
		`
		rows, err = p.Con.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get product(s): %w", err)
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Category,
			&product.Price,
			&product.Image.Thumbnail,
			&product.Image.Mobile,
			&product.Image.Tablet,
			&product.Image.Desktop,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	// Optional: check for errors from iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if len(products) == 0 {
		logger.Log.Infoln("No products found")
	}

	return products, nil
}

func (p *Postgres) UpsertOrder(sessionID, tableID string, items []models.OrderItem) (string, error) {
	// 1. Check if order exists
	var existingOrderID string
	err := p.Con.QueryRow(`
		SELECT id FROM orders 
		WHERE session_id = $1 AND table_id = $2
	`, sessionID, tableID).Scan(&existingOrderID)

	// Calculate total price from item prices
	totalPrice := 0.0
	for _, item := range items {
		totalPrice += float64(item.Quantity) * item.Price
	}

	if err == sql.ErrNoRows {
		// 2. No existing order, create new
		orderID := uuid.New().String()
		_, err := p.Con.Exec(`
			INSERT INTO orders (id, session_id, table_id, total_price,payment_status)
			VALUES ($1, $2, $3, $4,$5)
		`, orderID, sessionID, tableID, totalPrice, "PAYMENT_PENDING")
		if err != nil {
			return "", fmt.Errorf("failed to insert order: %w", err)
		}

		// Insert order items
		for _, item := range items {
			_, err = p.Con.Exec(`
				INSERT INTO order_items (order_id, product_id, quantity, price)
				VALUES ($1, $2, $3, $4)
			`, orderID, item.ProductID, item.Quantity, item.Price)
			if err != nil {
				return "", fmt.Errorf("failed to insert order items: %w", err)
			}
		}

		return orderID, nil
	} else if err != nil {
		// DB error
		return "", fmt.Errorf("failed to check existing order: %w", err)
	}

	// 3. Existing order found, update order total
	_, err = p.Con.Exec(`
		UPDATE orders 
		SET total_price = total_price + $1 
		WHERE id = $2
	`, totalPrice, existingOrderID)
	if err != nil {
		return "", fmt.Errorf("failed to update order total: %w", err)
	}

	// 4. Insert new items to order
	for _, item := range items {
		_, err = p.Con.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`, existingOrderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return "", fmt.Errorf("failed to insert order items: %w", err)
		}
	}

	return existingOrderID, nil
}

func (p *Postgres) UpdatePaymentStatus(orderID string) error {

	query := `
		UPDATE orders
		SET payment_status = $1
		WHERE id =$2 
	`

	result, err := p.Con.Exec(query, "PAYMENT_SUCCESS", orderID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to determine update status: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no order found with ID: %s", orderID)
	}

	logger.Log.Infof("Payment status updated to Success for order ID %s", orderID)
	return nil
}

func (p *Postgres) GetOrders(req models.OrderRequest) ([]models.Order, error) {
	logger.Log.Infoln("Fetching orders based on user type " + req.UserType + " and payment status " + req.PaymentStatus)

	var (
		rows  *sql.Rows
		err   error
		query string
	)
	if req.UserType == "CUSTOMER" {
		query = `
			SELECT id, total_price, created_at, table_id, session_id, payment_status 
			FROM orders where table_id = $1 AND session_id = $2
		`
		rows, err = p.Con.Query(query, req.TableId, req.SessionID)
		if err != nil {
			return nil, fmt.Errorf("query error: %w", err)
		}
	} else if req.UserType == "ADMIN" {
		if req.PaymentStatus == "" {
			query = `
			SELECT id, total_price, created_at, table_id, session_id, payment_status 
			FROM orders 
		`
			rows, err = p.Con.Query(query)
			if err != nil {
				return nil, fmt.Errorf("query error: %w", err)
			}
		} else {
			query = `
			SELECT id, total_price, created_at, table_id, session_id, payment_status 
			FROM orders where payment_status = $1
		`
			rows, err = p.Con.Query(query, req.PaymentStatus)
			if err != nil {
				return nil, fmt.Errorf("query error: %w", err)
			}
		}
	} else {
		return nil, fmt.Errorf("invalid user type: %s", req.UserType)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var ord models.Order
		if err := rows.Scan(&ord.ID, &ord.TotalPrice, &ord.CreatedAt, &ord.TableID, &ord.SessionID, &ord.PaymentStatus); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		orders = append(orders, ord)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return orders, nil
}
