#!/bin/bash
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
cd $(dirname "$0")
echo "Running tsc..."
npm run build > build_log.txt 2>&1
cat build_log.txt
