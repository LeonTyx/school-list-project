package school_list

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"school-list-project/database"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()

	//router.Delete("/{todoID}", DeleteSchool)
	//router.Post("/", CreateSchool)
	router.Get("/", GetSchools)
	return router
}

type School struct {
	Name           string `json:"name"`
	EducationStage string `json:"education_stage"`
	SchoolID       int    `json:"school_id"`
}

type SchoolList struct {
	Schools      []School `json:"schools"`
	DistrictName string   `json:"district"`
}

func GetSchools(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DBCon.Query("select S.name, S.education_stage, d.name, S.school_id from school S INNER JOIN district d on S.district_id = d.nces_id")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var schools []School
	var DistrictName string
	for rows.Next() {
		var s School
		err := rows.Scan(&s.Name, &s.EducationStage, &DistrictName, &s.SchoolID)
		if err != nil {
			log.Fatal(err)
		}

		schools = append(schools, s)
	}

	SchoolList := SchoolList{
		Schools:      schools,
		DistrictName: DistrictName,
	}

	render.JSON(w, r, SchoolList)
}
