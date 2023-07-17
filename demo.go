package main

import (
	"bytes"
	"context"
	"demo_day_8/signal"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func DeleteAPI() []byte {
	// Địa chỉ URL API
	url := "https://jsonplaceholder.typicode.com/todos/1" // Ví dụ: Xóa todo có id = 1

	// Tạo request DELETE
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Tạo HTTP client và thực thi request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Đọc nội dung phản hồi
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Xóa thành công")
	} else {
		fmt.Println("Xóa thất bại")
	}

	// In ra nội dung phản hồi
	fmt.Println(string(body))
	return body
}
func homeHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home page"))
}

func aboutHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About page"))
}
func contactHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Contact page"))
}
func returnHtml(w http.ResponseWriter, _ *http.Request) {
	result := GetAPIExample()
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
func returnPost(w http.ResponseWriter, _ *http.Request) {
	result := PostAPIExample()
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
func returnDelete(w http.ResponseWriter, _ *http.Request) {
	result := DeleteAPI()
	result = GetAPIExample()
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func GetAPIExample() []byte {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos")
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Failed==>", err)
		return []byte{}
	} else {
		body, _ := io.ReadAll(resp.Body)
		return body
	}
}
func DetailAPIExample() []byte {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos")
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Failed==>", err)
		return []byte{}
	} else {
		body, _ := io.ReadAll(resp.Body)
		return body
	}
}
func PostAPIExample() []byte {
	jsonBody := []byte(`{
	"userId": 202,
    "title": "delectus aut autem",
    "completed": false
}
`)
	bodyReader := bytes.NewReader(jsonBody)
	resp, err := http.Post("https://jsonplaceholder.typicode.com/todos", "application/json", bodyReader)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Failed==>", err)
		return []byte{}
	} else {
		fmt.Println("Success")
		body, _ := io.ReadAll(resp.Body)
		return body
	}
	//resp.Header.Add("Content-Type", "application/json")

}
func main() {
	//http.HandleFunc("/", homeHandle)
	//http.HandleFunc("/about", aboutHandle)
	//http.HandleFunc("/contact", contactHandle)
	//http.HandleFunc("/get", returnHtml)
	//http.HandleFunc("/post", returnPost)
	//http.HandleFunc("/delete", returnDelete)
	//http.HandleFunc("/get/", getTodoDetail)
	//http.HandleFunc("/todos/", contextDemo)
	//fmt.Println("Server listenning on port 3000 ...")
	//fmt.Println(http.ListenAndServe(":3000", nil))
	//signal.DemoSignal()
	signal.Demo2()
}

type Todo struct {
	UserId    int    `json:"userId"`
	ID        int    `json:"ID"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func getTodoDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/get/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todo, err := fetchTodoDetail(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func fetchTodoDetail(id int) (*Todo, error) {
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	var todo Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}
func contextDemo(w http.ResponseWriter, r *http.Request) {
	// Tạo một context với deadline
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	//defer cancel()
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// Gọi hàm lấy detail thông tin với context
	todo, err := getTodoDetailContext(ctx, id)
	if err != nil {
		log.Fatal(err)
	}
	// In ra thông tin todo
	fmt.Printf("Todo ID: %d\n", todo.ID)
	fmt.Printf("Todo Title: %s\n", todo.Title)
	fmt.Printf("Todo Completed: %t\n", todo.Completed)
	if context.Cause(ctx) != nil {
		fmt.Println(context.Cause(ctx))
		json.NewEncoder(w).Encode("Request timeout!")
	} else {
		json.NewEncoder(w).Encode(todo)

	}
}

func getTodoDetailContext(ctx context.Context, todoID int) (*Todo, error) {
	// Tạo một request mới với context
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", todoID), nil)
	if err != nil {
		return nil, err
	}

	// Tạo một HTTP client
	client := &http.Client{}

	// Gửi request và nhận phản hồi
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get todo detail: %s", resp.Status)
	}

	// Đọc phản hồi JSON và chuyển đổi thành struct Todo
	var todo Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	if err != nil {
		return nil, err
	}
	time.Sleep(7 * time.Second)
	return &todo, nil
}
