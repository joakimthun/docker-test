package endpoints

import (
    "net/http"
    "encoding/json"
    "io/ioutil"
    "log"
    "github.com/gorilla/mux"
    "github.com/gorilla/schema"
    "github.com/joakimthun/docker-test/db"
    "github.com/joakimthun/docker-test/redis"
    "fmt"
)

var http_requests_total_index int = 0

func init() {
    log.Println("Endpoints init")
    
    r := mux.NewRouter()
    
    r.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        http_requests_total_index += 1
        writeHTML(w, r, getView("views/index"))
    }).Methods("GET")
    
    r.HandleFunc("/users", getUsers).Methods("GET")
    
    r.HandleFunc("/createuser", func (w http.ResponseWriter, r *http.Request) {
        writeHTML(w, r, getView("views/user"))
    }).Methods("GET")
    
    r.HandleFunc("/createuser", postUser).Methods("POST")
    
    r.HandleFunc("/redisget/{key:(.)*}", redisGet).Methods("GET")
    
    r.HandleFunc("/redisset", func (w http.ResponseWriter, r *http.Request) {
        writeHTML(w, r, getView("views/redis"))
    }).Methods("GET")
    
    r.HandleFunc("/redisset", redisSet).Methods("POST")
    r.HandleFunc("/metrics", prometheusMetrics).Methods("GET")    
    
    http.Handle("/", r)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    users, err := db.Users()
    
    if err != nil {
        http.Error(w, "Server error", 500)
        return
    }
    
    writeJSON(w, users)
}

func postUser(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()

    if err != nil {
        http.Error(w, "Server error", 500)
        return
    }
    
    u := db.User{}
    decoder := schema.NewDecoder()
    err = decoder.Decode(&u, r.PostForm)

    if err != nil {
        http.Error(w, "Server error", 500)
        return
    }
    
    if u.Email == "" || u.Name == "" {
        http.Error(w, "Bad request", 400)
        return
    }
    
    err = db.Create(&u)
    
    if err != nil {
        http.Error(w, "Server error", 500)
        return
    }
    
    log.Println("Creating user:")
    log.Printf("Name: %s Email: %s", u.Name, u.Email)
    
    http.Redirect(w, r, "/users", 302)
}

type redisModel struct {
    Key string
    Value string
}

func redisGet(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
	key := vars["key"]
    
    log.Println("Fetching from redis with key:")
    log.Printf(key)
    
    m := redisModel{}
    err := redis.Get(key, m)
    
    if err != nil {
        log.Println(err)
        http.Error(w, "Could not get the value from redis", 500)
        return
    }
    
    writeJSON(w, m)
}

func prometheusMetrics(w http.ResponseWriter, r *http.Request) {
     fmt.Fprintf(w, "# HELP http_requests_total_index Total number of http requests made to the index page\n# TYPE http_requests_total_index COUNTER\nhttp_requests_total_index %v\n", http_requests_total_index)
}

func redisSet(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()

    if err != nil {
        http.Error(w, "Server error", 500)
        return
    }
    
    m := redisModel{}
    decoder := schema.NewDecoder()
    err = decoder.Decode(&m, r.PostForm)

    if err != nil {
        http.Error(w, "Server error", 500)
        return
    }
    
    if m.Key == "" || m.Value == "" {
        http.Error(w, "Bad request", 400)
        return
    }
    
    err = redis.Set(m.Key, &m)
   
    if err != nil {
        http.Error(w, "Server error", 500)
        return
    }
    
    log.Println("Redis set:")
    log.Printf("Key: %s value: %s", m.Key, m.Value)
   
    
    http.Redirect(w, r, "/redisget/" + m.Key, 302)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	response, _ := json.Marshal(v)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(response))
}

func writeHTML(w http.ResponseWriter, r *http.Request, v *view) {
	w.Header().Add("Content-Type", "text/html")
	w.Write(v.Content)
}

type view struct {
	Content []byte
}

func getView(name string) *view {
	b, err := ioutil.ReadFile(name + ".html")
	
	if err != nil {
		log.Fatal(err)
	}
	
    return &view{ Content: b }
}