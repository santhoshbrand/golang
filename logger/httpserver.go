package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
	"regexp"
	"strconv"
 	"myutils" 
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)



type Asset struct {
	ID      int    `json:"id"  xml:"ID"`
	Brand   string `json:"brand"  xml:"Model"`
	Value   int    `json:"value" xml:"Price"`
	BOMDate string `json:"-" xml:"BOM"`
}

type AssetList []Asset

var assets AssetList

func (al *AssetList) AddAsset(rw http.ResponseWriter, r *http.Request) {
	ml := myutils.NewMyLogger("HTTP")
	ml.Println(myutils.Debug,"Add a new element in Asset slice")
	//ml.Println(Debug, r.Body)
	dc := json.NewDecoder(r.Body)
	var asset Asset
	err := dc.Decode(&asset)
	//ml.Println(Debug, asset)
	if err != nil {
		ml.Println(myutils.Error, "Asset detail not added due to "+err.Error())
	}
	assets = append(assets, asset)
	ml.Println(myutils.Debug, "Asset detail added")
}

func (a *AssetList) ModifyAsset(rw http.ResponseWriter, r *http.Request) {
	ml := myutils.NewMyLogger("HTTP")
	ml.Println(myutils.Debug, "Modify asset to the Asset slice")
	re, _ := regexp.Compile("/([0-9]+)")
    id := re.FindAllStringSubmatch(r.URL.RequestURI(), -1)
    var ID int
    ID,_ =  strconv.Atoi(string(id[0][1])) //why two dimension array?
    ml.Println(myutils.Debug,string(id[0][1]))
   
    // get the modified value for the key from the body of the request
    dc := json.NewDecoder(r.Body)
    var asset Asset
    err := dc.Decode(&asset)
    if err != nil {
        ml.Println(myutils.Error,"Asset detail not modified due to "+err.Error()) 
        return   
    }
    
    // iterate the assets slice for the key recieved and once identified modify
    // the values based on the values in the body of the request
    // if given key not found in the entire assets then report error
    
    for i,_ := range assets {
        if assets[i].ID  == ID {
           assets[i].Brand = asset.Brand
           assets[i].Value = asset.Value  
        }

    }
}

func (a *AssetList) DeleteAsset(rw http.ResponseWriter, r *http.Request) {
	ml := myutils.NewMyLogger("HTTP")
	ml.Println(myutils.Debug,"Delte asset to the Asset slice")
	

	/* i ,_ := range assets {
	   if assets[i].ID =  a.ID {

	    assets = append(assets[:i],assets[i+1:len(i)])
	   }

	} */

}


func (a *AssetList) GetAsset(rw http.ResponseWriter,r *http.Request) {
	ml := myutils.NewMyLogger("HTTP")
	ml.Println(myutils.Debug, "Get the asset from Asset slice")

    ec:=json.NewEncoder(rw)
    re, _ := regexp.Compile("/([0-9]+)")
    id := re.FindAllStringSubmatch(r.URL.RequestURI(), -1)
           
    if id == nil {
       ml.Println(myutils.Info,"Asset list sent as response") 
       ec.Encode(&assets)
       return
     }
    
    var ID int
    ID,_ =  strconv.Atoi(string(id[0][1]))
    ml.Println(myutils.Debug,string(ID))   
    for i,_ := range assets {

        if assets[i].ID  == ID {
			ec.Encode(&assets[i])
			ml.Println(myutils.Info,"Asset details sent as response")       
        }  
    }    
} 


type AssetHandler struct {
	number int
}

func (a *AssetHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//rw.Write([]byte("This is Asset Hanlder"))
	ml := myutils.NewMyLogger("HTTP")
	if r.Method == http.MethodPost {
		ml.Println(myutils.Debug, "Received POST request..")
		assets.AddAsset(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		ml.Println(myutils.Debug, "Received PUT request..")
		assets.ModifyAsset(rw, r)
		return
	}
	if r.Method == http.MethodDelete {
		ml.Println(myutils.Debug, "Received DELETE request..")
		assets.DeleteAsset(rw, r)
		return

	}
	if r.Method == http.MethodGet {
		ml.Println(myutils.Debug, "Received GET request..")
		assets.GetAsset(rw,r)
		return
	}
}

func main() {
	ml := myutils.NewMyLogger("HTTP ")
	var ah AssetHandler
	sm := http.NewServeMux()
	sm.Handle("/", &ah)

	go TCPServer()

	fmt.Println("Starting HTTP server @ Localhost:9090")

	svr := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	func() {
		err := svr.ListenAndServe()
		if err != nil {
			ml.Println(myutils.Error, "Error in Listening"+err.Error())

		}
	}()

}

func TCPServer() {
	ml := myutils.NewMyLogger("HTTP")
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		ml.Println(myutils.Error, "Error listening:"+err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	ml.Println(myutils.Info, "Listening on "+CONN_HOST+":"+CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			ml.Println(myutils.Error, "Error accepting: "+err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	var req myutils.ChangeLogLevelRequest
	ml := myutils.NewMyLogger("HTTP")

	dc := gob.NewDecoder(conn)
	err := dc.Decode(&req)

	if err != nil {
		ml.Println(myutils.Error, "Unable to decode request..")
		conn.Close()
		return
	}
	ml.SetLogLevel(req.NewLogLevel)
	// Close the connection when you're done with it.
	conn.Close()
}
