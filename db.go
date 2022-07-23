package main

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var connection *mongo.Client

func Connect(url string) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return errors.New(err.Error()), false
	}
	connection = client
	fmt.Println("Successfully Connected to DB!")
	return nil, true
}

type EventDataWithTeam struct {
	TeamId    string   `bson:"teamId"`
	EventID   string   `bson:"id"`
	EventName string   `bson:"name"`
	Users     []string `bson:"users"`
	TxnId     string   `bson:"txnId"`
}

type Event struct {
	ID    string   `bson:"id"`
	Users []string `bson:"users"`
}

type EventDataWithID struct {
	ID    string  `bson:"id"`
	Event []Event `bson:"event"`
}

type User struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Roll        string             `bson:"roll"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	WhatsApp    string             `bson:"whatsapp"`
	Department  string             `bson:"department"`
	College     string             `bson:"college"`
	PremiumUser bool               `bson:"premiumUser"`
	Events      []string           `bson:"events"` // It is Team ID which can be Queried from EventDataWithTeam
}

func GetUser(id string) (User, error) {
	defer func() {                        //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Println("Some Unkown error occured ", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
	defer cancel()

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return User{}, err
	}

	user := User{}

	collection := connection.Database(dbname).Collection("users")
	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetEventDataFromId(id string) (EventDataWithID, error) {
	defer func() {                        //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Println("Some Unkown error occured ", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
	defer cancel()

	event := EventDataWithID{}

	collection := connection.Database(dbname).Collection("events")
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&event)
	if err != nil {
		return EventDataWithID{}, err
	}

	return event, nil
}

func AddEventTeam(team EventDataWithTeam) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Some Unknown error occurred", err)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
	defer cancel()
	collection := connection.Database(dbname).Collection("teams")
	_, err := collection.InsertOne(ctx, team)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func AddPremiumUser(id string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Some Unkown error occured ", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
	defer cancel()

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := connection.Database(dbname).Collection("users")
	_, err = collection.UpdateOne(ctx, bson.M{"_id": docId}, bson.M{"$set": bson.M{"premiumUser": true}})
	if err != nil {
		return err
	}

	return nil
}

func UpdateTeamUsers(team string, usersString []string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Some Unkown error occured ", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
	defer cancel()

	var users []primitive.ObjectID
	for _, user := range usersString {
		docId, err := primitive.ObjectIDFromHex(user)
		if err != nil {
			return err
		}
		users = append(users, docId)
	}

	collection := connection.Database(dbname).Collection("users")
	_, err := collection.UpdateMany(ctx, bson.M{
		"_id": bson.M{"$in": users},
	}, bson.M{
		"$push": bson.M{
			"events": team,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
