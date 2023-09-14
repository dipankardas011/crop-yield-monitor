# crop-yield-monitor

# Project Documents
- [Project Docs](https://docs.google.com/document/d/1OA2uJ1kn7ileYLgplpS6gbcuHSR6UofU6R5m0gyzHmE/edit?usp=sharing)

- **Design Doc** -> folder `design/`

![design image](design-plan.jpeg)


# Demo Try

```bash
cd src/

cd authentication
docker build -t auth .
docker run -p 8080:8080 auth
cd -

cd predict
docker build -t predict .
docker run -p 8082:8082 predict

cd -

cd images
docker build -t images .
docker run -p 8081:8081 images

cd -
```

to try the routes
go to webbroswer

0.0.0.0:8080/swaggerui
0.0.0.0:8081/swaggerui
0.0.0.0:8082/swaggerui
