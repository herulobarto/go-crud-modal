package mahasiswamodel

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/herulobarto/go-crud-modal/config"
	"github.com/herulobarto/go-crud-modal/entities"
)

type MahasiswaModel struct {
	db *sql.DB
}

func New() *MahasiswaModel {
	db, err := config.DBConnection()
	if err != nil {
		panic(err)
	}
	return &MahasiswaModel{db: db}
}

func (m *MahasiswaModel) FindAll(mahasiswa *[]entities.Mahasiswa) error {
	rows, err := m.db.Query("select * from mahasiswa")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var data entities.Mahasiswa
		rows.Scan(
			&data.Id,
			&data.NamaLengkap,
			&data.JenisKelamin,
			&data.TempatLahir,
			&data.TanggalLahir,
			&data.Alamat)

		*mahasiswa = append(*mahasiswa, data)
	}

	return nil
}

func (m *MahasiswaModel) Create(mahasiswa *entities.Mahasiswa) error {
	result, err := m.db.Exec("Insert into mahasiswa (nama_lengkap, jenis_kelamin, tempat_lahir, tanggal_lahir, alamat) values(?,?,?,?,?)",
		mahasiswa.NamaLengkap, mahasiswa.JenisKelamin, mahasiswa.TempatLahir, mahasiswa.TanggalLahir, mahasiswa.Alamat)

	if err != nil {
		return err
	}

	lastInsertId, _ := result.LastInsertId()
	mahasiswa.Id = lastInsertId
	return nil
}

func (m *MahasiswaModel) Find(id int64, mahasiswa *entities.Mahasiswa) error {
	err := m.db.QueryRow("SELECT * FROM mahasiswa WHERE id = ?", id).Scan(
		&mahasiswa.Id,
		&mahasiswa.NamaLengkap,
		&mahasiswa.JenisKelamin,
		&mahasiswa.TempatLahir,
		&mahasiswa.TanggalLahir,
		&mahasiswa.Alamat,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("mahasiswa dengan ID %d tidak ditemukan", id)
		}
		return err
	}

	return nil
}
