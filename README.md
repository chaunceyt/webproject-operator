# WebProject Operator

After reading this article [Create a Kubernetes Operator in Golang to automatically manage a simple, stateful application](https://developers.redhat.com/blog/2020/12/16/create-a-kubernetes-operator-in-golang-to-automatically-manage-a-simple-stateful-application/) I decided to develop this experimental operator. 

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


```
operator-sdk-v0.15.2 build operator:v0.0.1
kind load docker-image operator:v0.0.1


```

Create a webproject `kubectl apply -f deploy/crds/webproject-1.yaml`

View items created `kubectl get all,secrets,cm,pvc,ing -l release=issue-403-fixing-footer`


```
kubectl get all,secrets,cm,pvc,ing -l release=issue-403-fixing-footer
NAME                                  TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
service/issue-403-fixing-footer-svc   ClusterIP   10.96.127.26   <none>        80/TCP    12m

NAME                                      READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/issue-403-fixing-footer   1/1     1            1           12m

NAME                                        TYPE     DATA   AGE
secret/issue-403-fixing-footer-aws-secret   Opaque   4      12m
secret/issue-403-fixing-footer-secret       Opaque   1      12m

NAME                                              DATA   AGE
configmap/issue-403-fixing-footer-common-config   6      12m
configmap/issue-403-fixing-footer-env-config      2      12m

NAME                                                  STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
persistentvolumeclaim/issue-403-fixing-footer-data    Bound    pvc-873c10f8-7364-433c-8afd-fb3f621f84a7   10Gi       RWO            standard       12m
persistentvolumeclaim/issue-403-fixing-footer-files   Bound    pvc-e6f4447b-2482-4a75-9faa-cc3a59a4a79d   10Gi       RWO            standard       12m

NAME                                                 CLASS    HOSTS                                     ADDRESS   PORTS     AGE
ingress.extensions/issue-403-fixing-footer-ingress   <none>   issue-403-fixing-footer.kube.domain.tld             80, 443   12m
```