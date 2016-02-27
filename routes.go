package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Upload",
		"GET",
		"/upload",
		Upload,
	},
	Route{
		"Query",
		"POST",
		"/query",
		Query,
	},
	Route{
		"View",
		"POST",
		"/view",
		View,
	},
	Route{
		"ViewUpload",
		"Post",
		"/viewupload",
		ViewUpload,
	},
}
