package payment

import (
	"go_api_foundation/user"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type service struct{}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = ""
	midclient.ClientKey = ""
	midclient.APIEnvType = midtrans.Sandbox

	coreGateway := midtrans.CoreGateway{
		Client: midclient,
	}

	chargeReq := &midtrans.ChargeReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	resp, err := coreGateway.Charge(chargeReq)
	if err != nil {
		return "", err
	}

	return resp.ReURL, nil

}
