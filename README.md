# ToDo App

---

## Информация об сервисе
В репозитории будет 2 ветки:
 - main - в этой вертке сервис строго соблюдает методы ис сценарии, указанные в ТЗ
 - second - в этой ветке сервис не будет строго соблюдать ТЗ,
реализую модель пользователей, регистрацию, JWT и еще чего-нибудь(я очень хочу найти работу),
чтобы продемострировать свои навыки.
Отличия будут указаны в ветке. 

#### Прошу рассмотреть именно вторую ветку (если успею реализовать запланированное)

####  Связь со мной: https://t.me/BigBrainlittle

В этом проекте также реализованны:
- UNIT-тестирование эндпойнтов с mock services
- Миграции при запуске
- Swagger
- docker, docker compose, k8s manifest

#### API и модели данных описанны в Swagger

## START PROJECT
- **make app_dev.env and db_dev.env from example `example.env` file**

### Local:
- **install dependencies:**
```
go mod tidy
```
- **run project:**
```
go run ./cmd/main.go 
```
- **swagger available on http://localhost:8082/swagger/index.html **

### Docker (only app):
- **build image**
```
docker build -t obuhovskaia11/todoverba:latest .
```
- **or pull from Docker Hub**
```
docker pull obuhovskaia11/todoverba:latest
```
- **run container**
```
docker run -d --name serverToDo -p 8082:8082 --env-file app_dev.env obuhovskaia11/todoverba:latest
```

### Docker compose(app and database):
- **run docker-compose**
```
docker-compose rm --build
```

### Kubernetes (app and service):

- **copy manifests from `/k8s/`**
- **create database resource:**
```
kubectl create -f ./database/db-config.yaml
kubectl create -f ./database/db-deployment.yaml
kubectl create -f ./database/db-service.yaml
```
- **create app resources:**
- **set APP_HOST="your_node_external_ip" env in `app-config.yaml`**
```
kubectl create -f ./app/app-config.yaml
kubectl create -f ./app/app-deployment.yaml
kubectl create -f ./app/app-service.yaml
```
