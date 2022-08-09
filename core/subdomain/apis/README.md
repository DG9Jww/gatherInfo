## How to customize API?

###Step1. Add json file
There are some fields you should know:
- `needre` which means regular expression.For example,If your `subdomain` result still need `re`:
```json
    "needre":{
        "ip":false,
        "subdomain":true
    }
```
And we only support `ip` and `subdomain`

- `response_type` Only support `json` and `raw`.  
- `variables` Some API need key and secret and you should use symbol `{{variable}}` to set them.By the way, variable `{{domain}}` is necessary.Just metion the location it should be.



###Step2. Add struct
