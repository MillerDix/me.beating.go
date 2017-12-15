# 解决Nginx下使用React-router出现404问题

## 场景描述

用React + React-router做SPA项目，路由模式为BrowserHistory，在Nginx下访问项目，默认地址为`zlzkj.io`，点击路由`zlzkj.io/goods`，可以正常切换页面，一旦刷新页面就会报404错误。

## 问题分析

这是因为Nginx访问`zlzkj.io/goods`会去找`goods.html`，实际上我们是没有这个文件的，所有内容都是通过路由去渲染React组件，自然会报404错误。

## 解决方法

通过配置Nginx，访问任何URI都指向index.html，浏览器上的path，会自动被React-router处理，进行无刷新跳转。
配置文件参考：

```
server {
   listen 80;
   server_name zlzkj.io;
   index  index.html;
   root /Volumes/Mac/www/antd-admin/;
   location / {
       try_files $uri $uri/ /index.html;
   }
}
```

