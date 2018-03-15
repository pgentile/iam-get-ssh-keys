package main

/*
#include <stdarg.h>
#include <stdlib.h>
#include <syslog.h>

void go_syslog(int priority, _GoString_ message) {
    syslog(priority, "%s", _GoStringPtr(message));
}

*/
import "C"

import "fmt"

// OpenSyslog opens the syslog
func OpenSyslog() {
	C.openlog(C.CString("toto"), 0, C.LOG_LOCAL0)
}

// SyslogInfo generates an info message
func SyslogInfo(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	C.go_syslog(C.LOG_INFO, message)
}

// SyslogNotice generates an info message
func SyslogNotice(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	C.go_syslog(C.LOG_NOTICE, message)
}

// SyslogError generates an error message
func SyslogError(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	C.go_syslog(C.LOG_ERR, message)
}

// CloseSyslog closes the syslog
func CloseSyslog() {
	C.closelog()
}
