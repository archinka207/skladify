// Файл: cmd/server/main.go
// Файл: cmd/server/main.go (ИСПРАВЛЕННАЯ ВЕРСИЯ)
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"warehouse-api/config"
	"warehouse-api/generated"
	"warehouse-api/internal/api"
	"warehouse-api/internal/storage"

	"github.com/go-chi/chi/v5" // <-- Важно: импортируем chi
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Инициализация хранилища (подключение к БД)
	store, err := storage.NewPostgresStorage(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer store.Close()
	log.Println("Успешное подключение к базе данных")

	// 3. Создание экземпляра сервера с бизнес-логикой
	apiServer := api.NewServer(store)

	// --- ИЗМЕНЕНИЯ ЗДЕСЬ ---

	// 4. Сначала создаем экземпляр роутера chi
	chiRouter := chi.NewRouter()

	// 5. Добавляем к нему наши middleware
	chiRouter.Use(middleware.Logger)
	chiRouter.Use(middleware.Recoverer)

	// 6. Теперь передаем наш роутер в сгенерированную функцию.
	// Она добавит все API-маршруты в уже настроенный chiRouter.
	// Переменную `handler` можно использовать как основной обработчик для http.Server.
	handler := generated.HandlerFromMux(apiServer, chiRouter)

	// --- КОНЕЦ ИЗМЕНЕНИЙ ---

	// 7. Настройка и запуск HTTP-сервера
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler, // Используем итоговый handler
	}

	go func() {
		log.Printf("Сервер запускается на порту %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// 8. Graceful shutdown (корректное завершение работы)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Сервер останавливается...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	log.Println("Сервер успешно остановлен")
}
