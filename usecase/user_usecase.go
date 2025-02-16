package usecase

import (
	"os"
	"record-shop-rest-api/model"
	"record-shop-rest-api/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	Login(user model.User) (string, error)
	SignUp(user model.User) (model.UserResponse, error)
}

type userUsecase struct {
	ur repository.IUserRepository
}

func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	return &userUsecase{ur}
}

func (uu *userUsecase) Login(user model.User) (string, error) {
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	// テーブルのpwと、入力されたpw比較
	// bcryptでハッシュ化されたパスワードと、それに相当する可能性のある平文とを比較
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	// JSON Web Token生成
	// トークンに含める「クレーム」を指定する必要がある
	// 第1引数でJWTの署名アルゴリズムを指定
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// これらのれらの情報をペイロードというらしい
		// JWTのペイロードに入れる情報は、認証や認可に必要な最小限の情報に絞ることが望ましい
		// user_id: ユーザー認証を行うシステムは、JWTが発行されたユーザーを識別するために
		// user_id をペイロードに入れるのが一般的。
		// これにより、後でトークンを検証した際に、どのユーザーが認証されたのかを特定することが出来る
		// exp: トークンの有効期限を指定するための標準的なクレーム
		// これがないと、JWTが無期限に有効なトークンになってしまい、セキュリティリスクが増す
		// したがって、多くのアプリケーションでは exp を設定して、有効期限を過ぎたトークンを無効化する
		// 他に以下のようなClaimが使われることもある
		// iat (issued at): トークンが発行された日時を示す。サーバー側でトークンの発行日時を追跡するために使われる
		// aud (audience): トークンが意図する相手（例えば、特定のサービスやシステム）を示すためのクレーム
		// iss (issuer): トークンの発行者を示します。どのシステムがこのトークンを発行したかを識別する
		// scope: トークンに関連するアクセス権限（例えば、"read" や "write" など）を示すクレーム
		"user_id": storedUser.ID,
		// 現在時刻の12時間後をUNIXタイムスタンプで指定
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})
	// *jwt.Token.SignedString: JWTが改竄されていないことを保証するために、署名をトークンに追加
	// 署名は↑のjwt.SigningMethodHS256に基づいて行われる
	// os.Getenvで取得したSECRET(秘密鍵)を、[]byteに変換してSignedString()に渡す
	// tokenString は、認証や認可のためにクライアントに返されることが一般的
	// クライアントはこのトークンを使って、後続のリクエストでサーバーとやり取りを行う
	// JWTの署名に使用する「秘密鍵」は通常バイト列として扱う
	// ※受け手はinterface{}で任意の型を受けれるようになってる、なので型アサーションしてロジック書いてる
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		// エラー時は空のmodelを返す
		return model.UserResponse{}, err
	}
	// これで新しいstructが出来る、idはGormが自動で入れる？
	newUser := model.User{Email: user.Email, Password: string(hash)}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	// CreateUserが成功すれば、newUser、つまり引数が新しいユーザになっている、それを詰めて返す
	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	return resUser, nil
}
