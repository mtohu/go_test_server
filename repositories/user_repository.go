package repositories

import (
	"gotest/common"
	"gotest/datamodels"

)



// Query represents the visitor and action queries.
type Query func(datamodels.User) bool

type UserRepository struct {
	BaseRep IBaseRepository

}



// UserRepository handles the basic operations of a user entity/model.
// It's an interface in order to be testable, i.e a memory user repository or
// a connected to an sql database.
type IUserRepository interface {
	//Exec(query Query, action Query, limit int, mode int) (ok bool)

	Select(query Query) (user datamodels.User, found bool)
	//SelectMany(query Query, limit int) (results []datamodels.User)

	//InsertOrUpdate(user datamodels.User) (updatedUser datamodels.User, err error)
	//Delete(query Query, limit int) (deleted bool)
}

func NewUserRepository() IUserRepository {
	var base = NewBaseRepository()
	return &UserRepository{BaseRep:base}
}

func (ur *UserRepository) Select(query Query) (datamodels.User, bool) {
	sqls:=ur.BaseRep.Table("admin").Fields("*").Where("1=1").FetchSql(1)
	common.Log("ssssss===qqqq==",sqls)
    return datamodels.User{},true
}


