# WebProject Operator

The is Operator is to demostrate creating a WebProject kind of resource using Kubernetes controller pattern. The operator is to show how to build a custom controller that encapsulates specific domain/application level knowledge. The Operate is built using the operator-sdk framework.

The Operator manages the following objects.

- Deployment 
- Service
- Persitent Volume Claim(s)
- Ingress
- configmaps

## About the kind: Deployment

The deployment object managed by this operator creates a number of sidecar containers

- web (primary)
- cli (docksal default)
- database (mysql, mariadb)
- cache (memcache, redis)


