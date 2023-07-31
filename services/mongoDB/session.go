package mongoDB

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mapping struct {
	BoundingBox struct {
		X      int `bson:"x"`
		Y      int `bson:"y"`
		Width  int `bson:"width"`
		Height int `bson:"height"`
	} `bson:"boundingBox"`
	Data struct {
		Alias   string   `bson:"alias"`
		Content []string `bson:"content"`
	} `bson:"data"`
}

type Hotspot struct {
	UUID  string `bson:"uuid"`
	Label string `bson:"label"`
	Meta  struct {
		CreatedBy string `bson:"createdBy"`
		CreatedAt string `bson:"createdAt"`
		UpdatedAt string `bson:"updatedAt"`
		IsActive  bool   `bson:"isActive"`
	} `bson:"meta"`
	PartOf      []string  `bson:"partOf"`
	DirectedFor []string  `bson:"directedFor"`
	ImageBytes  string    `bson:"imageBytes"`
	Mapping     []Mapping `bson:"mapping"`
	ImageHeight int       `bson:"imageHeight"`
	ImageWidth  int       `bson:"imageWidth"`
}

func GetPanoramicImages(institutionUUID string, hotspotUUID string, directedFor []*string) []Hotspot {
	uri := "mongodb+srv://andresilmor:project-greenhealth@carexr-mongodb-cluster.mcnxmz7.mongodb.net/?retryWrites=true&w=majority"
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	fmt.Printf("%s\n", institutionUUID)
	coll := client.Database("SessionMaterial").Collection("360_hotspots")

	ids := []string{"597d38a7-2fd1-4686-a3ef-96861466eb36"}

	// Construct the filter
	filter := bson.M{
		"partOf": bson.M{
			"$in": ids,
		},
	}

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error querying documents:", err)
		return nil
	}
	defer cursor.Close(context.Background())

	var hotspotList []Hotspot

	// Iterate over the cursor and decode each document
	for cursor.Next(context.Background()) {
		var hotspot Hotspot
		err := cursor.Decode(&hotspot)
		if err != nil {
			fmt.Println("Error decoding document:", err)
			continue
		}
		hotspotList = append(hotspotList, hotspot)
		// Access the values of each document

	}

	/*	ONE RESULT
		title := "test"
		var result bson.M

		filter := bson.D{{"title", title}}

		err = coll.FindOne(context.TODO(), filter).Decode(&result)
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the title %s\n", title)
			return
		}
		if err != nil {
			panic(err)
		}
		jsonData, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", jsonData)
	*/

	return hotspotList
}
