# ==== CONFIGURE =====
# Use a Node 16 base image
FROM node:16-alpine As builder
# Set the working directory to /app inside the container
WORKDIR /app
# Copy app files
COPY ./web .
# ==== BUILD =====
# Install dependencies (npm ci makes sure the exact versions in the lockfile gets installed)
RUN npm ci 
# Build the app
RUN npm run build

# Bundle the app with nginx
FROM nginx:1.23.4-alpine
# Set the env to "development"
ENV NODE_ENV development
# Copy the build output from the builder stage
COPY --from=builder /app/build /usr/share/nginx/html
# Copy the nginx config
COPY ./build/web/nginx.conf /etc/nginx/conf.d/default.conf
# Expose port 80
EXPOSE 80
# Start nginx
CMD ["nginx", "-g", "daemon off;"]