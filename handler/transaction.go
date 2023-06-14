package handler

import (
	"go_api_foundation/helper"
	"go_api_foundation/transaction"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	transactions, err := h.service.GetTransactionsByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List campaign transaction", http.StatusOK, "success", transactions)
	c.JSON(http.StatusOK, response)
}

// ------- Handler -----------
// parameter dari uri
// tangkap parameter dan mapping ke struct input
// panggil service, struct di input menjadi paraameter nya

// ------- Service ---------
// Service memanggil repository dengan campaign id dari handler

// ----- Repository -------
// repo mencari  data trnasaction dari suatu campaign
