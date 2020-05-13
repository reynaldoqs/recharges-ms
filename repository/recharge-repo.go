package repository

import "audiman/projects/recharger-x/entity"

// RechargeRepository model for rechages repository
type RechargeRepository interface {
	Save(recharge *entity.Recharge) (*entity.Recharge, error)
	//Update(recharge *entity.Recharge) error
	UpdateRecharger(rmsg *entity.RechargeMsg) error
	FindAll() ([]*entity.Recharge, error)
}
