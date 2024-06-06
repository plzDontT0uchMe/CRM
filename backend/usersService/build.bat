docker stop auth-service
docker rm auth-service
docker build -t auth-service:latest .
docker run --name auth-service --link auth-postgres:auth-postgres --link redis:redis -p 3001:3001 -d auth-service:latest