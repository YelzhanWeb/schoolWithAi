FROM python:3.10-slim

# Рабочая директория
WORKDIR /app

# Копируем requirements.txt
COPY requirements.txt .

# Устанавливаем зависимости
RUN pip install --no-cache-dir -r requirements.txt

# Копируем приложение
COPY app ./app

# Порт
EXPOSE 5000

# Запуск
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "5000"]