###############################################################################
# Variables
###############################################################################
KIND_CLUSTER_NAME = local-k8s # context will be "kind-local-k8s"
APP_DOCKER_IMAGE = ibaiborodine/campsite-booking-go

###############################################################################
# Target: image-build
###############################################################################
.PHONY: image-build
image-build:
	@docker build --no-cache -t $(APP_DOCKER_IMAGE) -f ./docker/Dockerfile .
	@docker image ls | grep $(APP_DOCKER_IMAGE)

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
	@docker volume prune -f

################################################################################
# Target: all-deploy
################################################################################
.PHONY: all-deploy
all-deploy: db-deploy api-deploy proxy-deploy

################################################################################
# Target: all-remove
################################################################################
.PHONY: all-remove
all-remove: proxy-remove api-remove db-remove

################################################################################
# Target: services-list
################################################################################
.PHONY: services-list
services-list:
	@kubectl get services \
		--all-namespaces \
		--output=jsonpath='{range .items[*]}{.metadata.name}.{.metadata.namespace}.svc.cluster.local{"\n"}{end}'

################################################################################
# Target: db-deploy
################################################################################
.PHONY: db-deploy
db-deploy:
	@kubectl create secret generic postgres-secret --from-literal=POSTGRES_PASSWORD=postgres
	@kubectl create secret generic campgrounds-secret --from-literal=CAMPGROUNDS_PASSWORD=campgrounds_pass
	@kubectl create configmap initdb-config --from-file=./db/init/
#	@kubectl get configmap initdb-config -o yaml
	@kubectl apply -f ./k8s/postgres.yaml

################################################################################
# Target: db-remove
################################################################################
.PHONY: db-remove
db-remove:
	@kubectl get secret postgres-secret > /dev/null 2>&1 \
    		&& kubectl delete secret postgres-secret \
    		|| echo "secret 'postgres-secret' does not exist."
	@kubectl get secret campgrounds-secret > /dev/null 2>&1 \
			&& kubectl delete secret campgrounds-secret \
			|| echo "secret 'campgrounds-secret' does not exist."
	@kubectl get configmap initdb-config > /dev/null 2>&1 \
		&& kubectl delete configmap initdb-config \
		|| echo "configmap 'initdb-config' does not exist."
	@kubectl delete statefulset postgres
	@kubectl delete pvc -l app=postgres

################################################################################
# Target: secret-campgrounds-show
################################################################################
.PHONY: secret-campgrounds-show
secret-campgrounds-show:
	@kubectl get secret campgrounds-secret -o jsonpath="{.data.CAMPGROUNDS_PASSWORD}"; echo

################################################################################
# Target: api-deploy
################################################################################
.PHONY: api-deploy
api-deploy:
	@kind load docker-image $(APP_DOCKER_IMAGE) --name local-k8s
	@kubectl apply -f ./k8s/campgrounds.yaml

################################################################################
# Target: api-remove
################################################################################
.PHONY: api-remove
api-remove:
	@kubectl delete deployment campgrounds

################################################################################
# Target: proxy-deploy
################################################################################
.PHONY: proxy-deploy
proxy-deploy:
	@kubectl create configmap envoy-config --from-file=./k8s/envoy-config.yaml
#	@kubectl get configmap envoy-config -o yaml
	@kubectl apply -f ./k8s/envoy.yaml

################################################################################
# Target: proxy-remove
################################################################################
.PHONY: proxy-remove
proxy-remove:
	@kubectl get configmap envoy-config > /dev/null 2>&1 \
		&& kubectl delete configmap envoy-config \
		|| echo "configmap 'envoy-config' does not exist."
	@kubectl delete deployment envoy
