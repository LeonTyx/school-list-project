package supplies

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

	router.Get("/supplies/{districtID}", GetSupplies)

	return router
}

type Supplies struct {
	DistrictID   string   `json:"district_id"`
	DistrictName string   `json:"district_name"`
	Supplies     []Supply `json:"supply_lists"`
}

type Supply struct {
	ID   string `json:"supply_id"`
	Name string `json:"supply_name"`
	Desc string `json:"supply_desc"`
}

func GetSupplies(w http.ResponseWriter, r *http.Request) {
	districtID := chi.URLParam(r, "districtID")

	if IsNumeric(districtID) && len(districtID) < 5 {
		rows, err := database.DBCon.Query("SELECT supply_id, supply_name, supply_name FROM supply_item WHERE district_id=$1", districtID)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var supplyListItems []Supply
		for rows.Next() {
			var supply Supply
			err := rows.Scan(&supply.ID, &supply.Name, &supply.Desc)

			if err != nil {
				log.Fatal(err)
			}

			supplyListItems = append(supplyListItems, supply)
		}
		supplyList := Supplies{
			DistrictID:   districtID,
			DistrictName: "",
			Supplies:     nil,
		}

		render.JSON(w, r, supplyList)
	} else {
		render.Status(r, 414)
		render.JSON(w, r, nil)
	}
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
