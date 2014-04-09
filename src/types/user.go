package types

const (
	SYS_USR = 0
)

type User struct {
	Id             int32  // 用户id
	Domain         string // 玩家所在分服
	Name           string // 用户名
	Flag           int32  // 状态标记
	Pass           []byte // 密码(MD5 Hash)
	Score          int32  // 分数
	ProtectTimeout int64  // 护盾截止时间
	Mac            string // 玩家MAC地址
	CountryCode    string // 国家代码
	Language       string // 界面语言
	DeviceType     string // 设备类型
	LastSaveTime   int64  // 服务器最后一次刷入数据库的时间
	CreatedAt      int64  // 注册时间
}
