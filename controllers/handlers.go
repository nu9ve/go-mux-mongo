package writeon

import (
  "encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	// "github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"writeon/model"
	db "writeon/database"
)

var mh *db.MongoHandler

//NewMongoHandler uses MongoHandler constructor to create handler in package
func NewMongoHandler(address string) {
	mh = db.NewHandler(address)
}


// TOPIC HANDLING

// GetAllTopics query
func GetAllTopics(w http.ResponseWriter, r *http.Request) {
	topics := mh.GetTopics(bson.M{})
	json.NewEncoder(w).Encode(topics)
}

// AddTopic query
func AddTopic(w http.ResponseWriter, r *http.Request) {
	existingTopic := &model.Topic{}
	var topic model.TopicTemplate
	json.NewDecoder(r.Body).Decode(&topic)
	if topic.Name == "" {
		http.Error(w, fmt.Sprintf("Topic needs name"), 400)
		return
	}
	topic.CreatedAt = time.Now()
	err := mh.GetOneTopic(existingTopic, bson.M{"name": topic.Name})
	if err == nil {
		http.Error(w, fmt.Sprintf("Topic with essayID: %s already exist", topic.Name), 400)
		return
	} 
	inserted, insertErr := mh.AddOneTopic(&topic)
	if insertErr != nil {
		fmt.Println("InsertOne ERROR:", insertErr)
		// os.Exit(1) // safely exit script on error
		http.Error(w, fmt.Sprint(insertErr), 400)
		return
	}

	var finalTopic model.Topic
	finalTopic.ID = inserted.InsertedID.(primitive.ObjectID)
	finalTopic.Name = topic.Name
	finalTopic.CreatedAt = topic.CreatedAt
	json.NewEncoder(w).Encode(finalTopic)
}

// DeleteTopic query
func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	existingEssay := &model.Essay{}
	vars := mux.Vars(r)
	title := vars["topicID"]
	if title == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	err := mh.GetOneEssay(existingEssay, bson.M{"title": title})
	if err != nil {
		http.Error(w, fmt.Sprintf("Contact with topicID: %s does not exist", title), 400)
		return
	}
	_, err = mh.RemoveOneTopic(bson.M{"title": title})
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}
	w.Write([]byte("Contact deleted"))
	w.WriteHeader(200)
}


// ESSAY HANDLING

// GetEssay function
func GetEssay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	essayID := vars["essayID"]
	if essayID == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	essay := &model.Essay{}
	objID, _ := primitive.ObjectIDFromHex(essayID)
	err := mh.GetOneEssay(essay, bson.M{"_id": objID})
	if err != nil {
		fmt.Println("essay id not found")
		http.Error(w, fmt.Sprintf("Essay with id: %s not found", essayID), 404)
		return
	}
	json.NewEncoder(w).Encode(essay)
}

// GetAllEssays query
func GetAllEssays(w http.ResponseWriter, r *http.Request) {
	essays := mh.GetEssays(bson.M{})
	json.NewEncoder(w).Encode(essays)
}

// GetAllTopicEssays query
func GetAllTopicEssays(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicID"]
	essays := mh.GetEssays(bson.M{"topic": topicID})
	json.NewEncoder(w).Encode(essays)
}

// AddEssay query
func AddEssay(w http.ResponseWriter, r *http.Request) {
	existingEssay := &model.Essay{}
	var essay model.EssayTemplate
	json.NewDecoder(r.Body).Decode(&essay)
	if essay.Title == "" {
		json.NewEncoder(w).Encode("Essay needs title")
		// http.Error(w, fmt.Sprintf("Essay needs title"), 400)
		return
	}
	if essay.Body == "" {
		json.NewEncoder(w).Encode("Essay needs body")
		// http.Error(w, fmt.Sprintf("Essay needs body"), 400)
		return
	}
	if essay.Topic == "" {
		json.NewEncoder(w).Encode("Essay needs topic")
		// http.Error(w, fmt.Sprintf("Essay needs topic"), 400)
		return
	}
	essay.CreatedAt = time.Now()
	err := mh.GetOneEssay(existingEssay, bson.M{"title": essay.Title})
	if err == nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("Essay with essayId: %s already exist", essay.Title))
		// http.Error(w, fmt.Sprintf("Essay with essayId: %s already exist", essay.Title), 400)
		return
	}
	_, err = mh.AddOneEssay(&essay)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}
	
	json.NewEncoder(w).Encode("Essay created successfully")
		// w.Write([]byte("Essay created successfully"))
	// w.WriteHeader(201)
}

// DeleteEssay query
func DeleteEssay(w http.ResponseWriter, r *http.Request) {
	existingEssay := &model.Essay{}
	vars := mux.Vars(r)
	essayID := vars["essayID"]
	if essayID == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	err := mh.GetOneEssay(existingEssay, bson.M{"title": essayID})
	if err != nil {
		http.Error(w, fmt.Sprintf("Contact with essayID: %s does not exist", essayID), 400)
		return
	}
	_, err = mh.RemoveOneEssay(bson.M{"title": essayID})
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}
	w.Write([]byte("Contact deleted"))
	w.WriteHeader(200)
}

// UpdateEssay function
func UpdateEssay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["essayID"]
	if title == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	essay := &model.Essay{}
	json.NewDecoder(r.Body).Decode(essay)
	_, err := mh.UpdateEssay(essay, bson.M{"title": title}, essay)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}
	w.Write([]byte("Contact update successful"))
	w.WriteHeader(200)
}

