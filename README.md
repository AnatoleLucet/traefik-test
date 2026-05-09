# Traefik test

The proxy can be configured in a `greenlight.yml` file:

```yaml
server:
  host: 0.0.0.0
  ports:
    http: 8080
    https: 8443
  tls:
    key: key.pem
    cert: cert.pem

rules:
  - if: # when this rule is applied
      host: "example.com"
      path: "/"
      method: "GET"

    then: # what this rule does (only specify one)
      forward: internal.example.com
      redirect: internal.example.com
      respond:
        status: 200
        body: hello
        headers:
          X-MyHeader: my-value

    middleware: # extra features
      cache:
        ttl: 60s
```

<details>
<summary>Example config</summary>

```yaml
server:
  host: 0.0.0.0
  ports:
    http: 8080

rules:
  - if:
      path: "/hello"
    then:
      forward: localhost:8090
    middleware:
      cache:
        ttl: 10s

  # catch all
  - then:
      respond:
        status: 404
        body: not found
```
