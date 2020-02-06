package supply_list

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"school-list-project/database"
	"strconv"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/school/{schoolID}", GetSupplyLists)
	router.Get("/school/{schoolID}/grade/{grade}", GetSupplyList)

	return router
}


//A GradeList is a list of all supply lists
//Belonging to a school
type GradeList struct {
	SchoolName     string              `json:"school_name"`
	EducationStage string              `json:"education_stage"`
	SupplyLists    []SupplyListDetails `json:"supply_lists"`
}

type SupplyListDetails struct {
	ListID       int `json:"list_id"`
	Grade        int `json:"grade"`
	StartingYear int `json:"starting_year"`
	EndingYear   int `json:"ending_year"`
}

func GetSupplyLists(w http.ResponseWriter, r *http.Request) {
	schoolID := chi.URLParam(r, "schoolID")

	if IsNumeric(schoolID) && len(schoolID) < 5 {
		rows, err := database.DBCon.Query("SELECT list_id, grade, starting_year, ending_year FROM supply_list P INNER JOIN school S ON S.school_id = P.school_id WHERE P.school_id=$1 ORDER BY grade", schoolID)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var supplyLists []SupplyListDetails
		for rows.Next() {
			var sl SupplyListDetails
			err := rows.Scan(&sl.ListID, &sl.Grade, &sl.StartingYear, &sl.EndingYear)
			if err != nil {
				log.Fatal(err)
			}

			supplyLists = append(supplyLists, sl)
		}
		schoolDetails := GetSchoolName(schoolID)
		gradeList := GradeList{
			SchoolName:     schoolDetails[0],
			EducationStage: schoolDetails[1],
			SupplyLists:    supplyLists,
		}
		render.JSON(w, r, gradeList)
	} else {
		RespondWithError(w, r, 414, "School ID's must be numerical and smaller than 4 digits")
	}
}
func GetSchoolName(schoolID string) []string {
	rows, err := database.DBCon.Query("SELECT name, education_stage FROM school S WHERE S.school_id=$1", schoolID)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var schoolName string
	var educationStage string

	for rows.Next() {
		err := rows.Scan(&schoolName, &educationStage)
		if err != nil {
			log.Fatal(err)
		}
	}
	schoolDetails := make([]string, 2)
	schoolDetails[0] = schoolName
	schoolDetails[1] = educationStage

	return schoolDetails
}

type SupplyListItem struct {
	SupplyID string `json:"supply_id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Optional bool   `json:"optional"`
	Amount int   `json:"amount"`
}

type SupplyList struct {
	SupplyListItems []SupplyListItem `json:"supply_list"`
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func GetSupplyList(w http.ResponseWriter, r *http.Request) {
	grade := chi.URLParam(r, "grade")
	schoolID := chi.URLParam(r, "schoolID")


	if IsNumeric(grade) && len(grade) < 2 && IsNumeric(schoolID) && len(schoolID) < 2{
		rows, err := database.DBCon.Query(`SELECT P.supply_id, P.supply_name, P.supply_desc, B.optional, B.amount
													FROM supply_item P
													JOIN supply_list_bridge B ON P.supply_id = B.supply_id
													JOIN supply_list S ON S.list_id = B.list_id
													WHERE S.grade=$1 AND
													S.school_id=$2`, grade, schoolID)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var supplyListItems []SupplyListItem
		for rows.Next() {
			var sli SupplyListItem
			err := rows.Scan(&sli.SupplyID, &sli.Name, &sli.Desc, &sli.Optional, &sli.Amount)

			if err != nil {
				log.Fatal(err)
			}

			supplyListItems = append(supplyListItems, sli)
		}
		supplyList := SupplyList{
			SupplyListItems: supplyListItems,
		}

		render.JSON(w, r, supplyList)
	} else {
		RespondWithError(w, r, 414, "Grades must be numerical and smaller than 4 digits")
	}
}

type Error struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_msg"`
}

func RespondWithError(w http.ResponseWriter, r *http.Request, status_code int, error_msg string) {
	render.Status(r, status_code)
	render.JSON(w, r, Error{
		StatusCode:   status_code,
		ErrorMessage: error_msg,
	})
}
