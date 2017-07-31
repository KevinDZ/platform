IQIYI_LOG_PATH=/home/pengjianzhi/log/IqiyiLog
BAIDU_LOG_PATH=/home/pengjianzhi/log/BaiduLog
mkdir -p $IQIYI_LOG_PATH
mkdir -p $BAIDU_LOG_PATH
setsid ./iqiyiserver_dev -log_dir="$IQIYI_LOG_PATH" &
setsid ./baiduserver_dev -log_dir="$BAIDU_LOG_PATH" &

