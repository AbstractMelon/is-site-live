# Build stage
FROM node:18-alpine as build-stage

WORKDIR /app

# Install pnpm globally
RUN npm install -g pnpm

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN pnpm i

# Copy all files
COPY . .

# Build the app
RUN pnpm run build

# Production stage
FROM nginx:stable-alpine as production-stage

# Copy built files from build stage
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
