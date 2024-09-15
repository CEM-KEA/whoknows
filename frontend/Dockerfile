# Stage 1: Build the Vite React app
FROM node:20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build-time argument
ARG VITE_API_URL

# Expose environment variable
ENV VITE_API_URL=${VITE_API_URL}

# Build the Vite React app
RUN npm run build

# Stage 2: Serve the built files using a lightweight web server
FROM nginx:alpine

# Copy the built files from the previous stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose port 80 to the outside world
EXPOSE 80

# Command to run the web server
CMD ["nginx", "-g", "daemon off;"]
