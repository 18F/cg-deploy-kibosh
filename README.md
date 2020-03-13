# Kibosh

Deploy pipeline for the Cloud Foundry Kibosh service broker.

The Kibosh broker deploys pre-configured Helm charts to a Kubernetes cluster.
We're using it to deploy Elasticsearch and Redis service instances.

## Links

* [Kibosh](https://github.com/cf-platform-eng/kibosh)
* Pipeline (is this CUI?)
* Elasticsearch Helm chart - https://github.com/elastic/helm-charts/tree/master/elasticsearch
* Redis Helm charts - https://github.com/helm/charts/tree/master/stable/redis-ha and https://github.com/helm/charts/tree/master/stable/redis

## Deploying
### Kubernetes Requirements
Step 1 - Create kibosh user on targeted kbs

Create the template `kibosh.yml`:
``` 
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kibosh-broker-user
  namespace: default
secrets:
- name: kibosh-broker-user-secret
---
apiVersion: v1
kind: Secret
metadata:
  name: kibosh-broker-user-secret
  annotations:
    kubernetes.io/service-account.name: kibosh-broker-user
type: kubernetes.io/service-account-token
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: kkibosh-broker-user
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
subjects:
- kind: ServiceAccount
  name: kibosh-broker-user
  namespace: default
```

Step 2 - Create the user/role/secret:
```
kubectl apply -f kibosh.yml
```
Step 3 - Get kibosh-broker-user token for the broker:
```
kubectl describe secret kibosh-broker-user-secret
```
For token only
```
kubectl describe secret kibosh-broker-user-secret | grep token | awk '{print $2}'
```

### Push the Kibosh as a CF Application
From the `kibosh` folder. 

1. Update the `doc/sample-manifest.yaml` with required information.
```
---
applications:
  - name: kibosh_broker
    memory: 256M
    instances: 1
    buildpacks:
    - binary_buildpack
    command: ./kibosh-0.2.49.linux
    env:
      SECURITY_USER_NAME: <USERNAME_FOR_BROKER>
      SECURITY_USER_PASSWORD: <PASSWORD_FOR_BROKER>
      TILLER_NAMESPACE: kube-system
      CA_DATA: <CA_DATA_FROM_KUBE_SECRET>
      SERVER:  <ELB_FOR_K8S>
      TOKEN: <TOKEN_FROM_KUBE_SECRET>
```

CA_DATA, SERVER AND TOKEN can be taken from above requirements.

Charts are ready and will be used when pushed along with Kibosh.

2. CF Push with the following command
```
cf push -f docs/sample-manifest.yaml 
```

3. Register broker

Follow this doc: https://docs.cloudfoundry.org/services/managing-service-brokers.html

Once it's register, CF will talk to kibosh and will generate service plans automagically.


## Development
Note: Pretty much all of the development is done via the Helm Charts

Inside each folder udner `charts`, there are several files that is used to configure kibosh's usage of the Helm Chart

* `values.yaml` - This is where most of the configuration files for the service is done
* `plans.yml` - This is where you specify the plans for each service
* `bind.yml` - This is how you provide service credentials when the application is bound to the service
* `plans` folder - This is where you can configure how big and plan specific configs. 

## Architecture

![architecture](architecture.jpg)

### Kibosh

Kibosh itself is deployed as a Cloud Foundry application.

### Ingress Constraints

Applications need a way of connecting to their Elasticsearch or Redis service
instances.  These instances are StatefulSets running inside the Kubernetes
cluster, so unreachable from within the Cloud Foundry cluster.  To allow
connections, we need some form of Kubernetes Service (and possibly Ingress) in
between.

### Elasticsearch

Since Elasticsearch uses an HTTP-based protocol, we can use a single Nginx
Ingress controller to route to the right Elasticsearch service based on the
requested HTTP Host header.  Each Elasticsearch Helm chart includes an Ingress
resource with the correct host configuration.  All of these resources are
handled by the single Nginx controller, which is fronted by a single AWS
Network Load Balancer.

### Redis

Redis is more complicated, as it's [a custom TCP
protocol](https://redis.io/topics/protocol).  Because of this, we do not have a
Host header to use for routing.  We've considered a number of solutions:

#### NodePort Services

We could provision a NodePort service with each Redis instance.  This would
save money, but is non-viable due to changing node IP addresses.  We would need
to monitor the Kubernetes Node IP addresses as they change through
re-provisioning and cluster scaling, and update the service instances during
these events.

Worse, because Cloud Foundry doesn't (can't) update environment variables in
running application instances, the customer applications would still be
attempting to connect to non-existent Node IPs.

This is known as the "stale binding problem".

(Also, [think before you NodePort](https://oteemo.com/2017/12/12/think-nodeport-kubernetes/))

#### Round Robin DNS on the BOSH VMs

We could get around the stale binding problem by adding an agent on the BOSH
VMs that would monitor the Kubernetes Node IP addresses and maintain a set of
A-records for them in BOSH DNS.  Call this `k8s.cloud.gov.internal`, for
example.  Then we could point applications in Cloud Foundry at the
`k8s.cloud.gov.internal:NodePort` pair for their service instance.

However, this would require the authoring and maintenance of the agent.

Also, many applications _*cough*java*cough*_ don't play well with DNS TTLs,
meaning they'd experience prolonged downtime when the Kubernetes cluster
topology changes.

Finally, unless the BOSH DNS server randomizes the returned A-records,
applications would all pile up on the first Node.

#### Nginx TCP Ingress

(**This is our current path.**)

The Nginx controller [can actually act as a TCP
router](https://kubernetes.github.io/ingress-nginx/user-guide/exposing-tcp-udp-services/).
In this mode, it looks at ConfigMap entries to map separate external ports to
internal Services.

We could use this by injecting such a ConfigMap into each Helm deployment,
giving each Redis service a unique external port on the Nginx Network Load
Balancer.  We believe this configuration would work, but it does involve a bit
more configuration and management on our end.

#### Single Load Balancer per Redis

All of the above have been investigated in an attempt to reduce costs by not
requiring a separate AWS NLB for each Redis instance.  However, NLBs are
roughly $16 per month, and it looks as though we currently have fewer than 70
Redis instances deployed.

**If we believe this cost to be acceptable, then it's _much easier_ to
configure the Helm chart to provision a Kubernetes LoadBalancer, and call it a
day.**
