middlewares
===========

`middlewares` is a collection of Gorilla Mux middlewares. 

To use:

```
go get -u gitlab.com/amihan/common/libraries/go/middlewares.git
```

```
import (
    "github.com/gorilla/mux"

    "gitlab.com/amihan/common/libraries/go/middlewares.git"
)

func example() {
    router := mux.NewRouter()

    // this protects all routes with API Key
    apiKeyMiddleware := middlewares.APIKey(middlewares.APIKeyTypeHeader, "TEST-API-KEY", "secret")
    router.Use(apiKeyMiddleware)

    // ... do stuff with router
}
```
