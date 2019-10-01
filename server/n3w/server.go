// server.go

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	elog "github.com/labstack/gommon/log"
	n3context "github.com/nsip/n3-context"
)

var demoUser, demoContext = "user01", "mySchool"

//
// token secret for demonstration puroses only
//
var demoSecret = []byte("nias3demosecret")

var cm = n3context.NewN3ContextManager()

//
// wrapper type to capture graphql input
//
type GQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func main() {

	//
	// restore context manager from previous runs
	//
	err := cm.Restore()
	if err != nil && err != n3context.ErrContextFileNotFound {
		log.Fatal("Unable to restore saved contexts: ", err)
	}

	//
	// Create web server, and graceful shutdown handler
	//

	// Setup the web server
	e := echo.New()
	e.Logger.SetLevel(elog.INFO)

	// handler to create a new context for demonstrations only
	e.POST("/admin/newdemocontext", createContext)

	// demo query
	e.GET("/staffTraversal", staffTraversal)

	// Start server
	go func() {
		e.Use(middleware.Gzip())
		// e.Use(middleware.Recover())
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		})) // allow cors requests during testing

		// entry point for javascript/css/html resources etc.
		e.Static("/", "public")

		// protect n3 endpoints /publish /graphql with
		// requirement for jwt
		r := e.Group("/n3")
		r.Use(middleware.JWT(demoSecret))
		r.POST("/graphql", graphql)
		r.POST("/publish", publish)

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

	// shut down the contexts, persist details
	log.Println("Closing contexts, and saving...")
	err = cm.Close(true)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("...ContextManager closed")

}

//
// handler allows user to create a new context
// and returns a token that then allows them
// to publish and query that context.
//
// tokes are scoped to demo audience for now
// and do not expire.
//
func createContext(c echo.Context) error {

	uname := c.FormValue("userName")
	cname := c.FormValue("contextName")

	if uname == "" || cname == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Must have userName and contextName")
	}

	// create the context
	// add a context
	n3ctx, err := cm.AddContext(uname, cname)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create context", err)
	}
	// attach the context to the streaming server
	// to recieve updates
	err = n3ctx.Activate()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to activate context", err)
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["uname"] = uname
	claims["cname"] = cname
	claims["aud"] = "demo"
	// claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString(demoSecret)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token":   t,
		"message": "Context created and activated successfully. Use this token to publish: /n3/publish and query: /n3/graphql data. Token must be provided in Autorization: Bearer header.",
	})

}

//
// for any request with a valid jwt, returns the associated context
// for the user.
//
func getContextFromToken(c echo.Context) (*n3context.N3Context, *echo.HTTPError) {

	// get context details from token
	var u interface{}
	if u = c.Get("user"); u == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Missing token in request")
	}

	user, ok := u.(*jwt.Token)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Cannot read supplied jwt")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Cannot find claims in jwt")
	}

	uname := claims["uname"].(string)
	cname := claims["cname"].(string)
	if uname == "" || cname == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid claims in jwt")
	}

	// get the context
	n3ctx, err := cm.GetContext(uname, cname)
	if err != nil {
		log.Println("error fetching context: ", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve context", err)
	}

	return n3ctx, nil

}

//
// universal data publishing handler
//
func publish(c echo.Context) error {

	// get the context
	n3ctx, httpErr := getContextFromToken(c)
	if httpErr != nil {
		log.Println("error fetching context publish(): ", httpErr)
		return httpErr
	}

	// send in the data, via the crdt layer
	err := n3ctx.PublishFromHTTPRequest(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to publish data to context", err)
	}

	return c.JSON(http.StatusOK, "data published ok")
}

//
// universal gql query handler
//
func graphql(c echo.Context) error {

	// get the user's context from the request token
	n3ctx, httpErr := getContextFromToken(c)
	if httpErr != nil {
		log.Println("error fetching context graphql(): ", httpErr)
		return httpErr
	}

	// get the request graphql query element
	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to parse graphql query / params")
	}

	// handle the query
	results, err := n3ctx.GQLQuery(grq.Query, grq.Variables)
	if err != nil {
		log.Println("gqlHandler error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Query failed:", err)
	}

	return c.JSON(http.StatusOK, results)
}

//
// to be replaced
//
func staffTraversal(c echo.Context) error {
	// staffId := c.QueryParam("staffid")

	// // create traversal as json
	// traversalDef := `{"TraversalSpec":[
	// 	"StaffPersonal",
	// 	"TeachingGroup",
	// 	"GradingAssignment",
	// 	"Property.Link",
	// 	"XAPI",
	// 	"Property.Link",
	// 	"Subject",
	// 	"Unique.Link",
	// 	"Syllabus",
	// 	"Unique.Link",
	// 	"Lesson"
	// ]}`
	// // fmt.Println("Traversal Spec:\n", traversalDef)
	// var jsonTraversal deep6.Traversal
	// if err := json.Unmarshal([]byte(traversalDef), &jsonTraversal); err != nil {
	// 	panic(err)
	// }
	// // create filterspec as json
	// filterDef := `{
	// 	"XAPI":[{
	// 		"Predicate":"actor.name","TargetValue":"Albert Lombardi"
	// 	}],
	// 	"TeachingGroup":[{
	// 		"Predicate":".LocalId","TargetValue":"2018-History-8-1-A"
	// 	}]
	// }`
	// // fmt.Println("Filter Spec:\n", filterDef)
	// var jsonFilterSpec deep6.FilterSpec
	// if err := json.Unmarshal([]byte(filterDef), &jsonFilterSpec); err != nil {
	// 	panic(err)
	// }

	// results, err := n3c.Query(staffId, jsonTraversal, jsonFilterSpec)
	// if err != nil {
	// 	return err
	// }
	// return c.JSON(http.StatusOK, results)

	return c.JSON(http.StatusOK, map[string]interface{}{"deprecated": "do not use this method"})
}
