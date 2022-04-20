package mongocoba_test

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var defaultConfig = config{
	user:       "user",
	pass:       "pass",
	host:       "localhost",
	port:       27017,
	db:         "mongo_coba",
	authSource: "admin",
}

type config struct {
	user       string
	pass       string
	host       string
	port       int
	authSource string
	db         string
}

type User struct {
	ID    int    `bson:"_id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

type Group struct {
	GroupID string `bson:"group_id"`
	Name    string `bson:"name"`
	Email   string `bson:"email"`
}

func connect1(cnf config) (*mongo.Database, error) {
	const (
		funcName = `connect1`
		// refer : https://www.mongodb.com/docs/manual/reference/connection-string/
		// format from doc : mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
		// format from ex1 : mongodb://myDBReader:D1fficultP%40ssw0rd@mongodb0.example.com:27017/?authSource=admin
		formatURL = `mongodb://%s:%s@%s:%d/?authSource=%s`
	)
	clientOptions := options.Client()
	url := fmt.Sprintf(formatURL, cnf.user, cnf.pass, cnf.host, cnf.port, cnf.authSource)
	fmt.Println(funcName, `: url : `, url)
	clientOptions.ApplyURI(url)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	return client.Database(cnf.db), nil
}

func connect2(cnf config) (*mongo.Database, error) {
	const (
		funcName = `connect2`
		// refer : https://www.mongodb.com/docs/manual/reference/connection-string/
		// format from ex1 : mongodb://localhost:27017
		formatURL = `mongodb://%s:%d`
	)
	clientOptions := options.Client()
	url := fmt.Sprintf(formatURL, cnf.host, cnf.port)
	fmt.Println(funcName, `: url : `, url)
	clientOptions.ApplyURI(url).SetAuth(options.Credential{
		// AuthMechanism: "SCRAM-SHA-1",
		AuthSource: cnf.authSource,
		Username:   cnf.user,
		Password:   cnf.pass,
	})
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	return client.Database(cnf.db), nil
}

func Test_connect1(t *testing.T) {
	const funcName = `Test_connect1`
	db, err := connect1(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	t.Logf(`%s : db : %+v \n`, funcName, db)
	t.Logf(`%s : db.Client : %+v \n`, funcName, db.Client())
}

func Test_connect2(t *testing.T) {
	const funcName = `Test_connect2`
	db, err := connect2(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	t.Logf(`%s : db : %+v \n`, funcName, db)
	t.Logf(`%s : db.Client : %+v \n`, funcName, db.Client())
}

func Test_insertOne(t *testing.T) {
	const (
		funcName   = `Test_insertOne`
		collection = `mst_user`
	)
	db, err := connect2(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	resp, err := db.Collection(collection).InsertOne(context.Background(), User{
		ID:    1,
		Name:  "name",
		Email: "email",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(`%s : resp : %+v \n`, funcName, resp)
}

func Test_insertErrorDuplicate(t *testing.T) {
	const (
		funcName   = `Test_insertErrorDuplicate`
		collection = `mst_user`
	)
	db, err := connect2(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	usr := User{
		ID:    r.Intn(10000),
		Name:  "name",
		Email: "email",
	}

	// insert 1x
	_, err = db.Collection(collection).InsertOne(context.Background(), usr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("success insert pertama : ", usr)

	// insert lagi biar duplicate
	_, err = db.Collection(collection).InsertOne(context.Background(), usr)
	if err != nil {
		if !strings.Contains(err.Error(), `E11000 duplicate key error collection: mongo_coba.mst_user `) {
			t.Fatal(err)
		}
		fmt.Println("Error duplicate coy : ", err)
	}
}

func Test_insertBedaStruct(t *testing.T) {
	const (
		funcName   = `Test_insertBedaStruct`
		collection = `mst_group`
	)
	db, err := connect2(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	usr := User{
		ID:    r.Intn(100000),
		Name:  "name",
		Email: "email",
	}

	// insert 1x
	_, err = db.Collection(collection).InsertOne(context.Background(), usr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("success insert pertama : ", usr)

	grp := Group{
		GroupID: fmt.Sprint("group_id-", usr.ID),
		Name:    "group_name",
		Email:   "group_email",
	}

	// insert lagi biar duplicate
	_, err = db.Collection(collection).InsertOne(context.Background(), grp)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("success insert kedua : ", grp)
}

func Test_readData(t *testing.T) {
	const (
		funcName   = `Test_readData`
		collection = `mst_user`
		randNum    = 100000
	)
	db, err := connect2(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	id := r.Intn(randNum)

	usr := User{
		ID:    id,
		Name:  fmt.Sprint("nama-", id),
		Email: fmt.Sprint("email-", id),
	}

	// insert 1x
	_, err = db.Collection(collection).InsertOne(context.Background(), usr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("success insert pertama : ", usr)

	bsonM := bson.M{reflect.TypeOf(usr).Field(1).Tag.Get("bson"): usr.Name}
	fmt.Println("success insert pertama : ", bsonM)
	// read data
	cursor, err := db.Collection(collection).Find(context.Background(), bsonM)
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(context.Background())

	result := make([]User, 0)
	for cursor.Next(context.Background()) {
		var row User
		err := cursor.Decode(&row)
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, row)
	}

	t.Logf(`%s : result : %+v`, funcName, result)
}

func Test_readDatas(t *testing.T) {
	const (
		funcName   = `Test_readDatas`
		collection = `mst_user`
		randNum    = 100000
	)
	db, err := connect2(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	// pastikan data name itu ada banyak
	dt := reflect.TypeOf(User{}).Field(1).Tag.Get("bson")

	bsonM := bson.M{dt: dt}
	fmt.Println("success insert pertama : ", bsonM)
	// read data
	cursor, err := db.Collection(collection).Find(context.Background(), bsonM)
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(context.Background())

	result := make([]User, 0)
	for cursor.Next(context.Background()) {
		var row User
		err := cursor.Decode(&row)
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, row)
	}

	t.Logf(`%s : result : %+v`, funcName, result)
}

func Test_updateOne(t *testing.T) {
	const (
		funcName   = `Test_updateOne`
		collection = `mst_user`
		randNum    = 100000
	)

	ctx := context.Background()

	db, err := connect2(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(ctx)

	selector := bson.M{"name": "nama-11871-updated"}
	changes := User{
		ID:    11870, // id tidak bisa diubah
		Name:  "nama-11871-updated",
		Email: "email-11871-updated",
	}
	result, err := db.Collection(collection).UpdateOne(ctx, selector, bson.M{"$set": changes})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(`%s : result : %+v`, funcName, result)
}

func Test_deleteOne(t *testing.T) {
	const (
		funcName   = `Test_updateOne`
		collection = `mst_user`
		randNum    = 100000
	)

	ctx := context.Background()

	db, err := connect2(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(ctx)

	selector := bson.M{"name": "name"}
	result, err := db.Collection(collection).DeleteOne(ctx, selector)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(`%s : result : %+v`, funcName, result)
}
