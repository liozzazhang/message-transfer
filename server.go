package main

import (
	"cloopy/services"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"log"
	"net/http"
)

func main() {
	// use default container, create a new webservice
	ws := new(restful.WebService)
	restful.DefaultContainer.Add(ws)

	// add cloopyResource
	u := services.GroupChatResource{}
	restful.DefaultContainer.Add(u.WebService())

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("/Users/lio/Projects/GoProjects/cloopy/swagger/dist"))))

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer}
	restful.DefaultContainer.Filter(cors.Filter)
	log.Printf("Get the API using http://localhost:12345/apidocs.json")
	log.Printf("Open Swagger UI using http://localhost:12345/apidocs/?url=http://localhost:12345/apidocs.json")
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Ops API",
			Description: "An Internal Open API Platform",
			Contact: &spec.ContactInfo{
				Name:  "lizzano",
				Email: "zlprasy@gmail.com",
				URL:   "http://lizzano.cn",
			},
			License: &spec.License{
				Name: "MIT",
				URL:  "http://mit.org",
			},
			Version: "1.0.0",
		},
	}
}
