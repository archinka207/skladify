// Файл: internal/storage/postgres.go
package storage

import (
	"context"
	"fmt"

	"warehouse-api/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStorage реализует взаимодействие с базой данных PostgreSQL.
type PostgresStorage struct {
	pool *pgxpool.Pool
}

// NewPostgresStorage создает новый экземпляр PostgresStorage.
func NewPostgresStorage(ctx context.Context, connString string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать пул соединений: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	return &PostgresStorage{pool: pool}, nil
}

// Close закрывает пул соединений с БД.
func (s *PostgresStorage) Close() {
	s.pool.Close()
}

// GetSuppliersByMaterial находит поставщиков по ID материала.
func (s *PostgresStorage) GetSuppliersByMaterial(ctx context.Context, materialId int64) ([]generated.Supplier, error) {
	query := `
		SELECT s.supplier_id, s.name, s.inn, s.legal_zip_code, s.legal_city, s.legal_street_address,
		       s.bank_zip_code, s.bank_city, s.bank_street_address, s.bank_account
		FROM Suppliers s
		JOIN WarehouseReceipts wr ON s.supplier_id = wr.supplier_id
		WHERE wr.material_id = $1
		GROUP BY s.supplier_id;
	`
	rows, err := s.pool.Query(ctx, query, materialId)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var suppliers []generated.Supplier
	for rows.Next() {
		var sup generated.Supplier
		err := rows.Scan(
			&sup.SupplierId, &sup.Name, &sup.Inn, &sup.LegalZipCode, &sup.LegalCity, &sup.LegalStreetAddress,
			&sup.BankZipCode, &sup.BankCity, &sup.BankStreetAddress, &sup.BankAccount,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		suppliers = append(suppliers, sup)
	}

	return suppliers, nil
}

// CountSuppliersByMaterial считает количество поставщиков по ID материала.
func (s *PostgresStorage) CountSuppliersByMaterial(ctx context.Context, materialId int64) (int, error) {
	query := `
		SELECT COUNT(DISTINCT supplier_id) 
		FROM WarehouseReceipts 
		WHERE material_id = $1;
	`
	var count int
	err := s.pool.QueryRow(ctx, query, materialId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	return count, nil
}

// CreateReceipt создает новую запись о приходе.
func (s *PostgresStorage) CreateReceipt(ctx context.Context, receipt generated.NewReceipt) (int64, error) {
	query := `
		INSERT INTO WarehouseReceipts (order_number, receipt_date, supplier_id, balance_account, doc_type_id, 
		                             document_number, material_id, material_account, unit_id, quantity, unit_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING receipt_id;
	`
	var receiptId int64
	err := s.pool.QueryRow(ctx, query,
		receipt.OrderNumber, receipt.ReceiptDate.Time, receipt.SupplierId, receipt.BalanceAccount, receipt.DocTypeId,
		receipt.DocumentNumber, receipt.MaterialId, receipt.MaterialAccount, receipt.UnitId, receipt.Quantity, receipt.UnitPrice,
	).Scan(&receiptId)

	if err != nil {
		return 0, fmt.Errorf("ошибка вставки записи: %w", err)
	}
	return receiptId, nil
}

// CountSuppliersByBankAddress считает поставщиков по адресу банка.
func (s *PostgresStorage) CountSuppliersByBankAddress(ctx context.Context, params generated.CountSuppliersParams) (int, error) {
	// Мы строим запрос динамически, чтобы обрабатывать опциональные параметры
	query := "SELECT COUNT(*) FROM Suppliers WHERE 1=1"
	args := []interface{}{}
	argId := 1

	if params.BankCity != nil {
		query += fmt.Sprintf(" AND bank_city = $%d", argId)
		args = append(args, *params.BankCity)
		argId++
	}
	if params.BankStreetAddress != nil {
		query += fmt.Sprintf(" AND bank_street_address = $%d", argId)
		args = append(args, *params.BankStreetAddress)
		argId++
	}
	if params.BankZipCode != nil {
		query += fmt.Sprintf(" AND bank_zip_code = $%d", argId)
		args = append(args, *params.BankZipCode)
		argId++
	}

	var count int
	err := s.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	return count, nil
}
