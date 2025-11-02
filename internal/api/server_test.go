package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"warehouse-api/generated"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage - мок для хранилища, реализующий storage.Storage
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

func TestServer_GetSuppliersByMaterial(t *testing.T) {
	mockStore := new(MockStorage)
	server := NewServer(mockStore)

	tests := []struct {
		name           string
		materialId     int64
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:       "Успешное получение поставщиков",
			materialId: 1,
			mockSetup: func() {
				suppliers := []generated.Supplier{
					{
						SupplierId:         int64Ptr(1),
						Name:               strPtr("Поставщик 1"),
						Inn:                strPtr("1234567890"),
						LegalZipCode:       strPtr("123456"),
						LegalCity:          strPtr("Москва"),
						LegalStreetAddress: strPtr("ул. Ленина, 1"),
						BankZipCode:        strPtr("654321"),
						BankCity:           strPtr("Москва"),
						BankStreetAddress:  strPtr("ул. Банковая, 1"),
						BankAccount:        strPtr("40702810000000000001"),
					},
				}
				mockStore.On("GetSuppliersByMaterial", mock.Anything, int64(1)).Return(suppliers, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "Ошибка базы данных",
			materialId: 2,
			mockSetup: func() {
				mockStore.On("GetSuppliersByMaterial", mock.Anything, int64(2)).Return([]generated.Supplier{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest("GET", "/materials/1/suppliers", nil)
			w := httptest.NewRecorder()

			server.GetSuppliersByMaterial(w, req, tt.materialId)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestServer_CountSuppliersByMaterial(t *testing.T) {
	mockStore := new(MockStorage)
	server := NewServer(mockStore)

	tests := []struct {
		name           string
		materialId     int64
		mockSetup      func()
		expectedStatus int
		expectedCount  int
	}{
		{
			name:       "Успешный подсчет поставщиков",
			materialId: 1,
			mockSetup: func() {
				mockStore.On("CountSuppliersByMaterial", mock.Anything, int64(1)).Return(5, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  5,
		},
		{
			name:       "Ошибка базы данных при подсчете",
			materialId: 2,
			mockSetup: func() {
				mockStore.On("CountSuppliersByMaterial", mock.Anything, int64(2)).Return(0, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest("GET", "/materials/1/suppliers/count", nil)
			w := httptest.NewRecorder()

			server.CountSuppliersByMaterial(w, req, tt.materialId)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response generated.CountResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, *response.Count)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestServer_CreateReceipt(t *testing.T) {
	mockStore := new(MockStorage)
	server := NewServer(mockStore)

	// Создаем дату в UTC и обрезаем время, так как openapi_types.Date представляет только дату
	now := time.Now().UTC()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	testDate := openapi_types.Date{Time: date}

	validReceipt := generated.NewReceipt{
		OrderNumber:     "П-100/24",
		ReceiptDate:     testDate,
		SupplierId:      1,
		BalanceAccount:  "10.01",
		DocTypeId:       1,
		DocumentNumber:  "УПД-12345",
		MaterialId:      1001,
		MaterialAccount: "10.01.1",
		UnitId:          2,
		Quantity:        50.5,
		UnitPrice:       120.00,
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:        "Успешное создание приходного ордера",
			requestBody: validReceipt,
			mockSetup: func() {
				// Используем mock.MatchedBy для сравнения дат правильно
				mockStore.On("CreateReceipt", mock.Anything, mock.MatchedBy(func(receipt generated.NewReceipt) bool {
					return receipt.OrderNumber == validReceipt.OrderNumber &&
						receipt.ReceiptDate.Format("2006-01-02") == validReceipt.ReceiptDate.Format("2006-01-02") &&
						receipt.SupplierId == validReceipt.SupplierId &&
						receipt.BalanceAccount == validReceipt.BalanceAccount &&
						receipt.DocTypeId == validReceipt.DocTypeId &&
						receipt.DocumentNumber == validReceipt.DocumentNumber &&
						receipt.MaterialId == validReceipt.MaterialId &&
						receipt.MaterialAccount == validReceipt.MaterialAccount &&
						receipt.UnitId == validReceipt.UnitId &&
						receipt.Quantity == validReceipt.Quantity &&
						receipt.UnitPrice == validReceipt.UnitPrice
				})).Return(int64(123), nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Некорректный JSON",
			requestBody:    "invalid json",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Ошибка базы данных при создании",
			requestBody: validReceipt,
			mockSetup: func() {
				// Упрощаем матчинг для этого случая - проверяем только OrderNumber
				mockStore.On("CreateReceipt", mock.Anything, mock.MatchedBy(func(receipt generated.NewReceipt) bool {
					return receipt.OrderNumber == validReceipt.OrderNumber
				})).Return(int64(0), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var bodyBytes []byte
			switch v := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				var err error
				bodyBytes, err = json.Marshal(v)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/receipts", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.CreateReceipt(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// Для случая ошибки базы данных убедимся, что мок был вызван
			if tt.name == "Ошибка базы данных при создании" {
				// Проверяем, что был вызов метода CreateReceipt с ошибкой
				mockStore.AssertCalled(t, "CreateReceipt", mock.Anything, mock.AnythingOfType("generated.NewReceipt"))
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestServer_CountSuppliers(t *testing.T) {
	mockStore := new(MockStorage)
	server := NewServer(mockStore)

	// Используем латиницу для параметров запроса чтобы избежать проблем с кодировкой
	params := generated.CountSuppliersParams{
		BankCity:          strPtr("Moscow"),
		BankStreetAddress: strPtr("Bank Street"),
		BankZipCode:       strPtr("123456"),
	}

	tests := []struct {
		name           string
		queryString    string
		mockSetup      func()
		expectedStatus int
		expectedCount  int
	}{
		{
			name:        "Успешный подсчет поставщиков по адресу банка",
			queryString: "bankCity=Moscow&bankStreetAddress=Bank+Street&bankZipCode=123456",
			mockSetup: func() {
				mockStore.On("CountSuppliersByBankAddress", mock.Anything, params).Return(3, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			// Создаем URL с правильным кодированием параметров
			reqURL := "/suppliers/count?" + tt.queryString
			req := httptest.NewRequest("GET", reqURL, nil)
			w := httptest.NewRecorder()

			// Создаем параметры для вызова обработчика
			requestParams := generated.CountSuppliersParams{
				BankCity:          strPtr("Moscow"),
				BankStreetAddress: strPtr("Bank Street"),
				BankZipCode:       strPtr("123456"),
			}

			server.CountSuppliers(w, req, requestParams)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response generated.CountResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, *response.Count)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

// Альтернативный тест для CountSuppliers с использованием английских параметров
func TestServer_CountSuppliers_WithEnglishParams(t *testing.T) {
	mockStore := new(MockStorage)
	server := NewServer(mockStore)

	// Используем параметры на английском
	params := generated.CountSuppliersParams{
		BankCity:          strPtr("Moscow"),
		BankStreetAddress: strPtr("Bank Street"),
		BankZipCode:       strPtr("123456"),
	}

	mockStore.On("CountSuppliersByBankAddress", mock.Anything, params).Return(5, nil)

	// Создаем запрос с английскими параметрами
	req := httptest.NewRequest("GET", "/suppliers/count?bankCity=Moscow&bankStreetAddress=Bank+Street&bankZipCode=123456", nil)
	w := httptest.NewRecorder()

	server.CountSuppliers(w, req, params)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.CountResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 5, *response.Count)

	mockStore.AssertExpectations(t)
}

// Вспомогательные функции для создания указателей с правильными типами
func int64Ptr(i int64) *int64 {
	return &i
}

func strPtr(s string) *string {
	return &s
}
