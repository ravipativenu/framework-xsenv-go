docker build -t myapp:latest .
docker tag myapp:latest ravipativenu/myapp:latest
docker push ravipativenu/myapp:latest
kubectl apply -f deployment.yaml
kubectl get deployments
kubectl apply -f service.yaml
kubectl get services
kubectl logs myapp-deployment-657d6bb4b7-x8wwr



