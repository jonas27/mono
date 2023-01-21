# goserver
This is a go server for serving static files.  
It is super simple to use which comes with a few downsides.

## Args
`-port`: The port on which to serve (string, ex: *":8080"*).  
`-staticDir`: The root directory of the static content (string, ex *"/static"*).  
`-metrics`: If server should serve metrics (bool, ex: *true*). 


### Prometheus Metrics
`Promotheus style metrics for http requests can be 
When ever the service is called it increases a promotheus gauge by one.