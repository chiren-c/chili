package dao

import "database/sql"

type UserDAO interface {
}

type User struct {
	// 这三个字段表达为 sql.NullXXX 的意思，
	// 就是希望使用的人直到，这些字段在数据库中是可以为 NULL 的
	// 这种做法好处是看到这个定义就知道数据库中可以为 NULL，坏处就是用起来没那么方便
	// 大部分公司不推荐使用 NULL 的列
	// 所以你也可以直接使用 string, int64，那么对应的意思是零值就是每填写
	// 这种做法的好处是用起来好用，但是看代码的话要小心空字符串的问题
	Id       int64          `gorm:"primary_key,auto_increment"`
	Email    sql.NullString `gorm:"unique;comment:邮箱"`
	Password string         `gorm:"comment:密码（加密后）"`
	Phone    sql.NullString `gorm:"unique;comment:手机号"`
	Birthday sql.NullInt64  `gorm:"comment:生日"`
	NickName sql.NullString `gorm:"comment:昵称"`
	AboutMe  sql.NullString `gorm:"type:varchar(1024);comment:简介"`
	Status   uint8          `gorm:"comment:状态:1正常;default:1;not null"`
	// 创建时间
	Ctime int64 `gorm:"comment:创建时间"`
	// 更新时间
	Utime int64 `gorm:"comment:更新时间"`
}

// TableName UserM's table name
func (*User) TableName() string {
	return "user"
}
