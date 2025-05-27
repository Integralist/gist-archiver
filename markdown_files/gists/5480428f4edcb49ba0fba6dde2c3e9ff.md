# Local Memcache and ElastiCache

[View original Gist on GitHub](https://gist.github.com/Integralist/5480428f4edcb49ba0fba6dde2c3e9ff)

## Local Memcache and ElastiCache.md

## Single Node

- `brew install memcached` (or `docker run -d -p 11211:11211 memcached`)
- `memcached` (optional `-d` to background & `-p` to change port)
- `telnet localhost 11211`
- `stats`
- `quit`

> http://blog.elijaa.org/2010/05/21/memcached-telnet-command-summary

Notice the line break for the value (when using `set` or `add` etc) is required...

```
set foo 0 60 3
bar
STORED

get foo
VALUE foo 0 3
bar
END
```

## AWS ElastiCache Cluster Endpoint

- `gem install fake_elasticache`
- `fake_elasticache` (run in a separate shell as it's run in the foreground)
- `telnet localhost 11212` (fyi that's a non-standard port)
- `config get cluster`
- `quit`

## Go Client

Here's a Go client that interacts with the AWS ElastiCache endpoint and parses out the data.

To run this script you'll need `memcached` running locally along with `fake_elasticache` which looks up the locally running memcache at the default location of `localhost` and port `11211`.

```go
package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
)

// Node is a single ElastiCache instance node
type Node struct {
	URL  string
	Host string
	IP   string
	Port int
}

func main() {
	var response string
	var nodes []Node
	var urls []string

	conn, err := net.Dial("tcp", "127.0.0.1:11212")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer conn.Close()

	command := "config get cluster\r\n"

	fmt.Fprintf(conn, command)

	count := 0
	location := 3 // AWS docs suggest that nodes will always be listed on line 3

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		count++
		if count == location {
			response = scanner.Text()
		}
		if scanner.Text() == "END" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	items := strings.Split(response, " ")

	for _, v := range items {
		fields := strings.Split(v, "|") // ["host", "ip", "port"]

		port, err := strconv.Atoi(fields[2])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		node := Node{fmt.Sprintf("%s:%d", fields[1], port), fields[0], fields[1], port}
		nodes = append(nodes, node)
		urls = append(urls, node.URL)

		fmt.Printf("Host: %s\n", node.Host)
		fmt.Printf("IP: %s\n", node.IP)
		fmt.Printf("Port: %d\n", node.Port)
		fmt.Printf("URL: %s\n\n", node.URL)
	}

	mc := memcache.New(urls...)
	mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})

	it, err := mc.Get("foo")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%+v", it)
}
```

