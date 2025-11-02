package storage

import (
	"context"
	"warehouse-api/generated"
)

// Storage интерфейс определяет контракт для работы с данными
type Storage interface {
	GetSuppliersByMaterial(ctx context.Context, materialId int64) ([]generated.Supplier, error)
	CountSuppliersByMaterial(ctx context.Context, materialId int64) (int, error)
	CreateReceipt(ctx context.Context, receipt generated.NewReceipt) (int64, error)
	CountSuppliersByBankAddress(ctx context.Context, params generated.CountSuppliersParams) (int, error)
	Close()
}

// Убедимся, что PostgresStorage реализует интерфейс Storage
var _ Storage = (*PostgresStorage)(nil)
