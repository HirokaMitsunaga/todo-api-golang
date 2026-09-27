package repository

import (
	"errors"
	"todo-api/domain"
	"todo-api/domain/repository"
	"todo-api/infra/gorm/model"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) FindById(id ulid.ULID) (*domain.User, error) {
	var user model.User

	if err := ur.db.Where("id=?", id.String()).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return domain.Reconstructor(
		id, //値が見つからなかった場合errorが返るため、引数のidをそのまま使う。
		user.Name,
		user.Email,
		user.Password,
	)
}
func (ur *userRepository) Create(user *domain.User) error {

	record := model.User{
		ID:       user.ID().String(),
		Name:     user.Name(),
		Email:    user.Email(),
		Password: user.Password(),
	}

	if err := ur.db.Create(&record).Error; err != nil {
		//TODO:どのフィールドで重複しているのかわかるようにする
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrUserDuplicated
		}
	}
	return nil
}

func (ur *userRepository) Update(user *domain.User) error {

	record := model.User{
		ID:       user.ID().String(),
		Name:     user.Name(),
		Email:    user.Email(),
		Password: user.Password(),
	}

	if err := ur.db.Updates(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrUserDuplicated
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
	}
	return nil
}

func (ur *userRepository) Delete(id ulid.ULID) error {
	result := ur.db.
		Where("id = ?", id.String()).
		Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil

}

var ErrUserNotFound = errors.New(
	"user not Found",
)

var ErrUserDuplicated = errors.New(
	"user is Duplicated",
)
