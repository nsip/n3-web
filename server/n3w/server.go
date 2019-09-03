// server.go

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	elog "github.com/labstack/gommon/log"
	n3context "github.com/nsip/n3-context"
	deep6 "github.com/nsip/n3-deep6"
)

var demoUser, demoContext = "user01", "mySchool"

func main() {

	// create context manager
	cm := n3context.NewN3ContextManager()

	// restore from file

	// add a context
	n3c, err := cm.AddContext(demoUser, demoContext)
	if err != nil {
		log.Fatal(err)
	}
	// attach the context to the streaming server
	// to recieve updates
	err = n3c.Activate()
	if err != nil {
		log.Fatal(err)
	}

	// // send in some data, via the crdt layer
	// dataFile := "./sample_data/xapi/xapi.json"
	// err = n3c.PublishFromFile(dataFile)
	// if err != nil {
	// 	log.Fatal("PublishFromFile() Error: ", err)
	// }
	// dataFile = "./sample_data/sif/sif.json"
	// err = n3c.PublishFromFile(dataFile)
	// if err != nil {
	// 	log.Fatal("PublishFromFile() Error: ", err)
	// }

	//
	// Create web server, and graceful shutdown handler
	//

	// Setup the web server
	e := echo.New()
	e.Logger.SetLevel(elog.INFO)
	e.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

	// handler for publish
	e.POST("/publish", func(c echo.Context) error {
		// send in some data, via the crdt layer
		err = n3c.PublishFromHTTPRequest(c.Request())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, "data published ok")
	})

	e.GET("/staffTraversal", func(c echo.Context) error {
		staffId := c.QueryParam("staffid")

		// create traversal as json
		traversalDef := `{"TraversalSpec":[
		"StaffPersonal",
		"TeachingGroup",
		"GradingAssignment",
		"Property.Link",
		"XAPI",
		"Property.Link",
		"Subject",
		"Unique.Link",
		"Syllabus",
		"Unique.Link",
		"Lesson"
	]}`
		// fmt.Println("Traversal Spec:\n", traversalDef)
		var jsonTraversal deep6.Traversal
		if err := json.Unmarshal([]byte(traversalDef), &jsonTraversal); err != nil {
			panic(err)
		}
		// create filterspec as json
		filterDef := `{
		"XAPI":[{
			"Predicate":"actor.name","TargetValue":"Albert Lombardi"
		}],
		"TeachingGroup":[{
			"Predicate":".LocalId","TargetValue":"2018-History-8-1-A"
		}]
	}`
		// fmt.Println("Filter Spec:\n", filterDef)
		var jsonFilterSpec deep6.FilterSpec
		if err := json.Unmarshal([]byte(filterDef), &jsonFilterSpec); err != nil {
			panic(err)
		}

		results, err := n3c.Query(staffId, jsonTraversal, jsonFilterSpec)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, results)
	})

	// Start server
	go func() {
		e.Use(middleware.Gzip())
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		})) // allow cors requests during testing

		// entry point for javascript/css resources etc.
		e.Static("/", "public")

		if err := e.Start(":1323"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	//
	// post-webserver shutdown clean-up
	//

	// shut down the contexts, but persist details
	log.Println("Closing created contexts, and saving...")
	err = cm.Close(true)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("...ContextManager closed")

}
