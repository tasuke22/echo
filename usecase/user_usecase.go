package usecase

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/tasuke/udemy/model"
	"github.com/tasuke/udemy/repository"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (string, error)
}

// この構造体はなにに依存すべきか考える => 具象ではなく抽象に依存することによって疎結合となりテストしやすい
type userUsecase struct {
	ur repository.IUserRepository
}

// IUserRepositoryを依存性の注入して、IUserUsecaseを返す
// &userUsecase{}としてreturnするとどうなる？ => userUsecaseインスタンスはリポジトリとのやり取りができないため、データベースにアクセスする機能が失われます。
// つまり、依存性の注入ができていないことになる。
func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	return &userUsecase{ur}
}

// ポインタレシーバになっている意味はなに？ => 構造体の状態の変更を反映させる。厳密には、SignUpメソッドがuserUsecaseの状態を変更する必要はないかもしれませんが、内部で保持しているリポジトリ（uu.ur）への参照を通じてデータベース操作を行うために、ポインタレシーバが使われています。
func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	// パスワードを平文からハッシュ化している
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		// なぜゼロ値を返すのか => もしエラーだとしても返り値を成功した場合と一貫させることによって呼び出し元のエラーハンドリングがしやすい。
		return model.UserResponse{}, err
	}
	// &newUserとすることで参照渡しとしてポインタを渡しているわけなので、newUserの値も変更されている
	newUser := model.User{Email: user.Email, Password: string(hash)}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	// わかりやすいようにレスポンスの形の構造体を事前に定義していた。DTO。
	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	return resUser, nil
}

func (uu *userUsecase) Login(user model.User) (string, error) {
	storedUser := model.User{}
	// 登録されているユーザーをEmailで検索する
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	// パスワードの検証：クライアントのパスワードと、DBに保存されているパスワードを比較する
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	// jwtトークンの生成準備
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	// jwtトークンを実際に生成
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, err
}
