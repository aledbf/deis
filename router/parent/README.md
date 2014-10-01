[**NAXSI**](https://github.com/nbs-system/naxsi) is an open-source, high performance, low rules maintenance WAF for NGINX.

Why a firewall for nginx in deis?
Well in the last weeks [Shellshock](https://shellshocker.net) showed a vulnerability that some apps (mostly CGI's) inside a web server can be exploited like is explained here [Inside Shellshock: How hackers are using it to exploit systems](https://blog.cloudflare.com/inside-shellshock)

To reduce the contact surface of this attack and others (like sql injection and cross site scripting) this module is enabled by default.

Example:
```console
TODO
```

The rules included are taken from this project [doxi-rules](https://bitbucket.org/lazy_dogtown/doxi-rules)

Only this modules are enabled:


|---|---|
|---|---|
|web_app.rules       |detect exploit/misuse-attempts againts web-applications
|web_server.rules    |generic rules to protect a webserver from misconfiguration and known mistakes / exploit-vectors 
|active-mode.rules   |rules to configure active-mode (block)
|naxsi_core          |core naxsi rules
