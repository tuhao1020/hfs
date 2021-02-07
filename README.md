# hfs
a Http File Server based on [Iris](https://github.com/kataras/iris)

<img src="https://github.com/tuhao1020/filerepo/blob/master/images/hfs.png">

Usage:
```
./hfs
```

this will start a http file server on 8090 for current directory.

command line parameters
```
-dir        // specify root file directory, default: ./
-port       // specify http file server port, default: 8090
-showHidden // show hidden files or directories, default: false
-tls        // enable tls， default: false
-certFile   // cert file path
-keyFile    // key file path
-v          // show version
```
