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
ENV VITE_API_URL=http://52.169.32.176:8080/api

# Build the Vite React app
RUN npm run build

# Stage 2: Serve the built files using a lightweight web server
FROM nginx:alpine

# Copy the built files from the previous stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Remove the default Nginx configuration file
RUN rm /etc/nginx/conf.d/default.conf

# Copy the custom Nginx configuration file
COPY nginx.conf /etc/nginx/conf.d

# Expose port 80 to the outside world
EXPOSE 80

# Command to run the web server
CMD ["nginx", "-g", "daemon off;"]
