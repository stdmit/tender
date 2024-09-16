# Сервис проведения тендеров
## В сервисе реализован следующий функционал:

 GET    /api/ping
 
 POST   /api/tenders/new
 
 GET    /api/tenders/my
 
 GET    /api/tenders
 
 GET    /api/tenders/:tenderId/status
 
 PUT    /api/tenders/:tenderId/status
 
 PATCH  /api/tenders/:tenderId/edit
 
 PUT    /api/tenders/:tenderId/rollback/:ver

 POST   /api/bids/new

 GET    /api/bids/my

 GET    /api/bids/:Id/list

 GET    /api/bids/:Id/status

 PUT    /api/bids/:Id/status

 PATCH  /api/bids/:Id/edit

 PUT    /api/bids/:Id/rollback/:ver

 PUT    /api/bids/:Id/feedback

 GET    /api/bids/:Id/reviews
 


Тестирование проводилось при помощи утилит Postman и Swagger.

Комманды для сборки и запуска:

docker build -t myapp:latest .

docker run -d -p 8080:8080 myapp:latest