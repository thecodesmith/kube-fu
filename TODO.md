# TODO

## Define output spec

## Pull all relevant data into structs first, then print it

something like this:
- find each daemonset, statefulset, deployment (any others)
  - find each pod for each
    - find the node for each

struct looks like this:
```
    {
      Pod,
      Node,
      Type
    }
```

## Add node resource usage requests/limits

## Colorize the output

- Different colors for each resource type
- Node resources
- Labels matching requested - maybe?
