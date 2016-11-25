Reverse proxy for opening your local development to a public port on a server, setting a custom Host header and optionally accepting an insecure local cert. This solution requires you already have a publicly available server with SSH access.

## Another one?!

- [ngrok](ngrok.com) is fantastic, but charges for SSL forwarding.
- [localtunnel](https://localtunnel.github.io/) is also great, but I want to customise the Host header without changing the actual host being accessed.

When developing sites, I often configure my web server and/or application to redirect non-https requests to https. This will cause a lot of proxies to not work. 

## Usage

`proxy --help` will show the options available to you:

```
Usage of ./proxy:
  -hostname string
    	Value of the Host http header for which your vhost is configured (default "localhost")
  -insecure
    	If your local server uses a self-signed cert, set this to true
  -internal-port string
    	Internal port on which your reverse SSH tunnel is listening (default "9998")
  -public-port string
    	Public port on which to serve proxied content (default "9999")
  -scheme string
    	Use http or https locally (default "http")
```

1. Copy the proxy binary to your server, or somehow make it available on the server.
2. SSH into your server, providing a remote tunnel pointing the internal port above to the local port you wish to expose.
   `ssh -R 9998:localhost:443 myuser@myserver proxy -hostname mylocalsite.dev -insecure -scheme https`

The above command will expose the port 9999 on the server, and proxy it it through to your local machine, configuring the Host header and requesting https. Which means a request to http://myserver:9999 will return your local server as if you'd requested https://mylocalsite.dev

## Build

With docker:

```
docker run --rm -v "$PWD":/usr/src/proxy -w /usr/src/proxy golang:latest go build -v
```
