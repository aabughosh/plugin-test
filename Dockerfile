# Stage 1: Build the frontend
FROM registry.access.redhat.com/ubi8/nodejs-16:latest AS web-builder
USER root

# Install yarn if not already installed
RUN command -v yarn || npm i -g yarn
WORKDIR /opt/app-root

COPY web web

# Install dependencies and build frontend
RUN cd web && yarn install && yarn build

# Stage 2: Build the Go backend
FROM golang:1.22 AS go-builder
WORKDIR /opt/app-root

# Copy the Go module files
COPY ./src/components/go.mod ./
COPY ./src/components/go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the rest of the backend files
COPY ./cmd/plugin-backend.go ./cmd/plugin-backend.go

# Build the backend
RUN go build -o plugin-backend ./cmd/plugin-backend.go

# Stage 3: Combine frontend and backend
FROM registry.redhat.io/ubi9/ubi-minimal
WORKDIR /opt/app-root

COPY --from=go-builder /opt/app-root/plugin-backend /opt/app-root/plugin-backend
COPY --from=web-builder /opt/app-root/web/dist /opt/app-root/web/dist

EXPOSE 8080

CMD ["./opt/app-root/plugin-backend"]