package sqltool

import "time"

type Judge_Base struct {
	ID           uint      `` //评测编号
	Problem      Problem   `` //
	PID          uint      `` //题目编号
	PTitle       string    `` //题目标题
	PShowID      string    `` //题目显示ID
	UserInfo     UserInfo  `` //
	UID          uint      `` //用户编号
	SubmitTime   time.Time `` //提交时间
	Status       string    `` //评测状态(Pending、Compiling、Running、AC、CE、RE、WA、TLE、MLE、OLE)
	ShareCode    bool      `` //是否共享代码
	ErrorMessage string    `` //错误信息
	TimeUse      int       `` //耗时，单位ms
	MemoryUse    int       `` //使用内存，单位MB
	Length       int       `` //代码长度
	Code         string    `` //代码
	Language     string    `` //语言
	Judger       string    `` //评测机IP
	Ip           string    `` //用户IP
	CID          uint      `` //比赛ID
	Version      uint      `` //评测版本
	CreatedAt    time.Time `` //创建时间
	UpdatedAt    time.Time `` //修改时间
}
