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
	//http.HandleFunc("/api/users", middleware.WithCORS(handlers.GetUsers))
	http.HandleFunc("/api/users/", middleware.WithCORS(handlers.HandleUserByID)) // with trailing slash
	http.HandleFunc("/api/users/search", middleware.WithCORS(handlers.SearchUsers))
	//http.HandleFunc("/api/posts", middleware.WithCORS(handlers.GetPosts))
	http.HandleFunc("/api/posts/", middleware.WithCORS(handlers.HandlePostsByUserId))
	http.HandleFunc("/api/homefeed", middleware.WithCORS(handlers.GetFeedPosts))
	http.HandleFunc("/api/groups", middleware.WithCORS(handlers.GetGroups))
	http.HandleFunc("/api/groups/", middleware.WithCORS(handlers.HandleGroupsByUserId)) // groups with user id
	http.HandleFunc("/api/login", middleware.WithCORS(handlers.HandleLogin))
	http.HandleFunc("/api/register", middleware.WithCORS(handlers.HandleRegister))
	http.HandleFunc("/api/logout", middleware.WithCORS(handlers.HandleLogout))
	http.HandleFunc("/api/me", middleware.WithCORS(handlers.HandleMe))
	http.HandleFunc("/api/me/update", middleware.WithCORS(handlers.HandleUpdateMe))
	http.HandleFunc("/api/posts/create", middleware.WithCORS(handlers.HandleCreatePost))
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

	http.HandleFunc("/ws", middleware.WithCORS(handlers.HandleWSConnections))
	//http.HandleFunc("/ws", handlers.HandleWSConnections)
	http.HandleFunc("/api/comments/show", middleware.WithCORS(handlers.HandleCommentsForPost))
	http.HandleFunc("/api/comments/create", middleware.WithCORS(handlers.HandleCreateCommentsForPost))

	// Serve the avatars directory as static content with CORS
	fs := http.FileServer(http.Dir("./avatars"))
	avatarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.WithCORS(fs.ServeHTTP)(w, r)
	})
	http.Handle("/avatars/", http.StripPrefix("/avatars/", avatarHandler))
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
