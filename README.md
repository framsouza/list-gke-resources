# Listing gke resources

This is a simple code that will list the Kubernetes resources using client-go.
The resources available at the moment is, `pods`, `deployments`, `nodes`, `statefulsets` and `policysecurity`.

### Usage
You must specify the project and resource you want to check by adding the `-project` and `-resource` parameters, like this:
```
go run main.go -project=myproject -resource=pods
```

You will get an output similar to this:
```
GKE NAME                          RESOURCE NAME                          NAMESPACE            
migration-test                    event-exporter-gke                     kube-system          
migration-test                    konnectivity-agent                     kube-system          
migration-test                    konnectivity-agent-autoscaler          kube-system          
migration-test                    kube-dns                               kube-system          
migration-test                    kube-dns-autoscaler                    kube-system          
migration-test                    l7-default-backend                     kube-system          
migration-test                    metrics-server-v0.4.5                  kube-system          

GKE NAME                                             RESOURCE NAME                          NAMESPACE            
migration-test-twenty-two-version                    event-exporter-gke                     kube-system          
migration-test-twenty-two-version                    konnectivity-agent                     kube-system          
migration-test-twenty-two-version                    konnectivity-agent-autoscaler          kube-system          
migration-test-twenty-two-version                    kube-dns                               kube-system          
migration-test-twenty-two-version                    kube-dns-autoscaler                    kube-system          
migration-test-twenty-two-version                    l7-default-backend                     kube-system          
migration-test-twenty-two-version                    metrics-server-v0.4.5                  kube-system          
migration-test-twenty-two-version                    event-exporter-gke                     kube-system          
migration-test-twenty-two-version                    konnectivity-agent                     kube-system          
migration-test-twenty-two-version                    konnectivity-agent-autoscaler          kube-system          
migration-test-twenty-two-version                    kube-dns                               kube-system          
migration-test-twenty-two-version                    kube-dns-autoscaler                    kube-system          
migration-test-twenty-two-version                    l7-default-backend                     kube-system          
migration-test-twenty-two-version                    metrics-server-v0.4.5                  kube-system   
```
