// Файл: internal/api/server.go
package api

import (
	"encoding/json"
	"log"
	"net/http"

	"warehouse-api/generated"
	"warehouse-api/internal/storage"
)

// Server - наша реализация ServerInterface.
type Server struct {
	Store *storage.PostgresStorage
}

// NewServer создает новый экземпляр сервера.
func NewServer(store *storage.PostgresStorage) *Server {
	return &Server{Store: store}
}

// Убеждаемся, что Server соответствует интерфейсу во время компиляции.
var _ generated.ServerInterface = (*Server)(nil)

// --- Вспомогательные функции для ответов ---

func (s *Server) renderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка при кодировании JSON-ответа: %v", err)
	}
}

func (s *Server) renderError(w http.ResponseWriter, status int, message string) {
	errResponse := generated.Error{
		Code:    &status,
		Message: &message,
	}
	s.renderJSON(w, status, errResponse)
}

// --- Реализация методов ServerInterface ---

func (s *Server) GetSuppliersByMaterial(w http.ResponseWriter, r *http.Request, materialId int64) {
	suppliers, err := s.Store.GetSuppliersByMaterial(r.Context(), materialId)
	if err != nil {
		log.Printf("Ошибка при получении поставщиков: %v", err)
		s.renderError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}
	s.renderJSON(w, http.StatusOK, suppliers)
}

func (s *Server) CountSuppliersByMaterial(w http.ResponseWriter, r *http.Request, materialId int64) {
	count, err := s.Store.CountSuppliersByMaterial(r.Context(), materialId)
	if err != nil {
		log.Printf("Ошибка при подсчете поставщиков: %v", err)
		s.renderError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}
	response := generated.CountResponse{Count: &count}
	s.renderJSON(w, http.StatusOK, response)
}

func (s *Server) CreateReceipt(w http.ResponseWriter, r *http.Request) {
	var newReceipt generated.NewReceipt
	if err := json.NewDecoder(r.Body).Decode(&newReceipt); err != nil {
		s.renderError(w, http.StatusBadRequest, "Некорректное тело запроса")
		return
	}

	receiptId, err := s.Store.CreateReceipt(r.Context(), newReceipt)
	if err != nil {
		log.Printf("Ошибка при создании записи о приходе: %v", err)
		// Здесь можно добавить проверку на specific pg errors, например, foreign key violation
		s.renderError(w, http.StatusInternalServerError, "Не удалось создать запись")
		return
	}

	// Возвращаем полный объект Receipt
	createdReceipt := generated.Receipt{
		ReceiptId:       &receiptId,
		OrderNumber:     newReceipt.OrderNumber,
		ReceiptDate:     newReceipt.ReceiptDate,
		SupplierId:      newReceipt.SupplierId,
		BalanceAccount:  newReceipt.BalanceAccount,
		DocTypeId:       newReceipt.DocTypeId,
		DocumentNumber:  newReceipt.DocumentNumber,
		MaterialId:      newReceipt.MaterialId,
		MaterialAccount: newReceipt.MaterialAccount,
		UnitId:          newReceipt.UnitId,
		Quantity:        newReceipt.Quantity,
		UnitPrice:       newReceipt.UnitPrice,
	}

	s.renderJSON(w, http.StatusCreated, createdReceipt)
}

func (s *Server) CountSuppliers(w http.ResponseWriter, r *http.Request, params generated.CountSuppliersParams) {
	count, err := s.Store.CountSuppliersByBankAddress(r.Context(), params)
	if err != nil {
		log.Printf("Ошибка при подсчете поставщиков по адресу банка: %v", err)
		s.renderError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}
	response := generated.CountResponse{Count: &count}
	s.renderJSON(w, http.StatusOK, response)
}
