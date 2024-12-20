# Use the official Golang image as a parent image for the build stage
FROM golang:1.23.1-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy only necessary files to leverage Docker cache more effectively
COPY src/go.mod src/go.sum ./

# Download all dependencies first, leveraging Docker caching
RUN go mod download

# Copy the rest of the source code into the container
COPY src/ /app/src

# Set the working directory to where main.go is located
WORKDIR /app/src

# Build the application (output as "main" instead of "main.go")
RUN go build -o /app/main .

# Expose the port that the application listens on
EXPOSE 8443

# Use the chromedp/headless-shell image for the runtime environment
FROM chromedp/headless-shell:latest

# Install CA certificates and dumb-init
RUN apt-get update && apt-get install -y \
    ca-certificates \
    dumb-init

# Set entrypoint to use dumb-init for signal handling
ENTRYPOINT ["dumb-init", "--"]

# Set the working directory for the final runtime image
WORKDIR /app

# Copy the compiled Go binary AND the required text file from the build stage
COPY --from=build /app/main /app/
COPY --from=build /app/src/services/spot_symbols_sorted.txt /app/services/
COPY --from=build /app/src/services/futures_symbols_sorted.txt /app/services/
COPY --from=build /app/src/services/spot_symbols.txt /app/services/
COPY --from=build /app/src/services/futures_symbols.txt /app/services/

# Command to run the application
CMD ["./main"]

# Expose the port the application listens on
EXPOSE 8443
