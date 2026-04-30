package gin_helper

import (
	"core-ticket/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type RouterFunc func(router *Router)

type Router struct {
	ginEngine   *gin.Engine
	dbInstances map[string]*sqlx.DB
	middleware  middlewares.Middleware
}

/*
Engine initializes a new Router instance with the provided ginEngine, dbInstances, middleware, and s3Session.
It returns a pointer to the newly created Router.

Parameters:
- ginEngine: The gin.Engine instance to be used by the Router.
- dbInstances: A map containing database instances with string keys and *sqlx.DB values.
- middleware: An implementation of the Middleware interface for handling various middlewares.

Returns:
- *Router: A pointer to the newly created Router instance.

Example:

	router := Engine(myGinEngine, myDBInstances, myMiddleware)
*/
func Engine(ginEngine *gin.Engine, dbInstances map[string]*sqlx.DB, middleware middlewares.Middleware) *Router {
	return &Router{ginEngine, dbInstances, middleware}
}

/*
Group creates a new RouterGroup with the specified relative path and handlers.
It returns a pointer to the newly created RouterGroup.

Parameters:
- relativePath: a string representing the relative path for the new RouterGroup.
- handlers: variadic list of gin.HandlerFunc functions to be used as handlers for the new RouterGroup.

Returns:
- *RouterGroup: a pointer to the newly created RouterGroup.

Example:

	router := &Router{}
	group := router.Group("/api", handler1, handler2)
*/
func (router *Router) Group(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
	return &RouterGroup{
		ginRouterGroup: router.ginEngine.Group(relativePath, handlers...),
	}
}

/*
GinEngine returns the gin.Engine instance associated with the Router.
This allows access to the underlying gin.Engine for registering routes, middleware, and handling HTTP requests.

Returns:
- *gin.Engine: The gin.Engine instance used by the Router.
*/
func (router *Router) GinEngine() *gin.Engine {
	return router.ginEngine
}

/*
DBInstances returns a map containing references to all the sqlx.DB instances stored in the Router instance.
The keys of the map are strings representing the names of the database instances, and the values are pointers to the corresponding sqlx.DB instances.

Returns:
- map[string]*sqlx.DB: A map containing references to all the sqlx.DB instances stored in the Router instance.
*/
func (router *Router) DBInstances() map[string]*sqlx.DB {
	return router.dbInstances
}

/*
Middleware returns the middleware instance associated with the router.
This middleware instance provides methods for handling JWT tokens, authentication tokens, and whitelisting.
*/
func (router *Router) Middleware() middlewares.Middleware {
	return router.middleware
}

/*
RegisterRouter registers the provided router function by calling it with the current router instance.
The router function should accept a pointer to the current router instance as a parameter.

Parameters:
- routerFunc: The router function to be registered, which should have the signature func(router *Router).

Example:

	router.RegisterRouter(func(router *Router) {
		// Define routes and middleware using the provided router instance
	})
*/
func (router *Router) RegisterRouter(routerFunc RouterFunc) {
	routerFunc(router)
}
