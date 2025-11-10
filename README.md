执行流程 

安装打包工具
在终端中执行以下命令，安装Fyne的打包命令行工具：
go install fyne.io/fyne/v2/cmd/fyne@latest

应用应用图标
fyne package -icon icon.png

命令说明：
fyne package：核心打包命令。
-icon icon.png：指定应用图标。如果省略，会使用默认图标。


使用fyne-cross（为所有平台打包）

安装fyne-cross
在终端中执行以下命令，安装fyne-cross：
go install github.com/fyne-io/fyne-cross@latest

为所有平台打包
在项目根目录下执行以下命令，为所有平台打包：
fyne-cross windows -icon icon.png
fyne-cross darwin -icon icon.png
fyne-cross linux -icon icon.png
fyne-cross android

这会在当前目录下生成对应的可执行文件。

打包移动端app
生成一个调试密钥库（一次性的）  (可选)
keytool -genkey -v -keystore my-release-key.keystore -alias alias_name -keyalg RSA -keysize 2048 -validity 10000

它会问你一堆问题（姓名、组织等），你可以全部按回车跳过，或者随便填写。

最后会让你设置一个密码，请务必记住它。

运行 fyne-cross android。 (需要安装docker) 密钥会自动处理
把 fyne-cross/dist/android/ 目录下的APK文件发到手机安装。
完成！
