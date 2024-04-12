package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// model
type Courses struct {
	CourseId    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"-"`
	Author      *Author `json:"author"`
}
type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

// fakedb
var courseDB []Courses

// middleware
func (c *Courses) IsEmpty() bool {
	return c.CourseName == ""
}

func main() {
	fmt.Println("API Started")
	r := mux.NewRouter()

	// seeding
	courseDB = append(courseDB, Courses{CourseId: "20", CourseName: "react", CoursePrice: 299, Author: &Author{Fullname: "ram", Website: "loco.dev"}})
	courseDB = append(courseDB, Courses{CourseId: "30", CourseName: "mern", CoursePrice: 399, Author: &Author{Fullname: "vijay", Website: "vijay.dev"}})

	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")
	r.HandleFunc("/course", createOneCourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateOneCourse).Methods("PUT")
	r.HandleFunc("/course/{id}", deleteOneCourse).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":4000", r))
}

// controller
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to api</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all courses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courseDB)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, course := range courseDB {
		if course.CourseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("No course found")
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var course Courses
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if course.IsEmpty() {
		http.Error(w, "No data inside JSON", http.StatusBadRequest)
		return
	}

	rand.Seed(time.Now().UnixNano())
	course.CourseId = strconv.Itoa(rand.Intn(100))
	courseDB = append(courseDB, course)
	json.NewEncoder(w).Encode(course)
}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var updatedCourse Courses
	err := json.NewDecoder(r.Body).Decode(&updatedCourse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, course := range courseDB {
		if course.CourseId == params["id"] {
			courseDB[i] = updatedCourse
			updatedCourse.CourseId = params["id"]
			json.NewEncoder(w).Encode(updatedCourse)
			return
		}
	}
	http.Error(w, "Course not found", http.StatusNotFound)
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, course := range courseDB {
		if course.CourseId == params["id"] {
			courseDB = append(courseDB[:i], courseDB[i+1:]...)
			json.NewEncoder(w).Encode("Course deleted")
			return
		}
	}
	http.Error(w, "Course not found", http.StatusNotFound)
}
