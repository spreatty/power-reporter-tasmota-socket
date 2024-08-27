# Power Status app for Tasmota socket
An app that reports power status of a Tasmota socket(s) to specific URL. The app expects a socket to call an endpoint to indicate it went online. The app then pings the socket until the socket becomes unavailable.

## Tasmota rule for notifying that the socket went online
```
Rule1 ON Wifi#Connected DO WebQuery http://10.0.0.2:8085 POST ENDON
```
And don't forget to enable the rule
```
Rule1 1
```
