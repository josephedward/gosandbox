#!/bin/bash

# CKAD Cloud Study Session Script
# Integrates gosandbox + kubelingo for comprehensive cloud-native learning

set -e

STUDY_SESSION_DIR="./ckad-study-$(date +%Y%m%d-%H%M%S)"
GOSANDBOX_DIR="./tools/gosandbox"

echo "🎯 Starting CKAD Cloud Study Session"
echo "📁 Session directory: $STUDY_SESSION_DIR"

# Create session directory
mkdir -p "$STUDY_SESSION_DIR"
cd "$STUDY_SESSION_DIR"

# Initialize gosandbox
echo "🔧 Initializing AWS Sandbox Environment..."
go run "$GOSANDBOX_DIR/main.go" --init-session

# Wait for credentials
echo "⏳ Waiting for AWS credentials..."
while [ ! -f "./aws-credentials.json" ]; do
    sleep 5
    echo "   Still waiting for credentials..."
done

# Export AWS credentials
echo "🔑 Setting up AWS environment..."
export AWS_ACCESS_KEY_ID=$(jq -r '.aws_access_key_id' aws-credentials.json)
export AWS_SECRET_ACCESS_KEY=$(jq -r '.aws_secret_access_key' aws-credentials.json)
export AWS_SESSION_TOKEN=$(jq -r '.aws_session_token' aws-credentials.json)
export AWS_DEFAULT_REGION=$(jq -r '.aws_default_region' aws-credentials.json)

# Verify AWS access
echo "✅ Verifying AWS access..."
aws sts get-caller-identity

# Create EKS cluster for practice
echo "🚀 Creating EKS cluster for CKAD practice..."
cat > eks-cluster.yaml << EOF
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: ckad-practice
  region: ${AWS_DEFAULT_REGION}

nodeGroups:
  - name: worker-nodes
    instanceType: t3.medium
    desiredCapacity: 2
    minSize: 1
    maxSize: 3
    volumeSize: 20
    ssh:
      allow: true

addons:
  - name: aws-ebs-csi-driver
  - name: coredns
  - name: kube-proxy
  - name: vpc-cni
EOF

eksctl create cluster -f eks-cluster.yaml

# Setup kubeconfig
echo "⚙️  Configuring kubectl..."
aws eks update-kubeconfig --region $AWS_DEFAULT_REGION --name ckad-practice

# Verify cluster access
echo "🔍 Verifying cluster access..."
kubectl get nodes
kubectl get pods --all-namespaces

# Create CKAD practice namespaces and resources
echo "📚 Setting up CKAD practice environment..."
kubectl create namespace ckad-practice
kubectl create namespace web-tier
kubectl create namespace data-tier
kubectl create namespace monitoring

# Deploy sample workloads for practice
kubectl apply -f - << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: ckad-practice
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.21
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
  namespace: ckad-practice
spec:
  selector:
    app: nginx
  ports:
  - port: 80
    targetPort: 80
  type: LoadBalancer
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-type: "nlb"
EOF

# Launch kubelingo with cloud exercises
echo "🎮 Launching kubelingo for interactive CKAD practice..."
python3 ../kubelingo/cli_quiz.py --cloud-mode --cluster-context ckad-practice

# Cleanup function
cleanup() {
    echo "🧹 Cleaning up resources..."
    eksctl delete cluster --name ckad-practice --region $AWS_DEFAULT_REGION
    echo "✅ Cleanup completed"
}

# Register cleanup function
trap cleanup EXIT

echo "🎉 CKAD Cloud Study Session Complete!"
echo "📊 Session logs saved in: $STUDY_SESSION_DIR"