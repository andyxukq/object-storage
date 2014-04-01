package adapters

import (
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"

	//. "file-storage/core"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var DbSession *mgo.Session
var MONGO_HOST string

func SetConfig(mongo_host string) {
	MONGO_HOST = mongo_host
	log.Println("Set to connect to MongoDB at", MONGO_HOST)
	Init()
}

func GetSession() *mgo.Session {
	if DbSession == nil {
		DbSession = ConnectToDB()
	}

	return DbSession.Clone()
}

func ConnectToDB() *mgo.Session {
	var err error

	DbSession, err = mgo.Dial(MONGO_HOST)
	if err != nil {
		panic(err)
	}

	session := GetSession()
	defer session.Close()

	return DbSession
}

func Init() {
	ConnectToDB()
	SetupAdapterLocator()
}

var AdaptersLocator struct {
}

func SetupAdapterLocator() {}

func InsertFileGlobal(fs string, v multipart.File) (string, error) {
	session := DbSession.Copy()
	defer session.Close()

	f, _ := session.DB("tiup-file-storage").GridFS(fs).Create("")
	var id string

	switch x := f.Id().(type) {
	case bson.ObjectId:
		id = x.Hex()
		break
	default:
		return "", errors.New("Mongo GridFS failed")
	}

	data, _ := ioutil.ReadAll(v)

	if _, err := f.Write(data); err != nil {
		log.Println("Unable to Insert: ", err.Error())
		return "", err
	}

	f.Close()
	return id, nil
}

func FindFileGlobal(fs, id string) ([]byte, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, nil
	}

	session := DbSession.Copy()
	defer session.Close()

	mgoid := bson.ObjectIdHex(id)

	f, err := session.DB("tiup-file-storage").GridFS(fs).OpenId(mgoid)
	if err != nil {
		return nil, nil
	}

	data, _ := ioutil.ReadAll(f)

	f.Close()
	return data, nil
}
