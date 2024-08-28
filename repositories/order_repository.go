package repositories

import (
	"database/sql"
	"imooc-product/common"
	"imooc-product/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
	//
}
type OrderManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (o *OrderManagerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		} else {
			o.mysqlConn = mysql
		}
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (productID int64, err error) {
	if err = o.Conn(); err != nil {
		return 0, err
	}
	sql := "INSERT " + o.table + "SET userID=?,productID=?,orderStatus=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(order.UserId, order.ProductID, order.OrderStatus)
	if err != nil {
		return productID, err
	}
	productID, err = result.LastInsertId()
	return productID, err
}

func (o *OrderManagerRepository) Delete(i int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}
	sql := "DELETE FROM " + o.table + " WHERE id = ?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(i)
	if err != nil {
		return false
	}
	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) error {
	if err := o.Conn(); err != nil {
		return err
	}
	sql := "UPDATE " + o.table + " SET userID=?,productID=?,orderStatus=? WHERE id = ?" + strconv.FormatInt(order.ID, 10)
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.UserId, order.ProductID, order.OrderStatus)
	return err
}

func (o *OrderManagerRepository) SelectByKey(i int64) (*datamodels.Order, error) {
	if err := o.Conn(); err != nil {
		return &datamodels.Order{}, err
	}
	sql := "SELECT * FROM " + o.table + " WHERE id = ?" + strconv.FormatInt(i, 10)
	row, err := o.mysqlConn.Query(sql)
	if err != nil {
		return &datamodels.Order{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}
	order := &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return order, nil
}

func (o *OrderManagerRepository) SelectAll() ([]*datamodels.Order, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT * FROM " + o.table
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}
	orders := make([]*datamodels.Order, 0)
	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orders = append(orders, order)
	}
	return orders, nil
}

func (o *OrderManagerRepository) SelectAllWithInfo() (map[int]map[string]string, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT o.ID,p.productName,o.orderStatus From imooc.order as o left join product as p on o.productID=p.ID"
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	resultRows := common.GetResultRows(rows)
	return resultRows, err

}

func NewOrderManagerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderManagerRepository{table: table, mysqlConn: sql}
}
