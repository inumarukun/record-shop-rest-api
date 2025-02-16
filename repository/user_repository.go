package repository

import (
	"record-shop-rest-api/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(user *model.User, email string) error
}

type userRepository struct {
	db *gorm.DB
}

// constructor
func NewUserRepository(db *gorm.DB) IUserRepository {
	// {}はGoにおいて構造体の初期化を行う構文
	// recordRepository構造体↑はdbという*gorm.DB型のフィールドを持っており
	// この構造体を作成するには、そのフィールドに値を設定する必要がある
	// インターフェースは構造体のポインタ型で実装されることが多い
	// インターフェース内のメソッドリストを全て実装する必要がある
	return &userRepository{db}
}

func (ur *userRepository) CreateUser(user *model.User) error {
	// Gormはこのように引数を直接変更する、.Errorにチェインさせるため
	// 引数変えていいのか、違和感があるが変えたくなければ、recordRrepositoryの例見て
	// 引数受け取らず、戻り値で([]model.Record, error) を返してる
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}
