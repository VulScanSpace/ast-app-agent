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
Source7: 	jattach

%description
ast-app-agent is a agent of ast-app, used to monitor linux machine、install iast/sca/rasp to java process and upload report.

# 预备参数，通常为 %setup -q
%prep

# 编译参数 ./configure --user=nginx --group=nginx --prefix=/usr/local/nginx/……
%build
ls -l /home/runner/rpmbuild
cp /tmp/ast-app-agent %{_sourcedir}/
cp /tmp/jattach %{_sourcedir}/
cp /tmp/ast-agent.jar %{_sourcedir}/
cp /tmp/ast-http-client.jar %{_sourcedir}/
cp /tmp/ast-iast-engine.jar %{_sourcedir}/
cp /tmp/ast-rasp-engine.jar %{_sourcedir}/
cp /tmp/ast-servlet.jar %{_sourcedir}/
cp /tmp/ast-spy.jar %{_sourcedir}/
ls -l /home/runner/rpmbuild
ls -l %{_sourcedir}

# 安装步骤,此时需要指定安装路径，创建编译时自动生成目录，复制配置文件至所对应的目录中
%install
stat %{_sourcedir}
ls -l %{_sourcedir}
%{__install} -p -D -m 0755 %{_sourcedir}/%{source0} %{buildroot}%{_prefix}/bin/ast-app-agent
%{__install} -p -D -m 0755 %{_sourcedir}/jattach %{buildroot}%{_prefix}/bin/jattach
%{__install} -p -D %{_sourcedir}/ast-agent.jar %{buildroot}%{_prefix}/libs/ast-agent.jar
%{__install} -p -D %{_sourcedir}/ast-http-client.jar %{buildroot}%{_prefix}/libs/ast-http-client.jar
%{__install} -p -D %{_sourcedir}/ast-iast-engine.jar %{buildroot}%{_prefix}/libs/ast-iast-engine.jar
%{__install} -p -D %{_sourcedir}/ast-rasp-engine.jar %{buildroot}%{_prefix}/libs/ast-rasp-engine.jar
%{__install} -p -D %{_sourcedir}/ast-servlet.jar %{buildroot}%{_prefix}/libs/ast-servlet.jar
%{__install} -p -D %{_sourcedir}/ast-spy.jar %{buildroot}%{_prefix}/libs/ast-spy.jar

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
ps aux | grep -v grep | grep ast-app-agent | xargs -I {} kill -9 {}
rm -rf /tmp/ast-app-agent.sock
rm -rf %{_prefix}

%clean
[ "$RPM_BUILD_ROOT" != "/" ] && rm -rf "$RPM_BUILD_ROOT"
rm -rf $RPM_BUILD_DIR/%{name}-%{version}

%files
%dir %{_prefix}/
%dir %{_prefix}/bin
%dir %{_prefix}/libs
%{_prefix}/bin/ast-app-agent
%{_prefix}/bin/jattach
%{_prefix}/libs/ast-agent.jar
%{_prefix}/libs/ast-spy.jar
%{_prefix}/libs/ast-servlet.jar
%{_prefix}/libs/ast-http-client.jar
%{_prefix}/libs/ast-iast-engine.jar
%{_prefix}/libs/ast-rasp-engine.jar

%changelog
* Mon Mar 20 2023 owefsad <owefsad@gmail.com>
- 1.add ast-app-agent service
