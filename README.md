# hfs
a Http File Server based on [Iris](https://github.com/kataras/iris)

<img src="https://github.com/tuhao1020/filerepo/blob/master/images/hfs.png">

## Usage:
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
-tmplName   // custom template file name (should be in ./template directory)
-v          // show version
```

## Custom favicon
place your custom `favicon.ico` into the same directory of `hfs` executable file.

## Custom list template
1. copy the `template` directory of the source code and place it into the same directory of `hfs` executable file.
2. edit the `dirlist.tmpl` to change the view.




