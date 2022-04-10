package sqltool

import (
	"time"
)

type Problem_Description struct {
	ID          uint   `gorm:"primaryKey"`             //题面编号
	Markdown    bool   `gorm:"not null;default:false"` //是否存储MD文件，为true的时候在Description字段存储MD值，否则由以下六个字段共同组合
	BackGround  string ``                              //题目背景
	Description string ``                              //题目描述
	Input       string ``                              //输入描述
	Output      string ``                              //输出描述
	Examples    string ``                              //题目样例，json格式，[{input: , output: ,descrption:},{input:, output: ,descrption:},...] ，允许有多个样例
	Hint        string ``                              //题目提示
	CaseFiles   string ``                              //测试数据文件名
}

type Problem struct {
	ID               uint                `gorm:"primaryKey"`                     //题目编号
	Judge_mode       byte                `gorm:"not null;default:0"`             //测试模式 (0:default，1:spj，2:interactive,3:Leetcode)
	Show_id          string              ``                                      //题目显示ID
	Title            string              `gorm:"not null;index:,class:FULLTEXT"` //题目标题
	Time_limit       int                 `gorm:"not null;default:1000"`          //题目时限，单位ms
	Memory_limit     int                 `gorm:"not null;default:256"`           //题目内存限制，单位MB
	Stack_limit      int                 `gorm:"not null;default:128"`           //栈空间，默认128MB
	Description      Problem_Description `gorm:"foreignKey:DID"`                 //
	DID              uint                ``                                      //题面编号
	Source           string              ``                                      //题目来源(hduoj、vj...)
	Owner            int                 ``                                      //题目创建者，为0时任何用户都允许访问，为1时仅管理员用户可访问，其他数字则只有对应用户可以访问
	Code_share       bool                `gorm:"default:false"`                  //是否允许共享代码
	Spj_language     string              `gorm:"not null"`                       //Spj的代码语言
	Case_version     uint                `gorm:"default:1"`                      //测试用例版本，可用于题目重测
	Open_case_result bool                ``                                      //是否公开测试用例
	SubmitNumber     int                 ``                                      //提交次数
	SubmitACNumber   int                 ``                                      //AC数
	Modified_user    string              ``                                      //最后修改的用户
	CreatedAt        time.Time           ``                                      //创建时间
	UpdatedAt        time.Time           ``                                      //修改时间
}

type Judge struct {
	ID            uint      `gorm:"primaryKey"`     //评测编号
	Problem       Problem   `gorm:"foreignKey:PID"` //
	PID           uint      ``                      //题目编号
	PTitle        string    ``                      //题目标题
	PShow_id      string    ``                      //题目显示ID
	User_info     User_info `gorm:"foreignKey:UID"` //
	UID           uint      `gorm:"index"`          //用户编号
	Submit_time   time.Time `gorm:"autoCreateTime"` //提交时间
	Status        string    ``                      //评测状态(Pending、Compiling、Running、AC、CE、RE、WA、TLE、MLE、OLE)
	Share_code    bool      ``                      //是否共享代码
	Error_message string    ``                      //错误信息
	Time_use      int       ``                      //耗时，单位ms
	Memory_use    int       ``                      //使用内存，单位MB
	Length        int       ``                      //代码长度
	Code          string    ``                      //代码
	Language      string    ``                      //语言
	Judger        string    ``                      //评测机IP
	Ip            string    ``                      //用户IP
	CID           uint      ``                      //比赛ID
	Version       uint      ``                      //评测版本
	CreatedAt     time.Time ``                      //创建时间
	UpdatedAt     time.Time ``                      //修改时间
}

type Judge_case struct {
	ID         uint   `gorm:"primaryKey"`     //评测用例编号
	Judge      Judge  `gorm:"foreignKey:JID"` //
	JID        uint   ``                      //评测编号
	Status     string ``                      //评测状态
	Time_use   int    ``                      //耗时，单位ms
	Memory_use int    ``                      //使用内存，单位MB
	Case_id    int    ``                      //测试数据编号
}

type Contest struct {
	ID             uint      `gorm:"primaryKey"`             //比赛ID
	User_info      User_info `gorm:"foreignKey:UID"`         //
	UID            uint      ``                              //创建比赛的用户ID
	Authur         string    ``                              //创建比赛的用户名
	Title          string    `gorm:"index,class:FULLTEXT"`   //比赛标题
	Type           byte      `gorm:"not null;default:0"`     //比赛种类，0:ACM、1:OI、2:CF
	Source         string    ``                              //比赛来源
	Owner          uint      `gorm:"not null;default:0"`     //比赛访问权限，为0时公开，为1时管理员可见，否则仅创建者可见
	Visible        byte      `gorm:"not null;default:0"`     //参赛权限，为0时公开报名，为1时邀请制，为2时密码制，为3时邀请、密码混合制
	Password       string    ``                              //参赛密码
	Start_time     time.Time ``                              //开始时间
	End_time       time.Time ``                              //结束时间
	Profile        string    ``                              //简介
	Description    string    ``                              //比赛详细介绍
	Seal_rank      bool      `gorm:"not null;default:false"` //是否封榜
	Seal_rank_time time.Time ``                              //封榜时间
	Status         byte      `gorm:"not null;default:0"`     //比赛状态
	Open_print     bool      ``                              //开放打印
	Rank_show_name string    ``                              //rank显示名字：username，nikename，realname
	Open_out_rank  bool      `gorm:"not null;default:true"`  //开放外榜
	CreatedAt      time.Time ``                              //创建时间
	UpdatedAt      time.Time ``                              //修改时间
}

type Contest_problem struct {
	ID             uint      `gorm:"primaryKey"`     //比赛题目ID
	Show_ID        string    ``                      //显示ID
	Contest        Contest   `gorm:"foreignKey:CID"` //
	CID            uint      `gorm:"index"`          //比赛ID
	Problem        Problem   `gorm:"foreignKey:PID"` //
	PID            uint      ``                      //题目ID
	Title          string    ``                      //题目标题
	Color          string    ``                      //题目气球颜色
	Score          string    ``                      //题目分数，在OI赛制和CF赛制下有意义
	SubmitNumber   int       ``                      //提交数
	SubmitACNumber int       ``                      //AC数
	CreatedAt      time.Time ``                      //创建时间
	UpdatedAt      time.Time ``                      //修改时间
}

type Contest_register struct {
	ID        uint      `gorm:"primaryKey"`     //报名表ID
	Contest   Contest   `gorm:"foreignKey:CID"` //
	CID       uint      `gorm:"index"`          //比赛ID
	User_info User_info `gorm:"foreignKey:UID"` //
	UID       uint      ``                      //用户ID
	CreatedAt time.Time ``                      //创建时间
	UpdatedAt time.Time ``                      //修改时间
}

type Contest_record struct {
	ID              uint            `gorm:"primaryKey"`      //比赛提交ID
	Contest         Contest         `gorm:"foreignKey:CID"`  //
	CID             uint            `gorm:"index"`           //比赛ID
	User_info       User_info       `gorm:"foreignKey:UID"`  //
	UID             uint            ``                       //用户ID
	Contest_problem Contest_problem `gorm:"foreignKey:CPID"` //
	CPID            uint            ``                       //比赛题目ID
	Judge           Judge           `gorm:"foreignKey:JID"`  //
	JID             uint            ``                       //提交ID
	Submit_time     time.Time       `gorm:"autoCreateTime"`  //提交时间
	Score           int             ``                       //得分
	CreatedAt       time.Time       ``                       //创建时间
	UpdatedAt       time.Time       ``                       //修改时间
}

type Contest_judge struct {
	ID            uint      `gorm:"primaryKey"`     //评测编号
	Problem       Problem   `gorm:"foreignKey:PID"` //
	PID           uint      ``                      //题目编号
	PTitle        string    ``                      //题目标题
	PShow_id      string    ``                      //题目显示ID
	User_info     User_info `gorm:"foreignKey:UID"` //
	UID           uint      ``                      //用户编号
	Submit_time   time.Time `gorm:"autoCreateTime"` //提交时间
	Status        string    ``                      //评测状态(Pending、Compiling、Running、AC、CE、RE、WA、TLE、MLE、OLE)
	Share_code    bool      ``                      //是否共享代码
	Error_message string    ``                      //错误信息
	Time_use      int       ``                      //耗时，单位ms
	Memory_use    int       ``                      //使用内存，单位MB
	Length        int       ``                      //代码长度
	Code          string    ``                      //代码
	Language      string    ``                      //语言
	Judger        string    ``                      //评测机IP
	Ip            string    ``                      //用户IP
	Contest       Contest   `gorm:"foreignKey:CID"` //
	CID           uint      ``                      //比赛ID
	Score         int       ``                      //得分
	Version       uint      ``                      //评测版本
	CreatedAt     time.Time ``                      //创建时间
	UpdatedAt     time.Time ``                      //修改时间
}

type Contest_judge_case struct {
	ID            uint          `gorm:"primaryKey"`     //评测用例编号
	Contest_Judge Contest_judge `gorm:"foreignKey:JID"` //
	JID           uint          `gorm:"index"`          //评测编号
	Status        string        ``                      //评测状态
	Time_use      int           ``                      //耗时，单位ms
	Memory_use    int           ``                      //使用内存，单位MB
	Case_id       int           ``                      //测试数据编号
	Score         int           ``                      //得分
}

type User_info struct {
	ID                uint      `gorm:"primaryKey"`                    //用户编号
	Username          string    `gorm:"not null;index,class:FULLTEXT"` //用户名
	Password          string    `gorm:"not null"`                      //密码
	Nickname          string    `gorm:"index,class:FULLTEXT"`          //昵称
	School            string    ``                                     //学校
	Course            string    ``                                     //专业
	StudentID         string    ``                                     //学号
	Realname          string    ``                                     //真实姓名
	Email             string    ``                                     //邮箱
	Gender            byte      ``                                     //性别
	Avatar            string    ``                                     //头像地址
	Signature         string    ``                                     //个性签名
	Cf_username       string    ``                                     //Codeforces id
	Luogu_username    string    ``                                     //洛谷id
	Nowcoder_username string    ``                                     //牛客id
	Vj_username       string    ``                                     //vj id
	Blog              string    ``                                     //博客
	Github            string    ``                                     //github
	Title             string    ``                                     //头衔
	Title_color       string    ``                                     //头衔颜色
	Status            byte      ``                                     //状态
	CreatedAt         time.Time ``                                     //创建时间
	UpdatedAt         time.Time ``                                     //修改时间
}
