#!/bin/sh
curl -X "DELETE" -H "Content-Type: application/json" -d '{"year":2010}' http://localhost:1234/Student
