docker build -t api-gateway:latest .
docker stop api-gateway
docker rm api-gateway
docker run -d -p 3000:3000 --name api-gateway api-gateway:latest