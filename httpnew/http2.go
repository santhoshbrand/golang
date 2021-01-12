package main

import (
	"fmt"

	"encoding/json"

	"net/http"

	"time"

	"regexp"

	"strconv"
)

type Asset struct {
	ID int `json:"id"  xml:"ID"`

	Brand string `json:"brand"  xml:"Model"`

	Value int `json:"value" xml:"Price"`

	BOMDate string `json:"-" xml:"BOM"`
}

type AssetList []Asset

var assets AssetList

func (al *AssetList) AddAsset(rw http.ResponseWriter, r *http.Request) {

	//fmt.Println(r.Body)

	dc := json.NewDecoder(r.Body)

	var asset Asset

	err := dc.Decode(&asset)

	fmt.Println(asset)

	if err != nil {

		fmt.Println("Asset detail not added due to ", err)

		return

	}

	assets = append(assets, asset)

	fmt.Println("Asset detail added")

}

func (a *AssetList) ModifyAsset(rw http.ResponseWriter, r *http.Request) {

	// get the key that needs modification using regexp

	re, _ := regexp.Compile("/([0-9]+)")

	id := re.FindAllStringSubmatch(r.URL.RequestURI(), -1)

	var ID int

	ID, _ = strconv.Atoi(string(id[0][1]))

	fmt.Println("ID", ID)

	// get the modified value for the key from the body of the request

	dc := json.NewDecoder(r.Body)

	var asset Asset

	err := dc.Decode(&asset)

	fmt.Println(asset)

	if err != nil {

		fmt.Println("Asset detail not added due to ", err)

		return

	}

	// iterate the assets slice for the key recieved and once identified modify

	// the values based on the values in the body of the request

	// if given key not found in the entire assets then report error

	for i, _ := range assets {

		if assets[i].ID == asset.ID {

			assets[i].Brand = asset.Brand

			assets[i].Value = asset.Value

		}

	}

}

func (a *AssetList) GetAsset(rw http.ResponseWriter, r *http.Request) {

	ec := json.NewEncoder(rw)

	re, _ := regexp.Compile("/([0-9]+)")

	id := re.FindAllStringSubmatch(r.URL.RequestURI(), -1)

	if id == nil {

		fmt.Println("Encoding List")

		ec.Encode(&assets)

		//rw.Write([]byte(assets))

		return

	}

	var ID int

	ID, _ = strconv.Atoi(string(id[0][1]))

	fmt.Println("Encoding details of ID", ID)

	for i, _ := range assets {

		if assets[i].ID == ID {

			ec.Encode(&assets[i])

		}

	}

}

func (a *AssetList) DeleteAsset(rw http.ResponseWriter, r *http.Request) {

	fmt.Println("Delte asset to the Asset slice")

	fmt.Println(r.Body)

	/* i ,_ := range assets {

	   if assets[i].ID =  a.ID {



	    assets = append(assets[:i],assets[i+1:len(i)])

	   }



	} */

}

type AssetHandler struct {
	number int
}

func (a *AssetHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	//rw.Write([]byte("This is Asset Hanlder"))

	if r.Method == http.MethodPost {

		assets.AddAsset(rw, r)

		return

	}

	if r.Method == http.MethodPut {

		assets.ModifyAsset(rw, r)

		return

	}

	if r.Method == http.MethodDelete {

		assets.DeleteAsset(rw, r)

		return

	}

	if r.Method == http.MethodGet {

		//rw.Write([]byte(assets))

		assets.GetAsset(rw, r)

		return

	}

}

func main() {

	var ah AssetHandler

	sm := http.NewServeMux()

	sm.Handle("/", &ah)

	fmt.Println("Starting HTTP server...")

	svr := &http.Server{

		Addr: ":9090",

		Handler: sm,

		IdleTimeout: 120 * time.Second,

		ReadTimeout: 1 * time.Second,

		WriteTimeout: 1 * time.Second,
	}

	func() {

		err := svr.ListenAndServe()

		if err != nil {

			fmt.Println("Error in Listening")

		}

	}()

}
