# Stage 1: Build the frontend
FROM registry.access.redhat.com/ubi8/nodejs-16:latest AS web-builder
USER root

# Install yarn if not already installed
RUN command -v yarn || npm i -g yarn
WORKDIR /opt/app-root

COPY . .

# Install dependencies and build frontend
RUN yarn install && yarn build

# Stage 2: Build the Go backend
FROM golang:1.22 AS go-builder
WORKDIR /opt/app-root

# Copy the Go module files
COPY ./go.mod ./
COPY ./go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the rest of the backend files
COPY ./ ./

# Build the backend
RUN go build -o plugin-backend ./src/components/plugin-backend.go

# Stage 3: Combine frontend and backend
FROM registry.redhat.io/ubi9/ubi-minimal

COPY --from=go-builder /opt/app-root/plugin-backend /opt/app-root/plugin-backend
COPY --from=web-builder /opt/app-root/dist /opt/app-root/dist

EXPOSE 9443

CMD ["/opt/app-root/plugin-backend", "-static-path", "/opt/app-root/dist"]