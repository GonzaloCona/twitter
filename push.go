package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "os"
    )

//var apiKey string
//var customSearchEngineKey string
var searchURL string



func main() {
  //https://www.googleapis.com/customsearch/v1?key=AIzaSyBsm97GWF2WLuAGGcyugUMAu2c4Pm27j_I&cx=007380653580648478702:j5liwfex5-m&q=pablo%20v&alt=json&start=1&num=1
  //apiKey = "AIzaSyBsm97GWF2WLuAGGcyugUMAu2c4Pm27j_I"
  //customSearchEngineKey = "007380653580648478702:j5liwfex5-m"
  searchURL= "http://localhost:8080/fimi_v0/webapi/u/SPush"
  //tosearch:=searchURL+"key="+apiKey+"&cx="+customSearchEngineKey+"&q=pepito&alt=json&start=1&num=1"
  tosearch:=searchURL+";id=APA91bH2ZceuKL4QbZX06OpsI9RltBlaShO6mF4pjFUsxMJc63FCacAlEG6myw1MvcA7UlM9gnB9todjzMrmLY8OMtxi8znCNBYTyoQRfF-izl7bb4jXlWQ;cod=1;contenido=asdf"
    response,err := http.Get(tosearch)
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }
        fmt.Printf("%s\n", string(contents))
    }
}

/*
import (
  "fmt"
  "net/http"
  "github.com/alexjlockwood/gcm"
  )


func main() {


}*/
