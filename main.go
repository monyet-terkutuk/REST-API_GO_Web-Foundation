package main

import (
	"fmt"
	"go_api_foundation/auth"
	"go_api_foundation/campaign"
	"go_api_foundation/handler"
	"go_api_foundation/helper"
	"go_api_foundation/user"
	"log"
	"net/http"
	"strings"

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

	// memaping data dari input user untuk di masukan ke userRepository
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)

	// panggil service auth
	authService := auth.NewService()

	// test
	input := campaign.CreateCampaignInput{}

	input.Name = "Penggalangan Dana Banjir Cibaduyut"
	input.ShortDescription = "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Asperiores aliquam perspiciatis optio quod obcaecati possimus? Blanditiis possimus ipsam cumque culpa saepe commodi, fugiat neque quis quia minima, temporibus praesentium in."
	input.Description = "Lorem, ipsum dolor sit amet consectetur adipisicing elit. Aspernatur consectetur ad nam ex consequuntur, dicta, deleniti necessitatibus unde consequatur temporibus debitis voluptatem omnis ullam dolorum ea aperiam natus totam ratione? Esse cum sed molestiae fuga corporis! Reprehenderit rerum nam minus numquam ducimus accusamus ratione iste id cupiditate, rem dolorum fugit ut blanditiis tempora suscipit qui asperiores aspernatur assumenda, excepturi hic sint enim porro impedit? Dolore dignissimos accusamus eos debitis quae similique, numquam quod odio saepe commodi, dolores excepturi culpa tempore pariatur laboriosam sit inventore amet laudantium, aut error aliquam harum. Expedita, numquam cum iusto reiciendis voluptate iste ad eos molestias fugit rem odit quas, magnam quae, voluptatem cupiditate laboriosam reprehenderit recusandae laborum eum pariatur est sit id? Officiis, quis provident! Alias at suscipit, hic vel eaque blanditiis! Vitae velit esse delectus dignissimos magni maiores tempore, nisi exercitationem? Natus sapiente numquam dolores quidem laborum quia architecto laudantium, vel illo minus, vero doloribus velit tempora. Modi, sapiente suscipit illum amet eaque impedit eum provident explicabo quis obcaecati. Sapiente explicabo nihil eveniet consequuntur temporibus maiores atque a ducimus ut enim suscipit adipisci, reprehenderit fugit distinctio alias assumenda ea, dolorem nulla neque eos nam. Explicabo, veniam esse! Molestias, qui numquam quidem, perspiciatis ab cum impedit fugit atque ut eaque rerum provident incidunt illum et?"
	input.GoalAmount = 100000000
	input.Perks = "Pahala,Kebaikan,Ketenangan"

	inputUser, _ := userService.GetUserByID(7)
	input.User = inputUser

	_, err = campaignService.CreateCampaign(input)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err.Error())
	}

	// mengambil data mentah dari client dan di convert dari json ke userService
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()

	// buat routing static untuk gambar
	router.Static("/avatar", "./images/avatar")

	// grouping router
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.UserLogin)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleaware(authService, userService), userHandler.UploadAvatar)
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaign/:id", campaignHandler.GetCampaign)

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
