package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

type server struct {
	db *sql.DB
}

type Car struct {
	ID         int
	Model      string
	Year       int
	VehicleNum int
	MarketNum  int
	Vehicle    string
	Market     string
}

func database() server {
	database, _ := sql.Open("sqlite3", "lib/db.db")
	server := server{db: database}
	return server
}

func (s *server) updatemarket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		t, _ := template.ParseFiles("static/html/upm.html")
		t.Execute(w, nil)
		return
	}
	id := r.FormValue("id")
	name := r.FormValue("name")
	if _, err := s.db.Exec("update markets set market=$1 where id=$2", name, id); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *server) deletemarket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		t, _ := template.ParseFiles("static/html/delm.html")
		t.Execute(w, nil)
		return
	}
	name := r.FormValue("name")
	if _, err := s.db.Exec("delete from markets where market=$1", name); err != nil {
		log.Fatal("delete", err)
	}
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *server) addmarket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		t, _ := template.ParseFiles("static/html/addm.html")
		t.Execute(w, nil)
		return
	}
	name := r.FormValue("name")
	if _, err := s.db.Exec("insert into markets(market) values($1)", name); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *server) updatevehicle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		t, _ := template.ParseFiles("static/html/upv.html")
		t.Execute(w, nil)
	}
	id := r.FormValue("id")
	name := r.FormValue("name")
	if _, err := s.db.Exec("update vehicles set vehicle=$1 where id=$2", name, id); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *server) deletevehicle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		t, _ := template.ParseFiles("static/html/delv.html")
		t.Execute(w, nil)
		return
	}
	name := r.FormValue("name")
	if _, err := s.db.Exec("delete from vehicles where vehicle=$1", name); err != nil {
		log.Fatal("delete", err)
	}
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *server) addvehicle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		t, _ := template.ParseFiles("static/html/addv.html")
		t.Execute(w, nil)
		return
	}
	name := r.FormValue("name")
	if _, err := s.db.Exec("insert into vehicles(vehicle) values($1)", name); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *server) updatecar(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		t, _ := template.ParseFiles("static/html/upc.html")
		t.Execute(w, nil)
		return
	}
	id := r.FormValue("id")
	model := r.FormValue("model")
	year := r.FormValue("year")
	vehicle := r.FormValue("vehicle")
	market := r.FormValue("market")
	if _, err := s.db.Exec("update cars set model=$1, year=$2, vehicle=$3, market=$4 where id=$5", model, year, vehicle, market, id); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *server) deletecar(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		t, _ := template.ParseFiles("static/html/delc.html")
		t.Execute(w, nil)
		return
	}
	id := r.FormValue("id")
	if _, err := s.db.Exec("delete from cars where id=$1", id); err != nil {
		log.Fatal("delete", err)
	}
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *server) addcar(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		t, _ := template.ParseFiles("static/html/addc.html")
		t.Execute(w, nil)
		return
	}
	model := r.FormValue("model")
	year := r.FormValue("year")
	vehicle := r.FormValue("vehicle")
	market := r.FormValue("market")
	if _, err := s.db.Exec("insert into cars(model, year, vehicle, market) values($1, $2, $3, $4)", model, year, vehicle, market); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *server) listPage(w http.ResponseWriter, r *http.Request) {
	res, err := s.db.Query("select * from cars;")
	if err != nil {
		log.Fatal("Select", err)
	}

	var cars []Car
	for res.Next() {
		var car Car
		res.Scan(&car.ID, &car.Model, &car.Year, &car.VehicleNum, &car.MarketNum)
		if err := s.db.QueryRow("select vehicle from vehicles where id=$1", car.VehicleNum).Scan(&car.Vehicle); err != nil {
			log.Fatal(err)
		}
		if err := s.db.QueryRow("select market from markets where id=$1", car.MarketNum).Scan(&car.Market); err != nil {
			log.Fatal(err)
		}
		cars = append(cars, car)
	}
	fmt.Print(cars)

	t, _ := template.ParseFiles("static/html/list.html")
	t.Execute(w, cars)
}

func main() {
	s := database()
	defer s.db.Close()
	server := http.FileServer(http.Dir("./static"))
	http.Handle("/", server)
	http.HandleFunc("/list", s.listPage)

	http.HandleFunc("/addc", s.addcar)
	http.HandleFunc("/delc", s.deletecar)
	http.HandleFunc("/upc", s.updatecar)

	http.HandleFunc("/addv", s.addvehicle)
	http.HandleFunc("/delv", s.deletevehicle)
	http.HandleFunc("/upv", s.updatevehicle)

	http.HandleFunc("/addm", s.addmarket)
	http.HandleFunc("/delm", s.deletemarket)
	http.HandleFunc("/upm", s.updatemarket)

	http.ListenAndServe(":8000", nil)

}
