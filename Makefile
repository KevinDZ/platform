a:ysdkserver_dev
all:accountserver loginserver payserver backendserver ocpcserver ysdkserver
all_dev:accountserver_dev loginserver_dev payserver_dev backendserver_dev ocpcserver_dev ysdkserver_dev
#ysdkServer
run_dev:
	mkdir -p AccountLog
	setsid ./accountserver_dev -log_dir="./AccountLog"
	mkdir -p PayLog
	setsid ./payserver_dev -log_dir="./PayLog"
	mkdir -p LoginLog
	setsid ./loginserver_dev -log_dir="./LoginLog"

run:
	mkdir -p AccountLog
	setsid ./accountserver -log_dir="./AccountLog"
	mkdir -p PayLog
	setsid ./payserver -log_dir="./PayLog"
	mkdir -p LoginLog
	setsid ./loginserver -log_dir="./LoginLog"



clean:
	rm -fr accountserver cache loginserver payserver accountserver_dev loginserver_dev payserver_dev backendserver_dev backendserver ysdkserver ysdkserver_dev ocpcserver_dev ocpcserver
third_dev: ThirdServer.go
	go build -o $@ $^
accountserver:AccountServer.go Protocol.go AccountIDTable.go ReleaseConst.go GameIDTable.go Rand.go DeviceInfo.go Logger.go CommonRequestProcess.go
	go build -o $@ $^
accountserver_dev:AccountServer.go Protocol.go AccountIDTable.go DebugConst.go GameIDTable.go Rand.go DeviceInfo.go Logger.go CommonRequestProcess.go
	go build -o $@ $^


loginserver:LoginServer.go LoginGameTable.go Protocol.go ReleaseConst.go DeviceInfo.go GameIDTable.go Rand.go Logger.go CommonRequestProcess.go FirstAccessTable.go
	go build -o $@ $^
loginserver_dev:LoginServer.go LoginGameTable.go Protocol.go DebugConst.go DeviceInfo.go GameIDTable.go Rand.go Logger.go CommonRequestProcess.go FirstAccessTable.go
	go build -o $@ $^

payserver:PayServer.go PayCpNotify.go Protocol.go ShenfutongWeixin.go ShenfutongAlipay.go PayOrderTable.go ShenfutongNotify.go  ReleaseConst.go DeviceInfo.go GameIDTable.go Rand.go Logger.go CommonRequestProcess.go
	go build -o $@ $^
payserver_dev:PayServer.go PayCpNotify.go Protocol.go ShenfutongWeixin.go ShenfutongAlipay.go PayOrderTable.go ShenfutongNotify.go  DebugConst.go DeviceInfo.go GameIDTable.go Rand.go Logger.go CommonRequestProcess.go
	go build -o $@ $^


ysdkserver: Logger.go ysdkServer.go PayCpNotify.go PayOrderTable.go GameIDTable.go Protocol.go  Rand.go ReleaseConst.go CommonRequestProcess.go DeviceInfo.go
	go build -o $@ $^
ysdkserver_dev: Logger.go ysdkServer.go PayCpNotify.go PayOrderTable.go GameIDTable.go Protocol.go  Rand.go DebugConst.go CommonRequestProcess.go DeviceInfo.go
	go build -o $@ $^


backendserver: backendServer.go getDate.go ReleaseConst.go
	go build -o $@ $^
backendserver_dev:backendServer.go getDate.go DebugConst.go
	go build -o $@ $^

ocpcserver_dev: OcpcServer.go OcpcTable.go DebugConst.go Logger.go
	go build -o $@ $^

ocpcserver: OcpcServer.go OcpcTable.go ReleaseConst.go Logger.go
	go build -o $@ $^






account:
	go build -o $@ AccountIDTable.go ReleaseConst.go

account-dev:
	go build -o $@ AccountIDTable.go TestReleaseConst.go

logingame:
	go build -o $@ LoginGameTable.go ReleaseConst.go

logingame-dev:
	go build -o $@ LoginGameTable.go TestDatabaseConst.go

payorder:
	go build -o $@ PayOrderTable.go ReleaseConst.go

payorder_dev: PayOrderTable.go DebugConst.go
	go build -o $@ $^

getdate:getDate.go ReleaseConst.go
	go build -o $@  $^


getdate_dev:getDate.go DebugConst.go
	go build -o $@  $^


firstaccess:FirstAccessTable.go DebugConst.go
	go build -o $@ $^



cache:
	go build CacheServer.go 
runcache:
	./CacheServer &
