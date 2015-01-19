# fwatch

A small tool to watch a file for inactivity. I needed this  because we had 
a problem in our production environment. Our services (a tomcat app) hosts
a webservice which is called very often (multiple times a second). But 
sometimes the whole application stuck and our customers needed more information
about this situations.

So i wrote this tool to identify such situations: `fwatch` watches the
server log file for inactivity of a given time duration and if this situation
happened it mailed me a notification and wrote a heapdump for the current
application server.

You simply use 

```
fwatch <serverlog-file> 10sec 0 mail -s "Server hangs" service@example.com
```

and you will receive a mail notification if you server log is not changed for
at least 10 seconds.