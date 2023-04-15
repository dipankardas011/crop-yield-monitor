# how to run the docker build and run

```bash
cd src/classifer
mkdir container-models
```

```bash
docker build -t classifer .
```

## to train the model
```bash
docker run -t -v $(pwd)/container-models:/app/models classifer main.py
```

## to run the prediction
```bash
docker run -t -v $(pwd)/container-models:/app/models classifer predict.py
```
