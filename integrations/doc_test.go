package integrations_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bnkamalesh/webgo/v4"
	"github.com/gin-gonic/gin"

	"github.com/keploy/go-sdk/integrations/kecho/v4"
	"github.com/keploy/go-sdk/integrations/kgin/v1"
	"github.com/keploy/go-sdk/integrations/kgrpc"
	"github.com/keploy/go-sdk/integrations/khttpclient"
	"github.com/keploy/go-sdk/integrations/kmongo"
	"github.com/keploy/go-sdk/integrations/kwebgo/v4"
	"github.com/keploy/go-sdk/keploy"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func ExampleNewCollection() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	var collection *kmongo.Collection
	result, err := collection.InsertOne(context.TODO(), bson.D{{"x", 1}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("inserted ID: %v\n", result.InsertedID)
}

func ExampleSingleResult_Err() {
	var (
		sr         *kmongo.SingleResult
		collection *kmongo.Collection
	)
	filter := bson.M{"name": "Ash"}
	findOneOpts := options.FindOne()
	findOneOpts.SetComment("this is cool stuff")

	sr = collection.FindOne(context.TODO(), filter, findOneOpts)
	err := sr.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSingleResult_Decode() {
	var (
		sr         *kmongo.SingleResult
		collection *kmongo.Collection
	)
	filter := bson.M{"name": "Ash"}
	var result bson.D
	findOneOpts := options.FindOne()
	findOneOpts.SetComment("this is cool stuff")

	sr = collection.FindOne(context.TODO(), filter, findOneOpts)
	err := sr.Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found document: %v", fmt.Sprint(result))
}

func ExampleCursor_Err() {
	var (
		collection *kmongo.Collection
		cur        *kmongo.Cursor
		err        error
	)
	filter := bson.M{"name": "Misty"}
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{"age", -1}})
	cur, _ = collection.Find(context.TODO(), filter, findOpts)
	err = cur.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCursor_Next() {
	var (
		collection *kmongo.Collection
		cur        *kmongo.Cursor
		err        error
	)
	filter := bson.M{"name": "Misty"}
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{"age", -1}})
	cur, err = collection.Find(context.TODO(), filter, findOpts)
	if err != nil {
		log.Fatal(err)
	}
	var moreDocs bool
	moreDocs = cur.Next(context.TODO())
	fmt.Printf("More Docs: %v", moreDocs)
}

func ExampleCursor_TryNext() {
	var (
		collection *kmongo.Collection
		cur        *kmongo.Cursor
		err        error
	)
	filter := bson.M{"name": "Misty"}
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{"age", -1}})
	cur, err = collection.Find(context.TODO(), filter, findOpts)
	if err != nil {
		log.Fatal(err)
	}
	var moreDocs bool
	moreDocs = cur.TryNext(context.TODO())
	fmt.Printf("More Docs: %v", moreDocs)
}

func ExampleCursor_Close() {
	var (
		collection *kmongo.Collection
		cur        *kmongo.Cursor
		err        error
	)
	filter := bson.M{"name": "Misty"}
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{"age", -1}})
	cur, err = collection.Find(context.TODO(), filter, findOpts)
	if err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
}

func ExampleCursor_All() {
	var (
		collection *kmongo.Collection
		cur        *kmongo.Cursor
		err        error
	)
	filter := bson.M{"name": "Misty"}
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{"age", -1}})
	cur, err = collection.Find(context.TODO(), filter, findOpts)
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.D
	cur.All(context.TODO(), &results)
}

func ExampleCursor_Decode() {
	var (
		collection *kmongo.Collection
		cur        *kmongo.Cursor
		err        error
	)
	filter := bson.M{"name": "Misty"}
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{"age", -1}})
	cur, err = collection.Find(context.TODO(), filter, findOpts)
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.D
	for cur.Next(context.TODO()) {
		var elem bson.D
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, elem)
	}
}

func ExampleCollection_InsertOne() {
	var (
		collection      *kmongo.Collection
		err             error
		insertOneResult *mongo.InsertOneResult
	)
	ash := bson.D{{"name", "Alice"}}
	insertOneOpts := options.InsertOne()
	insertOneOpts.SetBypassDocumentValidation(false)
	insertOneResult, err = collection.InsertOne(context.TODO(), ash, insertOneOpts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("inserted document with ID %v\n", insertOneResult.InsertedID)
}

func ExampleCollection_FindOne() {
	var (
		sr         *kmongo.SingleResult
		collection *kmongo.Collection
		err        error
	)
	filter := bson.M{"name": "Ash"}
	var resulto bson.D
	findOneOpts := options.FindOne()
	findOneOpts.SetComment("this is cool stuff")
	sr = collection.FindOne(context.TODO(), filter, findOneOpts)
	err = sr.Err()
	if err != nil {
		log.Fatal(err)
	} else {
		sr.Decode(&resulto)
	}
}

func ExampleCollection_InsertMany() {
	var (
		insertManyResult *mongo.InsertManyResult
		collection       *kmongo.Collection
		err              error
	)
	docs := []interface{}{
		bson.D{{"name", "Alice"}},
		bson.D{{"name", "Bob"}},
	}
	insertManyOpts := options.InsertMany()
	insertManyOpts.SetOrdered(true)
	insertManyResult, err = collection.InsertMany(context.TODO(), docs, insertManyOpts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("inserted documents with IDs %v\n", insertManyResult.InsertedIDs...)
}

func ExampleCollection_Find() {
	var (
		collection *kmongo.Collection
		cur        *kmongo.Cursor
		err        error
	)
	filter := bson.M{"name": "Misty"}
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{"age", -1}})
	cur, err = collection.Find(context.TODO(), filter, findOpts)
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.D
	for cur.Next(context.TODO()) {
		var elem bson.D
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, elem)
	}
}

func ExampleCollection_UpdateOne() {
	var (
		result     *mongo.UpdateResult
		collection *kmongo.Collection
		err        error
	)
	filter := bson.M{"name": "Brock"}
	updateOpts := options.Update()
	updateOpts.SetBypassDocumentValidation(false)
	update := bson.D{{"$set", bson.D{{"name", "Brock"}, {"age", 22}, {"city", "Pallet Town"}}}}
	result, err = collection.UpdateOne(context.TODO(), filter, update, updateOpts)
	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}
}

func ExampleCollection_UpdateMany() {
	var (
		result     *mongo.UpdateResult
		collection *kmongo.Collection
		err        error
	)
	filter := bson.M{"name": "Brock"}
	updateOpts := options.Update()
	updateOpts.SetBypassDocumentValidation(false)
	update := bson.D{{"$set", bson.D{{"name", "Brock"}, {"age", 22}, {"city", "Pallet Town"}}}}
	result, err = collection.UpdateMany(context.TODO(), filter, update, updateOpts)
	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return
	}
}

func ExampleCollection_DeleteOne() {
	var (
		result     *mongo.DeleteResult
		collection *kmongo.Collection
		err        error
	)
	filter := bson.M{"name": "Brock"}
	deleteOpts := options.Delete()
	deleteOpts.SetHint("Go to cartoon network")
	result, err = collection.DeleteOne(context.TODO(), filter, deleteOpts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v document\n", result.DeletedCount)
}

func ExampleCollection_DeleteMany() {
	var (
		result     *mongo.DeleteResult
		collection *kmongo.Collection
		err        error
	)
	filter := bson.M{"name": "Brock"}
	deleteOpts := options.Delete()
	deleteOpts.SetHint("Go to cartoon network")
	result, err = collection.DeleteMany(context.TODO(), filter, deleteOpts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v documents\n", result.DeletedCount)
}


func ExampleWebGoV4() {
	port := "6060"
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "webgo-v4-app",
			Port: port,
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:8081/api",
		},
	})
	router := webgo.NewRouter(&webgo.Config{
		Host:         "",
		Port:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
	}, []*webgo.Route{})
	kwebgo.WebGoV4(k, router)
	router.Start()
}

func ExampleEchoV4() {
	e := echo.New()
	port := "6060"
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "echo-v4-app",
			Port: port,
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:8081/api",
		},
	})
	// Remember to call integrations.EchoV4 before route handling
	kecho.EchoV4(k, e)
	e.GET("/echo", func(c echo.Context) error {
		return nil
	})
	e.Start(":" + port)
}

func ExampleGinV1() {
	r := gin.New()
	port := "6060"
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "gin-v1-app",
			Port: port,
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:8081/api",
		},
	})
	//Call integration.GinV1 before routes handling
	kgin.GinV1(k, r)
	r.GET("/gin/:color/*type", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":" + port)
}

func ExampleWithClientUnaryInterceptor() {
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "my-app",
			Port: "8080",
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:8081/api",
		},
	})
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), kgrpc.WithClientUnaryInterceptor(k))
	if err != nil {
		log.Fatalf("Did not connect : %v", err)
	}
	defer conn.Close()
}

func ExampleWithClientStreamInterceptor() {
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "my-app",
			Port: "8080",
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:8081/api",
		},
	})

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), kgrpc.WithClientStreamInterceptor(k))
	if err != nil {
		log.Fatalf("Did not connect : %v", err)
	}
	defer conn.Close()
}

func ExampleNewHttpClient() {
	r := &http.Request{} // Here, r is for demo. You should use your handler's request as r.
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := khttpclient.NewHttpClient(&http.Client{
		Transport: tr,
	})

	// SetCtxHttpClient is called before mocked http.Client's Get method.
	client.SetCtxHttpClient(r.Context())
	resp, err := client.Get("https://example.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println("BODY : ", body)
}

func ExampleHttpClient_SetCtxHttpClient() {
	r := &http.Request{} // Here, r is for demo. You should use your handler's request as r.
	client := khttpclient.NewHttpClient(&http.Client{})

	// SetCtxHttpClient is called before mocked http.Client's Get method.
	client.SetCtxHttpClient(r.Context())
	resp, err := client.Get("https://example.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println("BODY : ", body)
}

func ExampleHttpClient_Get() {
	r := &http.Request{} // Here, r is for demo. You should use your handler's request as r.
	client := khttpclient.NewHttpClient(&http.Client{})

	// SetCtxHttpClient is called before mocked http.Client's Get method.
	client.SetCtxHttpClient(r.Context())
	resp, err := client.Get("https://example.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println("BODY : ", body)
}

func ExampleHttpClient_Do() {
	r := &http.Request{} // Here, r is for demo. You should use your handler's request as r.
	client := khttpclient.NewHttpClient(&http.Client{})

	// SetCtxHttpClient is called before mocked http.Client's Do method.
	client.SetCtxHttpClient(r.Context())
	req, err := http.NewRequestWithContext(r.Context(), "GET", "http://localhost:6060/getdocs?name=name&value=Ash", nil)
	if err != nil {
		log.Fatal("http client in webgo-v4 main.go")
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println("BODY : ", body)
}

func ExampleHttpClient_Post() {
	r := &http.Request{} // Here, r is for demo. You should use your handler's request as r.
	client := khttpclient.NewHttpClient(&http.Client{})

	// SetCtxHttpClient is called before mocked http.Client's Post method.
	client.SetCtxHttpClient(r.Context())
	postBody, _ := json.Marshal(map[string]interface{}{
		"name": "Toby",
		"age":  21,
		"city": "New York",
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:6060/createone", "application/json", responseBody)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println("BODY : ", body)
}
