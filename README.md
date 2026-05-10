# Traefik test

My little reverse-proxy called Greenlight, based on the assignment from Traefik Labs.

## Features

- Rule-based routing
- Requests forwarding and redirects
- Static responses
- Per-rule caching

## Introduction

Greenlight is an HTTP/HTTPS reverse-proxy that discriminates incoming requests and handles them in different ways.

You define a set of rules to identify incoming requests and specify what to do with each one.

#### Configuration

Greenlight can be configured in a `greenlight.yml` file placed in the current working directory.

Example config:

```yaml
server:
  host: 0.0.0.0
  ports:
    http: 8080

rules:
  - if: # when this rule is applied
      path: "/hello"
    then: # what this rule does
      forward: localhost:8090
    middleware: # extra features
      cache:
        ttl: 10s

  # catch all
  - then:
      respond:
        status: 404
        body: not found
```

<details>
<summary>Full config specs</summary>

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
  - if:
      host: "example.com"
      path: "/"
      method: "GET"

    then:
      forward: internal.example.com
      redirect: internal.example.com
      respond:
        status: 200
        body: hello
        headers:
          X-MyHeader: my-value

    middleware:
      cache:
        ttl: 60s
```

</details>

---

<details>
<summary><b>How AI was used on this project.</b></summary>

Because it feels increasingly important to specify how AI and LLMs where used by the person behind the keyboard:

- **Research:** the main use of AI by far. Mainly to understand standard reverse-proxy behaviors and also for exploring the Go stdlib a bit.
- **Code Generation:** a few tiny functions like [`matchSegments`](https://github.com/AnatoleLucet/traefik-test/blob/main/pkg/router/router.go#L109) were generated, but I tweaked them beyond recognition to better match my style and the behaviors I wanted ([⛵](https://en.wikipedia.org/wiki/Ship_of_Theseus)).
- **Code Review:** an agent reviewed the code for any obvious mistake I might have missed.

Apart from that, everything was done *old-school* by hand.

</details>
