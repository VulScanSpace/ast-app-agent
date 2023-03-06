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
Source0: 	%{name}-%{version}
URL: http://www.ast-app.com
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
ast-app-agent is a agent of ast-app, used to monitor linux machine、install iast/sca/rasp to java process and upload report.

%prep
rm -rf %{_prefix}/*

%install
%{_prefix}/bin/ast-app-agent-linux -s

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
