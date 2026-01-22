# Application Service Deployment Guide

This guide details the process for setting up and deploying the Application services (Student, Teacher, and Employee) using HashiCorp Vault for secret management and Docker for image building and pushing. The services will be deployed on Kubernetes using the provided Kubernetes manifest files.


## Step 1: Enable and Configure Vault Secrets(ON Vault VM)

### Enable Vault KV Secrets Engine

Enable the KV (Key-Value) secrets engine at the path `kindergarten` and configure it:

```bash
# Enable KV secrets engine at path "kindergarten" with version 2
vault secrets enable -path=kindergarten -version=2 kv

# Tune the "kindergarten" secrets engine with default TTL of 1 year and max TTL of 10 years
vault secrets tune -default-lease-ttl=8760h -max-lease-ttl=87600h kindergarten/

```

```bash
# MongoDB credentials
vault kv put kindergarten/mongodb \
  username="myUser" \
  password="myPassword" \
  database="kindergarten" \
  uri="mongodb://myUser:myPassword@mongo:27017/kindergarten?authSource=admin"

# Elastic APM configuration
vault kv put kindergarten/apm \
  server_url="http://192.168.121.224:8200" \
  secret_token="your_apm_access_tocken" \
  environment="production"

# Service-specific APM names
vault kv put kindergarten/services \
  student="student-service" \
  teacher="teacher-service" \
  employee="employee-service"

# Service port configurations
vault kv put kindergarten/ports \
  student=5001 \
  teacher=5002 \
  employee=5003

```


### Create Policies
Create Vault policies for  student service. These policie will give the services read-only access to the secrets in Vault.

```bash
# Define a policy for the student service that allows read access to specific paths in Vault
echo 'path "kindergarten/data/mongodb" {
  capabilities = ["read"]
}

path "kindergarten/data/apm" {
  capabilities = ["read"]
}

path "kindergarten/data/services" {
  capabilities = ["read"]
}

path "kindergarten/data/ports" {
  capabilities = ["read"]
}' > student-policy.hcl



# Write the student-policy to Vault using the configuration in student-policy.hcl
vault policy write student-policy student-policy.hcl


```


Create Vault policies for the teacher and employee services. These policie will give the services read-only access to the secrets in Vault.
```bash
# Define a policy for the teacher service that allows read access to specific paths in Vault
echo 'path "kindergarten/data/mongodb" {
  capabilities = ["read"]
}

path "kindergarten/data/apm" {
  capabilities = ["read"]
}

path "kindergarten/data/services" {
  capabilities = ["read"]
}

path "kindergarten/data/ports" {
  capabilities = ["read"]
}' > teacher-policy.hcl

# Write the teacher-policy to Vault using the configuration in teacher-policy.hcl
vault policy write teacher-policy teacher-policy.hcl

```
Create Vault policies for the employee service. These policie will give the services read-only access to the secrets in Vault.

```bash
# Define a policy for the employee service that allows read access to specific paths in Vault
echo 'path "kindergarten/data/mongodb" {
  capabilities = ["read"]
}

path "kindergarten/data/apm" {
  capabilities = ["read"]
}

path "kindergarten/data/services" {
  capabilities = ["read"]
}

path "kindergarten/data/ports" {
  capabilities = ["read"]
}' > employee-policy.hcl

# Write the employee-policy to Vault using the configuration in employee-policy.hcl
vault policy write employee-policy employee-policy.hcl
```


### Create Tokens for Each Service

```bash
# Student service token
vault token create -policy="student-policy" -ttl=8760h -display-name="student-service-token"

# Teacher service token
vault token create -policy="teacher-policy" -ttl=8760h -display-name="teacher-service-token"

# Employee service token  
vault token create -policy="employee-policy" -ttl=8760h -display-name="employee-service-token"

```


Save the tokens for later use:

* Student token: hvs.CA**YXo1Y0c

* Teacher token: hvs.CAE**6aFlwSkpQSGs

* Employee token: hvs.CAE**NKRjA





## If you want to create EKS Cluster Using eksctl
Before you begin, ensure you have the following tools and configurations in place:
 1. eksctl (EKS Cluster Management Tool)
 2. kubectl (Kubernetes Command-Line Tool)
 3. AWS CLI (AWS Command-Line Interface)

```bash
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: rnd-cluster
  region: us-east-1
  version: "1.34"
  tags:
    Environment: Dev
iam:
  withOIDC: true

vpc:
  cidr: 172.31.0.0/16
  clusterEndpoints:
    publicAccess: true   
    privateAccess: false 

availabilityZones:
  - us-east-1a
  - us-east-1b

managedNodeGroups:
  - name: ng-general
    instanceTypes: 
      - t3.large 
    minSize: 1
    maxSize: 6
    desiredCapacity: 1
    volumeSize: 40 
    volumeType: gp3
    tags:
      NodeGroup: General
      Environment: dev
    iam:
      attachPolicyARNs:
        - arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy
        - arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy
        - arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly
# iam:
#   withOIDC: true
#   podIdentityAssociations:
#   - namespace: kube-system
#     serviceAccountName: karpenter
#     roleName: karpenter-c1-karpenter
#     permissionPolicyARNs:
#     - arn:aws:iam::605134426044:policy/KarpenterControllerPolicy-karpenter-c1

# iamIdentityMappings:
# - arn: "arn:aws:iam::605134426044:role/KarpenterNodeRole-karpenter-c1"
#   username: system:node:{{EC2PrivateDNSName}}
#   groups:
#   - system:bootstrappers
#   - system:nodes
```


## Step 2: Build and Push Docker Images (ON Jump Host)


First, clone the repository from GitHub to your local environment.
```bash
git clone https://github.com/Nabil720/Cncf.git
cd Cncf
```

### Set Up Docker Build Script
Before running the build script, you need to update the following variables in the build_and_push.sh script:

* DOCKER_USERNAME: Your Docker Hub username.

* DOCKER_PASSWORD: Your Docker Hub personal access token.

* VAULT_ADDR: The Vault server address.

* VAULT_TOKEN: The Vault token (use the appropriate service token here).

* SERVICE_NAME: The name of the service (e.g., student, teacher, employee).

```bash
cd studentservice/
sudo chmod +x build_and_push.sh
./build_and_push.sh

cd ../teacherservice/
sudo chmod +x build_and_push.sh
./build_and_push.sh

cd ../employeeservice/
sudo chmod +x build_and_push.sh
./build_and_push.sh

```


## Step 3: Deploy Services to Kubernetes

### Apply Kubernetes Manifest Files
Navigate to the Kubernetes manifest directory and apply the manifest files for all services:
```bash
cd k8s_manifest/
kubectl apply -f mongo-deployment.yaml
kubectl apply -f employee-service-deployment.yaml
kubectl apply -f student-service-deployment.yaml
kubectl apply -f teacher-service-deployment.yaml
kubectl apply -f frontend-deployment.yaml
```