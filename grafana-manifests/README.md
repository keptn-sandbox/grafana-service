# Grafana manifests

This folder holds manifests to deploy Grafana into the `monitoring` namespace of your Kubernetes cluster. Please note that this is mainly for development and demo purposes and these manifests are not meant to be used in a production environment.

1. Deploy Grafana into your cluster:
    ```
    kubectl apply -f .
    ```

1. Login to Grafana by getting the EXTERNAL-IP:
    ```
    kubectl get service grafana -n monitoring
    ```

1. Login with the default credentials `admin` / `admin`.

1. Once logged in, navigate to `Configuration -> API Keys` and generate a new API key with the role *Admin*. If the automatic creation of a Prometheus datasource is not needed, the role *Editor* is sufficient. 

1. Use the API Key in the `deploy/service.yaml` file to configure the Grafana-service to talk to your Grafana instance.

