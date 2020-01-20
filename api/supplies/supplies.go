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

	router.Get("/{districtID}", GetSupplies)

	return router
}

type Supplies struct {
	DistrictID   string   `json:"district_id"`
	DistrictName string   `json:"district_name"`
	Supplies     []Supply `json:"supplies"`
}

type Supply struct {
	ID   string `json:"supply_id"`
	Name string `json:"supply_name"`
	Desc string `json:"supply_desc"`
}

func GetSupplies(w http.ResponseWriter, r *http.Request) {
	districtID := chi.URLParam(r, "districtID")

	if IsNumeric(districtID) && len(districtID) == 7 {
		rows, err := database.DBCon.Query("SELECT supply_id, supply_name, supply_desc, name FROM supply_item LEFT OUTER JOIN district d on supply_item.district_id = d.nces_id WHERE d.nces_id=$1", districtID)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var supplyListItems []Supply
		var districtName string
		for rows.Next() {
			var supply Supply //Possible bottleneck. Creates new struct for each row.
			err := rows.Scan(&supply.ID, &supply.Name, &supply.Desc, &districtName)

			if err != nil {
				log.Fatal(err)
			}

			supplyListItems = append(supplyListItems, supply)
		}
		supplyList := Supplies{
			DistrictID:   districtID,
			DistrictName: districtName,
			Supplies:     supplyListItems,
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
