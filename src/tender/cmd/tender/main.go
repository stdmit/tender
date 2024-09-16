package main

import (
	// "fmt"
	"log"
	// "net/http"

	// "example.com/tender/internal/database"
	"example.com/tender/internal/controllers"
	"example.com/tender/internal/database"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)
const port = "8080"

func main() {

	router := gin.New()

	router.Use(gin.Logger())

	database.PsqlConnect()

	router.GET("/api/ping",controllers.PingServer())

	router.POST("/api/tenders/new", controllers.CreateTender())
	router.GET("/api/tenders/my",controllers.ListMyTenders())
	router.GET("/api/tenders",controllers.ListTenders())	
	router.GET("/api/tenders/:tenderId/status",controllers.ShowStatusTender())
	router.PUT("/api/tenders/:tenderId/status",controllers.ChangeStatusTender())
	router.PATCH("/api/tenders/:tenderId/edit",controllers.EditTender())
	router.PUT("/api/tenders/:tenderId/rollback/:ver",controllers.RollbackVerTender())

	router.POST("/api/bids/new", controllers.CreateBid())
	router.GET("/api/bids/my", controllers.ListMyBids())
	router.GET("/api/bids/:Id/list", controllers.ListTenderBids())
	router.GET("/api/bids/:Id/status",controllers.ShowStatusBid())
	router.PUT("/api/bids/:Id/status",controllers.ChangeStatusBid())
	
	router.PATCH("/api/bids/:Id/edit",controllers.EditBid())
	router.PUT("/api/bids/:Id/rollback/:ver",controllers.RollbackVerBid())
	
	router.PUT("/api/bids/:Id/feedback",controllers.BidFeedback())
	router.GET("/api/bids/:Id/reviews",controllers.BidReviews())

	
	log.Fatal(router.Run(":" + port))


}