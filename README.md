<p align="center">


<h1 align="center">go-erltimer</h1>
<p align="center">
    <em> The adorable timer library for Golang</em>
</p>
<p align="center">
    <a href="https://godoc.org/github.com/botanyzh/go-erltimer">
        <img src="https://godoc.org/github.com/botanyzh/go-erltimer?status.svg" alt="GoDoc">
    </a>
    <a href="https://travis-ci.org/botanyzh/go-erltimer">
        <img src="https://travis-ci.org/botanyzh/go-erltimer.svg?branch=master" alt="Build Status">
    </a>
    <a href="https://ci.appveyor.com/project/botanyzh/go-erltimer">
        <img src="https://ci.appveyor.com/api/projects/status/kdxi0s1nc1t6dqn0?svg=true" alt="Build Status">
    </a>
    <a href="https://goreportcard.com/report/github.com/botanyzh/go-erltimer">
        <img src="https://goreportcard.com/badge/github.com/botanyzh/go-erltimer" alt="Go Report Card">
    </a>
    <a href="https://opensource.org/licenses/MIT">
        <img src="https://img.shields.io/badge/License-MIT-brightgreen.svg" alt="MIT License">
    </a>
</p>

# erltimer
## 注意 warning:
 **实验阶段项目尚存bug只阐述思想并非实际可用的代码，如感兴趣可与作者联系获取可生产级别代码 （the code is experimental not usable）**

## 特性 features:

**1.并发安全  concurrent safe**

**2.大规模可伸缩  massive scalable**

**3.linux 下 高精度 high precision for linux  but do not use this in windows for milli precision (windows下受原生api的系统限制无法提供高精度只能精确到毫秒)**

**4.小根堆保证高并发下的定时队列读写效率**

**5. 底层基于erlang gen_server 式的actor设计实现**
## 用例 examples:

请看 please see the  ./timer/timer_test.go
 

### 📃 LICENSE

MIT [©botanyzh](https://github.com/botanyzh)
