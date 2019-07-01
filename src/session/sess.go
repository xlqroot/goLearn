package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"sync"
	"time"
)

//定义session的接口
type Session interface{
	Set(key,value interface{})	
	Get(key interface{}) interface{}	
    Remove(key interface{}) error
    GetId() string
} 

//定义一个实现内存session接口的结构体
type SessionFromMemory struct{
	sid 	string 		//每个cookie对应的sessionID值
	lock	sync.Mutex 	//互斥锁
	lastAccessedTime	time.Time	//用户最后的访问时间
	maxAge	int64		//存活时间 单位秒
	data	map[interface{}]interface{} 	//每个用户对应的session值/列表
}

//内存seesion结构体实现session接口
func (s SessionFromMemory) Set(key,value interface{}){
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[key] = value
}
func (s SessionFromMemory) Get(key interface{}) interface{} {
	if value := s.data[key];value != nil {
		return value
	}
	return nil
}
func (s SessionFromMemory) Remove(key interface{}) error {
	if value := s.data[key];value != nil {
		delete(s.data,key)
	}
	return nil
}
func (s SessionFromMemory) GetId() string {
	return s.sid
}

//获取新的session
func newSessionFromMemory() *SessionFromMemory {
	return &SessionFromMemory{
		data:   make(map[interface{}]interface{}),
        maxAge: 1800, //默认30分钟
	}
}



//定义session存储的驱动的接口
type Storage interface {
	//初始化一个session，id根据需要生成后传入
	InitSession(sid string,maxAge int64) (Session,error)
	//销毁session
    DestroySession(sid string) error
    //回收
    GCSession()
}

//定义实现内存session接口的管理器---即所有session的集合
type FromMemory struct {
	lock sync.Mutex
	sessions map[string]Session
}
//实例化一个内存session实现
func newFromMemory() *FromMemory {
    return &FromMemory{
        sessions: make(map[string]Session, 0),
    }
}

//初始化一个session 并将其加入管理列表
func (ms FromMemory) InitSession(sid string,maxAge int64) (Session,error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	
	newSession := newSessionFromMemory()
	newSession.sid = sid
	if maxAge > 0 {
		newSession.maxAge = maxAge
	}
	newSession.lastAccessedTime = time.Now()
	ms.sessions[sid] = newSession
	return newSession,nil
}
//销毁某个用户的session
func (ms FromMemory) DestroySession(sid string) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	if session :=ms.sessions[sid];session != nil {
		delete(ms.sessions,sid)
	}
	return nil
}
//清除超时的session -- 一段时间内未使用的sesiion
func (ms FromMemory)GCSession()  {
	sessions := ms.sessions
	if len(sessions) < 1 {
		return
	}
	for k,v := range sessions {
		t := v.(*SessionFromMemory).lastAccessedTime.Unix() + v.(*SessionFromMemory).maxAge
		if t < time.Now().Unix() {
			delete(ms.sessions,k)
		}
	}

}

//控制session storage cookie的管理器
type SessionManager struct {
	cookieName string   //客户端cookie的名字
	storage Storage    //session驱动
	maxAge int64 		//存活时间
	lock sync.Mutex     //互斥锁
}

//定时清理数据
func (m *SessionManager) GC() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.storage.GCSession()
	//在多长时间后执行匿名函数，这里指在某个时间后执行GC
	time.AfterFunc(time.Duration(m.maxAge*10), func() {
		m.GC()
	})
}
//生成一定长度的随机数
func (m *SessionManager) randomId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	//加密
	return base64.URLEncoding.EncodeToString(b)
}

//通过ID获取session
func (m *SessionManager) GetSessionById(sid string) Session {
	session := m.storage.(*FromMemory).sessions[sid]
	return session
}

func NewSessionManager() *SessionManager {
	sessionManager := &SessionManager{
		cookieName: "lzy-cookie",
		storage:    newFromMemory(), //默认以内存实现
		maxAge:     1800,         //默认30分钟
	}
	go sessionManager.GC()

	return sessionManager
}




//https://studygolang.com/articles/12856





	
