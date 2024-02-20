# proxy_manager
Сервис управления проксями.

Функции:
- Хранит прокси;
- Отдает список проксей;
- Контролирует сколько клиентов на данный момент использует конкретную проксю;

Методы /api/v1:
- GET /proxies - возвращает список проксей:
    - Пагинация с помощью параметров offset и limit;
- POST /proxies - добавляет новую проксю;
- GET /proxies/:proxy_id/ - получение инфы по конкретной проксе;
- UPDATE /proxies/:proxy_id - обновление инфы о проксе;
- DELETE /proxies/:proxy_id - удаление прокси;
- POST /proxies/occupy - занять свободную проксю;
- POST /proxies/release - освободить проксю;

На /api/v1/swagger/index.html есть swagger.

TODO:
- [ ] Gin логирование;
- [ ] GRPC;
- [ ] Dockerfile && docker-compose.yaml
- [ ] ProxyCheck, на уровне repo
- [ ] фильтр на GetProxyList{enabled=[true, false]}
# cmds
## migrate
```
migrate -database postgresql://login:pass@ip:port/db?sslmode=disable -path ./migrations up/down/force
```
## swag
```
swag init -g internal/controller/http/v1/router.go
```

## golangci-lint
```
golangci-lint run
```