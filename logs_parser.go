package main

import (
	"strconv"
	"strings"
)

func (e *PostfixExporter) cleanupParser(matches []string) {
	if strings.Contains(matches[2], ": message-id=<") {
		e.cleanupProcesses.Inc()
	} else if strings.Contains(matches[2], ": reject: ") {
		e.cleanupRejects.Inc()
	} else {
		e.unsupportedLogEntries.WithLabelValues(matches[1]).Inc()
	}
}

func (e *PostfixExporter) lmtpParser(matches []string) {
	if lmtpMatches := lmtpPipeSMTPLine.FindStringSubmatch(matches[2]); lmtpMatches != nil {
		pdelay, err := strconv.ParseFloat(lmtpMatches[2], 64)
		parseError("Couldn't convert LMTP pdelay: %v", err)
		e.lmtpDelays.WithLabelValues("before_queue_manager").Observe(pdelay)
		adelay, err := strconv.ParseFloat(lmtpMatches[3], 64)
		parseError("Couldn't convert LMTP adelay: %v", err)
		e.lmtpDelays.WithLabelValues("queue_manager").Observe(adelay)
		sdelay, err := strconv.ParseFloat(lmtpMatches[4], 64)
		parseError("Couldn't convert LMTP adelay: %v", err)
		e.lmtpDelays.WithLabelValues("connection_setup").Observe(sdelay)
		xdelay, err := strconv.ParseFloat(lmtpMatches[5], 64)
		parseError("Couldn't convert LMTP xdelay: %v", err)
		e.lmtpDelays.WithLabelValues("transmission").Observe(xdelay)
	} else {
		e.unsupportedLogEntries.WithLabelValues(matches[1]).Inc()
	}
}

func (e *PostfixExporter) pipeParser(matches []string) {
	if pipeMatches := lmtpPipeSMTPLine.FindStringSubmatch(matches[2]); pipeMatches != nil {
		pdelay, err := strconv.ParseFloat(pipeMatches[2], 64)
		parseError("Couldn't convert PIPE pdelay: %v", err)
		e.pipeDelays.WithLabelValues(pipeMatches[1], "before_queue_manager").Observe(pdelay)
		adelay, err := strconv.ParseFloat(pipeMatches[3], 64)
		parseError("Couldn't convert PIPE adelay: %v", err)
		e.pipeDelays.WithLabelValues(pipeMatches[1], "queue_manager").Observe(adelay)
		sdelay, err := strconv.ParseFloat(pipeMatches[4], 64)
		parseError("Couldn't convert PIPE sdelay: %v", err)
		e.pipeDelays.WithLabelValues(pipeMatches[1], "connection_setup").Observe(sdelay)
		xdelay, err := strconv.ParseFloat(pipeMatches[5], 64)
		parseError("Couldn't convert PIPE xdelay: %v", err)
		e.pipeDelays.WithLabelValues(pipeMatches[1], "transmission").Observe(xdelay)
	} else {
		e.unsupportedLogEntries.WithLabelValues(matches[1]).Inc()
	}
}

func (e *PostfixExporter) qmgrParser(matches []string) {
	if qmgrInsertMatches := qmgrInsertLine.FindStringSubmatch(matches[2]); qmgrInsertMatches != nil {
		size, err := strconv.ParseFloat(qmgrInsertMatches[1], 64)
		parseError("Couldn't convert QMGR size: %v", err)
		e.qmgrInsertsSize.Observe(size)
		nrcpt, err := strconv.ParseFloat(qmgrInsertMatches[2], 64)
		parseError("Couldn't convert QMGR nrcpt: %v", err)
		e.qmgrInsertsNrcpt.Observe(nrcpt)
	} else if strings.HasSuffix(matches[2], ": removed") {
		e.qmgrRemoves.Inc()
	} else {
		e.unsupportedLogEntries.WithLabelValues(matches[1]).Inc()
	}
}

func (e *PostfixExporter) smtpParser(matches []string) {
	if smtpMatches := lmtpPipeSMTPLine.FindStringSubmatch(matches[2]); smtpMatches != nil {
		pdelay, err := strconv.ParseFloat(smtpMatches[2], 64)
		parseError("Couldn't convert SMTP pdelay: %v", err)
		e.smtpDelays.WithLabelValues("before_queue_manager").Observe(pdelay)
		adelay, err := strconv.ParseFloat(smtpMatches[3], 64)
		parseError("Couldn't convert SMTP adelay: %v", err)
		e.smtpDelays.WithLabelValues("queue_manager").Observe(adelay)
		sdelay, err := strconv.ParseFloat(smtpMatches[4], 64)
		parseError("Couldn't convert SMTP sdelay: %v", err)
		e.smtpDelays.WithLabelValues("connection_setup").Observe(sdelay)
		xdelay, err := strconv.ParseFloat(smtpMatches[5], 64)
		parseError("Couldn't convert SMTP xdelay: %v", err)
		e.smtpDelays.WithLabelValues("transmission").Observe(xdelay)
	} else if smtpTLSMatches := smtpTLSLine.FindStringSubmatch(matches[2]); smtpTLSMatches != nil {
		e.smtpTLSConnects.WithLabelValues(smtpTLSMatches[1:]...).Inc()
	} else {
		e.unsupportedLogEntries.WithLabelValues(matches[1]).Inc()
	}
}

func (e *PostfixExporter) smtpdParser(matches []string) {
	if strings.HasPrefix(matches[2], "connect from ") {
		e.smtpdConnects.Inc()
	} else if strings.HasPrefix(matches[2], "disconnect from ") {
		e.smtpdDisconnects.Inc()
	} else if smtpdFCrDNSErrorsLine.MatchString(matches[2]) {
		e.smtpdFCrDNSErrors.Inc()
	} else if smtpdLostConnectionMatches := smtpdLostConnectionLine.FindStringSubmatch(matches[2]); smtpdLostConnectionMatches != nil {
		e.smtpdLostConnections.WithLabelValues(smtpdLostConnectionMatches[1]).Inc()
	} else if smtpdProcessesSASLMatches := smtpdProcessesSASLLine.FindStringSubmatch(matches[2]); smtpdProcessesSASLMatches != nil {
		e.smtpdProcesses.WithLabelValues(smtpdProcessesSASLMatches[1]).Inc()
	} else if strings.Contains(matches[2], ": client=") {
		e.smtpdProcesses.WithLabelValues("").Inc()
	} else if smtpdRejectsMatches := smtpdRejectsLine.FindStringSubmatch(matches[2]); smtpdRejectsMatches != nil {
		e.smtpdRejects.WithLabelValues(smtpdRejectsMatches[1]).Inc()
	} else if smtpdSASLAuthenticationFailuresLine.MatchString(matches[2]) {
		e.smtpdSASLAuthenticationFailures.Inc()
	} else if smtpdTLSMatches := smtpdTLSLine.FindStringSubmatch(matches[2]); smtpdTLSMatches != nil {
		e.smtpdTLSConnects.WithLabelValues(smtpdTLSMatches[1:]...).Inc()
	} else {
		e.unsupportedLogEntries.WithLabelValues(matches[1]).Inc()
	}
}
