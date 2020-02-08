package supplies

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"school-list-project/auth/authorization"
	"school-list-project/database"
	"strconv"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		suppliesCtx,
		authorization.ValidSession,
		authorization.ResourceCtx,
	)

	router.With(
		authorization.CanView,
	).Get("/{districtID}", GetSupplies)

	router.With(
		authorization.CanView,
	).Get("/supply/{supplyID}", GetSupply)

	router.With(
		authorization.CanDelete,
	).Delete("/{todoID}", DeleteSupply)

	router.With(
		authorization.CanCreate,
		isUnique,
	).Post("/create", CreateSupply)
	//router.Post("/", EditSupply)

	return router
}

func suppliesCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "resource", "supplies")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetSupply(w http.ResponseWriter, r *http.Request) {

}

//func CreateSupply(w http.ResponseWriter, r *http.Request) {}

func DeleteSupply(w http.ResponseWriter, r *http.Request) {

}

func isUnique(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var supply Supply
		_ = json.NewDecoder(r.Body).Decode(&supply)

		rows, err := database.DBCon.Query(`SELECT COUNT(*) as count FROM supply_item WHERE supply_name = $1`, supply.Name)
		if err != nil {
			RespondWithError(w, r, 503, "There was an error contacting the database.")
			log.Fatal(err)
		}
		defer rows.Close()
		var count int

		for rows.Next() {
			err := rows.Scan(&count)
			if err != nil {
				RespondWithError(w, r, 503, "There was an error contacting the database.")
				panic(err)
			}
		}

		if count > 0 {
			RespondWithError(w, r, 400, "Supply item names must be unique")
			return
		}

		next.ServeHTTP(w, r)
	})
}

//Create a supply using the Supply struct and post body information
//TODO Add ability to add item to specific districts
func CreateSupply(w http.ResponseWriter, r *http.Request) {
	var supply Supply

	_ = json.NewDecoder(r.Body).Decode(&supply)

	rows, err := database.DBCon.Query(`INSERT INTO supply_item (supply_name, supply_desc, district_id) 
											  VALUES($1, $2, 5305400)`, supply.Name, supply.Desc)
	if err != nil {
		RespondWithError(w, r, 503, "There was an error contacting the database.")

		log.Fatal(err)
	}
	defer rows.Close()

	render.Status(r, 201)
	render.JSON(w, r, `"message":"successfully created"`)
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
