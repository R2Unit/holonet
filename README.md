# Holonet
> [!WARNING]
> Holonet is still under development so it is also still in the alpha phase! So do not use it for production purposes.

Holonet leverages Netbox to ensure your network and compute configurations always reflect your definitive single source of truth. 
This integration automates consistency, making your infrastructure reliable, auditable, and ready for advanced automation.
## Documentation

## Usage

### Logging Configuration

Holonet supports different logging levels to control the verbosity of log output. The available log levels are:

- **DEBUG**: Detailed troubleshooting information
- **INFO**: General operational information (default)
- **WARN**: Warning messages
- **ERROR**: Error messages
- **FATAL**: Critical errors that cause the program to exit

You can set the log level using the `-log-level` command-line flag:

```bash
# Run with default INFO level
./main

# Run with DEBUG level for more verbose output
./main -log-level=debug

# Run with ERROR level to show only errors and fatal messages
./main -log-level=error
```

When running in a container, you can pass the log level as an environment variable:

```bash
# Using Docker
docker run -e LOG_LEVEL=debug holonet/core

# Using Docker Compose
# In your docker-compose.yml file:
# services:
#   holonet:
#     image: holonet/core
#     environment:
#       - LOG_LEVEL=debug
```

### Step 1:

### Step 2: 
