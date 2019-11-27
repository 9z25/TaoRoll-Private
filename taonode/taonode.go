// taonode.go

package taonode

import (
    "fmt"
    "net/http"
    "io/ioutil"

)

const url = "http://34.223.110.241:8000/api/getaddress/"
const url2 = "https://taoexplorer.com/ext/getaddress/TfDJV4odVTsR8u7maQWg7yBTE4aghxwd4h"
const url3 = "https://taoexplorer.com/ext/getaddress/"
const url4 = "http://34.223.110.241:8000/api/getnewaddress/"
const addy = "TfDJV4odVTsR8u7maQWg7yBTE4aghxwd4h"

var client = & http.Client {}

// GetAddress : for Testing purposes ex. addr, _ := taonode.GetAddress()

func GetAddress() string {

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("X-Csrf-Token", "123")
    res, _ := client.Do(req)

    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
    }


    return string(body)
}


// GetNewAddress create new tao address ex. addr, _ := taonode.GetNewAddress()

func GetNewAddress() string {

    req, _ := http.NewRequest("GET", url4, nil)
    req.Header.Set("X-Csrf-Token", "123")
    res, _ := client.Do(req)

    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
    }


    return string(body)

}

// Balance get wallet information from address ex. json, _ := taonode.Balance(str).

func Balance(addr string) string {

    req1, _ := http.NewRequest("GET", url3 + addy, nil)

    //req1.Header.Set("X-Csrf-Token", "123")

    res1, _ := client.Do(req1)


    defer res1.Body.Close()
    body1, err := ioutil.ReadAll(res1.Body)

    if err != nil {
        fmt.Println(err)
    }

    return string(body1)
}




