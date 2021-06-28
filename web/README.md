# Simple url-shortener service

## About
This is a course project for golang backend lvl 1. Realization contains most part of basic features, such as
* web site;
* database connection (internal);
* api endpoints.
  
Since database connection described as an interface, potential users of the project can change the database they would like to use. Currently project uses **Redis** to store data.

## How to build and run
Use *Makefile* to pass lint checks and build docker image.

Run docker compose file. Since the database is redis it will pull its image and start services on **localhost:8080**.

## Web

![main-page](./docs/images/main.png)

Main page represents almost all functionality of the service. It awaits for the user to enter full link, select expiration date (can be omitted, default link timeout is 24 hours) and press **Generate** button.

![generated-link](./docs/images/generated-link.png)

After generation is complited, you will see a page with link information. It includes:

* Full link;
* Short link;
* Expiration time.

Now this short link will lead you to the web-page it was created for.
Confirmation can be found in docker container logs:

![redirect](./docs/images/redirect.png)

Or you can use statistics for specific link:

![stat](./docs/images/stat.png)

## API

API was tested with *Postman*, you can copy all further examples and test them yourself.

* **[host]/api/generate**

    This endpoint creates short link with given expiration time.

    Request:
    ```json
    {
        "link": "https://en.wikipedia.org/wiki/URL",
        "expiration_time": "2021-11-18"
    }
    ```

    Expiration can be set with a template: **YYYY-MM-DD**. In case if expiration is ommited or set as an empty string it will be default 24 hours.

    Response for the request is gonna be with the status code *201*:

    ```json
    {
        "full_url": "https://en.wikipedia.org/wiki/URL",
        "short_url": "localhost:8080/u0Sif4g7R",
        "expiration_time": 13566136803266447,
        "redirects": 0
    }
    ```

    The response represents link structure.

    ### Error status codes:

    * 400 - returns if given json object cannot be parsed. If link field is empty. Or if duration is in the past.

    Request:
    ```json
    {
        "link": "",
        "expiration_time": "2021-11-18"
    }
    ```
    Response:
    ```json
    {
        "error": "url is empty"
    }
    ```

    Request:
    ```json
    {
        "link": "https://en.wikipedia.org/wiki/URL",
        "expiration_time": "2020-11-18"
    }
    ```
    Response:
    ```json
    {
        "error": "expiration time is in the past"
    }
    ```

    * 500 - returns if there is a problem with data processing or database interaction. 

    Error returns a json-object with error message, which can describe a problem more clearly:

    Request:
    ```json
    {
        "link": "https://en.wikipedia.org/wiki/URL",
        "expiration_time": "20-11-18"
    }
    ```
    Response:
    ```json
    {
    "error": "parsing time \"20-11-18\" as \"2006-01-02\": cannot parse \"1-18\" as \"2006\""
    }
    ```

* **[host]/api/statistics**

    This endpoint retrieves short link statistics.

    Request:
    ```json
    {
        "link": "localhost:8080/u0Sif4g7R"
    }
    ```

    Response for the request is gonna be with the status code *200*:

    ```json
    {
        "full_url": "https://en.wikipedia.org/wiki/URL",
        "short_url": "localhost:8080/u0Sif4g7R",
        "expiration_time": 13565446000000000,
        "redirects": 0
    }
    ```

    ### Error status codes:

    * 400 - returns if given json object cannot be parsed. Or link field is empty.

    Request:
    ```json
    {
        "link": ""
    }
    ```
    Response:
    ```json
    {
        "error": "url is empty"
    }
    ```

    * 500 - returns if there is a problem with getting data from a database.

