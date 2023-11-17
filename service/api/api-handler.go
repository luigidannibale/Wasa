package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.PUT("/users/:userID/followed", rt.followUser)
	rt.router.DELETE("/users/:userID/followed/:followedID", rt.unfollowUser)
	
	rt.router.PUT("/users/:userID/banned", rt.banUser)
	rt.router.DELETE("/users/:userID/banned/:bannedID", rt.unbanUser)
	
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
