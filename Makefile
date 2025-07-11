.PHONY: img img_push img_rm sqlc

IMAGE_NAME = mdmitrym/food_delivery_registration
TAG ?= latest

img:
	docker build -t ${IMAGE_NAME}:${TAG} .

img_push:
	docker push ${IMAGE_NAME}:${TAG}

img_rm:
	docker image rm ${IMAGE_NAME}:${TAG}

sqlc:
	sqlc generate