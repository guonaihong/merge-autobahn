### 一、安装

```bash
go install github.com/guonaihong/merge-autobahn/cmd/merge-autobahn@latest
```

### 二、例子

```bash
merge-autobahn -f ~/reports/servers -f ~/reports/servers2 -o ./output
// 输出 ./output/merge_index.html
// 查看 open ./output/merge_index.html
```

### 三、merge-autobahn  usage

-f, --from: 可以指定多个输入目录

```console
Usage:
    merge-autobahn [Options] 

Options:
    -f,--from      input directory
    -h,--help      print the help information
    -o,--output    output directory

```
