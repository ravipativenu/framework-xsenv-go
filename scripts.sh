docker build -t myapp:latest .
docker tag myapp:latest ravipativenu/myapp:latest
docker push ravipativenu/myapp:latest
kubectl apply -f deployment.yaml
kubectl get deployments
kubectl apply -f service.yaml
kubectl get services



