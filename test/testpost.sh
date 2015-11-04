#!/bin/sh
curl -H "Content-Type: application/json" -d '{"NetID":"147001235", "Name":"Paul","Major":"Computer Science","Year":2015,"Grade":90,"Rating":"D"}' http://localhost:1234/Student
