package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"wb_test/internal/model"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

type OrderStorage struct {
	db *Storage
}

func (s *Storage) Orders() *OrderStorage {
	return &OrderStorage{db: s}
}

func (os *OrderStorage) Upsert(ctx context.Context, order *model.Order) error {
	tx, err := os.db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if order.DateCreated.IsZero() {
		order.DateCreated = time.Now().UTC()
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id,
		                    delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO UPDATE SET 
		    track_number = EXCLUDED.track_number,
		    entry = EXCLUDED.entry,
		    locale = EXCLUDED.locale,
		    internal_signature = EXCLUDED.internal_signature,
		    customer_id = EXCLUDED.customer_id,
		    delivery_service = EXCLUDED.delivery_service,
		    shardkey = EXCLUDED.shardkey,
		    sm_id = EXCLUDED.sm_id,
		    date_created = EXCLUDED.date_created,
		    oof_shard = EXCLUDED.oof_shard
	`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO deliveries (name, phone, zip, city, address, region, email, order_uid)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (order_uid) DO UPDATE SET
		    name = EXCLUDED.name,
		    phone = EXCLUDED.phone,
		    zip = EXCLUDED.zip,
		    city = EXCLUDED.city,
		    address = EXCLUDED.address,
		    region = EXCLUDED.region,
		    email = EXCLUDED.email
	`, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email, order.OrderUID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, order_uid)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO UPDATE SET
		    transaction = EXCLUDED.transaction,
		    request_id = EXCLUDED.request_id,
		    currency = EXCLUDED.currency,
		    provider = EXCLUDED.provider,
		    amount = EXCLUDED.amount,
		    payment_dt = EXCLUDED.payment_dt,
		    bank = EXCLUDED.bank,
		    delivery_cost = EXCLUDED.delivery_cost,
		    goods_total = EXCLUDED.goods_total,
		    custom_fee = EXCLUDED.custom_fee
	`, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee, order.OrderUID)
	if err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, `DELETE FROM items WHERE order_uid=$1`, order.OrderUID); err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, it := range order.Items {
		batch.Queue(`
			INSERT INTO items (chrt_id, price, rid, name, sale, size, total_price, nm_id, brand, status, track_number, order_uid)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
			ON CONFLICT (chrt_id) DO UPDATE SET
			    price = EXCLUDED.price,
			    rid = EXCLUDED.rid,
			    name = EXCLUDED.name,
			    sale = EXCLUDED.sale,
			    size = EXCLUDED.size,
			    total_price = EXCLUDED.total_price,
			    nm_id = EXCLUDED.nm_id,
			    brand = EXCLUDED.brand,
			    status = EXCLUDED.status,
			    track_number = EXCLUDED.track_number,
			    order_uid = EXCLUDED.order_uid
		`, it.ChrtID, it.Price, it.RID, it.Name, it.Sale, it.Size, it.TotalPrice,
			it.NmID, it.Brand, it.Status, it.TrackNumber, order.OrderUID)
	}

	br := tx.SendBatch(ctx, batch)
	for range order.Items {
		if _, err := br.Exec(); err != nil {
			_ = br.Close()
			return err
		}
	}
	if err := br.Close(); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (os *OrderStorage) Get(ctx context.Context, orderUID string) (*model.Order, error) {
	order := model.Order{}
	err := os.db.pool.QueryRow(ctx, `
		SELECT order_uid, track_number, entry, locale, COALESCE(internal_signature,''), customer_id,
		       delivery_service, COALESCE(shardkey,''), COALESCE(sm_id,0), date_created, COALESCE(oof_shard,'')
		FROM orders
		WHERE order_uid=$1
	`, orderUID).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	err = os.db.pool.QueryRow(ctx, `
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries
		WHERE order_uid=$1
	`, orderUID).Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
		&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		return nil, err
	}

	err = os.db.pool.QueryRow(ctx, `
		SELECT transaction, COALESCE(request_id,''), currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments
		WHERE order_uid=$1
	`, orderUID).Scan(&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider,
		&order.Payment.Amount, &order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal, &order.Payment.CustomFee)
	if err != nil {
		return nil, err
	}

	rows, err := os.db.pool.Query(ctx, `
		SELECT chrt_id, price, rid, name, sale, size, total_price, nm_id, brand, status, track_number
		FROM items
		WHERE order_uid=$1
	`, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var it model.Item
		if err := rows.Scan(&it.ChrtID, &it.Price, &it.RID, &it.Name, &it.Sale, &it.Size, &it.TotalPrice,
			&it.NmID, &it.Brand, &it.Status, &it.TrackNumber); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, it)
	}
	return &order, nil
}
