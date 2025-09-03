#!/bin/bash

OUTPUT="README.md"
OS=(win win64 mac mac_arm64 linux android webview ios lacros)

truncate -s 0 $OUTPUT
for i in "${OS[@]}"; do
  echo "## $i" >> $OUTPUT
  echo "\`\`\`" >> $OUTPUT
  echo $(cat $i) >> $OUTPUT
  echo "\`\`\`" >> $OUTPUT
done
