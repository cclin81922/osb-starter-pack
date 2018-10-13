# Deploy broker using Helm

Deploy with Helm and pass custom image and tag name.
Note: This also pushes the generated image with docker.

```console
$ make deploy-sc
$ IMAGE=cclin81922/osbapi-broker TAG=latest PULL=Never make deploy-broker
```

Keep watch by `svcat get brokers broker-skeleton` until its status becomes ready.

# Use broker

```console
$ make provision
$ make bind
```

Verify

```console
$ svcat get instances -n test-ns
$ svcat get bindings -n test-ns
$ kubectl get secret -n test-ns
```

# Adding your business logic

To implement your broker, you fill out just a few methods and types in
`pkg/broker` package:

- The `Options` type, which holds options for the broker
- The `AddFlags` function, which adds CLI flags for an Options
- The methods of the `BusinessLogic` type, which implements the broker's
  business logic
- The `NewBusinessLogic` function, which creates a BusinessLogic from the
  Options the program is run with