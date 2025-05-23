# Запуск

`go run cmd/main.go`

## Запросы для ручек

Создать:
<pre>
curl -X POST http://localhost:8080/quotes \
     -H "Content-Type: application/json" \
     -d '{"author":"Confucius", "text":"Life is simple, but we insist on making it complicated."}'
</pre>

Получить все записи:
<pre>
curl http://localhost:8080/quotes
</pre>

Получить все записи по автору:
<pre>
curl http://localhost:8080/quotes?author=Confucius
</pre>

Получить рандомную запись:
<pre>
curl http://localhost:8080/quotes/random
</pre>

Удалить по ID:
<pre>
curl -X DELETE http://localhost:8080/quotes/1
</pre>

# Тесты

Они есть, но только на основной логике
<pre>
$ go test -cover ./...
quotes/cmd coverage: 0.0% of statements
quotes/config coverage: 0.0% of statements
quotes/internal/controller coverage: 0.0% of statements
? quotes/internal/entity [no test files] 
quotes/internal/repository coverage: 0.0% of statements
quotes/internal/server  coverage: 0.0% of statements       
ok quotes/internal/service 0.360s  coverage: 88.1% of statements      
? quotes/pkg/error [no test files]
quotes/pkg/logs  coverage: 0.0% of statements
quotes/pkg/middleware coverage: 0.0% of statements       
quotes/pkg/utils coverage: 0.0% of statements 
</pre>