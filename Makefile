####### HTTP Server
recommend:
	@echo "building recommendation"
	cd recommendation && docker build -t recommend .

image:
	@echo "building image"
	cd imageprocessing && docker build -t image .

auth:
	@echo "building auth"
	cd authentication && docker build -t auth .

####### Database
image-db:
	@echo "building image-db"
	cd db/image && docker build -t image-db .

auth-db:
	@echo "building auth-db"
	cd db/auth && docker build -t auth-db .

####### Build all containers
build: recommend image auth auth-db image-db
	@echo "Building done"

####### Run the backend containers(each)
recommend_run:
	@echo "running recommend"
	docker run -dp 8100:8100 --name recommend recommend

image_run:
	@echo "running image"
	docker run -dp 8090:8090 --name image image

auth_run:
	@echo "running auth"
	docker run -dp 8080:8080 --name auth auth

######## Start the docker compose
run: 
	docker compose up --detach
	@echo "Running done"

run-watch: 
	docker compose up
	@echo "Running done"

destroy:
	@echo "deleting containers"
	docker compose down

destroy_all:
	@echo "deleting containers and images"
	make destroy
	docker rmi -f recommend auth image auth-db image-db
