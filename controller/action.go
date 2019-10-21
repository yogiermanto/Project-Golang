package controller

import (
	"database/sql"
	"net/http"
	"text/template"

	c "project/config"
	m "project/model"
)

func Index(w http.ResponseWriter, r *http.Request) {

	var tmpl = template.Must(template.ParseFiles(
		"views/index.tmpl",
		"views/_header.tmpl",
		"views/_navbar.tmpl",
		"views/_sidebar.tmpl",
		"views/_footer.tmpl",
		"views/_control-sidebar.tmpl",
		"views/_script.tmpl",
	))

	// fmt.Println(ResMeja)
	var errs = tmpl.ExecuteTemplate(w, "index", nil)
	if errs != nil {
		http.Error(w, errs.Error(), http.StatusInternalServerError)
	}
}

func CustQueue(w http.ResponseWriter, r *http.Request) {
	// var data = M{"name": "Batman"}
	var tmpl = template.Must(template.ParseFiles(
		"views/cust-queue.tmpl",
		"views/_header.tmpl",
		"views/_navbar.tmpl",
		"views/_sidebar.tmpl",
		"views/_footer.tmpl",
		"views/_control-sidebar.tmpl",
		"views/_script.tmpl",
	))

	var err = tmpl.ExecuteTemplate(w, "cust-queue", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AdminQueue(w http.ResponseWriter, r *http.Request) {
	db := c.DbConn()
	MejaDB, err := db.Query("SELECT * FROM meja")
	var tmpl = template.Must(template.ParseFiles(
		"views/admin-queue.tmpl",
		"views/_header.tmpl",
		"views/_navbar.tmpl",
		"views/_sidebar.tmpl",
		"views/_footer.tmpl",
		"views/_control-sidebar.tmpl",
		"views/_script.tmpl",
	))

	Meja := m.Meja{}
	ResMeja := []m.Meja{}

	for MejaDB.Next() {
		var id_meja int
		var status bool
		err = MejaDB.Scan(&id_meja, &status)
		if err != nil {
			panic(err.Error())
		}
		Meja.IdMeja = id_meja
		Meja.Status = status
		ResMeja = append(ResMeja, Meja)
	}

	// untuk loop customer
	result, errS := db.Query("Select COUNT(*) FROM customer where status = false")
	TotalAntrian := CheckCount(result)
	if errS != nil {
		panic(errS.Error())
	}
	var no_antrian int
	resultMax, _ := db.Query("Select max(no_antrian) FROM customer")
	for resultMax.Next() {
		resultMax.Scan(&no_antrian)
	}

	resAdmin := m.ResAdmin{}
	resAdmin.Meja = ResMeja
	resAdmin.Total = TotalAntrian
	resAdmin.AntrianSekarang = no_antrian
	var errs = tmpl.ExecuteTemplate(w, "admin-queue", resAdmin)
	if errs != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CheckCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err.Error())
		}
	}
	return count
}

func TambahAntrian(w http.ResponseWriter, r *http.Request) {
	db := c.DbConn()
	if r.Method == "POST" {
		Nama := r.FormValue("nama")
		NoHp := r.FormValue("no_hp")
		result, errS := db.Query("Select COUNT(no_antrian) FROM customer")
		NoAntrian := CheckCount(result)
		if errS != nil {
			panic(errS.Error())
		}
		NoAntrian++
		statusmeja := CekMeja()
		_, err := db.Exec("INSERT INTO customer(nama, no_hp, no_antrian, status) VALUES(?,?,?,?)", Nama, NoHp, NoAntrian, statusmeja)
		if err != nil {
			panic(err.Error())
		}
		// log.Println("INSERT: Nama: " + Nama + " | No. Hp: " + NoHp + "| No. Antrian: " + NoAntrian)
	}
	defer db.Close()
	http.Redirect(w, r, "/admin-queue", 301)
}

func CekMeja() bool {
	db := c.DbConn()
	var StatusMeja bool
	CekMejaDB, err := db.Query("SELECT * FROM meja")
	CekMeja := m.Meja{}
	CekResMeja := []m.Meja{}

	for CekMejaDB.Next() {
		var id_meja int
		var status bool
		err = CekMejaDB.Scan(&id_meja, &status)
		if err != nil {
			panic(err.Error())
		}
		CekMeja.IdMeja = id_meja
		CekMeja.Status = status
		CekResMeja = append(CekResMeja, CekMeja)
	}

	for _, cekResMeja := range CekResMeja {
		if cekResMeja.Status == false {
			// fmt.Println("meja ke", cekResMeja.IdMeja, "berubah")
			_, err := db.Exec("UPDATE meja set status = 1 where id_meja = ?", cekResMeja.IdMeja)
			if err != nil {
				panic(err.Error())
			}
			StatusMeja = true
			break
		}
	}

	return StatusMeja
}

func UpdateMeja(w http.ResponseWriter, r *http.Request) {
	db := c.DbConn()
	emp := r.URL.Query().Get("id")
	result, errS := db.Query("Select COUNT(no_antrian) FROM customer where status = false")
	TotalAntrian := CheckCount(result)
	if errS != nil {
		panic(errS.Error())
	}
	if TotalAntrian > 0 {
		for i := 0; i < TotalAntrian; i++ {
			_, err := db.Exec("UPDATE customer set status = 1 where no_antrian = (select min(no_antrian) from customer where status = 0) ")
			if err != nil {
				panic(err.Error())
			}
			defer db.Close()
			http.Redirect(w, r, "/admin-queue", 301)
			break
		}
	} else {
		_, err := db.Exec("UPDATE meja set status = 0 where id_meja = ?", emp)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		http.Redirect(w, r, "/admin-queue", 301)
	}

}
