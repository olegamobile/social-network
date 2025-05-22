package main

import (
	"backend/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func setHandlers() {
	// CORS to allow developement on same address
	http.HandleFunc("/api/users/", middleware.WithCORS(handlers.HandleUserByID))
	http.HandleFunc("/api/users/search", middleware.WithCORS(handlers.SearchUsers))
	http.HandleFunc("/api/posts/", middleware.WithCORS(handlers.HandlePostsByUserId))
	http.HandleFunc("/api/posts/create", middleware.WithCORS(handlers.HandleCreatePost))
	http.HandleFunc("/api/group/posts/", middleware.WithCORS(handlers.HandlePostsByGroupId))
	http.HandleFunc("/api/group/members/", middleware.WithCORS(handlers.HandleMembersByGroupId))
	http.HandleFunc("/api/group/events/", middleware.WithCORS(handlers.HandleEventsByGroupId))
	http.HandleFunc("/api/homefeed", middleware.WithCORS(handlers.GetFeedPosts))

	http.HandleFunc("/api/suggestgroups", middleware.WithCORS(handlers.HandleSuggestGroups))
	http.HandleFunc("/api/groups/search", middleware.WithCORS(handlers.SearchGroups))
	http.HandleFunc("/api/groups/user/", middleware.WithCORS(handlers.HandleGroupsByUserId))            // groups with user id
	http.HandleFunc("/api/groups/requested", middleware.WithCORS(handlers.HandleGroupRequests))         // active user group requests
	http.HandleFunc("/api/groups/invitations", middleware.WithCORS(handlers.HandleGroupInvitations))    // active user group invitations
	http.HandleFunc("/api/groups/administered", middleware.WithCORS(handlers.HandleGroupsAdministered)) // active user group invitations
	http.HandleFunc("/api/group/", middleware.WithCORS(handlers.HandleGroupById))                       // group with group id
	http.HandleFunc("/api/group/join", middleware.WithCORS(handlers.HandleGroupMembership))

	http.HandleFunc("/api/group-posts/create", middleware.WithCORS(handlers.CreateGroupPostHandler))

	http.HandleFunc("/api/login", middleware.WithCORS(handlers.HandleLogin))
	http.HandleFunc("/api/register", middleware.WithCORS(handlers.HandleRegister))
	http.HandleFunc("/api/logout", middleware.WithCORS(handlers.HandleLogout))
	http.HandleFunc("/api/me", middleware.WithCORS(handlers.HandleMe))
	http.HandleFunc("/api/me/update", middleware.WithCORS(handlers.HandleUpdateMe))
	http.HandleFunc("/api/following/", middleware.WithCORS(handlers.HandleFollowing))
	http.HandleFunc("/api/follow", middleware.WithCORS(handlers.HandleFollowAction))
	http.HandleFunc("/api/followers/", middleware.WithCORS(handlers.GetFollowers))
	http.HandleFunc("/api/followed/", middleware.WithCORS(handlers.GetFollowedUsers))
	http.HandleFunc("/api/follow/requests/sent", middleware.WithCORS(handlers.GetSentFollowRequests))
	http.HandleFunc("/api/follow/requests/received", middleware.WithCORS(handlers.GetReceivedFollowRequests))
	http.HandleFunc("/api/follow/requests/{id}/accept", middleware.WithCORS(handlers.HandleFollowRequestApprove))
	http.HandleFunc("/api/follow/requests/{id}/decline", middleware.WithCORS(handlers.HandleFollowRequestApprove))
	http.HandleFunc("/api/suggest/users", middleware.WithCORS(handlers.GetSuggestedUsers))

	http.HandleFunc("/api/notifications", middleware.WithCORS(handlers.GetNotifications))
	http.HandleFunc("/api/notifications/{id}", middleware.WithCORS(handlers.GetNotificationByID))
	http.HandleFunc("/api/notifications/{id}/read", middleware.WithCORS(handlers.ReadNotification))
	http.HandleFunc("/api/notifications/new", middleware.WithCORS(handlers.GetNewNotifications))
	http.HandleFunc("/api/notifications/{id}/joingroup", middleware.WithCORS(handlers.HandleJoinReqsByGroupId))

	http.HandleFunc("/ws", middleware.WithCORS(handlers.HandleWSConnections))
	//http.HandleFunc("/ws", handlers.HandleWSConnections)

	// Serve the avatars directory as static content with CORS
	avatarsFS := http.FileServer(http.Dir("./avatars"))
	avatarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.WithCORS(avatarsFS.ServeHTTP)(w, r)
	})
	http.Handle("/avatars/", http.StripPrefix("/avatars/", avatarHandler))

	postsFS := http.FileServer(http.Dir("./uploads/posts"))
	postImageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.WithCORS(postsFS.ServeHTTP)(w, r)
	})
	http.Handle("/uploads/posts/", http.StripPrefix("/uploads/posts/", postImageHandler))

}

func main() {
	config.InitConfig()

	err := database.NewDatabase(config.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	setHandlers()
	fmt.Printf("Backend running on port %s, allowing requests from %s\n", config.Port, config.FrontendURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
