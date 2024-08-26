###############################################################################
# Variables
###############################################################################
KIND_CLUSTER_NAME = local-k8s # context will be "kind-local-k8s"
HELM_RELEASE_NAME = campgrounds

###############################################################################
# Target: cluster-create
###############################################################################
.PHONY: cluster-create
cluster-create:
	@kind create cluster --name $(KIND_CLUSTER_NAME) --config ./k8s/kind-config.yaml

################################################################################
# Target: cluster-delete
################################################################################
.PHONY: cluster-delete
cluster-delete:
	@kind delete cluster --name $(KIND_CLUSTER_NAME)

################################################################################
# Target: postgres-install
################################################################################
.PHONY: postgres-install
postgres-install:
	@kubectl create configmap init-db-config --from-file=./db/init/
	@kubectl get configmap init-db-config -o yaml
	@helm install $(HELM_RELEASE_NAME) bitnami/postgresql --version 15.2.10 -f ./k8s/postgres/values.yaml

################################################################################
# Target: postgres-uninstall
################################################################################
.PHONY: postgres-uninstall
postgres-uninstall:
	@kubectl get configmap init-db-config > /dev/null 2>&1 \
		&& kubectl delete configmap init-db-config \
		|| echo "ConfigMap 'init-db-config' does not exist."
	@helm uninstall $(HELM_RELEASE_NAME)
