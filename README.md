# os_manage

为什么会有这个项目呢？因为我经常会有需要电脑挂着任务跑的情况，而我又希望在睡前屏幕能在我关灯后稍微照亮一下房间~~以便我躺在床上继续玩手机~~

当我真正准备睡觉时，只需要一个命令就能关闭电脑屏幕电源而不影响后台任务，~~台灯？送人了。~~ 于是就有了这个项目。

其实最开始做的是计划关机，屏幕电源虽然是后来者，但反而是用得最多的功能

## 功能

+ Windows计划关机
+ Windows屏幕电源控制
+ VirtualNetworkComputing(VNC)远程操作

## 基本介绍

### 计划关机

Windows提供的ui目前是没有定时关机的（至少据我所知是没有的）

大学的时候室友小志有天告诉我计划关机的命令，当时试了下感觉挺有意思的所以不自觉就记住了

项目的计划关机其实就是对该命令的封装调用，下面是命令基本使用

```shell
# Windows cmd
# 180秒后关机
shutdown -s -t 180

# 取消上面命令设置的定时关机任务（如果有的话）
shutdown -a

# 更多命令帮助
shutdown -help
```

### 屏幕电源控制

功能上类似控制显示器电源键，不影响电脑本身任务的运行，通过Windows提供的`user32.dll`向屏幕电源发送信号实现

但我自己的电脑在屏幕熄灭后机箱风扇会猛转1分钟左右，之前问过AI，可能是触发主板电源策略调整啥的，后续抽空再具体排查

### 远程桌面

桌面远程控制，非常轻量，Server端使用的是[UltraVNC](https://github.com/ultravnc/UltraVNC)

为了在手机上也能远程控制桌面，项目使用网页版来实现，具体是通过[noVNC](https://github.com/novnc/noVNC)提供的网页文件提供

在第一次使用VNC功能时，会检查上面的UltraVNC和noVNC，没有会自动下载（一共约20M）

> 注：使用本项目的VNC功能时**默认同意遵守上面的[UltraVNC](https://github.com/ultravnc/UltraVNC)和[noVNC](https://github.com/novnc/noVNC)相关协议**。

## 声明

> 本项目是我自己在局域网内部控制电脑使用，实际使用过程中注意联网控制，自行甄别不安全来源的控制请求，一切损失均与本项目无关。

另外，因为noVNC使用的是WebSocket协议，因此使用VNC功能时本项目也会充当WS-TCP协议的代理转发站，源代码来自
[novnc/websockify-other](https://github.com/novnc/websockify-other/blob/master/golang/websockify.go)

## 使用

### 直链下载

.exe可执行文件[前往下载](https://github.com/vantinfy/os_manage/releases)

### 自行编译

``` shell
git clone https://vantinfy/os_manage.git
cd os_manage
# 需要确保安装了go环境
./build.bat
```

## 计划

- [ ] **系统声音控制**
- [ ] 排查屏幕电源熄灭后机箱风扇猛转一分钟左右的问题
- [ ] 浏览远程文件 ~~我直接被窝里看电脑下载到本地的动漫！~~
- [ ] 应该不会做的跨操作系统
