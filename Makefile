recommend:
	@echo "building recommendation"
	cd recommendation && docker build -t recommend .

image:
	@echo "building image"
	cd imageprocessing && docker build -t image .

auth:
	@echo "building auth"
	cd authentication && docker build -t auth .

build: recommend image auth
	@echo "Building done"

recommend_run:
	@echo "running recommend"
	docker run -dp 8100:8100 --name recommend recommend

image_run:
	@echo "running image"
	docker run -dp 8090:8090 --name image image

auth_run:
	@echo "running auth"
	docker run -dp 8080:8080 --name auth auth

run: recommend_run image_run auth_run
	@echo "Running done"

destroy:
	@echo "deleting containers"
	docker rm -f auth image recommend

destroy_all:
	@echo "deleting containers and images"
	docker rmi -f recommend auth image
	make destroy
