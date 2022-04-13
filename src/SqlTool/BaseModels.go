package sqltool

import "time"

type Judge_Base struct {
	ID            uint      `` //评测编号
	Problem       Problem   `` //
	PID           uint      `` //题目编号
	PTitle        string    `` //题目标题
	PShow_id      string    `` //题目显示ID
	User_info     User_info `` //
	UID           uint      `` //用户编号
	Submit_time   time.Time `` //提交时间
	Status        string    `` //评测状态(Pending、Compiling、Running、AC、CE、RE、WA、TLE、MLE、OLE)
	Share_code    bool      `` //是否共享代码
	Error_message string    `` //错误信息
	Time_use      int       `` //耗时，单位ms
	Memory_use    int       `` //使用内存，单位MB
	Length        int       `` //代码长度
	Code          string    `` //代码
	Language      string    `` //语言
	Judger        string    `` //评测机IP
	Ip            string    `` //用户IP
	CID           uint      `` //比赛ID
	Version       uint      `` //评测版本
	CreatedAt     time.Time `` //创建时间
	UpdatedAt     time.Time `` //修改时间
}
