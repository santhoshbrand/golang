package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Employee struct {
    Sapid   int    `json:"Sapid"  xml:"ID"`
    Empname string `json:"Empname"  xml:"Model"`
    Salary  int    `json:"value" xml:"Price"`
    
}
type EmployeeList []Employee

var employees EmployeeList

func (el *EmployeeList) AddEmployee(rw http.ResponseWriter, r *http.Request) {

    dc := json.NewDecoder(r.Body)
    var employee Employee
    err := dc.Decode(&employee)
    fmt.Println(employee)
    if err != nil {
        fmt.Println("Employee detail not added due to ", err)
        return
    }
    employees = append(employees, employee)
    fmt.Println("Employee detail added")
}

type AssetHandler struct {
    number int
}

func (a *AssetHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
    //rw.Write([]byte("This is Asset Hanlder"))
    if r.Method == http.MethodPost {
        employees.AddEmployee(rw, r)
        return
    }

}

func main() {
    var ah AssetHandler
    sm := http.NewServeMux()
    sm.Handle("/", &ah)

    fmt.Println("Starting HTTP server...")
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
            fmt.Println("Error in Listening")
        }
    }()
}
