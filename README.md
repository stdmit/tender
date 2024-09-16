# Сервис проведения тендеров
## В сервисе реализован следующий функционал:

 GET    /api/ping                 --> main.main.PingServer.func1 (2 handlers)
POST   /api/tenders/new          --> main.main.CreateTender.func2 (2 handlers)
 GET    /api/tenders/my           --> main.main.ListMyTenders.func3 (2 handlers)
 GET    /api/tenders              --> main.main.ListTenders.func4 (2 handlers)
 GET    /api/tenders/:tenderId/status --> main.main.ShowStatusTender.func5 (2 handlers)
 PUT    /api/tenders/:tenderId/status --> main.main.ChangeStatusTender.func6 (2 handlers)
 PATCH  /api/tenders/:tenderId/edit --> main.main.EditTender.func7 (2 handlers)
 PUT    /api/tenders/:tenderId/rollback/:ver --> main.main.RollbackVerTender.func8 (2 handlers)
 POST   /api/bids/new             --> main.main.CreateBid.func9 (2 handlers)
 GET    /api/bids/my              --> main.main.ListMyBids.func10 (2 handlers)
 GET    /api/bids/:Id/list        --> main.main.ListTenderBids.func11 (2 handlers)
 GET    /api/bids/:Id/status      --> main.main.ShowStatusBid.func12 (2 handlers)
 PUT    /api/bids/:Id/status      --> main.main.ChangeStatusBid.func13 (2 handlers)
 PATCH  /api/bids/:Id/edit        --> main.main.EditBid.func14 (2 handlers)
 PUT    /api/bids/:Id/rollback/:ver --> main.main.RollbackVerBid.func15 (2 handlers)
 PUT    /api/bids/:Id/feedback    --> main.main.BidFeedback.func16 (2 handlers)
 GET    /api/bids/:Id/reviews     --> main.main.BidReviews.func17 (2 handlers)


Тестирование проводилось при помощи утилит Postman и Swagger.

Комманды для сборки и запуска:

docker build -t myapp:latest .

docker run -d -p 8080:8080 myapp:latest