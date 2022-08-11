## How to customize API?
We can use *Json* file to customize `http request` for API.  
We support some field according `APIRequest` struct in `request.go`  
- `baseurl`
- `path`
- `method`
- `headers`
- `postbody`
- `variables`
- `needre`
- `response_type`


There are some field in json file you should be careful:
- `variables` Define variables in this field such as some API key or secret and use them like `{{variable}}`.By the way,the `{{domain}}` variable must be exsitent as it is your target and placing it in the correct location according to the API request.
- `needre` We want to result like this `www.example.com`. But some API's response is like this `result is www.example.com xxxxx` even if it's json format. So we need regular expression after getting API result.  
- `response_type` We need use diiferent ways to process response body according to this field.

## Function
We provide some practical functions.You can check `funcList` in `run.go`  
- `base64`

## response_type details
three type : `raw` | `json` | `special`
### Raw Response
```json
{"response_type":"raw"}
```
It's pretty easy when it comes to `raw response`.You just need build `http request` in a *json file*. (remember to set `response_type` to `raw`)

### Json Response
```json
{"response_type":"json"}
```
If the response is json format,you have to do following steps:  
1. Create a `Json` file in `scripts` folder which is for build `request` and build structs in file `apistruct.go` according to API response.  
    *Attention:* The file name should be *the same* as the API struct  

2. Append `struct` to map `SpecialRespMap` in file `apistruct.go`.The key is file name,and the value is the `struct` defined for API.  

We don't support some json format. For instance:`a big list` because it can't be parsed into a struct normally. So we have the third `response_type` SpecialResp

### Special Response
```json
{"response_type":"special"}
```

`special response` indicate that you can't process the response content through the two ways above. And you need a function to process it.  
1. build `request` in `scripts` folder
2. Create a file with `go` suffix. You should write a struct for API and write a function named `SpecailProcess` which for implementing `SpecialResp` interface.The function `SpecailProcess` just need a `response data`.Just process response data in this function and return the subdomain slice you need in API.
