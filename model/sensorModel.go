package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"os"
	"smhome/database"
)

type Sensor struct {
	Id           string `json:"id"`
	Value        string `json:"value"`
	FeedID       int    `json:"feed_id"`
	FeedKey      string `json:"feed_key"`
	CreatedAt    string `json:"created_at"`
	CreatedEpoch int    `json:"created_epoch"`
	Expiration   string `json:"expiration"`
}

type Sensors struct {
	Type    string   `json:"type"`
	Payload []Sensor `json:"payload"`
}

func (s *Sensors) SetElement(typ string, value interface{}) error {
	switch typ {
	case "type":
		s.Type = value.(string)
		return nil
	}
	return nil
}

func (s *Sensors) GetEntity(param string) (interface{}, error) {
	errEnv := godotenv.Load()
	if errEnv != nil {
		return nil, errEnv
	}

	var api string
	typ, _ := s.GetElement("type")
	switch *typ {
	case "temperature":
		api = os.Getenv("API_TEMP")
	case "humidity":
		api = os.Getenv("API_HUMID")
	default:
		return nil, errors.New(fmt.Sprintf("no type in entity:%s", *typ))
	}
	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}

	//We Read the response body on the line below.
	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		return nil, errBody
	}

	var sensors Sensors
	errSen := json.Unmarshal(body, &sensors.Payload)
	if errSen != nil {
		return nil, errSen
	}

	s.Payload = sensors.Payload
	s.Type = sensors.Payload[0].FeedKey
	return sensors, nil
}

func (s *Sensors) DeleteEntity(param string) error {
	return nil
}

func (s *Sensors) UpdateData(payload interface{}) error {
	sensor, ok := payload.(Sensors)
	if !ok {
		return errors.New("InitField: Require a Sensors")
	}
	filter := bson.D{{"type", s.Type}}
	update := bson.D{{"$set", bson.D{{"payload", sensor.Payload}}}}
	collection := database.GetConnection().Database("SmartHomeDB").Collection("Sensors")
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sensors) InsertData(payload interface{}) error {
	typ, _ := s.GetElement("type")
	instanceSensor, _ := s.FindDocument("type", *typ)
	if instanceSensor != nil {
		sensors, ok := payload.(Sensors)
		if !ok {
			return errors.New("InitField: Require a Sensors")
		}
		err := s.UpdateData(sensors)
		if err != nil {
			return err
		}
		return nil
	}

	collection := database.GetConnection().Database("SmartHomeDB").Collection("Sensors")

	_, err := collection.InsertOne(context.TODO(), s)
	if err != nil {
		return err
	}
	return nil
}
func (s *Sensors) FindDocument(key string, val string) (interface{}, error) {

	filter := bson.D{{key, val}}

	var res Sensors
	collection := database.GetConnection().Database("SmartHomeDB").Collection("Sensors")

	err := collection.FindOne(context.TODO(), filter).Decode(&res)

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Sensors) GetElement(msg string) (*string, error) {
	switch msg {
	case "type":
		return &s.Type, nil
	default:
		return nil, errors.New("no element in user entity")
	}
}
