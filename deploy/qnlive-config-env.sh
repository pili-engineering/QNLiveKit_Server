#!/bin/bash

file_env() {
	local var="$1"
	local s='$'"{$var}"
	local val="${!var}"

	sed -i "s/$s/$val/g" /etc/qnlive.yaml
}

file_env QINIU_ACCESS_KEY
file_env QINIU_SECRET_KEY
file_env IM_ADMIN_TOKEN
file_env PILI_PUB_KEY