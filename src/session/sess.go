package main

import "sync"

//session 管理器
type Manager struct {
	cookieName string  //private 客户端cookie名称 如jsessionID
	lock sync.Mutex    //互斥锁
	provider Provider  //session数据存储驱动  文件、内存、数据库
	maxlifetime int64  //session 垃圾回收时间
}

//session 数据驱动接口
type Provider interface {
	SessionInit(sid string) (Session,error)  //初始化session驱动
	SessionRead(sid string) (Session,error)  //读取session
	SessionDestroy(sid string) error         //删除session
	SessionGC(maxlifetime int64)             //设置session 定时过期
}

//session 存取接口
type Session interface {
	Set(key,value interface{}) error   //存数据
	Get(key interface{}) interface{}  //读取数据
	Delete(key interface{}) 		  //删除数据
	SessionId() string				  //返回当前的sessionID
	GC() 							  //删掉用户的所有json
}

//驱动列表
var providers = make(map[string]Provider)

//驱动注册
func Register(name string,provide Provider)  {
	
}

//创建新的session管理器
func NewManager(providerName,cookieName string,maxlifetime int64)  {
	
}

func main()  {

}
