# 设置GitHub使用ssh拉取
git config --global url."git@github.com:".insteadof "https://github.com/"

# check QBOXROOT 是否为空，空则设置上级目录为QBOXROOT
if [ ! $QBOXROOT ]; then
	QBOXROOT=$(cd ../; pwd)
	export QBOXROOT
fi

export GO111MODULE="on"

