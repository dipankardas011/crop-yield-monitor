
cd recommendation ; docker build -t recommend . ; cd 

cd imageprocessing ; docker build -t image . ; cd ..

cd authentication ; docker build -t auth . ; cd ..

cd lb-backend ; docker build -t lb-backend . ; cd ..

cd frontend ; docker build -t frontend . ; cd ..

cd db\image ; docker build -t image-db . ; cd ..\..

cd db\auth ; docker build -t auth-db . ; cd ..\..

cd db\recommend ; docker build -t recommend-db . ; cd ..\..

docker compose up

# then visit the localhost:8080

