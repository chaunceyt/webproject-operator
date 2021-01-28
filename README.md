# WebProject Operator

The WebProject operator is a Kubernetes Native Operator for managing a development team's testing environment.


### Purpose

I'm the cluster administrator for a multi-tenant GKE cluster. The workloads that runs in this cluster follow the sidecar pattern. Each workload running in the cluster is created after a `git push` of a `bug/feature/issue` branch that triggers a Gitlab pipeline, during the **deploy** stage of the pipeline we use `helm template` to generate the manifests needed to create the Kubernetes workload for the code being tested. This helm chart is introduced via the forking of a "starter kit" when a project is starting sprint zero and defining what components they need for their workloads. There are times when projects who forked "months ago" have a fair amount of drift and keeping the helm chart in sync adds toil.


After reading this article [Create a Kubernetes Operator in Golang to automatically manage a simple, stateful application](https://developers.redhat.com/blog/2020/12/16/create-a-kubernetes-operator-in-golang-to-automatically-manage-a-simple-stateful-application/) and reviewing this [git repo](https://github.com/priyanka19-98/Wordpress-Operator). I decided to write this operator to better understand the Kubernetes operator pattern, to create an "automated site reliability engineer" for webprojects by removing the toil of maintaining a complex chart across a number of projects.

### What is a webproject?

A managed Kubernetes environment for development teams to test their web applications. 

The core Kubernetes components managed are:


- ConfigMaps
- CronJobs
- Deployment
- Ingress
- PersistentVolumeClaims
- Secrets
- Service
- Volumesnapshots/Volumecloning (TODO)


The core containers in a webproject pod:

- Webcontainer (can be almost anything that listens on port 80)
- Database sidecar (Mariadb or Mysql)
- Cache sidecar (Memcached or Redis)
- CLI container sidecar (same mountpaths as webcontainer)
- Search sidecar (Solr or ElasticSearch)


Managed Cronjobs

- Cronjobs for each sidecar (WIP)
- Backup for database sidecar (WIP)




### Version

- v1alpha1



