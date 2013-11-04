# Enable Client Hints with User Agent sniffing

Nowadays many web applications use User Agent sniffing to serve
appropriate content based on device or browser capabilities. e.g:
big websites use the UA to decide whether to serve a mobile or a desktop
site.

It's common knowledge that doing User Agent sniffing is considered a bad
practice. However, whether we like it or not, all big companies rely on it to
perform content negotiation.

Ideally we would enrich the information browsers send to web servers so
that we don't need to rely on User Agent sniffing so much. HTTP Client
Hints proposal is a good example.

I thought it would be cool to develop a proof of concept where you could
start using Client Hints in your server even when the browser doesn't
provide them.

I created a simple proxy that would inject the Client Hints when these are
not provided by the browser by doing User Agent sniffing and gathering the
device information from the [DeviceAtlas][deviceatlas] API.

I think, as we get more proposals to add HTTP headers for content
negotiation this approach of enriching HTTP requests server side
would have different benefits:

 * Code for content negotiation/optimization can rely only on the new
   HTTP Headers
 * Old browsers that don't support the new header will benefit from
 * New browser versions supporting the headers would automatically
   benefit from the content negotation with no need to wait for the User
   Agent database to be updated.
 * Once a critical mass of browsers support such headers, header
   enrichment could be removed and there wouldn't be any changes.


## What are Client Hints?

Client Hints are a couple of new proposed HTTP headers to allow for
content negotiation. The headers in question are:

 * ```CH-DPR``` (Client Hint - Device Pixel Ratio): the ratio of real
   pixels for each device independent pixel. e.g: a retina iPhone has a
   DRP of 2.0 since it has 2 real pixels for every virtual pixel.
 * ```CH-RW``` (Client Hint - Resource Width): the width of the expected
   resource in device independent pixels.

You can find out more in [igrigorik/http-client-hints][ch_repo] and [read
the draft spec][ch_draft]

## Why is this a proxy?

This solution could be implemented at many different levels:
Apache/Nginx module, Rack middleware, server-side library, etc.

I decided to go with a proxy for two main reasons:

  * It doesn't depend on your stack. You can try it with any web
    server/language you like
  * I wanted to learn some Go :)

This is not a new approach. [Google's PageSpeed Service][pagespeed]
basically works by proxying your web server and rewrites your server's
responses to optimize your webages.

## What is DeviceAtlas?

[DeviceAtlas][deviceatlas] maintains a database with hundreds of user
agents.

They offer a paid API where you can send them a User Agent string and
they would return all the properties they know about the device:
manufacturer, screen size, pixel density, whether it's
mobile/desktop/tv/..., etc.

## Caveats

 * __This is not intended for production usage__. It's just a proof of
   concept.
 * This proxy doesn't do any content optimization, it just adds the
   Client Hints so that your server/app can perform the appropriate
   optimizations.
 * I did set the CH-RW as the device's width since that's the only
   relevant information i had.

## Installing and running the proxy

In order to install and run the proxy you will need:

 * [Go][golang]
 * A [DeviceAtlas Cloud Premium][deviceatlas] license key (you can get a 30 day
   trial key for free without entering your CC details).

Getting the proxy:

```bash
go get github.com/ernesto-jimenez/ch_deviceatlas
go install github.com/ernesto-jimenez/ch_deviceatlas
```

Running the server:

```bash
ch_deviceatlas --listen=[proxy port] --deviceatlas_key=[key] --proxy_to=[(host|ip):port]
```

## Deploying in heroku

```bash
git clone https://github.com/ernesto-jimenez/ch_deviceatlas.git
cd ch_deviceatlas
heroku create -b https://github.com/kr/heroku-buildpack-go.git
heroku config:set PROXY_TO=[(host|ip):port]
heroku config:set DEVICEATLAS_KEY=[key]
git push heroku master
```

## What about SSL?

You will notice that the proxy is HTTP only. We can only modify the
request if it's not encrypted and rather than adding SSL support to the
proxy I would recommend terminating SSL in the load balancer or in the
machine with [stunnel][stunnel].

## Questions/Comments?

If you have any questions or comments you drop me a mail/tweet! :)

 * [me@ernesto-jimenez.com](me@ernesto-jimenez.com)
 * [@ernesto_jimenez](http://twitter.com/ernesto_jimenez)

[ch_repo]: https://github.com/igrigorik/http-client-hints
[ch_draft]: https://github.com/igrigorik/http-client-hints/blob/master/draft.md
[pagespeed]: https://www.youtube.com/watch?v=FCyExI6Blfo
[deviceatlas]: https://deviceatlas.com
[golang]: http://golang.org
[stunnel]: https://www.stunnel.org

