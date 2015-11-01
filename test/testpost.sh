#!/bin/sh
curl -H "Content-Type: application/json" -d '{"NetID":"147001234", "Name":"Mike","Major":"Computer Science","Year":2015,"Grade":90,"Rating":"D"}' http://localhost:1234/Student
