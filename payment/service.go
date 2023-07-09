package payment

import (
	"go_api_foundation/transaction"
	"go_api_foundation/user"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type service struct{}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
	ProcessPayment(input transaction.TransactionNotificationInput) error
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = ""
	midclient.ClientKey = ""
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	resp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return resp.RedirectURL, nil

}

func (s *service) ProcessPayment(input transaction.TransactionNotificationInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID)
}
