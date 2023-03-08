%define _prefix    /opt/ast-app
%define _user      root
%define _user_uid  0
%define _group     root
%define _group_uid 0

Name:		ast-app-agent
Version:    0.1.0
Release: 	1
Summary:    ast-app-agent is a agent of ast-app, used to monitor linux machine、install iast/sca/rasp to java process and upload report
License: 	AGPL-3.0 license
URL: http://www.ast-app.com
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

Source0: 	ast-app-agent
Source1: 	ast-agent.jar
Source2: 	ast-http-client.jar
Source3: 	ast-iast-engine.jar
Source4: 	ast-rasp-engine.jar
Source5: 	ast-servlet.jar
Source6: 	ast-spy.jar
Source7: 	jattach-linux

%description
ast-app-agent is a agent of ast-app, used to monitor linux machine、install iast/sca/rasp to java process and upload report.

# 预备参数，通常为 %setup -q
%prep


# 编译参数 ./configure --user=nginx --group=nginx --prefix=/usr/local/nginx/……
%build
cp /tmp/ast-app-agent %{_sourcedir}
curl -sL -o "%{_sourcedir}/jattach-linux" "https://github.com/jattach/jattach/releases/download/$(curl -sL "https://api.github.com/repos/jattach/jattach/releases/latest" | grep -E 'tag_name\": \"' | head -n 1 | tr -d 'tag_name\": ' | tr -d ',')/jattach"
cp /tmp/%{Source0} %{_sourcedir}
cp /tmp/%{Source1} %{_sourcedir}
cp /tmp/%{Source2} %{_sourcedir}
cp /tmp/%{Source3} %{_sourcedir}
cp /tmp/%{Source4} %{_sourcedir}
cp /tmp/%{Source5} %{_sourcedir}
cp /tmp/%{Source6} %{_sourcedir}

# 安装步骤,此时需要指定安装路径，创建编译时自动生成目录，复制配置文件至所对应的目录中
%install
echo %buildroot
rm -rf %{buildroot}
%{__install} -p -D -m 0755 %{_sourcedir}/%{Source0} %{buildroot}/bin/ast-app-agent
%{__install} -p -D -m 0755 %{_sourcedir}/%{Source7} %{buildroot}/bin/jattach-linux
%{__install} -p -D %{_sourcedir}/%{Source1} %{buildroot}/libs/ast-agent.jar
%{__install} -p -D %{_sourcedir}/%{Source2} %{buildroot}/libs/ast-http-client.jar
%{__install} -p -D %{_sourcedir}/%{Source3} %{buildroot}/libs/ast-iast-engine.jar
%{__install} -p -D %{_sourcedir}/%{Source4} %{buildroot}/libs/ast-rasp-engine.jar
%{__install} -p -D %{_sourcedir}/%{Source5} %{buildroot}/libs/ast-servlet.jar
%{__install} -p -D %{_sourcedir}/%{Source6} %{buildroot}/libs/ast-spy.jar

# 安装前需要做的任务，如：创建用户
%pre
rm -rf %{_prefix}/*

# 安装后需要做的任务 如：自动启动的任务
%post

# 卸载前需要做的任务 如：停止任务
%preun
ps aux | grep -v grep | grep ast-app-agent | awk '{$2}' | xargs -I {} kill -9 {}

#
%postun
ps aux | grep -v grep | grep ast-app-agent-linux | xargs -I {} kill -9 {}
rm -rf /tmp/ast-app-agent.sock
rm -rf %{_prefix}

%clean
[ "$RPM_BUILD_ROOT" != "/" ] && rm -rf "$RPM_BUILD_ROOT"
rm -rf $RPM_BUILD_DIR/%{name}-%{version}

%files
%attr(-,_user,_group,-)
%dir %{_prefix}/
%dir %{_prefix}/bin
%dir %{_prefix}/libs
%dir %{_prefix}/logs
%{_prefix}/bin/ast-app-agent-linux
%{_prefix}/bin/jattach-linux
%{_prefix}/libs/ast-agent.jar
%{_prefix}/libs/ast-spy.jar
%{_prefix}/libs/ast-servlet.jar
%{_prefix}/libs/ast-http-client.jar
%{_prefix}/libs/iast-engine.jar
%{_prefix}/libs/rasp-engine.jar

%changelog
* Mon Mar 20 2023 owefsad <owefsad@gmail.com>
- 1.add ast-app-agent service
