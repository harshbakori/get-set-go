client needs server that can process 10K requests per minit.
we can use echo server to server request assuming hardware is capable of getting request

steps 
- create a get endpoint to be used by client
- use log file for logging do not log in console only.
- create a swagger ui for ease of use.(Not nessesory but better to mantain)

When the endpoint is provided, the service should fire an HTTP POST request to the provided endpoint with count of unique requests in the current minute as a query parameter. Also log the HTTP status code of the response. 

- create logic to send post request 
- handle empty endpoint and not reachable endpoint.
- while it just post request to endpoint it will be a seprate thread running asyncronously and will not affect output of get request endpoint.

Extension 1:
while any database will be suficiant to store and share unique request on multipal instance running we will use redis for if i spinup 100 of instances redis will be able to handle it while any db we migth need to do some adjustment.

Extension 2:
to share unique count on any distributed service i'm using kafka just becauase it can handle high throughput. it will be helpful while there will be 100s of instances.

NOTE: running kafka on docker is hard so i will be using redpanda usign docker compose to mimic kafka. it is just easear to setup for current synario.

what a mess i should cleanup some code and restructure some config.

some screenshots attached 

multi server setup:
![alt text](<Screenshot 2025-02-07 004242.png>)

insted of kafka creating topic in redpanda is way easy
![alt text](<Screenshot 2025-02-07 011043.png>)

got my first msg on kafka bus
![alt text](<Screenshot 2025-02-07 011453.png>)

