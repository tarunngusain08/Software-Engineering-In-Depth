# Docker

Docker simplifies application deployment by packaging code, dependencies, and runtime into containers, ensuring consistency across environments. Here's what you need to know:

---

## Key Concepts
- **Image**: Blueprint of your application.
- **Container**: Running instance of an image.
- **Dockerfile**: Instructions to build an image.
- **Volume**: Persistent storage for containers.
- **Registry**: Repository for storing and sharing images (e.g., Docker Hub).

---

## Essential Commands

### Containers
- List running containers: `docker ps`
- Stop/start a container: `docker stop/start <container_id>`
- Remove a container: `docker rm <container_id>`

### Images
- List images: `docker images`
- Build an image: `docker build -t <image_name>:<tag> .`
- Run a container: `docker run -d -p <host_port>:<container_port> <image_name>`
- Remove an image: `docker rmi <image_id>`

### Debugging
- View logs: `docker logs <container_id>`
- Access shell: `docker exec -it <container_id> /bin/bash`

### Cleanup
- Remove stopped containers: `docker container prune`
- Remove unused images: `docker image prune`
- Remove all unused resources: `docker system prune -a`

---

### Pulling Postgres Image
<img width="756" alt="Screenshot 2024-12-26 at 8 53 12 PM" src="https://github.com/user-attachments/assets/09113be3-3814-4e7b-ac7c-5775e934aa37" />

---

### Running Postgres Image on a container
<img width="1160" alt="Screenshot 2024-12-26 at 8 53 42 PM" src="https://github.com/user-attachments/assets/dea05791-7838-4a23-a970-0a245e549d24" />

---

### Removing Postgres Image and the container
<img width="1051" alt="Screenshot 2024-12-26 at 8 56 03 PM" src="https://github.com/user-attachments/assets/1cbf536b-c919-41d5-8244-396e7385ee0f" />

---

### Running Mysql Image
<img width="978" alt="Screenshot 2024-12-26 at 9 04 31 PM" src="https://github.com/user-attachments/assets/31f0751b-6087-442a-8a0b-7ab368afd2f1" />

<img width="586" alt="Screenshot 2024-12-26 at 9 04 43 PM" src="https://github.com/user-attachments/assets/2bede068-4d94-4a51-8127-ca0f10b3a846" />

---

### Inspecting containers and Network
<img width="1265" alt="Screenshot 2024-12-26 at 9 05 57 PM" src="https://github.com/user-attachments/assets/733eca19-c8a2-4bfd-8268-b7d1ade8a3ce" />

<img width="698" alt="Screenshot 2024-12-26 at 9 07 23 PM" src="https://github.com/user-attachments/assets/b926255c-6a05-4342-bb54-c7597cd25d4e" />

