# Этап 1: Сборка (Node.js)
FROM node:18-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
# ВАЖНО: Если у вас Vite, команда 'npm run build' создаст папку 'dist'.
# Если Create-React-App, то папку 'build'. Проверьте это!
RUN npm run build

# Этап 2: Раздача (Nginx)
FROM nginx:alpine

# Копируем собранные файлы в Nginx
# ВАЖНО: Если у вас Create-React-App, замените /app/dist на /app/build
COPY --from=build /app/dist /usr/share/nginx/html

# Создаем конфиг Nginx прямо на лету, чтобы работал React Router (перезагрузка страниц)
RUN echo 'server { \
    listen 80; \
    location / { \
    root /usr/share/nginx/html; \
    index index.html index.htm; \
    try_files $uri $uri/ /index.html; \
    } \
    }' > /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]