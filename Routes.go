package main

// TODO: perhaps we do not need the pattern?

import (
  "net/http"
  "github.com/gorilla/mux"
)

type Route struct {
  Name string
  Method string
  Pattern string
  HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
  router := mux.NewRouter().StrictSlash(true)
  for _, route := range routes {
    var handler http.Handler
    handler = route.HandlerFunc
    handler = Logger(handler, route.Name)
    // attach the routes to their handlers
    router.
      Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(handler)
  }
  return router
}

/*** THIS IS WHERE WE PUT ALL THE ROUTING INFORMATION ***/
/*** COMMUNICATES WITH HANDLERS.GO ***/

var routes = Routes{
  Route {
    "POST_Operation",
    "POST",
    "/Student",
    PostHandler,
  },
  Route {
    "GET_Operation",
    "GET",
    "/Student/getstudent?name={Name}",
    GetHandler,
  },
  Route {
    "DELETE_Operation",
    "DELETE",
    "/Student",
    DeleteHandler,
  },
  Route {
    "UPDATE_Operation",
    "PUT",
    "/Student",
    UpdateHandler,
  },
  Route {
    "LIST_Operation",
    "GET",
    "/Student/listall",
    ListHandler,
  },
}
