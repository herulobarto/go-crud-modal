package mahasiswacontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/herulobarto/go-crud-modal/entities"
	"github.com/herulobarto/go-crud-modal/models/mahasiswamodel"
)

var mahasiswaModel = mahasiswamodel.New()

func Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"data": template.HTML(GetData()),
	}

	temp, _ := template.ParseFiles("views/mahasiswa/index.html")
	temp.Execute(w, data)
}

func GetData() string {

	buffer := &bytes.Buffer{}

	temp, _ := template.New("data.html").Funcs(template.FuncMap{
		"increment": func(a, b int) int {
			return a + b
		},
	}).ParseFiles(("views/mahasiswa/data.html"))

	var mahasiswa []entities.Mahasiswa
	err := mahasiswaModel.FindAll(&mahasiswa)
	if err != nil {
		panic(err)
	}

	fmt.Println(mahasiswa)

	data := map[string]interface{}{
		"mahasiswa": mahasiswa,
	}

	temp.ExecuteTemplate(buffer, "data.html", data)

	return buffer.String()
}

func GetForm(w http.ResponseWriter, r *http.Request) {

	queryString := r.URL.Query()
	id, err := strconv.ParseInt(queryString.Get("id"), 10, 64)

	var data map[string]interface{}

	if err != nil {
		data = map[string]interface{}{
			"title": "Tambah Data Mahasiwa",
		}
	} else {
		var mahasiswa entities.Mahasiswa
		err := mahasiswaModel.Find(id, &mahasiswa)
		if err != nil {
			panic(err)
		}

		data = map[string]interface{}{
			"title":     "Edit Data Mahasiwa",
			"mahasiswa": mahasiswa,
		}
	}

	temp, _ := template.ParseFiles("views/mahasiswa/form.html")
	temp.Execute(w, data)
}

func Store(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		r.ParseForm()
		var mahasiswa entities.Mahasiswa

		mahasiswa.NamaLengkap = r.Form.Get("nama_lengkap")
		mahasiswa.JenisKelamin = r.Form.Get("jenis_kelamin")
		mahasiswa.TanggalLahir = r.Form.Get("tanggal_lahir")
		mahasiswa.TempatLahir = r.Form.Get("tempat_lahir")
		mahasiswa.Alamat = r.Form.Get("alamat")

		err := mahasiswaModel.Create(&mahasiswa)
		if err != nil {
			ResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		data := map[string]interface{}{
			"message": "Data berhasil disimpan",
			"data":    template.HTML(GetData()),
		}

		ResponseJson(w, http.StatusOK, data)
	}
}

func ResponseError(w http.ResponseWriter, code int, message string) {

	ResponseJson(w, code, map[string]string{"error": message})
}

func ResponseJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
