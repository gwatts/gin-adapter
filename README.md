# Gin Middleware Adapter

[![GoDoc](https://godoc.org/github.com/gwatts/gin-adapter?status.svg)](https://godoc.org/github.com/gwatts/gin-adapter)
![Build Status](https://github.com/gwatts/gin-adapter/actions/workflows/go.yml/badge.svg)


This package provides a simple adapter to make it easier to 
use common Go HTTP middleware packages with the
[Gin Web Framework](https://github.com/gin-gonic/gin/).

For example, if you have some middleware that uses the common pattern of
having a function that takes an `http.Handler` to wrap and returning another
`http.Handler`, you can use `adapter.Wrap`.  

Using the [nosurf CSRF middleware](https://github.com/justinas/nosurf/):

```golang
engine.Use(adapter.Wrap(nosurf.NewPure)) 
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
