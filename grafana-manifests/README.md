# Grafana manifests

This folder holds manifests to deploy Grafana into the `monitoring` namespace of your Kubernetes cluster. Please note that this is mainly for development and demo purposes and these manifests are not meant to be used in a production environment.

Deploy Grafana into your cluster:
```
kubectl apply -f .
```

Login to Grafana by getting the EXTERNAL-IP:
```
kubectl get service grafana -n monitoring
```

Login with the default credentials `admin` / `admin`.

