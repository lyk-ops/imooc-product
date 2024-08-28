package services

import (
	"imooc-product/datamodels"
	"imooc-product/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*datamodels.Order) error
	InsertOrder(*datamodels.Order) (int64, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	//
}
type OrderService struct {
	OrderReopository repositories.IOrderRepository
}

func (o *OrderService) GetOrderByID(i int64) (*datamodels.Order, error) {
	order, err := o.OrderReopository.SelectByKey(i)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *OrderService) DeleteOrderByID(i int64) bool {
	return o.OrderReopository.Delete(i)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) error {
	return o.OrderReopository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (int64, error) {
	insert, err := o.OrderReopository.Insert(order)
	if err != nil {
		return 0, err
	}
	return insert, nil
}

func (o *OrderService) GetAllOrder() ([]*datamodels.Order, error) {
	all, err := o.OrderReopository.SelectAll()
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	info, err := o.OrderReopository.SelectAllWithInfo()
	if err != nil {
		return nil, err
	}
	return info, nil
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService {
	return &OrderService{
		OrderReopository: repository,
	}
}
