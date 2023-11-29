package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	//rt.router.POST("/users", rt.doLogin)
	//rt.router.GET("/users", rt.getUserProfile)
	//rt.router.PUT("/users/:userID", rt.setMyUserName)
	rt.router.GET("/users/:userID/getMyStream", rt.getMyStream)

	rt.router.PUT("/users/:userID/followed", rt.followUser)
	rt.router.DELETE("/users/:userID/followed/:followedID", rt.unfollowUser)

	//rt.router.PUT("/users/:userID/banned", rt.banUser)
	//rt.router.DELETE("/users/:userID/banned/:bannedID", rt.unbanUser)

	rt.router.POST("/photos", rt.uploadPhoto)
	rt.router.DELETE("/photos/:photoID", rt.deletePhoto)

	rt.router.PUT("/photos/:photoID/likes", rt.likePhoto)
	//rt.router.GET("/photos/:photoID/likes/:likeID", rt.getLike)
	rt.router.DELETE("/photos/:photoID/likes/:likeID", rt.unlikePhoto)

	//rt.router.PUT("/photos/:photoID/comments", rt.commentPhoto)
	//rt.router.GET("/photos/:photoID/comments/:likeID", rt.getComment)
	//rt.router.DELETE("/photos/:photoID/comments/:likeID", rt.uncommentPhoto)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
