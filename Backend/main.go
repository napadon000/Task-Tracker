package main

import (
	// "log"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"project/configs"
	"project/routes"
	"time"

	"project/models"

	// dotenv "github.com/dsh2dsh/expx-dotenv"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/cors"
)

var client *http.Client

func init() {
	// Create a cookie jar to automatically handle cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	// Create HTTP client with cookie jar
	client = &http.Client{
		Jar: jar,
	}
}

const URL = "http://localhost:3222/"

func main() {
	configs.ConnectDatabase()
	app := fiber.New()

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "http://localhost:3221",
	// 	AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	// 	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	// }))

	userGroup := app.Group("/api")
	routes.UserRoute(userGroup)
	routes.TaskRoute(userGroup)

	go app.Listen(":3222")
	time.Sleep(2 * time.Second) // Wait for the server to start
	for {
		//receive input command
		var input string
		fmt.Print("Enter command (type 'exit' to quit):")
		fmt.Scanln(&input)

		switch input {
		case "exit":
			return
		case "register":
			var email, password string
			fmt.Print("Enter email: ")
			fmt.Scanln(&email)
			fmt.Print("Enter password: ")
			fmt.Scanln(&password)

			rsp, err := client.PostForm(URL+"api/user/register", url.Values{
				"email":    {email},
				"password": {email},
			})
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			defer rsp.Body.Close()

		case "login":
			var email, password string
			fmt.Print("Enter email: ")
			fmt.Scanln(&email)
			fmt.Print("Enter password: ")
			fmt.Scanln(&password)

			rsp, err := client.PostForm(URL+"api/user/login", url.Values{
				"email":    {email},
				"password": {password},
			})
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			defer rsp.Body.Close()

			body, err := io.ReadAll(rsp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				continue
			}

			fmt.Println(string(body))

		case "getTasks":
			rsp, err := client.Get(URL + "api/task/gettasks")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			defer rsp.Body.Close()

			var tasks []models.Task

			//read the response body
			body, err := io.ReadAll(rsp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				continue
			}

			err = json.Unmarshal(body, &tasks)
			if err != nil {
				fmt.Println("Error unmarshalling response body:", err)
				continue
			}

			for i, task := range tasks {
				fmt.Println("Task", i+1)
				fmt.Println("ID:", task.ID.Hex())
				fmt.Println("Title:", task.Title)
				fmt.Println("Description:", task.Description)
				fmt.Println("Status:", task.Status)
				fmt.Println("============================")
			}

		case "addTask":
			var title, description, status string
			fmt.Print("Enter task title: ")
			fmt.Scanln(&title)
			fmt.Print("Enter task description: ")
			fmt.Scanln(&description)
			fmt.Print("Enter task status (todo, in_progress, done): ")
			fmt.Scanln(&status)

			rsp, err := client.PostForm(URL+"api/task/create", url.Values{
				"title":       {title},
				"description": {description},
				"status":      {status},
			})
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			defer rsp.Body.Close()

			body, err := io.ReadAll(rsp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				continue
			}

			fmt.Println("Task added:", string(body))

		case "deleteTask":
			var taskID string
			fmt.Print("Enter task ID to delete: ")
			fmt.Scanln(&taskID)

			req, err := http.NewRequest("DELETE", URL+"api/task/delete/"+taskID, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}

			rsp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			defer rsp.Body.Close()

			body, err := io.ReadAll(rsp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				continue
			}

			fmt.Println("Task deleted:", string(body))

		case "changeTaskStatus":
			var taskID, status string
			fmt.Print("Enter task ID to change status: ")
			fmt.Scanln(&taskID)
			fmt.Print("Enter new status (todo, in_progress, done): ")
			fmt.Scanln(&status)

			data := map[string]string{
				"status": status,
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				continue
			}

			req, err := http.NewRequest("PATCH", URL+"api/task/updatestatus/"+taskID, bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}

			req.Header.Set("Content-Type", "application/json")

			rsp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			defer rsp.Body.Close()

			body, err := io.ReadAll(rsp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				continue
			}
			fmt.Println("Status of respond:", rsp.Status)
			fmt.Println("Task status changed:", string(body))
			fmt.Println(taskID)

		case "changeTaskDescription":
			var taskID, description string
			fmt.Print("Enter task ID to change description: ")
			fmt.Scanln(&taskID)
			fmt.Print("Enter new description: ")
			fmt.Scanln(&description)

			data := map[string]string{
				"description": description,
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				continue
			}

			req, err := http.NewRequest("PATCH", URL+"api/task/updatetask/"+taskID, bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}

			req.Header.Set("Content-Type", "application/json")

			rsp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			defer rsp.Body.Close()

			body, err := io.ReadAll(rsp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				continue
			}

			fmt.Println("Task description changed:", string(body))

		default:
			fmt.Println("Unknown command:")
		}
	}

}
