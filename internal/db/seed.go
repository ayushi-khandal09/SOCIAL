package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"math/rand"

	"github.com/ayushi-khandal09/social/internal/store"
)

var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve",
	"frank", "grace", "heidi", "ivan", "judy",
	"karl", "laura", "mike", "nancy", "oscar",
	"peggy", "quinn", "ruth", "sam", "tina",
	"ursula", "victor", "wendy", "xavier", "yvonne",
	"zach", "abby", "ben", "cindy", "dan",
	"ella", "finn", "gina", "harry", "irene",
	"jack", "kate", "leo", "maya", "nina",
	"oliver", "paula", "quincy", "ryan", "sara",
	"tom", "una", "vince", "will", "zoe",
}

var titles = []string{
	"Hello World", "Learning Go", "Web Dev Tips", "AI Basics",
	"My First Project", "Why I Code", "Daily Routine", "Tech Trends",
	"Go vs Python", "Build and Deploy", "UI Design", "Clean Code",
	"Git Essentials", "Docker Guide", "Code Faster", "Cloud 101",
	"Simple Life", "React Hooks", "Debug Tricks", "Career Tips",
	"Study Plan", "REST API", "Go for Web", "Design Patterns",
	"Tiny Projects", "Fast Backend", "UX Myths", "Remote Work",
	"Daily Standup", "Weekend Build", "Happy Coding", "Go Modules",
	"Next.js Intro", "TypeScript Fun", "AI for All",
	"Minimal UI", "Write Better", "Code Review", "Learning Curve",
	"Big O", "Test Driven", "Bug Fixes", "My Portfolio", "Launch Day",
	"Build Tools", "Open Source", "Learn by Doing",
	"Weekly Notes", "Quick Tips", "From Zero",
}

var contents = []string{
	"Just started coding!",
	"Coffee and code ☕",
	"Bug fixed 🚀",
	"New project dropped!",
	"Learning Go today.",
	"Design looks clean.",
	"What a day!",
	"Deploy done ✅",
	"Trying dark mode.",
	"Build. Fail. Repeat.",
	"Tiny win today!",
	"Code. Eat. Sleep.",
	"Version 1.0 live!",
	"Docs updated!",
	"API is working.",
	"Nice UI tweak.",
	"Dev life 😅",
	"Refactoring code.",
	"Let’s ship it!",
	"More to come...",
}

var tags = []string{
	"go",
	"webdev",
	"ai",
	"productivity",
	"coding",
	"startup",
	"design",
	"cloud",
	"javascript",
	"react",
	"python",
	"opensource",
	"backend",
	"frontend",
	"devlife",
	"tutorial",
	"machinelearning",
	"tech",
	"daily",
	"tips",
}

var comments = []string{
	"Great post!",
	"Thanks for sharing.",
	"Very helpful.",
	"I learned something new.",
	"Can you explain more?",
	"Awesome content!",
	"This is gold 🔥",
	"Totally agree.",
	"Nice and simple.",
	"Saved it for later.",
	"Mind blown 🤯",
	"I needed this today.",
	"Keep it up!",
	"Really well written.",
	"Interesting point!",
	"Loved this!",
	"Super useful.",
	"Thanks a ton!",
	"Where can I find more?",
	"Following for more.",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}
	log.Println("Seeding Complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		var pw store.Password
		if err := pw.Set("123123"); err != nil {
			log.Fatalf("hash password: %v", err)
		}

		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: pw, // ✅ correct type
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserId:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))], // ✅ use len(comments)
		}
	}
	return cms
}
