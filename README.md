# Gin Middleware Adapter

This package provides a simple adapter to make it easier to 
use common Go HTTP middleware packages with the
[Gin Web Framework](https://github.com/gin-gonic/gin/).

For example, if you have some middleware that uses the common pattern of
having a function that takes an `http.Handler` to wrap and returning another
`http.Handler`, you can use `adapter.Wrap`

```golang
engine.Use(adapter.Wrap(nosurf.New)) 
```

If you need to pass in the wrapped handler explictly (eg. via configuration),
or use a different signature, you can use `adapter.New`, which returns a
handler that can be passed to the middleware, and a function that wraps
the handler the middleware returns into a gin.HandlerFunc:

```golang
nextHandler, wrapper := adapter.New()
ns := nosurf.New(nextHandler)
ns.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "Custom error about invalid CSRF tokens", http.StatusBadRequest)
}))
engine.Use(wrapper(ns)) 
```

The wrapper will ensure that the request's context is correctly passed down
the stack and that the next middleware in the chain is called.

It will additionally call the Abort method on Gin's context if the middleware
does not call the wrapped handler.
