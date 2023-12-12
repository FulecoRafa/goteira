# Goteira

Simple cli tool to open multiple ping servers.

One good usage is to check a service that may call other services to check if they are hitting their targets.

## Usage

The application is pretty simple, just use the `-p` flag to list the ports to be listening

### Example

```bash
goteira -p 6969 -p 7979
```

```
Serving the following ping servers:
http://localhost:6969
http://localhost:7979
```
