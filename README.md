# XM



## Build
```azure
go build
```

## Run the app
```azure
./xm 
```

## Run the tests
```azure
go test -v ./...
```

# REST API

The REST APIs to the xm company entity is described below.

## Create a new company

### Request
- Request origin should be Cyprus otherwise it will return 401 error with `Key_InvalidRequestOrigin` as response
```azure
    HTTP Method: POST
    Request URL: http://localhost:8080/api/companies
    Payload:
    {
        "name":"abc",
        "code": "123",
        "country": "india",
        "website": "https://www.abc.com/",
        "phone": "900000000"
    }
```

### Response

    HTTP/1.1 201 Created

    {
        "id": "21af21ba-dc2e-4994-aabc-e4d497a479b2",
        "name": "abc",
        "code": "123",
        "country": "india",
        "website": "https://www.abc.com/",
        "phone": "900000000"
    }

## Update company

### Request

```azure
    HTTP Method: PUT
    Request URL: http://localhost:8080/api/companies/21af21ba-dc2e-4994-aabc-e4d497a479b2
    Payload:
        {
            "name":"abc",
            "code": "001",
            "country": "india",
            "website": "https://www.abc.com/",
            "phone": "900000000"
        }
```

### Response

    HTTP/1.1 200 OK

    {
        "id": "21af21ba-dc2e-4994-aabc-e4d497a479b2",
        "name": "abc",
        "code": "001",
        "country": "india",
        "website": "https://www.abc.com/",
        "phone": "900000000"
    }

## Get list of companies
### Request
```azure
    HTTP Method: GET
    Request URL: http://localhost:8080/api/companies
    List of query parameters: [name, code, country, website, phone]
```

### Response

    HTTP/1.1 200 OK
    [
        {
            "id": "21af21ba-dc2e-4994-aabc-e4d497a479b2",
            "name": "abc",
            "code": "123",
            "country": "india",
            "website": "https://www.abc.com/",
            "phone": "900000000"
        }
    ]

## Get a specific company
### Request
```azure
    HTTP Method: GET
    Request URL: http://localhost:8080/api/companies/21af21ba-dc2e-4994-aabc-e4d497a479b2
```

### Response

    HTTP/1.1 200 OK
    {
        "id": "21af21ba-dc2e-4994-aabc-e4d497a479b2",
        "name": "abc",
        "code": "123",
        "country": "india",
        "website": "https://www.abc.com/",
        "phone": "900000000"
    }

## Delete company
### Request
- Request origin should be Cyprus otherwise it will return 401 error with `Key_InvalidRequestOrigin` as response
```azure
    HTTP Method: DELETE
    Request URL: http://localhost:8080/api/companies/21af21ba-dc2e-4994-aabc-e4d497a479b2
```

### Response

    HTTP/1.1 200 OK

