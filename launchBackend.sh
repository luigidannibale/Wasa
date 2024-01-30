docker build -f Dockerfile.backend -t backend:latest .
docker run -it --rm -p 3000:3000 backend:latest 
