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

lb:
	@echo "building lb-backend"
	cd lb-backend && docker build -t lb-backend .

ui:
	@echo "building frontend"
	cd frontend && docker build -t frontend .

model:
	@echo "building ml-model"
	cd ml && docker build -t ml .

####### Database
image-db:
	@echo "building image-db"
	cd db/image && docker build -t image-db .

auth-db:
	@echo "building auth-db"
	cd db/auth && docker build -t auth-db .

recommend-db:
	@echo "building recommend-db"
	cd db/recommend && docker build -t recommend-db .

####### Build all containers
build: recommend image auth lb ui model auth-db image-db recommend-db
	@echo "Building done"


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
	docker rmi -f recommend image auth lb-backend frontend auth-db image-db recommend-db
