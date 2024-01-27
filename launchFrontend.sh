docker build -f Dockerfile.frontend -t frontend:latest .
docker run -it --rm -p 8080:80 frontend:latest
