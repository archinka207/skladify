package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"warehouse-api/generated"
	"warehouse-api/internal/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestIntegration_BasicAPI(t *testing.T) {
	// Создаем мок хранилища для интеграционных тестов
	mockStore := &MockStorage{}
	server := api.NewServer(mockStore)

	// Настраиваем роутер
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.Logger)
	chiRouter.Use(middleware.Recoverer)
	handler := generated.HandlerFromMux(server, chiRouter)

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("CountSuppliersByMaterial - успешный запрос", func(t *testing.T) {
		// Настраиваем мок
		mockStore.On("CountSuppliersByMaterial", mock.Anything, int64(1)).Return(3, nil)

		req, err := http.NewRequest("GET", testServer.URL+"/materials/1/suppliers/count", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var countResponse generated.CountResponse
		err = json.NewDecoder(resp.Body).Decode(&countResponse)
		assert.NoError(t, err)
		assert.Equal(t, 3, *countResponse.Count)

		mockStore.AssertExpectations(t)
	})

	t.Run("CreateReceipt - успешное создание", func(t *testing.T) {
		// Создаем дату в UTC и обрезаем время для согласованности
		now := time.Now().UTC()
		date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		testDate := openapi_types.Date{Time: date}

		newReceipt := generated.NewReceipt{
			OrderNumber:     "TEST-001",
			ReceiptDate:     testDate,
			SupplierId:      1,
			BalanceAccount:  "10.01",
			DocTypeId:       1,
			DocumentNumber:  "TEST-DOC-001",
			MaterialId:      1,
			MaterialAccount: "10.01.1",
			UnitId:          1,
			Quantity:        100.0,
			UnitPrice:       50.0,
		}

		// Настраиваем мок с использованием MatchedBy для корректного сравнения дат
		mockStore.On("CreateReceipt", mock.Anything, mock.MatchedBy(func(receipt generated.NewReceipt) bool {
			return receipt.OrderNumber == newReceipt.OrderNumber &&
				receipt.ReceiptDate.Format("2006-01-02") == newReceipt.ReceiptDate.Format("2006-01-02") &&
				receipt.SupplierId == newReceipt.SupplierId &&
				receipt.BalanceAccount == newReceipt.BalanceAccount &&
				receipt.DocTypeId == newReceipt.DocTypeId &&
				receipt.DocumentNumber == newReceipt.DocumentNumber &&
				receipt.MaterialId == newReceipt.MaterialId &&
				receipt.MaterialAccount == newReceipt.MaterialAccount &&
				receipt.UnitId == newReceipt.UnitId &&
				receipt.Quantity == newReceipt.Quantity &&
				receipt.UnitPrice == newReceipt.UnitPrice
		})).Return(int64(123), nil)

		body, err := json.Marshal(newReceipt)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", testServer.URL+"/receipts", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var receipt generated.Receipt
		err = json.NewDecoder(resp.Body).Decode(&receipt)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), *receipt.ReceiptId)

		mockStore.AssertExpectations(t)
	})

	t.Run("GetSuppliersByMaterial - успешный запрос", func(t *testing.T) {
		// Настраиваем мок
		suppliers := []generated.Supplier{
			{
				SupplierId:         int64Ptr(1),
				Name:               strPtr("Тестовый поставщик"),
				Inn:                strPtr("1234567890"),
				LegalZipCode:       strPtr("123456"),
				LegalCity:          strPtr("Москва"),
				LegalStreetAddress: strPtr("ул. Тестовая, 1"),
				BankZipCode:        strPtr("654321"),
				BankCity:           strPtr("Москва"),
				BankStreetAddress:  strPtr("ул. Банковая, 1"),
				BankAccount:        strPtr("40702810000000000001"),
			},
		}
		mockStore.On("GetSuppliersByMaterial", mock.Anything, int64(1)).Return(suppliers, nil)

		req, err := http.NewRequest("GET", testServer.URL+"/materials/1/suppliers", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseSuppliers []generated.Supplier
		err = json.NewDecoder(resp.Body).Decode(&responseSuppliers)
		assert.NoError(t, err)
		assert.Len(t, responseSuppliers, 1)
		assert.Equal(t, "Тестовый поставщик", *responseSuppliers[0].Name)

		mockStore.AssertExpectations(t)
	})

	t.Run("CountSuppliersByBankAddress - успешный запрос", func(t *testing.T) {
		// Настраиваем мок
		mockStore.On("CountSuppliersByBankAddress", mock.Anything, mock.AnythingOfType("generated.CountSuppliersParams")).Return(2, nil)

		req, err := http.NewRequest("GET", testServer.URL+"/suppliers/count?bankCity=Moscow", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var countResponse generated.CountResponse
		err = json.NewDecoder(resp.Body).Decode(&countResponse)
		assert.NoError(t, err)
		assert.Equal(t, 2, *countResponse.Count)

		mockStore.AssertExpectations(t)
	})
}

// MockStorage для интеграционных тестов
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetSuppliersByMaterial(ctx context.Context, materialId int64) ([]generated.Supplier, error) {
	args := m.Called(ctx, materialId)
	return args.Get(0).([]generated.Supplier), args.Error(1)
}

func (m *MockStorage) CountSuppliersByMaterial(ctx context.Context, materialId int64) (int, error) {
	args := m.Called(ctx, materialId)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) CreateReceipt(ctx context.Context, receipt generated.NewReceipt) (int64, error) {
	args := m.Called(ctx, receipt)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorage) CountSuppliersByBankAddress(ctx context.Context, params generated.CountSuppliersParams) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) Close() {
	m.Called()
}

// Вспомогательные функции для создания указателей (нужны только для полей Supplier)
func int64Ptr(i int64) *int64 {
	return &i
}

func strPtr(s string) *string {
	return &s
}
