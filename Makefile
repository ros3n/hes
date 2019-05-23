GOOS?=linux

build-api:
	rm -f ./hes-api
	GOOS=${GOOS} go build -o hes-api github.com/ros3n/hes/api

build-mailer:
	rm -f ./hes-mailer
	GOOS=${GOOS} go build -o hes-mailer github.com/ros3n/hes/mailer

dockerize-api: build-api
	docker build -t ${DOCKER_USER}/hes-api:latest -f k8s/docker/hes-api.dockerfile .

dockerize-mailer: build-mailer
	docker build -t ${DOCKER_USER}/hes-mailer:latest -f k8s/docker/hes-mailer.dockerfile .

push-api: dockerize-api
	docker push ${DOCKER_USER}/hes-api:latest

push-mailer: dockerize-mailer
	docker push ${DOCKER_USER}/hes-mailer:latest
