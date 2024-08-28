###############################################################################
# Variables
###############################################################################
KIND_CLUSTER_NAME = local-k8s # context will be "kind-local-k8s"
HELM_RELEASE_NAME = campgrounds

###############################################################################
# Target: cluster-deploy
###############################################################################
.PHONY: cluster-deploy
cluster-deploy:
	@kind create cluster --name $(KIND_CLUSTER_NAME) --config ./k8s/kind-config.yaml
	@kubectl cluster-info --context kind-$(KIND_CLUSTER_NAME)

################################################################################
# Target: cluster-remove
################################################################################
.PHONY: cluster-remove
cluster-remove:
	@kind delete cluster --name $(KIND_CLUSTER_NAME)

################################################################################
# Target: postgres-deploy
################################################################################
.PHONY: postgres-deploy
postgres-deploy:
	@kubectl create secret generic postgres-secret --from-literal=POSTGRES_PASSWORD=postgres
	@kubectl create secret generic campgrounds-secret --from-literal=CAMPGROUNDS_PASSWORD=campgrounds_pass
	@kubectl create configmap initdb-config --from-file=./db/init/
	@kubectl get configmap initdb-config -o yaml
	@kubectl apply -f ./k8s/postgres.yaml

################################################################################
# Target: postgres-remove
################################################################################
.PHONY: postgres-remove
postgres-remove:
	@kubectl get secret postgres-secret > /dev/null 2>&1 \
    		&& kubectl delete secret postgres-secret \
    		|| echo "secret 'postgres-secret' does not exist."
	@kubectl get secret campgrounds-secret > /dev/null 2>&1 \
			&& kubectl delete secret campgrounds-secret \
			|| echo "secret 'campgrounds-secret' does not exist."
	@kubectl get configmap initdb-config > /dev/null 2>&1 \
		&& kubectl delete configmap initdb-config \
		|| echo "configmap 'initdb-config' does not exist."
	@kubectl delete deployment postgres

################################################################################
# Target: api-deploy
################################################################################
.PHONY: api-deploy
api-deploy:
	@kubectl apply -f ./k8s/campgrounds/api.yaml

################################################################################
# Target: api-remove
################################################################################
.PHONY: api-remove
api-remove:
	@kubectl delete deployment campgrounds-api
