<h1 style="color: #1bf89cff; font-size: 48px; font-weight: bold;">K8s based hosting</h1>


## Kindergarten Registry

A full-stack React + Node.js + MongoDB application deployed using K8s.


## Step-by-Step Deployment Guide

### Launch EC2 Instance
1. Create an EC2 instance using **Ubuntu 22.04 (t3.micro)**.  
2.  Connect via SSH:
   ```bash
   ssh -i key.pem ubuntu@<EC2_PUBLIC_IP>
```

3. Clone the repository

```bash
cd ~
https://github.com/Nabil720/Linux_Docker_k8s_Deployment.git
cd kindergarten-registry_GitOps
```

```
Linux_Docker_k8s_Deployment/kindergarten-registry_Docker/
├── docker-compose.yml
├── employeeservice
│   ├── database
│   │   └── db.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── handlers
│   │   └── employee.go
│   ├── main.go
│   └── models
│       └── employee.go
├── frontend
│   ├── Dockerfile
│   ├── package.json
│   ├── package-lock.json
│   ├── public
│   │   ├── favicon.ico
│   │   ├── index.html
│   │   ├── logo192.png
│   │   ├── logo512.png
│   │   ├── manifest.json
│   │   └── robots.txt
│   ├── README.md
│   └── src
│       ├── App.css
│       ├── App.js
│       ├── App.test.js
│       ├── components
│       │   ├── EmployeeForm.js
│       │   ├── EmployeeList.js
│       │   ├── StudentForm.js
│       │   ├── StudentList.js
│       │   ├── TeacherForm.js
│       │   └── TeacherList.js
│       ├── index.css
│       ├── index.js
│       ├── logo.svg
│       ├── reportWebVitals.js
│       └── setupTests.js
├── RIDEME.md
├── studentservice
│   ├── database
│   │   └── db.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── handlers
│   │   └── student.go
│   ├── main.go
│   └── models
│       └── student.go
└── teacherservice
    ├── database
    │   └── db.go
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── handlers
    │   └── teacher.go
    ├── main.go
    └── models
        └── teacher.go





# GitOps Files

kindergarten-registry_GitOps$ tree
.
├── employee-service-deployment.yaml
├── frontend-deployment.yaml
├── mongo-deployment.yaml
├── student-service-deployment.yaml
└── teacher-service-deployment.yaml


```


## kubectl apply -f .




##  Applicatio architecture


```
                        +---------------------+
                        |     Web Browser     |
                        | (React Frontend UI) |
                        +-------+-------------+
                                |
                                | (http://<node-ip>:30000) - Frontend
                                | (http://<node-ip>:30001) - Student API  
                                | (http://<node-ip>:30002) - Teacher API
                                | (http://<node-ip>:30003) - Employee API
                                v
                        +---------------------------------------------------+
                        |               Kubernetes Cluster                  |
                        |                                                   |
                        |  +---------------------+  +---------------------+ |
                        |  |  Frontend Service   |  | Backend Services    | |
                        |  |   (NodePort:30000)  |  | (NodePorts)         | |
                        |  +----------+----------+  | • Student:30001     | |
                        |             |             | • Teacher:30002     | |
                        |             |             | • Employee:30003    | |
                        |  +----------v----------+  +----------+----------+ |
                        |  |   Frontend Pod      |             |            |
                        |  | (React Container)   |             |            |
                        |  +----+----------+-----+  +----------v----------+ |
                        |       |          |        |   Backend Pods      | |
                        |       |          |        | • Student (5001)    | |
                        |       |          |        | • Teacher (5002)    | |
                        |       |          |        | • Employee (5003)   | |
                        |       |          |        +----------+----------+ |
                        |       |          |                   |            |
                        |       v          v                   |            |
                        |  +-----------+  +-------------------+ |            |
                        |  | React     |  | Internal API     | |            |
                        |  | Static    |  | Calls via Service| |            |
                        |  | Files     |  | Names            | |            |
                        |  |           |  | • student-service| |            |
                        |  +-----------+  | • teacher-service| |            |
                        |                 | • employee-service| |            |
                        |                 +-------------------+ |            |
                        |                                      |            |
                        |                            +---------v----------+ |
                        |                            |   MongoDB Service  | |
                        |                            |   (mongo:27017)    | |
                        |                            +---------+----------+ |
                        |                                      |            |
                        |                            +---------v----------+ |
                        |                            |   MongoDB Pod     | |
                        |                            |  (Data Storage)   | |
                        |                            +-------------------+ |
                        +---------------------------------------------------+



```
![Website View](./Images/Screenshot%20from%202025-11-11%2014-08-23.png)


![Instance View](./Images/Screenshot%20from%202025-11-11%2014-08-49.png)


![Instance View](./Images/Screenshot%20from%202025-11-11%2014-08-58.png)

![Instance View](./Images/Screenshot%20from%202025-11-11%2014-09-07.png)

