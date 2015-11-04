package main

import (
//  "gopkg.in/mgo.v2/bson"
)

type Student struct {
  NetID string `json: "id" bson: "id"`
  Name string `json: "name" bson: "name"`
  Major string `json: "major" bson: "major"`
  Year int `json: "year" bson: "year"`
  Grade int `json: "grade" bson: "grade"`
  Rating string `json: "rating" bson: "rating"`
}

type Students []Student
