---
title: API Reference

language_tabs: # must be one of https://git.io/vQNgJ
  - bash

toc_footers:
  - <a href='#'>Sign Up for a Developer Key</a>
  - <a href='https://github.com/slatedocs/slate'>Documentation Powered by Slate</a>

search: true

code_clipboard: true

meta:
  - name: description
    content: Documentation for the REST API
---


# Introduction

Welcome to the Kittn API! You can use our API to access Kittn API endpoints, which can get information on various cats, kittns, and breeds in our database.

We have language bindings in Shell. You can view code examples in the dark area to the right.

This example API documentation page was created with [Slate](https://github.com/slatedocs/slate). Feel free to edit it and use it as a base for your own API's documentation.


# Authentication

> To authorize, use this code:

```shell
# With bash, you can just pass the correct header with each request
curl "api_endpoint_here" \
  -H "Authorization: meowmeowmeow"
```

> Make sure to replace `meowmeowmeow` with your API key.

Kittn uses API keys to allow access to the API. You can register a new Kittn API key at our [developer portal](http://example.com/developers).

Kittn expects for the API key to be included in all API requests to the server in a header that looks like the following:

`Authorization: meowmeowmeow`

<aside class="notice">
You must replace <code>meowmeowmeow</code> with your personal API key.
</aside>



# Domain
## Root Path

```shell
curl "http://example.com/" \
  -H "Authorization: Bearer {token}" 
```


> The above command returns JSON structured like this:

```json
{
 "app_name": "rest-api",
 "build_info": {
  "build": "go1.18.2",
  "version": "fbb340b",
  "build_time": "2022-05-30T13:35:28"
 }
}
```
The root path returns the api name and version.

### HTTP Request

`GET http://example.com/`



# Kittns V1
## Delete a Specific Kitten

```shell
curl -X DELETE "http://example.com/v1/kittns/{id}" \
  -H "Authorization: Bearer {token}" 
```

This endpoint deletes a specific kittn

### HTTP Request

`DELETE http://example.com/v1/kittns/{id}`

### URL Parameters

| Parameter    | Description                                                                      |
| ------------ | -------------------------------------------------------------------------------- |
| id |  the id for a kittn |





## Add a New Kittn

```shell
curl -X POST "http://example.com/v1/kittns" \
  -H "Authorization: Bearer {token}" \
  --header 'Content-Type: application/json' --raw-data '{
 "name": "Stealth",
 "breed": "Siamese",
 "fluffiness": 2,
 "cuteness": 3
}'
```


> The body of the request should be JSON and structured like this:

```json
{
 "name": "Stealth",
 "breed": "Siamese",
 "fluffiness": 2,
 "cuteness": 3
}
```

> The above command returns JSON structured like this:

```json
{
 "id": 3,
 "name": "Stealth",
 "breed": "Siamese",
 "fluffiness": 2,
 "cuteness": 3
}
```
This endpoint deletes a specific kittn

### HTTP Request

`POST http://example.com/v1/kittns`

## Get All Kittens

```shell
curl "http://example.com/v1/kittns" \
  -H "Authorization: Bearer {token}" 
```


> The above command returns JSON structured like this:

```json
[
 {
  "id": 1,
  "name": "Fluffums",
  "breed": "calico",
  "fluffiness": 6,
  "cuteness": 7
 },
 {
  "id": 2,
  "name": "Max",
  "breed": "calico",
  "fluffiness": 5,
  "cuteness": 10
 }
]
```
This endpoint retrieves all kittns

### HTTP Request

`GET http://example.com/v1/kittns`

## Get a Specific Kitten

```shell
curl "http://example.com/v1/kittns/{id}" \
  -H "Authorization: Bearer {token}" 
```


> The above command returns JSON structured like this:

```json
{
 "id": 1,
 "name": "Fluffums",
 "breed": "calico",
 "fluffiness": 6,
 "cuteness": 7
}
```
This endpoint retrieves a specific kittn

### HTTP Request

`GET http://example.com/v1/kittns/{id}`

### URL Parameters

| Parameter    | Description                                                                      |
| ------------ | -------------------------------------------------------------------------------- |
| id |  the id for a kittn |








# Errors

This API uses the following error codes:

Error Code | Meaning
---------- | -------
400 | Bad Request -- Your request is invalid.
401 | Unauthorized -- Your API key is wrong.
403 | Forbidden -- Unauthorized
404 | Not Found -- The specified kittn could not be found.
405 | Method Not Allowed -- You tried to access an endpoint with an invalid method.
406 | Not Acceptable -- You requested a format that isn't json.
410 | Gone -- The resource requested has been removed from our servers.
418 | I'm a teapot.
429 | Too Many Requests -- Rate limits applied ! Slow down!
500 | Internal Server Error -- We had a problem with our server. Try again later.
503 | Service Unavailable -- We're temporarily offline for maintenance. Please try again later.











