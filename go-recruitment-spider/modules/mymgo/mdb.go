package mymgo

import (
	"fmt"
	"go-recruitment-spider/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mdbSession *MdbSession

func init() {
	mongoConfig := config.GetTomlConfig().Mongo
	session, err := mgo.Dial(mongoConfig.Addr)
	if err != nil {
		panic(err.Error())
	}
	mdbSession = &MdbSession{
		session: session,
		db:      mongoConfig.Db,
	}
}

func GetDb()  *MdbSession{
	return mdbSession
}

type Field struct {
	Id         string
	Collection string
}

type MdbSession struct {
	session *mgo.Session
	db      string
}

func (mdb *MdbSession) Session() *mgo.Session {
	return mdb.session.New()
}

var (
	AutoIncIdCollection = "amb_config"
	field               = &Field{
		Id:         "seq",
		Collection: "_id",
	}
)

func (mdb *MdbSession) AutoIncId(name string) (id int) {
	s := mdb.Session()
	id, err := autoIncr(s.DB(mdb.db).C(AutoIncIdCollection), name)
	s.Close()
	if err != nil {
		panic("Get next id of [" + name + "] fail:" + err.Error())
	}
	return
}

func autoIncr(c *mgo.Collection, name string) (id int, err error) {
	return incr(c, name, 1)
}

func incr(c *mgo.Collection, name string, step int) (id int, err error) {
	result := make(map[string]interface{})
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{field.Id: step}},
		Upsert:    true,
		ReturnNew: true,
	}
	_, err = c.Find(bson.M{field.Collection: name}).Apply(change, result)
	if err != nil {
		return
	}
	id, ok := result[field.Id].(int)
	if ok {
		return
	}
	id64, ok := result[field.Id].(int64)
	if !ok {
		err = fmt.Errorf("%s is ont int or int64", field.Id)
		return
	}
	id = int(id64)
	return
}
