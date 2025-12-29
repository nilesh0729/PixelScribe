#!/bin/bash
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
cd $(dirname "$0")
echo "Starting dev server..."
# Run with --host to ensure network binding if needed, though localhost is fine for test
# processgroup handling might be needed to kill it later, but run_command kills the process group
timeout 10s npm run dev
