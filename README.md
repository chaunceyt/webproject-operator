# WebProject Operator

After reading this article [Create a Kubernetes Operator in Golang to automatically manage a simple, stateful application](https://developers.redhat.com/blog/2020/12/16/create-a-kubernetes-operator-in-golang-to-automatically-manage-a-simple-stateful-application/) I decided to developer this experimental operator. 

The is Operator is to demostrate creating a WebProject kind of resource using Kubernetes controller pattern. The operator is to show how to build a custom controller that encapsulates specific domain/application level knowledge. The Operator is built using the operator-sdk framework.

The Operator manages the following objects.

- Deployment
- ConfigMaps
- PersitentVolumeClaim(s)
- Ingress
- Service
- Secrets

## About the kind: Deployment

The deployment object managed by this operator creates a number of sidecar containers

- web (primary)
- cli (docksal default)
- database (mysql, mariadb)
- cache (memcache, redis)

## About the kind: PersitentVolumeClaim

The operator manages a pvc for the database `/var/lib/mysql` and static files. i.e. `/var/www/build/html/sites/default/files`

## About the kind: Secret

The operator manages a secret for the `MYSQL_PASSWORD` at the moment.

## About the kind: ConfigMap

The operator manages a secret for `env` and `common` environment variables 

## About the kind: Service

The operator manages one service for the webproject.

## About the kind: Ingress

The operator manages one ingress for the webproject. At the moment the pattern only supports one domain per workload.


## Notes

After updating `webproject_types.go` run 

```
operator-sdk-v0.15.2 generate crds
operator-sdk-v0.15.2 generate k8s
kubectl apply -f deploy/crds/wp.com_webprojects_crd.yaml
operator-sdk-v0.15.2 run --local
```

Run in local mode

`operator-sdk-v0.15.2 run --local`

### Install and setup for development

```
wget https://github.com/operator-framework/operator-sdk/releases/download/v0.15.2/operator-sdk-v0.15.2-x86_64-apple-darwin
chmod +x ./operator-sdk-v0.15.2-x86_64-apple-darwin
mv operator-sdk-v0.15.2-x86_64-apple-darwin ~/bin/operator-sdk-v0.15.2
operator-sdk-v0.15.2 new webproject-operator --type go --repo github.com/chaunceyt/webproject-operator
operator-sdk-v0.15.2 add api --kind WebProject --api-version wp.com/v1
operator-sdk-v0.15.2 generate crds
operator-sdk-v0.15.2 generate k8s

kubectl apply -f deploy/crds/wp.com_webprojects_crd.yaml
# run locally for development
operator-sdk-v0.15.2 run --local

```