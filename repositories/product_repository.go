package repositories

import (
	"database/sql"
	"errors"
	"imooc-product/common"
	"imooc-product/datamodels"
	"strconv"
)

// 第一步 先开发对应接口

// 第二步 实现对应接口
type IProduct interface {
	Conn() error
	Insert(product *datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(product *datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}
type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

// 初始化连接
func (p *ProductManager) Conn() error {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		} else {
			p.mysqlConn = mysql
		}
		// 设置默认表名称
		if p.table == "" {
			p.table = "product"
		}
	}
	return nil
}

// 插入数据方法
func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	// 判断是否连接成功
	if err := p.Conn(); err != nil {
		return 0, err
	}
	sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	// 执行插入语句
	result, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return 0, err
	}
	// 获取插入数据的主键id
	productId, err = result.LastInsertId()
	if err != nil {
		return 0, err
	} else {
		return productId, nil
	}
}

// 删除数据方法
func (p *ProductManager) Delete(i int64) bool {
	// 判断是否连接成功
	if err := p.Conn(); err != nil {
		return false
	}
	sql := "DELETE FROM product WHERE id=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(i)
	if err != nil {
		return false
	} else {
		return true
	}
}

// 更新数据方法
func (p *ProductManager) Update(product *datamodels.Product) error {
	// 判断是否连接成功
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "UPDATE product SET productName=?,productNum=?,productImage=?,productUrl=? WHERE id=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, strconv.FormatInt(product.ID, 10))
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductManager) SelectByKey(i int64) (*datamodels.Product, error) {
	// 判断是否连接成功
	if err := p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "SELECT * FROM" + p.table + "WHERE ID=?" + strconv.FormatInt(i, 10)
	row, err := p.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return &datamodels.Product{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, errors.New("查询数据为空")
	}
	//映射数据库中的数据到结构体中并且转换类型
	common.DataToStructByTagSql(result, &datamodels.Product{})
	return &datamodels.Product{}, nil
}

// 查询所有数据
func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, err error) {
	// 判断是否连接成功
	if err := p.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT * FROM" + p.table
	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, errors.New("查询数据为空")
	}
	// 映射数据库中的数据到结构体中并且转换类型
	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return productArray, nil
}

// NewProductManager 构造函数
func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}
