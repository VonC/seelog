// Copyright (c) 2012 - Cloud Instruments Co. Ltd.
// 
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met: 
// 
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer. 
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution. 
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package seelog

import (
	"errors"
)

type loggerTypeFromString uint8

const (
	syncloggerTypeFromString = iota
	asyncLooploggerTypeFromString
	asyncTimerloggerTypeFromString
	adaptiveLoggerTypeFromString
	DefaultloggerTypeFromString = asyncLooploggerTypeFromString
)

const (
	syncloggerTypeFromStringStr = "sync"
	AsyncloggerTypeFromStringStr = "asyncloop"
	asyncTimerloggerTypeFromStringStr = "asynctimer"
	adaptiveLoggerTypeFromStringStr = "adaptive"
)

// asyncTimerLoggerData represents specific data for async timer logger
type asyncTimerLoggerData struct {
	AsyncInterval uint32
}

// adaptiveLoggerData represents specific data for adaptive timer logger
type adaptiveLoggerData struct {
	MinInterval 		uint32
	MaxInterval 		uint32
	CriticalMsgCount 	uint32
}

var loggerTypeToStringRepresentations = map[loggerTypeFromString]string{
	syncloggerTypeFromString:    		syncloggerTypeFromStringStr,
	asyncLooploggerTypeFromString:    AsyncloggerTypeFromStringStr,
	asyncTimerloggerTypeFromString:   asyncTimerloggerTypeFromStringStr,
	adaptiveLoggerTypeFromString:	  adaptiveLoggerTypeFromStringStr,
}

// loggerTypeFromStringFromString parses a string and returns a corresponding logger type, if sucessfull. 
func loggerTypeFromStringFromString(logTypeString string) (level loggerTypeFromString, found bool) {
	for logType, logTypeStr := range loggerTypeToStringRepresentations {
		if logTypeStr == logTypeString {
			return logType, true
		}
	}

	return 0, false
}

// logConfig stores logging configuration. Contains messages dispatcher, allowed log level rules 
// (general constraints and exceptions), and messages formats (used by nodes of dispatcher tree)
type logConfig struct {
	Constraints    logLevelConstraints      // General log level rules (>min and <max, or set of allowed levels)
	Exceptions     []*logLevelException     // Exceptions to general rules for specific files or funcs
	RootDispatcher dispatcherInterface // Root of output tree
	LogType        loggerTypeFromString
	LoggerData     interface{}
}

func newConfig(
	constraints logLevelConstraints, 
	exceptions []*logLevelException, 
	rootDispatcher dispatcherInterface,
	logType loggerTypeFromString,
	logData interface{}) (*logConfig, error) {
	if constraints == nil {
		return nil, errors.New("Constraints can not be nil")
	}
	if rootDispatcher == nil {
		return nil, errors.New("RootDispatcher can not be nil")
	}
	
	config := new(logConfig)
	config.Constraints = constraints
	config.Exceptions = exceptions
	config.RootDispatcher = rootDispatcher
	config.LogType = logType
	config.LoggerData = logData
	
	return config, nil
}

// IsAllowed returns true if logging with specified log level is allowed in current context.
// If any of exception patterns match current context, then exception constraints are applied. Otherwise,
// the general constraints are used.
func (config *logConfig) IsAllowed(level LogLevel, context logContextInterface) bool {
	allowed := config.Constraints.IsAllowed(level) // General rule

	// Exceptions:
	if context.IsValid() {
		for _, exception := range config.Exceptions {
			if exception.MatchesContext(context) {
				return exception.IsAllowed(level)
			}
		}
	}

	return allowed
}
