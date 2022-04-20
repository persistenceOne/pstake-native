IMAGE_NAME=anmol1696/$(PROCESS)
TAG_NAME=latest
CONTAINER_NAME=container-$(PROCESS)
FILE=Dockerfile
PROCESS=pstake

DOCKER_ENV=\
	-e CHAIN_ID=pstake \
	-e PSTAKE_DIR=/pstaked \
	-e COINS="100000000000000000stake"
DOCKER_VOLUME=\
    -v $(PWD):/pstake \
	-v $(PWD)/k8s/pstake/configs:/configs \
	-v /tmp:/tmp
DOCKER_ARGS=\
	-p 26657:26657 \
	--network=host

# Command for dockerization
docker-build:
	docker build . -f docker/$(PROCESS)/$(FILE) -t $(IMAGE_NAME):$(TAG_NAME)

docker-run:
	docker run --rm -it --name=$(CONTAINER_NAME) \
		$(DOCKER_VOLUME) $(DOCKER_ARGS) $(DOCKER_ENV) \
		$(IMAGE_NAME):$(TAG_NAME) /bin/bash

docker-push:
	docker push $(IMAGE_NAME):$(TAG_NAME)

docker-build-push: docker-build docker-push

docker-clean-container:
	-docker stop $(CONTAINER_NAME)
	-docker rm $(CONTAINER_NAME)

docker-clean: clean-container
	-docker rmi $(IMAGE_NAME):$(TAG_NAME)
