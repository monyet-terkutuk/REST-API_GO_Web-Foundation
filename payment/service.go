package payment

import (
	"go_api_foundation/transaction"
	"go_api_foundation/user"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
	transactionRepository transaction.Repository
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
	ProcessPayment(input transaction.TransactionNotificationInput) error
}

func NewService(transactionRepository transaction.Repository) *service {
	return &service{transactionRepository}
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

func (s *service) ProcessPayment(input transaction.TransactionNotificationInput) (Transaction,error) {
	transaction_id, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.transactionRepository.GetByID(transaction_id)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	}

	if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	}

	if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel"{\
		transaction.Status = "cancelled"
	}

	updatedTransaction,err := s.transaction.Repository.Update(transaction)
	if err != nil {
		return err
	}
}
