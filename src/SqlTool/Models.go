package sqltool

import (
	"fmt"
	"time"

	"gorm.io/plugin/optimisticlock"
)

type SqlModel interface {
	GetID() string
}

type ProblemDescription struct {
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

func (m *ProblemDescription) GetID() string {
	return fmt.Sprint(m.ID)
}

type Problem struct {
	ID             uint                   `gorm:"primaryKey"`                     //题目编号
	JudgeMode      byte                   `gorm:"not null;default:0"`             //测试模式 (0:default，1:spj，2:interactive,3:Leetcode)
	ShowID         string                 ``                                      //题目显示ID
	Title          string                 `gorm:"not null;index:,class:FULLTEXT"` //题目标题
	TimeLimit      int                    `gorm:"not null;default:1000"`          //题目时限，单位ms
	MemoryLimit    int                    `gorm:"not null;default:256"`           //题目内存限制，单位MB
	StackLimit     int                    `gorm:"not null;default:128"`           //栈空间，默认512MB
	Description    ProblemDescription     `gorm:"foreignKey:DID"`                 //
	DID            uint                   ``                                      //题面编号
	Source         string                 ``                                      //题目来源(hduoj、vj...)
	Owner          int                    ``                                      //题目创建者，为0时任何用户都允许访问，为1时仅管理员用户可访问，其他数字则只有对应用户可以访问
	CodeShare      bool                   `gorm:"default:false"`                  //是否允许共享代码
	SpjLanguage    string                 `gorm:"not null"`                       //Spj的代码语言
	CaseVersion    uint                   `gorm:"default:1"`                      //测试用例版本，可用于题目重测
	OpenCaseResult bool                   ``                                      //是否公开测试用例
	SubmitNumber   int                    ``                                      //提交次数
	SubmitACNumber int                    ``                                      //AC数
	ModifiedUser   string                 ``                                      //最后修改的用户
	CreatedAt      time.Time              ``                                      //创建时间
	UpdatedAt      time.Time              ``                                      //修改时间
	Version        optimisticlock.Version ``                                      //版本控制
}

func (m *Problem) GetID() string {
	return fmt.Sprint(m.ID)
}

type Judge struct {
	ID           uint      `gorm:"primaryKey"`     //评测编号
	Problem      Problem   `gorm:"foreignKey:PID"` //
	PID          uint      ``                      //题目编号
	PTitle       string    ``                      //题目标题
	PShowID      string    ``                      //题目显示ID
	UserInfo     UserInfo  `gorm:"foreignKey:UID"` //
	UID          uint      `gorm:"index"`          //用户编号
	SubmitTime   time.Time `gorm:"autoCreateTime"` //提交时间
	Status       string    ``                      //评测状态(Pending、Compiling、Running、AC、CE、RE、WA、TLE、MLE、OLE)
	ShareCode    bool      ``                      //是否共享代码
	ErrorMessage string    ``                      //错误信息
	TimeUse      int       ``                      //耗时，单位ms
	MemoryUse    int       ``                      //使用内存，单位MB
	Length       int       ``                      //代码长度
	Code         string    ``                      //代码
	Language     string    ``                      //语言
	Judger       string    ``                      //评测机
	Ip           string    ``                      //用户IP
	CID          uint      ``                      //比赛ID
	Version      uint      ``                      //评测版本
	CreatedAt    time.Time ``                      //创建时间
	UpdatedAt    time.Time ``                      //修改时间
}

func (m *Judge) GetID() string {
	return fmt.Sprint(m.ID)
}

type JudgeCase struct {
	ID        uint   `gorm:"primaryKey"`     //评测用例编号
	Judge     Judge  `gorm:"foreignKey:JID"` //
	JID       uint   ``                      //评测编号
	Status    string ``                      //评测状态
	TimeUse   int    ``                      //耗时，单位ms
	MemoryUse int    ``                      //使用内存，单位MB
	CaseId    int    ``                      //测试数据编号
}

func (m *JudgeCase) GetID() string {
	return fmt.Sprint(m.ID)
}

type Contest struct {
	ID           uint      `gorm:"primaryKey"`             //比赛ID
	UserInfo     UserInfo  `gorm:"foreignKey:UID"`         //
	UID          uint      ``                              //创建比赛的用户ID
	Authur       string    ``                              //创建比赛的用户名
	Title        string    `gorm:"index,class:FULLTEXT"`   //比赛标题
	Type         byte      `gorm:"not null;default:0"`     //比赛种类，0:ACM、1:OI、2:CF
	Source       string    ``                              //比赛来源
	Owner        uint      `gorm:"not null;default:0"`     //比赛访问权限，为0时公开，为1时管理员可见，否则仅创建者可见
	Visible      byte      `gorm:"not null;default:0"`     //参赛权限，为0时公开报名，为1时邀请制，为2时密码制，为3时邀请、密码混合制
	Password     string    ``                              //参赛密码
	StartTime    time.Time ``                              //开始时间
	EndTime      time.Time ``                              //结束时间
	Profile      string    ``                              //简介
	Description  string    ``                              //比赛详细介绍
	SealRank     bool      `gorm:"not null;default:false"` //是否封榜
	SealRankTime time.Time ``                              //封榜时间
	Status       byte      `gorm:"not null;default:0"`     //比赛状态
	OpenPrint    bool      ``                              //开放打印
	RankShowName string    ``                              //rank显示名字：username，nikename，realname
	OpenOutRank  bool      `gorm:"not null;default:true"`  //开放外榜
	CreatedAt    time.Time ``                              //创建时间
	UpdatedAt    time.Time ``                              //修改时间
}

func (m *Contest) GetID() string {
	return fmt.Sprint(m.ID)
}

type ContestProblem struct {
	ID             uint      `gorm:"primaryKey"`     //比赛题目ID
	ShowID         string    ``                      //显示ID
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

func (m *ContestProblem) GetID() string {
	return fmt.Sprint(m.ID)
}

type ContestRegister struct {
	ID        uint      `gorm:"primaryKey"`     //报名表ID
	Contest   Contest   `gorm:"foreignKey:CID"` //
	CID       uint      `gorm:"index"`          //比赛ID
	UserInfo  UserInfo  `gorm:"foreignKey:UID"` //
	UID       uint      ``                      //用户ID
	CreatedAt time.Time ``                      //创建时间
	UpdatedAt time.Time ``                      //修改时间
}

func (m *ContestRegister) GetID() string {
	return fmt.Sprint(m.ID)
}

type ContestRecord struct {
	ID             uint           `gorm:"primaryKey"`      //比赛提交ID
	Contest        Contest        `gorm:"foreignKey:CID"`  //
	CID            uint           `gorm:"index"`           //比赛ID
	UserInfo       UserInfo       `gorm:"foreignKey:UID"`  //
	UID            uint           ``                       //用户ID
	ContestProblem ContestProblem `gorm:"foreignKey:CPID"` //
	CPID           uint           ``                       //比赛题目ID
	Judge          Judge          `gorm:"foreignKey:JID"`  //
	JID            uint           ``                       //提交ID
	SubmitTime     time.Time      `gorm:"autoCreateTime"`  //提交时间
	Score          int            ``                       //得分
	CreatedAt      time.Time      ``                       //创建时间
	UpdatedAt      time.Time      ``                       //修改时间
}

func (m *ContestRecord) GetID() string {
	return fmt.Sprint(m.ID)
}

type ContestJudge struct {
	ID           uint      `gorm:"primaryKey"`     //评测编号
	Problem      Problem   `gorm:"foreignKey:PID"` //
	PID          uint      ``                      //题目编号
	PTitle       string    ``                      //题目标题
	PShowID      string    ``                      //题目显示ID
	UserInfo     UserInfo  `gorm:"foreignKey:UID"` //
	UID          uint      ``                      //用户编号
	SubmitTime   time.Time `gorm:"autoCreateTime"` //提交时间
	Status       string    ``                      //评测状态(Pending、Compiling、Running、AC、CE、RE、WA、TLE、MLE、OLE)
	ShareCode    bool      ``                      //是否共享代码
	ErrorMessage string    ``                      //错误信息
	TimeUse      int       ``                      //耗时，单位ms
	MemoryUse    int       ``                      //使用内存，单位MB
	Length       int       ``                      //代码长度
	Code         string    ``                      //代码
	Language     string    ``                      //语言
	Judger       string    ``                      //评测机IP
	Ip           string    ``                      //用户IP
	Contest      Contest   `gorm:"foreignKey:CID"` //
	CID          uint      ``                      //比赛ID
	Score        int       ``                      //得分
	Version      uint      ``                      //评测版本
	CreatedAt    time.Time ``                      //创建时间
	UpdatedAt    time.Time ``                      //修改时间
}

func (m *ContestJudge) GetID() string {
	return fmt.Sprint(m.ID)
}

type ContestJudgeCase struct {
	ID           uint         `gorm:"primaryKey"`     //评测用例编号
	ContestJudge ContestJudge `gorm:"foreignKey:JID"` //
	JID          uint         `gorm:"index"`          //评测编号
	Status       string       ``                      //评测状态
	TimeUse      int          ``                      //耗时，单位ms
	MemoryUse    int          ``                      //使用内存，单位MB
	CaseId       int          ``                      //测试数据编号
	Score        int          ``                      //得分
}

func (m *ContestJudgeCase) GetID() string {
	return fmt.Sprint(m.ID)
}

type UserInfo struct {
	ID               uint      `gorm:"primaryKey"`                                    //用户编号
	Username         string    `gorm:"not null;index,class:FULLTEXT;uniqueIndex(20)"` //用户名
	Password         string    `gorm:"not null"`                                      //密码
	Nickname         string    `gorm:"index,class:FULLTEXT;uniqueIndex(20)"`          //昵称
	School           string    ``                                                     //学校
	Course           string    ``                                                     //专业
	StudentID        string    ``                                                     //学号
	Realname         string    ``                                                     //真实姓名
	Email            string    `gorm:"uniqueIndex(320)"`                              //邮箱
	Gender           byte      ``                                                     //性别
	Avatar           string    ``                                                     //头像地址
	Signature        string    ``                                                     //个性签名
	CfUsername       string    ``                                                     //Codeforces id
	LuoguUsername    string    ``                                                     //洛谷id
	NowcoderUsername string    ``                                                     //牛客id
	VjUsername       string    ``                                                     //vj id
	Blog             string    ``                                                     //博客
	Github           string    ``                                                     //github
	Title            string    ``                                                     //头衔
	TitleColor       string    ``                                                     //头衔颜色
	Status           byte      ``                                                     //状态
	CreatedAt        time.Time ``                                                     //创建时间
	UpdatedAt        time.Time ``                                                     //修改时间
}

func (m *UserInfo) GetID() string {
	return fmt.Sprint(m.ID)
}
