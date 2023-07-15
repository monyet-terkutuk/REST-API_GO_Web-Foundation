package main

import (
	"go_api_foundation/auth"
	"go_api_foundation/campaign"
	"go_api_foundation/handler"
	"go_api_foundation/helper"
	"go_api_foundation/payment"
	"go_api_foundation/transaction"
	"go_api_foundation/user"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// bisa buat logging dan helper
	dsn := "root:@tcp(127.0.0.1:3306)/go_foundation?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	// menghubungkan struct ke database
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	// memaping data dari input user untuk di masukan ke userRepository
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	// panggil service auth
	authService := auth.NewService()

	// mengambil data mentah dari client dan di convert dari json ke userService
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := gin.Default()

	// set cors
	router.Use(cors.Default())

	// buat routing static untuk gambar
	router.Static("/avatar", "./images/avatar")

	// grouping router
	api := router.Group("/api/v1")

	// user endpoint
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.UserLogin)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleaware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleaware(authService, userService), userHandler.FetchUser)
	// campaign endpoint
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleaware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleaware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleaware(authService, userService), campaignHandler.UploadImage)
	// transaction endpoint
	api.GET("/campaign/:id/transactions", authMiddleaware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleaware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleaware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.Run()
}

func authMiddleaware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil isi header authorization
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		// ambil token nya saja
		var tokenString string
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		// validasi token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		// ambil user id dan ambil user dengan service

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// set context
		c.Set("currentUser", user)

	}

}
