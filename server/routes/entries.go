package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mayurlalwani/go-react-calorie-tracker/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)
var validate = validator.New()
var entryCollection *mongo.Collection = OpenCollection(Client, "calories")

func throwError(errorMessage error, c *gin.Context) error{
	if errorMessage != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage.Error()})
		fmt.Println(errorMessage)
	}
	return errorMessage
}

func AddEntry(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry models.Entry
	defer cancel()
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validationErr := validate.Struct(entry)
	throwError(validationErr, c)
	entry.ID = primitive.NewObjectID()
	result, insertErr := entryCollection.InsertOne(ctx, entry)
	if insertErr != nil {
		msg := fmt.Sprintf("order item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func GetEntries(c *gin.Context){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entries []bson.M
	defer cancel()
	cursor, err := entryCollection.Find(ctx, bson.M{})
	throwError(err,c)
	if err = cursor.All(ctx, &entries); err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		fmt.Println(err)
		return
	}
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries)
}

func GetEntryById(c *gin.Context){
	EntryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(EntryID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var entry bson.M
	if err := entryCollection.FindOne(ctx, bson.M{"_id":docID}).Decode(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		fmt.Println(err)
		return
	}
	fmt.Println(entry)
	c.JSON(http.StatusOK, entry)
}

func GetEntriesByIngredient(c *gin.Context) {
	ingredient := c.Params.ByName("id")
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entries []bson.M
	defer cancel()
	cursor, err := entryCollection.Find(ctx, bson.M{"ingredients": ingredient})
	throwError(err,c)
	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries)
}

func UpdateEntry(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry models.Entry
	defer cancel()
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validationErr := validate.Struct(entry)
	throwError(validationErr,c)
	result, err := entryCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docID},
		bson.M{
			"dish":        entry.Dish,
			"fat":         entry.Fat,
			"ingredients": entry.Ingredients,
			"calories":    entry.Calories,
		},
	)
	throwError(err,c)
	defer cancel()
	c.JSON(http.StatusOK, result.ModifiedCount)
}

func UpdateIngredient(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	type Ingredient struct {
		Ingredients *string `json:"ingredients"`
	}
	defer cancel()
	var ingredient Ingredient
	if err := c.BindJSON(&ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	result, err := entryCollection.UpdateOne(ctx, bson.M{"_id": docID},
		bson.D{{"$set", bson.D{{"ingredients", ingredient.Ingredients}}}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, result.ModifiedCount)
}


func DeleteEntry(c *gin.Context){
	entryId := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryId)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	result, err := entryCollection.DeleteOne(ctx, bson.M{"_id":docID})
	throwError(err,c)
	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)
}
