# Configure Wrk https://github.com/wg/wrk (brew install wrk) with Lua to execute against multiple URLs

## Dockerfile

```dockerfile
FROM debian:jessie
MAINTAINER Michał Czeraszkiewicz <czerasz.hosting@gmail.com>

# Set the reset cache variable
ENV REFRESHED_AT 2015-06-15

ENV DEBIAN_FRONTEND noninteractive

# Update packages list
RUN apt-get update -y

# Install useful packages
# RUN apt-get install -y strace procps tree vim git curl wget gnuplot

# Install required software
RUN apt-get install -y git make build-essential libssl-dev

# Install wrk - benchmarking software
# Resource: https://github.com/wg/wrk/wiki/Installing-Wrk-on-Linux
WORKDIR /tmp

RUN git clone https://github.com/wg/wrk.git &&\
    cd wrk &&\
    make &&\
    mv wrk /usr/local/bin

# Install Luarocks dependencies
RUN apt-get install -y curl \
                       make \
                       unzip \
                       lua5.1 \
                       liblua5.1-dev

# Install Luarocks - a lua package manager
RUN curl http://keplerproject.github.io/luarocks/releases/luarocks-2.2.2.tar.gz -O &&\
    tar -xzvf luarocks-2.2.2.tar.gz &&\
    cd luarocks-2.2.2 &&\
    ./configure &&\
    make build &&\
    make install

# Install the cjson package
RUN luarocks install lua-cjson

WORKDIR /

# Cleanup
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Raise the limits to successfully run benchmarks
RUN ulimit -c -m -s -t unlimited

ENV DEBIAN_FRONTEND=newt
```

## multi-request-json.lua

```lua
-- http://czerasz.com/2015/07/19/wrk-http-benchmarking-tool-example/
-- https://github.com/czerasz/docker-wrk-json
-- docker run --rm -v "$(pwd)/multi-request-json.lua":/multi-request-json.lua -v "$(pwd)/requests.json":/requests.json czerasz/wrk-json wrk -c1 -t1 -d5s -s /multi-request-json.lua https://www.example.com

-- TRIED LOCAL SETUP ON MAC BUT MODULE STILL COULDN'T BE FOUND...
-- sudo luarocks install lua-cjson
-- package.path = package.path .. ';/usr/local/lib/luarocks/rocks-5.2/?/?.lua'
-- wrk -c1 -t1 -d5s -s multi-request-json.lua https://www.example.com

-- Module instantiation
local cjson = require "cjson"
local cjson2 = cjson.new()
local cjson_safe = require "cjson.safe"

-- Initialize the pseudo random number generator
-- Resource: http://lua-users.org/wiki/MathLibraryTutorial
math.randomseed(os.time())
math.random(); math.random(); math.random()

-- Shuffle array
-- Returns a randomly shuffled array
function shuffle(paths)
  local j, k
  local n = #paths

  for i = 1, n do
    j, k = math.random(n), math.random(n)
    paths[j], paths[k] = paths[k], paths[j]
  end

  return paths
end

-- Load URL paths from the file
function load_request_objects_from_file(file)
  local data = {}
  local content

  -- Check if the file exists
  -- Resource: http://stackoverflow.com/a/4991602/325852
  local f=io.open(file,"r")
  if f~=nil then
    content = f:read("*all")

    io.close(f)
  else
    -- Return the empty array
    return lines
  end

  -- Translate Lua value to/from JSON
  data = cjson.decode(content)

  return shuffle(data)
end

-- Load URL requests from file
requests = load_request_objects_from_file("/requests.json")

-- Check if at least one path was found in the file
if #requests <= 0 then
  print("multiplerequests: No requests found.")
  os.exit()
end

print("multiplerequests: Found " .. #requests .. " requests")

-- Initialize the requests array iterator
counter = 1

request = function()
  -- Get the next requests array element
  local request_object = requests[counter]

  -- Increment the counter
  counter = counter + 1

  -- If the counter is longer than the requests array length then reset it
  if counter > #requests then
    counter = 1
  end

  -- Return the request object with the current URL path
  return wrk.format(request_object.method, request_object.path, request_object.headers, request_object.body)
end
```

## requests.json

```json
[
  {
    "path": "/bar",
    "body": "some content",
    "method": "GET",
    "headers": {
      "X-Custom-Header-1": "test 1",
      "X-Custom-Header-2": "test 2"
    }
  },
  {
    "path": "/baz",
    "body": "some content",
    "method": "GET",
    "headers": {
      "X-Custom-Header-1": "test 3",
      "X-Custom-Header-2": "test 4"
    }
  }
]
```

## 1. Usage.sh

```shell
$ docker run --rm -v "$(pwd)/multi-request-json.lua":/multi-request-json.lua -v "$(pwd)/requests.json":/requests.json czerasz/wrk-json wrk -c1 -t1 -d5s -s /multi-request-json.lua https://www.example.com

multiplerequests: Found 2 requests
multiplerequests: Found 2 requests

Running 5s test @ https://www.example.com

  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   887.09ms  415.48ms   1.36s    60.00%
    Req/Sec     0.60      0.55     1.00     60.00%
  5 requests in 5.10s, 1.54MB read

Requests/sec:      0.98
Transfer/sec:    309.12KB
```

