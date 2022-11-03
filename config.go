package gredis

const (
	Client   = "Client"
	Sentinel = "Sentinel"
	Cluster  = "Cluster"
)

// defaultConfiguration 默认配置
var _defaultConfiguration = configuration{}

type configuration struct {
	RedisType   string
	Hosts       []string
	MasterName  string
	PassWord    string
	PoolSize    int
	MinIdleCons int
}

//Option 配置参数接口
type Option interface {
	apply(*configuration)
}

//funcOption 配置参数函数
type funcOption struct {
	f func(*configuration)
}

//apply 执行赋值操作
func (fdo *funcOption) apply(do *configuration) {
	fdo.f(do)
}

//newFuncOption 创建配置参数函数
func newFuncOption(f func(*configuration)) *funcOption {
	return &funcOption{
		f: f,
	}
}

//WithHosts 设置连接地址
func WithHosts(hosts []string) Option {
	return newFuncOption(func(c *configuration) {
		c.Hosts = hosts
	})
}

//WithMasterName 设置主节点名称
func WithMasterName(masterName string) Option {
	return newFuncOption(func(c *configuration) {
		c.MasterName = masterName
	})
}

//WithPassWord 设置连接密码
func WithPassWord(passWord string) Option {
	return newFuncOption(func(c *configuration) {
		c.PassWord = passWord
	})
}

//WithPoolSize 设置连接池大小
func WithPoolSize(poolSize int) Option {
	return newFuncOption(func(c *configuration) {
		c.PoolSize = poolSize
	})
}

//WithHosts 设置最小存活连接数
func WithMinIdleCons(minIdleCons int) Option {
	return newFuncOption(func(c *configuration) {
		c.MinIdleCons = minIdleCons
	})
}

func WithClient() Option {
	return newFuncOption(func(c *configuration) {
		c.RedisType = Client
	})
}

func WithSentinel() Option {
	return newFuncOption(func(c *configuration) {
		c.RedisType = Sentinel
	})
}

func WithCluster() Option {
	return newFuncOption(func(c *configuration) {
		c.RedisType = Cluster
	})
}
