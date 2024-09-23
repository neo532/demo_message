package biz

import (
	"context"

	"demo_message/internal/biz/entity"
	"demo_message/internal/biz/repo"
)

type CampaignUsecase struct {
	tx repo.TransactionMessageRepo
	cp repo.CampaignRepo
}

func NewCampaignUsecase(
	tx repo.TransactionMessageRepo,
	cp repo.CampaignRepo,
) *CampaignUsecase {
	return &CampaignUsecase{
		tx: tx,
		cp: cp,
	}
}

func (uc *CampaignUsecase) Create(c context.Context, cp *entity.Campaign) (insID int64, err error) {

	if insID, err = uc.cp.Create(c, cp); err != nil {
		return
	}

	return
}

func (uc *CampaignUsecase) UpdateStatus(c context.Context, ID int64, status int) (err error) {
	return uc.cp.UpdateStatus(c, ID, status)
}
