package main

import (
        "fmt"
        "github.com/gorilla/sessions"
      _ "github.com/lib/pq"
        "database/sql"
        "html/template"
        "net/http"
        //"os"
        )
type  rdetail struct {
	Id int64
	Firstname string
	Lastname string
	Email string
	Password string
 }

func main() {

	http.HandleFunc("/" , Regestration)
	http.HandleFunc("/login" , Login)
	http.HandleFunc("/confm" , Confirm)
	http.HandleFunc("/hm" , Home)
	http.HandleFunc("/slogin" , SessionLogin)
	http.HandleFunc("/slogout" , SessionLogout)
	http.ListenAndServe(":8888" , nil)
}

func Dbconn() *sql.DB {

	const (
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  password = "test123"
  dbname   = "postgres"
)

	psqlInfo := fmt.Sprintf("host=%s  port=%d  user=%s "+
        "password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }
   // defer db.Close()

    err = db.Ping()
    if err != nil {
        panic(err)
    }

    fmt.Println("Successfully connected!")
    return db


}

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func Regestration(w http.ResponseWriter, r *http.Request) {

	t,_ := template.ParseFiles("Register.html")
	t.Execute(w , nil)
}

func Login(w http.ResponseWriter, r *http.Request) {

  t,_ := template.ParseFiles("login.html")
	t.Execute(w , nil)
}

func SessionLogin(w http.ResponseWriter,r *http.Request) {
	 session, _ := store.Get(r, "session")
	 session.Values["authenticated"] = true
	 session.Save(r,w)
	 http.Redirect(w,r,"/hm",307)
}

func SessionLogout(w http.ResponseWriter,r *http.Request) {
	session, _ := store.Get (r, "session")
	session.Values["authenticated"] = false
	session.Save(r,w)
	http.Redirect(w,r,"/login",307)
}

func Confirm(w http.ResponseWriter, r *http.Request) {

  db2 := Dbconn()

	data := new(rdetail)
	data.Firstname = r.FormValue("fname")
	data.Lastname = r.FormValue("lname")
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("pwd")

  var dta rdetail
  var total []rdetail
	// fmt.Println(data.Firstname)
	// fmt.Println(data.Lastname)
	// fmt.Println(data.Email)
	// fmt.Println(data.Password)
	sql := "select email from rgstr_detail"
	rows,_ := db2.Query(sql)

	for rows.Next() {
		er := rows.Scan(&dta.Email)
		if er != nil {
			fmt.Println(er)
		}
		total =append(total, dta)

	}

	flag := true

	for i := 0; i < len(total); i++ {
		if total[i].Email==data.Email {
			flag = false
		}
	}

	//fmt.Println(total)

	if flag==false {
   fmt.Fprintf(w, "duplicate data not entered")

	}else {
  stmt := "INSERT INTO rgstr_detail( fname, lname, email, password) VALUES ( $1, $2, $3, $4)"
   _, err := db2.Exec(stmt, data.Firstname, data.Lastname, data.Email, data.Password)
    if err != nil {
        panic(err)
    }
    t,_ := template.ParseFiles("confirm.html")
	  t.Execute(w , nil)
  }
	//fmt.Println(data)

}

func Home(w http.ResponseWriter, r *http.Request) {
  session,_ := store.Get(r , "session")
  if auth ,ok := session.Values["authenticated"].(bool) ; !ok || !auth{
  	http.Redirect(w,r,"/login",307)
  }
  db2 := Dbconn()

  Email :=r.FormValue("email")
  pwd := r.FormValue("pwd")

  fmt.Println("without database>>>>>>", Email)
  fmt.Println("without database>>>>>>>>", pwd)

  var d rdetail

  stmt := "select email,password from rgstr_detail where email = $1"
  rows , _ := db2.Query(stmt, Email)

  for rows.Next() {
  	err := rows.Scan(&d.Email,&d.Password)
    if err!= nil {
    	panic(err)
    }
  }

  fmt.Println(">>>>>>>>>>>>>", d.Email)
  fmt.Println(">>>>>>>>>>>>",d.Password)
  if Email==d.Email && pwd == d.Password {
  	t,_ := template.ParseFiles("home.html")
	  t.Execute(w , nil)
  }else {
     session.Values["authenticated"] = false
     session.Save(r,w)
     http.Redirect(w,r,"/login",307)
  }

}
