# envd 

envd is an daemon for starting and stopping processess using configuration from 
[etcd](https://github.com/coreos/etcd).  You can use it to either run once and 
exit or to watch a key field for updates and restart your application.

The goal of this project is to provide highly available configuration data
independent of application deployment.


## Configuration

```yaml
server: 
    - https://localhost:4001
key: /etc/envd/myclient.key 
cert: /etc/envd/myclient.crt 
cacert: /etc/envd/clientCA.crt
```


## Usage

```bash

$ envd -e my-environment-variables /start/my/app

```
This will run once and exit when the application exits.


```bash

$ envd -e my-environment-variables -d /start/my/app


```
This will start the app and daemonize in the background.


```bash

$ envd -e my-environment-variables -w /start/my/app

```
This will start the app, watch the variables for updates and restart the app
if any variables are changed in real time.


```bash

$ envd -e my-environment-variables -w -d /start/my/app

```
This will start the app, daemonize and restart the app if any variables are
changed.


## Author

Steve Rude <srude@riotgames.com>

