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

%prep
rm -rf %{_prefix}/*

%install
rm -rf %{buildroot}
%{__install} -p -D -m 0755 %{Source0} %{_prefix}/bin/ast-app-agent
%{__install} -p -D -m 0755 %{Source7} %{_prefix}/bin/jattach-linux
%{__install} -p -D %{Source1} %{_prefix}/libs/ast-agent.jar
%{__install} -p -D %{Source2} %{_prefix}/bin/ast-http-client.jar
%{__install} -p -D %{Source3} %{_prefix}/bin/ast-iast-engine.jar
%{__install} -p -D %{Source4} %{_prefix}/bin/ast-rasp-engine.jar
%{__install} -p -D %{Source5} %{_prefix}/bin/ast-servlet.jar
%{__install} -p -D %{Source6} %{_prefix}/bin/ast-spy.jar

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
