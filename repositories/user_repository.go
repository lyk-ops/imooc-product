package repositories

import (
	"database/sql"
	"errors"
	"imooc-product/common"
	"imooc-product/datamodels"
	"strconv"
)

type IUserRepository interface {
	Conn() error
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userID int64, err error)

	//
}
type UserManagerRepository struct {
	table     string
	myslqConn *sql.DB
}

func (u *UserManagerRepository) Conn() error {
	if u.myslqConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.myslqConn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}
	return nil
}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return nil, err
	}
	if userName == "" {
		return &datamodels.User{}, errors.New("用户名不能为空")
	}
	sqlStr := "select * from " + u.table + " where userName = ?"
	rows, err := u.myslqConn.Query(sqlStr, userName)
	defer rows.Close()
	if err != nil {
		return &datamodels.User{}, err
	}
	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return user, nil
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userID int64, err error) {
	if err = u.Conn(); err != nil {
		return 0, err
	}
	sqlStr := "insert into " + u.table + "set nickName=?,userName=?,password=?"
	stmt, err := u.myslqConn.Prepare(sqlStr)
	if err != nil {
		return userID, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if err != nil {
		return userID, err
	}
	userID, err = result.LastInsertId()
	return userID, err
}

func (u *UserManagerRepository) SelectByID(userID int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sqlStr := "select * from " + u.table + " where id=?" + strconv.FormatInt(userID, 10)
	rows, err := u.myslqConn.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		return &datamodels.User{}, err
	}
	row := common.GetResultRow(rows)
	if len(row) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(row, user)
	return user, nil
}
func NewUserRepository(table string, db *sql.DB) IUserRepository {
	return &UserManagerRepository{}
}
