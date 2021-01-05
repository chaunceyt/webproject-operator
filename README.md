# WebProject Operator

I currently manage a small GKE cluster. This cluster has two types of nodes, builder and worker. The builder nodes run gitlab-runners on `shared` and `project specific` builder node pools. The worker nodes are used by all of the projects using the cluster. Each project gets a namespace and can have a workload per branch. Currently the workload is generated via a custom helm chart, embeded in a "project starter" tool. Keeping the "project starter" and projects that started from the "project starter" in sync introduces some toil.

After reading this article [Create a Kubernetes Operator in Golang to automatically manage a simple, stateful application](https://developers.redhat.com/blog/2020/12/16/create-a-kubernetes-operator-in-golang-to-automatically-manage-a-simple-stateful-application/) and reviewing this [git repo](https://github.com/priyanka19-98/Wordpress-Operator). I decided to develop this experimental operator to see if I could reduce the toil of maintaining the helm chart.

The Operator is built using the `operator-sdk-v0.15.2` framework to demostrate

- creating a WebProject kind of resource using Kubernetes controller pattern. 
- creating a custom controller that encapsulates specific domain/application level knowledge of running mostly Drupal development projects in the cluster.

The Operator manages the following objects.

- Deployment
- ConfigMaps
- PersistentVolumeClaim(s)
- Ingress
- Service
- Secrets
 

## About the kind: Deployment

The deployment object managed by this operator creates a `web` container and a number of other sidecar containers

- cli (docksal default)
- database (mysql, mariadb)
- cache (memcache, redis)

## About the kind: PersistentVolumeClaim

The operator manages a pvc for the database `/var/lib/mysql` and static files. i.e. `/var/www/build/html/sites/default/files`

## About the kind: Secret

The operator manages a secret for the `MYSQL_PASSWORD` at the moment.

## About the kind: ConfigMap

The operator manages a configmap for `env` and `common` environment variables and `init-container` script. 

## About the kind: Service

The operator manages one service for the webproject.

## About the kind: Ingress

The operator manages one ingress for the webproject. At the moment the pattern only supports one domain per workload.


