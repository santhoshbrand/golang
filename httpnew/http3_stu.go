package main

import (
        "fmt"
        "encoding/json"
        "net/http"
        "time"
        "regexp"
)

type Student struct {
        Name string `json:"name"  xml:"NAME"`
        Rollno int64 `json:"rollno"  xml:"ROLLNO"`
        Grade string `json:"grade" xml:"GRADE"`
}

type Studentinfo []Student
var data Studentinfo

func (stu *Studentinfo) Addinfo(rw http.ResponseWriter, r *http.Request) {

        //fmt.Println(r.Body)
        decoder := json.NewDecoder(r.Body)
        var st Student
        err := decoder.Decode(&st)
        fmt.Println(st)
        if err != nil {
                fmt.Println("stdent's information is  not included due to ", err)
                return
        }
        data = append(data, st)
        fmt.Println("Student information is included")

}

func (stu *Studentinfo) Modifyinfo(rw http.ResponseWriter, r *http.Request) {

        // get the key that needs modification using regexp
        regular, _ := regexp.Compile("/([A-Za-z]+)")
	ID := regular.FindAllStringSubmatch(r.URL.RequestURI(), -1)
        // get the modified value for the key from the body of the request
        dc := json.NewDecoder(r.Body)
        var st Student
        err := dc.Decode(&st)
        fmt.Println(st)
        if err != nil {
                fmt.Println("Asset detail not added due to ", err)
                return
        }

        // iterate the assets slice for the key recieved and once identified modify
        // the values based on the values in the body of the request
        // if given key not found in the entire assets then report error
	 for i, _ := range data {

                if data[i].Name == ID[0][1] {
	              data[i].Rollno = st.Rollno
                      data[i].Grade = st.Grade
                }
	}
}
func (stu *Studentinfo) Getinfo(rw http.ResponseWriter, r *http.Request) {
	var st Student
        encoder := json.NewEncoder(rw)
        regular, _ := regexp.Compile("/([A-Za-z]+)")
        ID := regular.FindAllStringSubmatch(r.URL.RequestURI(), -1)
        if ID == nil {
                fmt.Println("Encoding List")
                encoder.Encode(&st)
                //rw.Write([]byte(assets))
                return
        }
        for i, _ := range data {
                if data[i].Name == ID[0][0] {
                        encoder.Encode(&data[i])

                }
	}

}
 
/*func (stu *Studentinfo) DeleteAsset(rw http.ResponseWriter, r *http.Request) {

        fmt.Println("Delte student info to the Asset slice")

        fmt.Println(r.Body)

        /* i ,_ := range assets {

           if assets[i].ID =  a.ID {



            assets = append(assets[:i],assets[i+1:len(i)])

           }



        } */

type StudentHandler struct {
        name string
}

func (stu *Studentinfo) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

        //rw.Write([]byte("This is Asset Hanlder"))

        if r.Method == http.MethodPost {

                data.Addinfo(rw, r)

                return

        }

        if r.Method == http.MethodPut {

                data.Modifyinfo(rw, r)

                return

        }


  /*      if r.Method == http.MethodDelete {

		   assets.DeleteAsset(rw, r)

                return

        }*/

        if r.Method == http.MethodGet {

                //rw.Write([]byte(assets))

                data.Getinfo(rw, r)

                return

        }

}

func main() {

        var stud Studentinfo

        httphandler := http.NewServeMux()

        httphandler.Handle("/", &stud)

        fmt.Println("Starting HTTP server...")
	svr := &http.Server{

                Addr: ":3333",

                Handler: httphandler,

                IdleTimeout: 60 * time.Second,

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

                                                                                                                             






