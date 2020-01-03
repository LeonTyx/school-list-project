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
	router.Get("/{listID}", GetASupplyList)

	return router
}

type GradeList struct {
	GradeList []SupplyListDetails `json:"grade_list"`
	School    string              `json:"school"`
}

type SupplyListDetails struct {
	ListID       int `json:"list_id"`
	Grade        int `json:"grade"`
	StartingYear int `json:"starting_year"`
	EndingYear   int `json:"ending_year"`
}

func GetSupplyLists(w http.ResponseWriter, r *http.Request) {
	schoolID := chi.URLParam(r, "schoolID")

	var supplyLists []SupplyListDetails
	supplyList := SupplyListDetails{
		ListID:       0,
		Grade:        0,
		StartingYear: 0,
		EndingYear:   0,
	}
	supplyLists = append(supplyLists, supplyList)
	gradeList := GradeList{
		GradeList: supplyLists,
		School:    schoolID,
	}

	render.JSON(w, r, gradeList) // A chi router helper for serializing and returning json
}

type SupplyListItem struct {
	SupplyID string `json:"list_id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Optional bool   `json:"optional"`
}

type SupplyList struct {
	SupplyListItems []SupplyListItem `json:"supply_list"`
	Grade           int8             `json:"grade"`
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func GetASupplyList(w http.ResponseWriter, r *http.Request) {
	listID := chi.URLParam(r, "listID")

	var grade int8
	grade = -1

	if IsNumeric(listID) && len(listID) < 5 {
		rows, err := database.DBCon.Query("SELECT S.grade, P.supply_id, P.supply_name, P.supply_desc, B.optional FROM supply_item P JOIN supply_list_bridge B ON P.supply_id = B.supply_id JOIN supply_list S ON S.list_id = B.list_id WHERE B.list_id=$1 ORDER BY grade ASC", listID)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var supplyListItems []SupplyListItem
		for rows.Next() {
			var sli SupplyListItem
			err := rows.Scan(&grade, &sli.SupplyID, &sli.Name, &sli.Desc, &sli.Optional)

			if err != nil {
				log.Fatal(err)
			}

			supplyListItems = append(supplyListItems, sli)
		}
		supplyList := SupplyList{
			SupplyListItems: supplyListItems,
			Grade:           grade,
		}

		render.JSON(w, r, supplyList)
	} else {
		render.Status(r, 414)
		render.JSON(w, r, nil)
	}

}
