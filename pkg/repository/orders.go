package repository

import (
	"fmt"

	"github.com/TemaStatham/OrderService/pkg/model"
	"github.com/jmoiron/sqlx"
)

type OrdersPostgres struct {
	db *sqlx.DB
}

func NewOrdersPostgres(db *sqlx.DB) *OrdersPostgres {
	return &OrdersPostgres{db: db}
}

// GetOrder : Запрос для получения основной информации о заказе
func (o *OrdersPostgres) GetOrder(orderID string) (*model.OrderClient, error) {
	orderQuery := "SELECT * FROM orders WHERE order_uid = $1"

	var order model.OrderClient
	if err := o.db.Get(&order, orderQuery, orderID); err != nil {
		return nil, err
	}

	payment, err := o.GetPayment(orderID)
	if err != nil {
		return nil, err
	}
	order.Payment = *payment

	items, err := o.GetItems(orderID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	delivery, err := o.GetDelivery(orderID)
	if err != nil {
		return nil, err
	}
	order.Delivery = *delivery

	return &order, nil
}

// GetPayment : Запрос для получения информации о платеже
func (o *OrdersPostgres) GetPayment(orderID string) (*model.Payment, error) {
	paymentQuery := "SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM order_payment WHERE order_uid = $1"

	var payment model.Payment
	if err := o.db.Get(&payment, paymentQuery, orderID); err != nil {
		return nil, err
	}

	return &payment, nil
}

// GetItems : Запрос для получения информации о товарах
func (o *OrdersPostgres) GetItems(orderID string) ([]model.Item, error) {
	itemsQuery := "SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM order_items WHERE order_uid = $1;"

	var items []model.Item
	if err := o.db.Select(&items, itemsQuery, orderID); err != nil {
		return nil, err
	}

	return items, nil
}

// GetDelivery : Запрос для получения информации о доставке
func (o *OrdersPostgres) GetDelivery(orderID string) (*model.Delivery, error) {
	deliveryQuery := "SELECT name, phone, zip, city, address, region, email FROM order_delivery WHERE order_uid = $1;"

	var delivery model.Delivery
	if err := o.db.Get(&delivery, deliveryQuery, orderID); err != nil {
		return nil, err
	}

	return &delivery, nil
}

// AddOrder : добавляет заказ в бд
func (o *OrdersPostgres) AddOrder(order *model.OrderClient) (string, error) {
	tx, err := o.db.Beginx()
	if err != nil {
		return "", fmt.Errorf("begin transaction failed %s", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	fmt.Print("4\n")
	orderID, err := o.addOrder(tx, order)
	if err != nil {
		return "", fmt.Errorf("failed to add order %s", err)
	}

	fmt.Print("1\n")
	_, err = o.addDelivery(tx, order.Delivery, order.OrderUID)
	if err != nil {
		return "", fmt.Errorf("failed to add delivery %s", err)
	}

	fmt.Print("2\n")
	_, err = o.addPayment(tx, order.Payment, order.OrderUID)
	if err != nil {
		return "", fmt.Errorf("failed to add payment %s", err)
	}

	fmt.Print("3\n")
	// itemIDs := make([]string, 0)
	for _, item := range order.Items {
		_, err = o.addItem(tx, item, order.OrderUID)
		if err != nil {
			return "", fmt.Errorf("failed to add items %s", err)
		}
		// itemIDs = append(itemIDs, itemID)
	}

	fmt.Print("5\n")
	return orderID, nil
}

func (o *OrdersPostgres) addDelivery(tx *sqlx.Tx, delivery model.Delivery, orderID string) (string, error) {
	query := `
		INSERT INTO order_delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING order_uid
	`
	row := tx.QueryRow(query, orderID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email)

	var deliveryID string
	if err := row.Scan(&deliveryID); err != nil {
		return "", fmt.Errorf("failed to insert delivery %s", err)
	}

	return deliveryID, nil
}

func (o *OrdersPostgres) addPayment(tx *sqlx.Tx, payment model.Payment, orderID string) (string, error) {
	query := `
		INSERT INTO order_payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING order_uid
	`
	row := tx.QueryRow(query, orderID, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount,
		payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)

	var paymentID string
	if err := row.Scan(&paymentID); err != nil {
		return "", fmt.Errorf("failed to insert payment %s", err)
	}

	return paymentID, nil
}

func (o *OrdersPostgres) addItem(tx *sqlx.Tx, item model.Item, orderID string) (string, error) {
	query := `
		INSERT INTO order_items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING order_uid
	`
	row := tx.QueryRow(query, orderID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice,
		item.NmID, item.Brand, item.Status)

	var itemID string
	if err := row.Scan(&itemID); err != nil {
		return "", fmt.Errorf("failed to insert item %s", err)
	}

	return itemID, nil
}

func (o *OrdersPostgres) addOrder(tx *sqlx.Tx, order *model.OrderClient) (string, error) {
	// Вставка в таблицу orders
	query := `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING order_uid
	`
	row := tx.QueryRow(query, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)

	var orderID string
	if err := row.Scan(&orderID); err != nil {
		return "", fmt.Errorf("failed to insert order %s", err)
	}

	return orderID, nil
}
