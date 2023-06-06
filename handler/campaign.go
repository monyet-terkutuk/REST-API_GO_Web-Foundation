package handler

import (
	"go_api_foundation/campaign"
	"go_api_foundation/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
- tangkap params di handler
- handler ke service
- service memanggil repository sesuai kebutuhan
- fungsi di repository : GetAll, GetByUserIdD
- dari database
*/

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaigns)
	c.JSON(http.StatusOK, response)

}
