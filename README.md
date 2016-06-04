##HTTP  proxy
##### this  Is actually a http  proxy 

####vanish cache clear agnet
please config the file servers.json  all data in one line ,end by '\n'
####build

install the go runtime

cd ./clearvanish
[root@07Node ~]#  go build
###run the agent
./start.sh

###测试
<pre><code>
[root@bogon 桌面]# curl -G -d "spid=10002001&epgid=100106"  -H"Host:interface5.v5.com" -H"Purge-ID:12" http://127.0.0.1:8000/



[root@bogon 桌面]# nc -l -k -p 8003
PURGE /?spid=10002001&epgid=100106 HTTP/1.1
Host: interface5.vooler.com
User-Agent: curl/7.29.0
Accept: */*
Purge-Id: 12
Accept-Encoding: gzip



</code></pre>
