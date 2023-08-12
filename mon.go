package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"karmabot2.0/internal/model"
)

var collection *mongo.Collection

const (
	connectionString = "mongodb://localhost:27017/"
	dbName           = "gokarmas"
	colName          = "karmas"
)

func init() {
	// client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongonnected")

	collection = client.Database(dbName).Collection(colName)

	// collection instance
	fmt.Println("Mongollection is ready")
}

func newKarma(karma model.Karma) {
	_, err := collection.InsertOne(context.TODO(), karma)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted!")
}

func modifyKarma(karmaID string, plus bool) {
	objectId, _ := primitive.ObjectIDFromHex(karmaID)
	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"karma": true, "lastGived": true}}

	if _, err := collection.UpdateOne(context.TODO(), filter, update); err != nil {
		log.Fatal(err)
	}
}

func getKarmas(sort int) {
	filter := bson.M{}

	// Opcional: Define opciones de ordenamiento (ascendente) y limitar resultados a 10 karmaumentos.
	options := options.Find().SetSort(bson.M{"_id": sort}).SetLimit(10)

	// Realiza la consulta con filtro y opciones.
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// Itera sobre los resultados.
	var karmas []model.Karma
	for cursor.Next(context.TODO()) {
		var karma model.Karma
		err := cursor.Decode(&karma)
		if err != nil {
			log.Fatal(err)
		}
		karmas = append(karmas, karma)
	}

	// Manejo de los karmaumentos obtenidos.
	for _, karma := range karmas {
		fmt.Println(karma)
	}
}
