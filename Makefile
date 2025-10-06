# ==== Makefile ====

# Пути к скриптам
PYTHON_SERVER = ai-service/app.py
PYTHON_ANALYZE = ai-service/analyze.py
GO_REC = go_recommender/main.go

# ===== Команды =====

# Запуск Python-сервера
server:
	@echo "🚀 Запуск Python-сервера..."
	python3 $(PYTHON_SERVER)

# Запуск анализа данных
analyze:
	@echo "📊 Запуск анализа данных..."
	python3 $(PYTHON_ANALYZE)

# Запуск Go-рекомендателя
rec:
	@echo "🤖 Запуск Go-рекомендательной системы..."
	go run $(GO_REC)

# Очистка кеша/файлов (опционально)
clean:
	@echo "🧹 Очистка временных файлов..."
	find . -type d -name "__pycache__" -exec rm -rf {} +

# Помощь
help:
	@echo "📘 Доступные команды:"
	@echo "  make server   — запустить Python сервер (ai-service/app.py)"
	@echo "  make analyze  — выполнить анализ (ai-service/analyze.py)"
	@echo "  make rec      — запустить Go рекомендатель (go_recommender/main.go)"
	@echo "  make clean    — очистить временные файлы"
