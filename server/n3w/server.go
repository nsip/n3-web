// server.go
//
// NOTE: On External Static Servers
//	DC-UI and dc-dynamic are hosted here on port 8002 and 8003 respectively.
//	n3/build brings these into public/dc-ui and public/dc-dynamic.
//	You can override the default location with ENVIRONENT varilables or .env file
//		N3_DCUI_PATH, N3_DCUI_PORT, N3_DCDYNAMIC_PATH, N3_DCDYNAMIC_PORT
//		e.g.
//			export N3_DCUI_PATH=~/nsip/DC-UI/dist/spa-mat/
//			export N3_DCDYNAMIC_PATH=~/nsip/dc-dynamic/dist/spa-mat/

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/digisan/gotk"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
	n3context "github.com/nsip/n3-context"
)

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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func makeStatic(keyPart string, portDefault string, pathDefault string) func() {
	return func() {
		le := echo.New()

		le.Use(middleware.Gzip())
		le.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		})) // allow cors requests during testing

		port := getEnv(keyPart+"_PORT", portDefault)
		path := getEnv(keyPart+"_PATH", pathDefault)

		// le.Use(middleware.Logger())
		// le.Use(middleware.Recover())

		le.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Index:  "index.html",
			Root:   path,
			Browse: true,
			HTML5:  true,
		}))
		if err := le.Start(":" + port); err != nil {
			le.Logger.Info("shutting down the server: " + path)
		}
	}
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

	// handler to create a temporary object model config for demonstrations only
	e.POST("/admin/addtmpobjconf", addTmpObjConf)

	// handler to clear all temporary object model configs for demonstrations only
	e.POST("/admin/clrtmpobjconf", clrTmpObjConf)

	// Start server
	go func() {
		e.Use(middleware.Gzip())
		// e.Use(middleware.Recover())
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		})) // allow cors requests during testing

		// // use for debugging only
		// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		// 	log.Println("reqBody:\n", string(reqBody))
		// }))

		// entry point for javascript/css/html resources etc.
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Index:  "index.html",
			Root:   "public",
			Browse: true,
			HTML5:  true,
		}))

		// protect n3 endpoints /publish /graphql with
		// requirement for jwt
		r := e.Group("/n3")
		r.Use(middleware.JWT(demoSecret))

		r.POST("/graphql", graphql)
		r.POST("/publish", publish)

		// 1323 - Main server
		if err := e.Start(":1323"); err != nil {
			e.Logger.Info("shutting down the server:", err)
		}
	}()

	// 8002 - dc-ui
	go makeStatic("N3_DCUI", "8002", "public/dc-ui")()
	go makeStatic("N3_DCDYNAMIC", "8003", "public/dc-dynamic")()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
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
// handler allows user to add a temporary object model classifier config
// which will be appended to n3-crdt default classifier config.
// temporary config is inferred from csv file via request body.
//
func addTmpObjConf(c echo.Context) error {
	bytes, err := io.ReadAll(c.Request().Body)
	if err == nil {
		if csvstr := string(bytes); gotk.IsCSV(csvstr) {
			model := c.QueryParam("modelName")
			if model == "" {
				model = fmt.Sprintf("tempmodel%v", time.Now().Unix())
			}
			if tc, err := n3context.AddTempContextConfigFromCSV(model, csvstr); err == nil {
				return c.String(http.StatusOK, tc)
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, "valid csv file is needed in request body")
	}
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

//
// handler allows user to clear all temporary object model classifier config
//
func clrTmpObjConf(c echo.Context) error {
	n3context.ClearTempContextConfig()
	return c.String(http.StatusOK, "Cleared")
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
	// to receive updates
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
		"message": "Context created and activated successfully. Use this token to publish: /n3/publish and query: /n3/graphql data. Token must be provided in Authorization: Bearer header.",
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
		// panic("Failed")
		return httpErr
	}

	// send in the data, via the crdt layer
	err := n3ctx.PublishFromHTTPRequest(c.Request())
	if err != nil {
		// XXX SCOTT
		// panic(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to publish data to context", err)
		// fmt.Println(err.(*errors.Error).ErrorStack())
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

	// log.Printf("\n\nquery:\n\n$%#v\n\n%#v\n", grq.Query, grq.Variables)

	// handle the query
	results, err := n3ctx.GQLQuery(grq.Query, grq.Variables)
	if err != nil {
		log.Println("gqlHandler error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Query failed:", err)
	}

	return c.JSON(http.StatusOK, results)
}
