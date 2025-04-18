# Kashyap - Site

A simple, all-in-one personal portfolio website written in Go.

## Overview

Kashyap - Site provides a self-hosted personal portfolio with essential services built into a single deployable project. It is containerized using Docker and follows the DRAG stack philosophy.

## Features

- Email (SMTP)
- Admin server
- Portfolio server
- Chat with AI
- Website monitoring

## Tech Stack

This project uses the following technologies:

**DRAG Stack:**

- Docker / Docker Compose
- Redis
- Alpine.js
- Go

**Additional Tools:**

- Pico.css (optional, minimal CSS framework)

## Running the Project

To start the project, use the included `Makefile` targets. Docker and Docker Compose must be installed on your system.

### Start the project

```bash
make run
```

This will:

- Start all services using Docker Compose (`config/compose.yml`)
- Pull the `smollm2` model into the `ollama-ai` container

### Stop the project

```bash
make stop
```

This will:

- Shut down all containers
- Remove any orphaned containers

## Performance

Tested using [`wrk`](https://github.com/wg/wrk):

```bash
wrk -t20 -c1000 -d10m http://localhost:11000/
```

### Results

```
Running 10m test @ http://localhost:11000/
  20 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    77.47ms  101.70ms   1.43s    83.69%
    Req/Sec     1.60k   348.42     4.90k    69.62%
  19,144,904 requests in 10.00m, 382.21GB read

Requests/sec:  31,902.80
Transfer/sec:     652.19MB
```
