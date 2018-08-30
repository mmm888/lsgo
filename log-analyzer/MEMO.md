# Memo

* task

https://github.com/pfnet/intern-coding-tasks/tree/master/2018/be

* Download log file

```bash
mkdir log
cd log
wget https://preferredjp.box.com/shared/static/mc01xtkcn0qmdgbb1r53uzljeva6e9tq.zip -O log_s.zip
wget https://preferredjp.box.com/shared/static/gpmix7flrrl4badutdqwo4hr9q2xnasp.zip -O log_l.zip
unzip log_s.zip
unzip log_l.zip
```

* Insert data

```bash
make build
./bin/log-analyzer clean
LOGDIR="./log/log_s"; for i in $(ls $LOGDIR); do for j in $(ls $LOGDIR/$i); do ./bin/log-analyzer add -f $LOGDIR/$i/$j; done; done
```
