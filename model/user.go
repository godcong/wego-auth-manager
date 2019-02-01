package model

import (
	"github.com/godcong/wego-auth-manager/util"
	"golang.org/x/exp/xerrors"
	"net/url"
)

// UserLogin ...
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User ...
type User struct {
	Model         `xorm:"extends" json:",inline"`
	Block         bool   `json:"block" xorm:"block"`                     //禁止访问
	Nickname      string `json:"nickname" xorm:"nickname"`               //名称
	Username      string `json:"username" xorm:"username"`               //用户名
	Email         string `json:"email" xorm:"email"`                     //邮件
	Mobile        string `json:"mobile" xorm:"mobile"`                   //移动电话
	IDCardFacade  string `json:"id_card_facade" xorm:"id_card_facade"`   //身份证(正)
	IDCardObverse string `json:"id_card_obverse" xorm:"id_card_obverse"` //身份证(反)
	Password      string `json:"password" xorm:"password"`               //密码
	Certificate   string `json:"certificate" xorm:"certificate"`         //证书
	PrivateKey    string `json:"private_key" xorm:"private_key"`         //私钥
	LoginIP       string `json:"login_ip" xorm:"login_ip"`               //本次登录IP
	Token         string `json:"token" xorm:"token"`                     //Token
	Salt          string `json:"-" xorm:"slat"`                          //盐值
}

// NewUser ...
func NewUser(id string) *User {
	return &User{Model: Model{
		ID: id,
	}}
}

// Get ...
func (obj *User) Get() (bool, error) {
	return Get(nil, obj)
}

// Update ...
func (obj *User) Update(cols ...string) (int64, error) {
	return Update(nil, obj.ID, obj)
}

// Paginate ...
func (obj *User) Paginate(v url.Values) (*Paginate, error) {
	return &Paginate{}, nil
}

// Users ...
func (obj *User) Users() ([]*User, error) {
	var users []*User
	err := DB().Table(obj).Find(&users)
	if err != nil {
		return nil, xerrors.Errorf("find: %w", err)
	}
	return users, nil
}

// Permissions ...
func (obj *User) Permissions() ([]*Permission, error) {
	var permissions []*Permission
	session := DB().Table(&Permission{}).Select("permission.*").
		Join("left", &PermissionUser{}, "permission_user.user_id = user.id")

	if obj.ID != "" {
		session = session.Where("user.id = ? ", obj.ID)
	}

	err := session.Find(&permissions)
	if err != nil {
		return nil, xerrors.Errorf("relate: %w", err)
	}

	return permissions, nil
}

// CheckPermission ...
func (obj *User) CheckPermission(funcName string) bool {
	session := DB().Table(&Permission{}).Select("permission.*").
		Join("left", &PermissionUser{}, "permission_user.user_id = user.id").
		Where("permission.slug = ?", funcName)

	if obj.ID != "" {
		session = session.Where("user.id = ? ", obj.ID)
	}

	b, err := session.Exist()
	if err != nil || !b {
		return false
	}
	return true
}

// Roles ...
func (obj *User) Roles() ([]*Role, error) {
	var roles []*Role
	session := DB().Table(&Role{}).Select("role.*").
		Join("left", &RoleUser{}, "role_user.user_id = user.id")

	if obj.ID != "" {
		session = session.Where("user.id = ? ", obj.ID)
	}

	err := session.Find(&roles)
	if err != nil {
		return nil, xerrors.Errorf("relate: %w", err)
	}

	return roles, nil
}

// CheckPermission ...
func (obj *User) Validate(u *UserLogin, key string) bool {
	u.Password = util.SHA256(u.Password, key, obj.Salt)
	session := DB().Table(obj).Where("username = ?", u.Username).And("password = ?", u.Password)

	b, err := session.Exist()
	if err != nil || !b {
		return false
	}
	return true
}
