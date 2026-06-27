package domain

// 例　User はたぶん使わない気がする。これを参考に他のpostとかを作る

type User struct {
	// フィールドはprivate
	name string
}

// コンストラクタをNew*みたいな感じで用意する
func NewUser(name string) *User {
	return &User{
		name: name,
	}
}

// 取得はメソッドを作る
func (u *User) GetName() string {
	return u.name
}
