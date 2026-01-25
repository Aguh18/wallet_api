# Kubernetes Deployment Guide

Complete Kubernetes manifests untuk deploy Wallet API ke production.

## 📁 Manifests

- `namespace.yaml` - Namespace `wallet` untuk isolasi
- `configmap.yaml` - Configuration (DB host, port, log level)
- `secret.yaml` - Credentials (DB user/pass, JWT secret) ⚠️
- `postgres.yaml` - PostgreSQL + PersistentVolume (5Gi)
- `deployment.yaml` - Wallet API (3 replicas, rolling updates)
- `service.yaml` - ClusterIP & LoadBalancer
- `ingress.yaml` - Nginx ingress + SSL/TLS
- `hpa.yaml` - Horizontal Pod Autoscaler (3-10 pods)

## 🚀 Quick Deploy

```bash
# Deploy all
kubectl apply -f k8s/

# Check status
kubectl get all -n wallet

# Port forward
kubectl port-forward -n wallet svc/wallet-api-service 8080:80

# Test
curl http://localhost:8080/healthz
```

## 📋 Step by Step

### 1. Create Namespace & Config

```bash
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
```

### 2. Deploy Database

```bash
kubectl apply -f postgres.yaml

# Wait for ready
kubectl wait --for=condition=ready pod -l app=postgres -n wallet --timeout=300s
```

### 3. Deploy Application

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# Wait for ready
kubectl wait --for=condition=ready pod -l app=wallet-api -n wallet --timeout=300s
```

### 4. Enable Auto-scaling

```bash
kubectl apply -f hpa.yaml
```

### 5. Setup Ingress (Optional)

```bash
# Update domain di ingress.yaml
kubectl apply -f ingress.yaml
```

## ⚙️ Configuration

### Update Secrets (Production)

```bash
kubectl create secret generic wallet-api-secret \
  --from-literal=DB_USER=prod_user \
  --from-literal=DB_PASSWORD=secure_pass \
  --from-literal=JWT_SECRET=very-secure-key \
  -n wallet --dry-run=client -o yaml | kubectl apply -f -
```

### Update ConfigMap

```bash
kubectl edit configmap wallet-api-config -n wallet
```

## 📊 Monitoring

### Check Status

```bash
# All resources
kubectl get all -n wallet

# Pods
kubectl get pods -n wallet

# HPA
kubectl get hpa -n wallet

# Services
kubectl get svc -n wallet
```

### View Logs

```bash
# API logs
kubectl logs -f deployment/wallet-api -n wallet

# Postgres logs
kubectl logs -f deployment/postgres -n wallet

# Specific pod
kubectl logs <pod-name> -n wallet
```

### Describe Resources

```bash
kubectl describe deployment wallet-api -n wallet
kubectl describe pod <pod-name> -n wallet
kubectl describe hpa wallet-api-hpa -n wallet
```

## 🌐 Access Application

### Via Port Forward (Local)

```bash
kubectl port-forward -n wallet svc/wallet-api-service 8080:80
curl http://localhost:8080/healthz
```

### Via LoadBalancer

```bash
# Get external IP
kubectl get svc wallet-api-lb -n wallet

# Access
curl http://<EXTERNAL-IP>/healthz
```

### Via Ingress

```bash
# Configure DNS
# wallet-api.yourdomain.com → LoadBalancer IP

# Access
curl https://wallet-api.yourdomain.com/healthz
```

## 🔧 Troubleshooting

### Pod Not Starting

```bash
# Check events
kubectl get events -n wallet --sort-by='.lastTimestamp'

# Describe pod
kubectl describe pod <pod-name> -n wallet

# Check logs
kubectl logs <pod-name> -n wallet
```

### Database Connection Issues

```bash
# Test connectivity
kubectl run -it --rm psql-client --image=postgres:18.1-alpine -n wallet -- \
  psql -h postgres-service -U wallet_user -d wallet_db

# Check postgres pod
kubectl logs deployment/postgres -n wallet
```

### HPA Not Scaling

```bash
# Check metrics
kubectl top pods -n wallet
kubectl top nodes

# Describe HPA
kubectl describe hpa wallet-api-hpa -n wallet

# Ensure metrics-server installed
kubectl get deployment metrics-server -n kube-system
```

## 🗑️ Cleanup

```bash
# Delete all resources
kubectl delete namespace wallet

# Or delete individually
kubectl delete -f k8s/
```

## 🔐 Security Notes

⚠️ **Before Production:**

1. Update `secret.yaml` dengan production credentials
2. Jangan commit secrets ke git
3. Gunakan sealed-secrets atau external secret manager
4. Enable HTTPS di Ingress
5. Configure network policies
6. Review resource limits

## 📈 Scaling

### Manual Scaling

```bash
# Scale deployment
kubectl scale deployment wallet-api --replicas=5 -n wallet

# Check status
kubectl get pods -n wallet
```

### Auto-scaling (HPA)

HPA sudah configured:
- Min: 3 replicas
- Max: 10 replicas
- CPU target: 70%
- Memory target: 80%

```bash
# Check HPA
kubectl get hpa -n wallet

# Load test to trigger scaling
kubectl run -it load-test --rm --image=busybox -n wallet -- \
  sh -c "while true; do wget -q -O- http://wallet-api-service/healthz; done"
```

## 🏗️ Architecture

```
┌─────────────────────────────────────┐
│          Ingress (SSL/TLS)          │
│     wallet-api.yourdomain.com       │
└──────────────┬──────────────────────┘
               │
               ▼
┌──────────────────────────────────────┐
│      wallet-api-service (ClusterIP)  │
│              Port 80                 │
└──────────────┬───────────────────────┘
               │
       ┌───────┴────────┐
       ▼                ▼
┌─────────────┐  ┌─────────────┐
│ wallet-api  │  │ wallet-api  │  ... (3-10 pods)
│   Pod 1     │  │   Pod 2     │
└──────┬──────┘  └──────┬──────┘
       │                │
       └────────┬───────┘
                ▼
    ┌──────────────────────┐
    │  postgres-service    │
    │    (ClusterIP)       │
    └──────────┬───────────┘
               ▼
    ┌──────────────────────┐
    │   postgres Pod       │
    │ + PersistentVolume   │
    └──────────────────────┘
```

## 📚 Resources

- [Kubernetes Docs](https://kubernetes.io/docs/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [HPA Guide](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
