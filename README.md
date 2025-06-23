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

### NetBox Integration Configuration

Holonet integrates with NetBox to ensure your network and compute configurations always reflect your definitive single source of truth. To connect to your NetBox instance, you need to configure the following environment variables:

- **NETBOX_HOST**: The URL of your NetBox instance (e.g., `http://netbox.example.com:8000`)
- **NETBOX_API_TOKEN**: Your NetBox API token for authentication

These environment variables are required for the NetBox integration to work. If they are not set, the NetBox integration will be disabled, but Holonet will continue to function with limited capabilities.

You can set these environment variables in various ways:

```bash
# Setting environment variables directly
export NETBOX_HOST="http://netbox.example.com:8000"
export NETBOX_API_TOKEN="your_api_token_here"
./main

# Or in one line
NETBOX_HOST="http://netbox.example.com:8000" NETBOX_API_TOKEN="your_api_token_here" ./main
```

When running in a container, you can pass these variables as environment variables:

```bash
# Using Docker
docker run -e NETBOX_HOST="http://netbox.example.com:8000" -e NETBOX_API_TOKEN="your_api_token_here" holonet/core

# Using Docker Compose
# In your docker-compose.yml file:
# services:
#   holonet:
#     image: holonet/core
#     environment:
#       - NETBOX_HOST=http://netbox.example.com:8000
#       - NETBOX_API_TOKEN=your_api_token_here
```

To obtain a NetBox API token:
1. Log in to your NetBox instance
2. Go to your user profile
3. Navigate to the "API Tokens" tab
4. Create a new token with appropriate permissions

### Step 1:

### Step 2: 
