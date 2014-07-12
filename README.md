yxorp
=====

A tiny reverse-proxy

*See what I did there?*

Install with `go get github.com/gomicroprojects/yxorp`

## The idea

Suppose you have now finished a couple of Go projects and you want to deploy them on a server. But you only have one IP address available and want to serve multiple projects under different domains on port 80.

yxorp to the rescue.

It:

* is configurable with a simple JSON file
* will act as a normal HTTP reverse proxy
* bonus: will optionally GZip encode the response

## Starting

Start with [yxorp.go](yxorp.go)

## www.gomicroprojects.com
