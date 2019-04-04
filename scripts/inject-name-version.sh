#!/bin/bash

PACKAGE="github.com/ncr-devops-platform/nagiosfoundation/cmd/initcmd"

echo "-ldflags"
echo -n "-X $PACKAGE.cmdName=$PRODUCT "
echo "-X $PACKAGE.cmdVersion=$VERSION"

