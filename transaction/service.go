package transaction

import (
	"errors"
	"go_api_foundation/campaign"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}
type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {

	// dapatkan campaign
	// lalu cek yg akses sama tidak dengan user id di campaign
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != input.User.Id {
		return []Transaction{}, errors.New("Not authorized to perform this action.")
	}

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
