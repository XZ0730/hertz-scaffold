package db

import "time"

type UserInfo struct {
	Id        int64
	Telephone string
	Name      string
	Age       int
	Sex       string
	CardId    string
	Account   string
	Password  string
	Role      int
}

func NewUserInfo() *UserInfo {
	return &UserInfo{}
}

// 判断账号是否已经注册
func JudgeAccount(account string) error {
	ui := NewUserInfo()
	return _db.Table("user_info").Where("account=?", account).First(ui).Error
}

// 判断身份证是否被绑定 - 普通用户
func JudgeIDCard(card_id string) error {
	ui := NewUserInfo()
	return _db.Table("user_info").Where("card_id=? AND role =?", card_id, 0).First(ui).Error
}

// 判断手机号是否被绑定 - 普通用户
func JudgePhoneNumber(phone string) error {
	ui := NewUserInfo()
	return _db.Table("user_info").Where("telephone=? AND role =?", phone, 0).First(ui).Error
}

// 更新用户个人信息
func UpdateUserInfo(user_info *UserInfo) error {
	return _db.Table("user_info").Updates(&user_info).Error
}

func GetUserInfoByUserName(user_name string) (*UserInfo, error) {
	user := new(UserInfo)
	err := _db.Table("user_info").Where("account=?", user_name).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(username, pwd string) error {
	ui := NewUserInfo()
	ui.Account = username
	ui.Password = pwd
	return _db.Table("user_info").Create(&ui).Error
}

func JudgeViolation(uid int64) bool {
	var cnt int64
	_db.Table("user_info").Where("id=? AND violation=1", uid).Count(&cnt)
	return cnt == 1
}

func JudgeParkCnt(vnumber string, gid string) bool {
	var cnt int64
	now := time.Now().UnixNano()
	now -= (time.Hour.Nanoseconds())
	first := time.Unix(now, 0).Format(time.DateTime)
	second := time.Unix(now+2*time.Hour.Nanoseconds(), 0).Format(time.DateTime)
	_db.Table("park_record").Where("vehicle_number=? AND garage_id=? AND entries>? AND entries<?", vnumber, gid, first, second).Count(&cnt)
	return cnt >= 3
}

func UpdateViolation(uid int64) error {
	return _db.Table("user_info").Where("id=?", uid).UpdateColumn("violation", 1).Error
}
