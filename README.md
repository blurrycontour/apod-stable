# APOD's stable API

Uses cache and scrapes the astropix webpage to create api.

## Environment Variables
- `LISTEN_ADDR`=:8080
- `USER_AGENT`="apod-stable"
- `APOD_URL`="https://apod.nasa.gov/apod/astropix.html"

## Building and Running
1. Build the Docker image:
    ```bash
    docker build .
    # (optional) to push to Docker Hub:
    docker tag apod-stable <yourusername>/apod-stable:latest
    docker push <yourusername>/apod-stable:latest
    ```
2. Run the container:
    ```bash
    docker run -p 8080:8080 apod-stable
    ```
3. Access the API at `http://localhost:8080/apod`
