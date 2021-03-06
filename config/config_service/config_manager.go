/*
 *  Licensed to Wikifeat under one or more contributor license agreements.
 *  See the LICENSE.txt file distributed with this work for additional information
 *  regarding copyright ownership.
 *
 *  Redistribution and use in source and binary forms, with or without
 *  modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *  this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright
 *  notice, this list of conditions and the following disclaimer in the
 *  documentation and/or other materials provided with the distribution.
 *  * Neither the name of Wikifeat nor the names of its contributors may be used
 *  to endorse or promote products derived from this software without
 *  specific prior written permission.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */

package config_service

import (
	"errors"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/registry"
	"strconv"
	"strings"
)

type ConfigManager struct{}

var configLocation = registry.EtcdPrefix + "/config/"

//Get a config parameter
func (cm *ConfigManager) getConfigParam(section string,
	paramName string) (string, error) {
	serviceSection, err := config.ServiceSectionFromString(section)
	if err != nil {
		return "", err
	}
	switch serviceSection {
	case config.AuthService:
		return cm.getAuthParam(strings.ToLower(paramName))
	default:
		return "", errors.New("Invalid config section requested")
	}

}

//Get an Auth config parameter
func (cm *ConfigManager) getAuthParam(paramName string) (string, error) {
	switch paramName {
	case "allowguestaccess":
		return strconv.FormatBool(config.Auth.AllowGuest), nil
	case "allownewuserregistration":
		return strconv.FormatBool(config.Auth.AllowNewUserRegistration), nil
	default:
		return "", errors.New("Invalid auth config parameter requested")
	}
}
