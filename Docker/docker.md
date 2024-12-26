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

