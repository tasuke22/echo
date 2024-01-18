package repository

import (
	"github.com/tasuke/udemy/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}

// Go言語では、インターフェースの実装は明示的に宣言されません。
// userRepositoryはdbに依存している
type userRepository struct {
	db *gorm.DB
}

// userRepository構造体は、IUserRepositoryインターフェースで定義された全てのメソッド（関数）を持っているため、IUserRepositoryとして返すことができます。
// これにより、この関数を使用するコードは、どの具体的な構造体（userRepository）が背後にあるかを知る必要はなく、インターフェースを通じて定義されたメソッドにアクセスできます。
// NewUserRepository関数は、userRepository構造体にデータベース接続（db *gorm.DB）を渡します。これは依存性注入の一例です。
// なんで&userRepository{db}アドレスを指定して返さなきゃいけないの？ => 値(コピー)として渡されると、メソッドはその構造体のコピー上で操作され、元の構造体には影響しなくなってしまうため。
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	// .Errorはなに？ => この部分は、上記のクエリ操作がエラーを発生させたかどうかをチェックして、エラーがあればエラーを返す。その後のエラーハンドリングはGoの慣習である。
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	// なぜnilを返さないといけないの？ => Goの慣習であり、使用している側でのエラーハンドリングを考えればnilを返すべきなのは想像がつく。
	return nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
